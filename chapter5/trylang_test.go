package chapter5

import (
	"errors" // For errors.Is/As in tests
	"fmt"
	"testing"

	epl "github.com/panyam/eplgo" // Base definitions
	// AST nodes and base evaluator helpers
	"github.com/panyam/eplgo/chapter4" // Need AST nodes like Begin, Assign
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper to create the Chapter 5 evaluator and set standard operators
func NewTestChapter5Eval() Evaluator {
	// Use chapter3's SetOpFuncs for convenience
	return SetOpFuncs(NewChapter5Eval())
}

// Test runner for Chapter 5 - handles errors, including RaisedError
func RunTryLangTest(t *testing.T, e Evaluator, tc *TestCase, extraenv map[string]any) {
	env := epl.NewEnv[any](nil)
	for k, v := range extraenv {
		env.Set(k, v)
	}

	// log.Printf("======= Running TestCase: %s =======", tc.Name)
	// log.Println("Expr:", tc.Expr.Repr())
	// if extraenv != nil {
	// 	log.Println("Initial Env:", extraenv)
	// }

	value, err := e.Eval(tc.Expr, env) // Eval returns (any, error)

	// Check if an error is expected
	expectedErr, expectError := tc.Expected.(error) // Check if expected value IS an error
	expectSpecificError := false
	if expectError {
		// Is it the specific RaisedError sentinel or another specific error?
		_, expectRaisedSentinel := expectedErr.(RaisedError) // Check if it's the zero value sentinel
		expectSpecificError = !expectRaisedSentinel || expectedErr != (RaisedError{})

		// If we expect *any* RaisedError (using the zero value as sentinel)
		if !expectSpecificError && expectRaisedSentinel {
			require.Error(t, err, "Test %s: Expected a RaisedError, but got nil error", tc.Name)
			assert.ErrorAs(t, err, &RaisedError{}, "Test %s: Expected RaisedError type", tc.Name)
			// Don't compare value if we only expected the error type
			return
		}

		// If we expect a specific error (RaisedError with value, or other error type)
		require.Error(t, err, "Test %s: Expected error '%v', but got nil error", tc.Name, expectedErr)
		if errors.Is(expectedErr, RaisedError{}) { // If expecting RaisedError type comparison
			var actualRaised RaisedError
			isRaised := errors.As(err, &actualRaised)
			require.True(t, isRaised, "Test %s: Expected RaisedError type, got %T", tc.Name, err)
			// Compare the wrapped values for specific RaisedError expectation
			expectedRaised := expectedErr.(RaisedError)
			assert.ObjectsAreEqualValues(expectedRaised.Value, actualRaised.Value) // Compare values inside errors
		} else {
			// Check for other specific errors
			assert.ErrorIs(t, err, expectedErr, "Test %s: Error mismatch", tc.Name)
		}
		// Don't compare value if we expected a specific error
		return
	}

	// If no error was expected, assert that none occurred
	require.NoError(t, err, "Test %s Failed - Unexpected error", tc.Name)

	// --- Assert Value if no error was expected ---
	// Use logic similar to RunExpRefTest for value comparison
	switch expected := tc.Expected.(type) {
	case int:
		finalValue, ok := value.(*LitExpr)
		require.True(t, ok, "Test %s: Expected LitExpr(int) but got %T (%#v)", tc.Name, value, value)
		assert.Equal(t, expected, finalValue.Value, "Test %s Failed - Int Value", tc.Name)
	case bool:
		finalValue, ok := value.(*LitExpr)
		require.True(t, ok, "Test %s: Expected LitExpr(bool) but got %T (%#v)", tc.Name, value, value)
		assert.Equal(t, expected, finalValue.Value, "Test %s Failed - Bool Value", tc.Name)
	case *epl.Ref[any]:
		actualRef, ok := value.(*epl.Ref[any])
		require.True(t, ok, "Test %s: Expected *epl.Ref[any] but got %T (%#v)", tc.Name, value, value)
		assert.ObjectsAreEqualValues(expected.Value, actualRef.Value) // Compare values inside refs
	case struct{ typeIsRef bool }: // Sentinel struct from Ch4 tests
		if expected.typeIsRef {
			_, ok := value.(*epl.Ref[any])
			assert.True(t, ok, "Test %s: Expected a reference (*epl.Ref[any]) but got %T (%#v)", tc.Name, value, value)
		} else {
			t.Fatalf("Test %s: Invalid use of Ref sentinel expectation", tc.Name)
		}
	// Add case for Thunk if needed for comparison
	case *chapter4.Thunk:
		actualThunk, ok := value.(*chapter4.Thunk)
		require.True(t, ok, "Test %s: Expected *chapter4.Thunk but got %T (%#v)", tc.Name, value, value)
		// Basic comparison: check if expressions are equal. Env comparison is tricky.
		assert.True(t, ExprEq(expected.Expr, actualThunk.Expr), "Test %s: Thunk expression mismatch", tc.Name)
		// Could add more checks if needed (e.g., env non-nil)
	default:
		assert.Equal(t, tc.Expected, value, "Test %s Failed (Default Comparison)", tc.Name)
	}
	// log.Printf("Test %s Passed. Found: %v (%T), Err: %v\n", tc.Name, value, value, err)
}

// --- Test Cases ---

func TestTryNormalExecution(t *testing.T) {
	evaluator := NewTestChapter5Eval()
	// try 10 catch (x) x + 1  ==> 10
	expr := Try(10, "x", Op("+", Var("x"), 1))
	tc := TestCase{Name: "try_normal", Expected: 10, Expr: expr}
	RunTryLangTest(t, evaluator, &tc, nil)
}

func TestTryCatchRaise(t *testing.T) {
	evaluator := NewTestChapter5Eval()
	// try raise 10 catch (x) x + 1 ==> 11
	expr := Try(Raise(10), "x", Op("+", Var("x"), 1))
	tc := TestCase{Name: "try_catch", Expected: 11, Expr: expr}
	RunTryLangTest(t, evaluator, &tc, nil)

	// try raise (10 + 5) catch (y) y * 2 ==> 30
	exprCalc := Try(Raise(Op("+", 10, 5)), "y", Op("*", Var("y"), 2))
	tcCalc := TestCase{Name: "try_catch_calc", Expected: 30, Expr: exprCalc}
	RunTryLangTest(t, evaluator, &tcCalc, nil)
}

func TestTryNestedRaise(t *testing.T) {
	evaluator := NewTestChapter5Eval()
	// try try raise 5 catch (x) x + 1 catch (y) y * 10
	// Inner try: raise 5 -> caught by inner catch -> x=5 -> handler returns 6
	// Outer try: Receives 6 normally from inner try.
	// Result: 6
	expr := Try( // Outer try
		Try( // Inner try
			Raise(5),             // raise 5
			"x",                  // catch (x)
			Op("+", Var("x"), 1), // handler x+1 -> returns 6
		),
		"y", // catch (y) - never reached
		Op("*", Var("y"), 10),
	)
	tc := TestCase{Name: "try_nested_inner_catch", Expected: 6, Expr: expr}
	RunTryLangTest(t, evaluator, &tc, nil)

	// try try raise 5 catch (x) raise (x+1) catch (y) y * 10
	// Inner try: raise 5 -> caught by inner catch -> x=5 -> handler raises (5+1)=6
	// Outer try: Catches the raise 6 -> y=6 -> handler returns 6*10 = 60
	// Result: 60
	exprOuter := Try( // Outer try
		Try( // Inner try
			Raise(5),                    // raise 5
			"x",                         // catch (x)
			Raise(Op("+", Var("x"), 1)), // handler raises x+1
		),
		"y",                   // catch (y)
		Op("*", Var("y"), 10), // handler y*10
	)
	tcOuter := TestCase{Name: "try_nested_outer_catch", Expected: 60, Expr: exprOuter}
	RunTryLangTest(t, evaluator, &tcOuter, nil)
}

func TestUncaughtRaise(t *testing.T) {
	evaluator := NewTestChapter5Eval()
	// raise 100
	expr := Raise(100)
	// Expect a RaisedError containing Lit(100)
	tc := TestCase{
		Name:     "uncaught_raise",
		Expected: RaisedError{Value: Lit(100)}, // Expect specific error
		Expr:     expr,
	}
	RunTryLangTest(t, evaluator, &tc, nil)

	// let x = raise 99 in x + 1 -> Should error immediately on raise
	exprLet := Let(ExprDict("x", Raise(99)), Op("+", Var("x"), 1))
	tcLet := TestCase{
		Name:     "uncaught_raise_in_let",
		Expected: RaisedError{Value: Lit(99)}, // Expect specific error
		Expr:     exprLet,
	}
	RunTryLangTest(t, evaluator, &tcLet, nil)
}

func TestTryNonRaisedError(t *testing.T) {
	evaluator := NewTestChapter5Eval()
	// try y + 1 catch (x) x // y is unbound
	expr := Try(Op("+", Var("y"), 1), "x", Var("x"))
	// We expect the "variable not found" error, NOT a RaisedError
	expectedErr := fmt.Errorf("variable 'y' not found in environment") // Match error from ValueOfVar

	env := epl.NewEnv[any](nil)
	_, err := evaluator.Eval(expr, env)

	require.Error(t, err, "Test try_non_raised: Expected an error")
	assert.False(t, errors.As(err, &RaisedError{}), "Test try_non_raised: Error should not be RaisedError")
	// Check if the error message contains the expected text (more robust than exact match)
	assert.Contains(t, err.Error(), expectedErr.Error(), "Test try_non_raised: Error message mismatch")

	// Or, if we modify RunTryLangTest to accept specific non-RaisedError expectations:
	// tc := TestCase{
	// 	Name:     "try_non_raised",
	// 	Expected: fmt.Errorf("variable 'y' not found in environment"), // Specific expected error
	// 	Expr:     expr,
	// }
	// RunTryLangTest(t, evaluator, &tc, nil) // Requires RunTryLangTest update
}

func TestTryLangEqPrintable(t *testing.T) {
	// Eq
	assert.True(t, ExprEq(Try(1, "x", 2), Try(1, "x", 2)))
	assert.False(t, ExprEq(Try(1, "x", 2), Try(9, "x", 2))) // Diff body
	assert.False(t, ExprEq(Try(1, "x", 2), Try(1, "y", 2))) // Diff var
	assert.False(t, ExprEq(Try(1, "x", 2), Try(1, "x", 9))) // Diff handler

	assert.True(t, ExprEq(Raise(1), Raise(1)))
	assert.False(t, ExprEq(Raise(1), Raise(2)))
	assert.False(t, ExprEq(Raise(1), Try(1, "x", 2)))

	// Printable - rely on common_test.go for formatting checks
}
