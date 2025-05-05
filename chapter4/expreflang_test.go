package chapter4

import (
	"testing"

	epl "github.com/panyam/eplgo" // Base definitions
	// AST nodes and base evaluator helpers
	"github.com/stretchr/testify/assert"
)

// Helper to create the evaluator and set standard operators
func NewTestExpRefLangEval() Evaluator {
	// Use chapter3's setOpFuncs for convenience
	return SetOpFuncs(NewExpRefLangEval())
}

// Test runner - adapted from RunTest
// It needs to handle the fact that evaluation can return non-LitExpr values
// (like *epl.Ref[any] for newref). Assertions must check types.
func RunExpRefTest(t *testing.T, e Evaluator, tc *TestCase, extraenv map[string]any) {
	env := epl.NewEnv[any](nil)
	for k, v := range extraenv {
		// Wrap initial env values in refs, as that's how variables are stored now
		// Correction: Env.Set already wraps in Ref. Env values should be the actual values.
		env.Set(k, v) // Set already handles wrapping in Ref
	}

	// log.Printf("======= Running TestCase: %s =======", tc.Name)
	// log.Println("Expr:", tc.Expr.Repr())
	// if extraenv != nil {
	// 	log.Println("Initial Env:", extraenv)
	// }

	value, err := e.Eval(tc.Expr, env) // Eval returns any

	// Check for unexpected errors first
	// TODO: Modify tests later to expect errors when needed
	assert.NoError(t, err, "Test %s Failed - Unexpected error", tc.Name)

	// Assert based on the *expected* type and value
	// For Chapter 4, expected might be a primitive OR could indicate structure (e.g., expect a ref)
	switch expected := tc.Expected.(type) {
	case int:
		finalValue, ok := value.(*LitExpr)
		if !ok {
			t.Fatalf("Test %s: Expected LitExpr(int) but got %T (%#v)", tc.Name, value, value)
		}
		assert.Equal(t, expected, finalValue.Value, "Test %s Failed - Int Value", tc.Name)
	case bool:
		finalValue, ok := value.(*LitExpr)
		if !ok {
			t.Fatalf("Test %s: Expected LitExpr(bool) but got %T (%#v)", tc.Name, value, value)
		}
		assert.Equal(t, expected, finalValue.Value, "Test %s Failed - Bool Value", tc.Name)
	// --- Add specific checks for Chapter 4 ---
	case *epl.Ref[any]:
		// Used when the test expects the *result* to be a reference
		actualRef, ok := value.(*epl.Ref[any])
		if !ok {
			t.Fatalf("Test %s: Expected *epl.Ref[any] but got %T (%#v)", tc.Name, value, value)
		}
		// Compare the *contents* of the reference
		// This requires the expected ref to also contain a comparable value
		assert.ObjectsAreEqualValues(expected.Value, actualRef.Value) // Compare values inside refs

	// Use a sentinel type/value if we just need to check if it's *any* reference
	case struct{ typeIsRef bool }: // Sentinel struct
		if expected.typeIsRef {
			_, ok := value.(*epl.Ref[any])
			assert.True(t, ok, "Test %s: Expected a reference (*epl.Ref[any]) but got %T (%#v)", tc.Name, value, value)
		} else {
			// Handle case where sentinel is used for non-ref expectation (if needed)
			t.Fatalf("Test %s: Invalid use of Ref sentinel expectation", tc.Name)
		}

	default:
		// Fallback: Use standard testify equality which might handle pointers well enough sometimes
		assert.Equal(t, tc.Expected, value, "Test %s Failed (Default Comparison)", tc.Name)
		// Or add more specific type checks as needed
	}
	// log.Printf("Test %s Passed. Found: %v (%T)\n", tc.Name, value, value)
}

// --- Test Cases ---

// Sentinel value to indicate we expect the result to be a reference, without caring about its content
var ExpectRef = struct{ typeIsRef bool }{true}

