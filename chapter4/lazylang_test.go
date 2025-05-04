package chapter4

import (
	"testing"

	"github.com/panyam/eplgo/chapter3" // AST nodes and base evaluator helpers
	"github.com/stretchr/testify/assert"
	// No specific epl import needed unless using ExpectRef etc.
)

// Helper to create the evaluator and set standard operators
func NewTestLazyLangEval() chapter3.Evaluator {
	// Use chapter3's setOpFuncs for convenience
	return chapter3.SetOpFuncs(NewLazyLangEval())
}

// Use the same test runner as expreflang/impreflang
var RunLazyTest = RunExpRefTest // Alias for clarity

// --- Test Cases ---

// Porting Python Tests from tests/chapter4/cases.py (lazy)

func TestLazy_InfiniteLoopAvoidance(t *testing.T) {
	/* Python:
	   "infinite": ("""
	       letrec infinite-loop(x) = lazy ( infinite-loop(x) ) // Simplified Python slightly
	           in let f = proc(z) 11
	               in (f (infinite-loop 0)) // Pass the lazy computation to f
	   """, 11)
	   The key is that `infinite-loop 0` evaluates to a Thunk immediately.
	   `f` receives the Thunk but never forces it (since `z` isn't used), so it just returns 11.
	*/
	evaluator := NewTestLazyLangEval()
	expr := chapter3.LetRec(
		chapter3.ProcMap(
			"infinite-loop", chapter3.Proc([]string{"x"},
				// Body creates the lazy expression containing the recursive call
				Lazy(chapter3.Call(chapter3.Var("infinite-loop"), chapter3.Var("x"))),
			),
		),
		// Body of LetRec
		chapter3.Let(
			chapter3.ExprDict("f", chapter3.Proc([]string{"z"}, chapter3.Lit(11))), // f = proc(z) 11
			// Call f, passing the result of (infinite-loop 0), which should be a Thunk
			chapter3.Call(chapter3.Var("f"), chapter3.Call(chapter3.Var("infinite-loop"), 0)),
		),
	)

	tc := chapter3.TestCase{Name: "lazy_infinite", Expected: 11, Expr: expr}
	RunLazyTest(t, evaluator, &tc, nil)
}

func TestLazy_Forcing(t *testing.T) {
	evaluator := NewTestLazyLangEval()
	// let x = lazy (1 + 2) in thunk x == 3
	expr := chapter3.Let(
		chapter3.ExprDict("x", Lazy(chapter3.Op("+", 1, 2))), // x = Ref(Thunk(1+2, env))
		ForceThunk(chapter3.Var("x")),                        // Force evaluation of the thunk stored in x
	)
	tc := chapter3.TestCase{Name: "lazy_forcing", Expected: 3, Expr: expr}
	RunLazyTest(t, evaluator, &tc, nil)

	// Test forcing multiple times (should evaluate only once if memoized - but ours isn't memoized yet)
	// let y = lazy -(10, 1) in begin (thunk y); (thunk y) end == 9
	// Our current implementation re-evaluates each time `thunk` is called.
	exprMultiForce := chapter3.Let(
		chapter3.ExprDict("y", Lazy(chapter3.Op("-", 10, 1))), // y = Ref(Thunk(10-1, env))
		Begin(
			ForceThunk(chapter3.Var("y")), // Eval -> 9 (discarded)
			ForceThunk(chapter3.Var("y")), // Eval -> 9 (returned)
		),
	)
	tcMultiForce := chapter3.TestCase{Name: "lazy_multiple_force", Expected: 9, Expr: exprMultiForce}
	RunLazyTest(t, evaluator, &tcMultiForce, nil)

}

// Test state interaction with lazy
func TestLazy_State(t *testing.T) {
	evaluator := NewTestLazyLangEval()
	// let r = newref(0)
	// in let lz = lazy (setref(r, 10); deref(r) + 5) // Capture r=Ref(0)
	// in begin
	//      setref(r, 100); // Change r *before* forcing
	//      thunk lz        // Force using captured env -> sets Ref(0) to Ref(10), returns 10+5=15
	// end == 15
	expr := chapter3.Let(chapter3.ExprDict("r", NewRef(0)), // r = Ref(0)
		chapter3.Let(chapter3.ExprDict("lz", // lz = Ref(Thunk(...))
			Lazy(
				Begin(
					SetRef(chapter3.Var("r"), 10),                 // uses captured r
					chapter3.Op("+", DeRef(chapter3.Var("r")), 5), // uses captured r
				),
			),
		),
			// Body
			Begin(
				SetRef(chapter3.Var("r"), 100), // Change r *outside* lazy expr to Ref(100)
				ForceThunk(chapter3.Var("lz")), // Force lazy evaluation using captured env
			),
		),
	)
	tc := chapter3.TestCase{Name: "lazy_state_capture", Expected: 15, Expr: expr}
	RunLazyTest(t, evaluator, &tc, nil)

	// Check final value of r after the test above
	// let r = newref(0)
	// in let lz = lazy (setref(r, 10); deref(r) + 5)
	// in begin
	//      setref(r, 100);
	//      thunk lz;
	//      deref(r) // Should be 10, because thunk evaluation modified it
	// end == 10
	exprFinalR := chapter3.Let(chapter3.ExprDict("r", NewRef(0)),
		chapter3.Let(chapter3.ExprDict("lz",
			Lazy(
				Begin(
					SetRef(chapter3.Var("r"), 10),
					chapter3.Op("+", DeRef(chapter3.Var("r")), 5),
				),
			),
		),
			Begin(
				SetRef(chapter3.Var("r"), 100),
				ForceThunk(chapter3.Var("lz")), // Force evaluation, returns 15 (discarded)
				DeRef(chapter3.Var("r")),       // Get current value of r
			),
		),
	)
	tcFinalR := chapter3.TestCase{Name: "lazy_state_final_value", Expected: 10, Expr: exprFinalR}
	RunLazyTest(t, evaluator, &tcFinalR, nil)
}

// Add Eq and Printable tests for LazyExpr, ThunkExpr
func TestLazyExprEqPrintable(t *testing.T) {
	// Eq
	assert.True(t, chapter3.ExprEq(Lazy(1), Lazy(1)))
	assert.False(t, chapter3.ExprEq(Lazy(1), Lazy(2)))
	assert.True(t, chapter3.ExprEq(ForceThunk("x"), ForceThunk("x")))
	assert.False(t, chapter3.ExprEq(ForceThunk("x"), ForceThunk("y")))
	assert.False(t, chapter3.ExprEq(Lazy(1), ForceThunk(1)))

	// Printable - rely on common_test.go for formatting checks
}