func TestNewRef(t *testing.T) {
	evaluator := NewTestExpRefLangEval()
	// newref(10)
	expr := NewRef(10)
	tc := TestCase{Name: "newref_int", Expected: ExpectRef, Expr: expr}
	// We expect the *result* of newref to be a reference
	RunExpRefTest(t, evaluator, &tc, nil)

	// Check the value inside (requires deref)
	// deref(newref(20)) == 20
	exprDeref := DeRef(NewRef(20))
	tcDeref := TestCase{Name: "deref_newref_int", Expected: 20, Expr: exprDeref}
	RunExpRefTest(t, evaluator, &tcDeref, nil)
}

func TestDeref(t *testing.T) {
	evaluator := NewTestExpRefLangEval()
	// let x = newref(5) in deref(x) == 5
	expr := Let(ExprDict("x", NewRef(5)), DeRef(Var("x")))
	tc := TestCase{Name: "deref_var", Expected: 5, Expr: expr}
	RunExpRefTest(t, evaluator, &tc, nil)
}

func TestSetRef(t *testing.T) {
	evaluator := NewTestExpRefLangEval()
	// let x = newref(5) in begin setref(x, 15); deref(x) end == 15
	expr := Let(
		ExprDict("x", NewRef(5)), // x bound to Ref(5)
		Begin( // Sequence
			SetRef(Var("x"), 15), // Set Ref(5) to Ref(15), returns 15 (discarded)
			DeRef(Var("x")),      // Deref Ref(15), returns 15
		),
	)
	tc := TestCase{Name: "setref_var", Expected: 15, Expr: expr}
	RunExpRefTest(t, evaluator, &tc, nil)

	// Check return value of setref itself
	// setref(newref(1), 99) == 99
	exprRet := SetRef(NewRef(1), 99)
	tcRet := TestCase{Name: "setref_return_val", Expected: 99, Expr: exprRet}
	RunExpRefTest(t, evaluator, &tcRet, nil)
}

func TestBlock(t *testing.T) {
	evaluator := NewTestExpRefLangEval()
	// begin 1; 2; 3 end == 3
	expr := Begin(1, 2, 3)
	tc := TestCase{Name: "block_simple", Expected: 3, Expr: expr}
	RunExpRefTest(t, evaluator, &tc, nil)

	// Test block with state changes
	// let x = newref(0) in begin setref(x, 1); setref(x, deref(x) + 10); deref(x) end == 11
	exprState := Let(
		ExprDict("x", NewRef(0)), // x = Ref(0)
		Begin(
			SetRef(Var("x"), 1), // x = Ref(1), returns 1 (discarded)
			SetRef(Var("x"), // Set x = ...
				Op("+", // Value is...
					DeRef(Var("x")), // deref(x) = 1
					10,
				), // Op result is Lit(11)
			), // x = Ref(11), setref returns 11 (discarded)
			DeRef(Var("x")), // deref(x) = 11 (final result)
		),
	)
	tcState := TestCase{Name: "block_state", Expected: 11, Expr: exprState}
	RunExpRefTest(t, evaluator, &tcState, nil)

	// Empty block - Python returned Lit(0), our eval mimics this
	exprEmpty := Begin()
	tcEmpty := TestCase{Name: "block_empty", Expected: 0, Expr: exprEmpty}
	RunExpRefTest(t, evaluator, &tcEmpty, nil)
}

// Porting Python Tests from tests/chapter4/cases.py (exprefs)

func TestExpref_OddEven(t *testing.T) {
	/* Python:
	   "oddeven": ("""
	       let x = newref(0) in
	           letrec
	               even(dummy)
	                   = if isz(deref(x))
	                     then 1
	                     else begin
	                       setref(x, -(deref(x), 1));
	                       (odd 888)
	                     end
	               odd(dummy)
	                   = if isz(deref(x))
	                     then 0
	                     else begin
	                       setref(x, -(deref(x), 1));
	                       (even 888)
	                     end
	           in begin setref(x, 13) ; (odd 888) end
	           """, 1),
	*/
	evaluator := NewTestExpRefLangEval()
	expr := Let(
		ExprDict("x", NewRef(0)), // x = Ref(0)
		LetRec( // Mutually recursive procs
			ProcMap(
				"even", Proc([]string{"dummy"},
					If(IsZero(DeRef(Var("x"))), // if isz(deref(x))
						Lit(1), // then 1
						Begin( // else begin
							SetRef(Var("x"), Op("-", DeRef(Var("x")), 1)), // setref(x, deref(x)-1)
							Call(Var("odd"), 888),                         // call odd (dummy arg)
						), // end
					),
				), // end even proc
				"odd", Proc([]string{"dummy"},
					If(IsZero(DeRef(Var("x"))), // if isz(deref(x))
						Lit(0), // then 0
						Begin( // else begin
							SetRef(Var("x"), Op("-", DeRef(Var("x")), 1)), // setref(x, deref(x)-1)
							Call(Var("even"), 888),                        // call even (dummy arg)
						), // end
					),
				), // end odd proc
			), // End ProcMap
			// Body of LetRec
			Begin(
				SetRef(Var("x"), 13),  // setref(x, 13)
				Call(Var("odd"), 888), // call odd (initial call)
			),
		), // End LetRec
	) // End Let

	tc := TestCase{Name: "expref_oddeven", Expected: 1, Expr: expr}
	RunExpRefTest(t, evaluator, &tc, nil)
}

func TestExpref_Counter(t *testing.T) {
	/* Python:
	   "counter": ("""
	       let g = let counter = newref(0)
	               in proc(dummy)
	                   begin
	                       setref(counter, -(deref(counter), -1)) ;
	                       deref(counter)
	                   end
	       in let a = (g 11)
	           in let b = (g 11)
	               in -(a,b)
	   """, -1)
	*/
	evaluator := NewTestExpRefLangEval()
	expr := Let(
		ExprDict( // let g = ...
			"g", Let( // Inner let for counter
				ExprDict("counter", NewRef(0)), // counter = Ref(0)
				// Proc body is the value bound to 'g'
				Proc([]string{"dummy"},
					Begin(
						// Equivalent to counter = counter + 1
						SetRef(Var("counter"), Op("-", DeRef(Var("counter")), -1)),
						DeRef(Var("counter")), // Return updated counter value
					),
				),
			), // End inner let
		), // End binding for g
		// Body of outer let
		Let(ExprDict("a", Call(Var("g"), 11)), // let a = (g 11)
			Let(ExprDict("b", Call(Var("g"), 11)), // let b = (g 11)
				Op("-", Var("a"), Var("b")), // Body: -(a, b)
			),
		),
	) // End outer let

	tc := TestCase{Name: "expref_counter", Expected: -1, Expr: expr}
	RunExpRefTest(t, evaluator, &tc, nil)
}

// TODO: Add tests for Printable and Eq for the new Chapter 4 types
func TestChapter4ExprEqPrintable(t *testing.T) {
	// Test Eq
	assert.True(t, ExprEq(NewRef(1), NewRef(1)))
	assert.False(t, ExprEq(NewRef(1), NewRef(2)))
	assert.True(t, ExprEq(DeRef("x"), DeRef("x")))
	assert.False(t, ExprEq(DeRef("x"), DeRef("y")))
	assert.True(t, ExprEq(SetRef("x", 1), SetRef("x", 1)))
	assert.False(t, ExprEq(SetRef("x", 1), SetRef("y", 1)))
	assert.False(t, ExprEq(SetRef("x", 1), SetRef("x", 2)))
	assert.True(t, ExprEq(Begin(1, 2), Begin(1, 2)))
	assert.False(t, ExprEq(Begin(1, 2), Begin(1, 3)))
	assert.False(t, ExprEq(Begin(1, 2), Begin(1)))
}

// Helper for simplified printable check (add to common.go?)
// func (p *Printable) Leaf string { return p.Leaf }
