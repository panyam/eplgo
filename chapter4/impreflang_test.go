package chapter4

import (
	"testing"

	// epl "github.com/panyam/eplgo"          // Base definitions (Not directly needed here)
	"github.com/panyam/eplgo/chapter3" // AST nodes and base evaluator helpers
	"github.com/stretchr/testify/assert"
)

// Helper to create the evaluator and set standard operators
func NewTestImpRefLangEval() chapter3.Evaluator {
	// Use chapter3's setOpFuncs for convenience
	return chapter3.SetOpFuncs(NewImpRefLangEval())
}

// Use the same test runner as expreflang, as the assertion logic should be similar
var RunImpRefTest = RunExpRefTest // Alias for clarity

// --- Test Cases ---

func TestAssign(t *testing.T) {
	evaluator := NewTestImpRefLangEval()
	// let x = 10 in begin set x = 20; x end == 20
	// Note: `let x = 10` binds x to Ref(Lit(10)) because Env.Set wraps it.
	expr := chapter3.Let(
		chapter3.ExprDict("x", chapter3.Lit(10)), // x bound to Ref(Lit(10))
		Begin( // Sequence
			Assign("x", 20),   // Set x = 20 -> Modifies Ref to Ref(Lit(20)), returns Lit(20) (discarded)
			chapter3.Var("x"), // Evaluate x -> Looks up Ref(Lit(20)), returns Lit(20)
		),
	)
	tc := chapter3.TestCase{Name: "assign_simple", Expected: 20, Expr: expr}
	RunImpRefTest(t, evaluator, &tc, nil)

	// Check return value of assign itself
	// let x = 0 in set x = 99 == 99
	exprRet := chapter3.Let(
		chapter3.ExprDict("x", chapter3.Lit(0)), // x = Ref(Lit(0))
		Assign("x", 99),                         // Set x = 99, returns Lit(99)
	)
	tcRet := chapter3.TestCase{Name: "assign_return_val", Expected: 99, Expr: exprRet}
	RunImpRefTest(t, evaluator, &tcRet, nil)
}

// Porting Python Tests from tests/chapter4/cases.py (imprefs)

func TestImpref_OddEven(t *testing.T) {
	/*
	   "oddeven": ("""
	       let x = 0 in
	           letrec
	               even(dummy)
	                   = if isz(x)
	                     then 1
	                     else begin
	                       set x = -(x,1) ;
	                       (odd 888)
	                     end
	               odd(dummy)
	                   = if isz(x)
	                     then 0
	                     else begin
	                       set x = -(x,1) ;
	                       (even 888)
	                     end
	           in begin set x = 13 ; (odd 888) end
	       """, 1),
	*/
	evaluator := NewTestImpRefLangEval()
	expr := chapter3.Let(
		chapter3.ExprDict("x", chapter3.Lit(0)), // x = Ref(Lit(0))
		chapter3.LetRec( // Mutually recursive procs
			chapter3.ProcMap(
				"even", chapter3.Proc([]string{"dummy"},
					chapter3.If(chapter3.IsZero(chapter3.Var("x")), // if isz(x)
						chapter3.Lit(1), // then 1
						Begin( // else begin
							Assign("x", chapter3.Op("-", chapter3.Var("x"), 1)), // set x = x - 1
							chapter3.Call(chapter3.Var("odd"), 888),             // call odd
						), // end
					),
				), // end even proc
				"odd", chapter3.Proc([]string{"dummy"},
					chapter3.If(chapter3.IsZero(chapter3.Var("x")), // if isz(x)
						chapter3.Lit(0), // then 0
						Begin( // else begin
							Assign("x", chapter3.Op("-", chapter3.Var("x"), 1)), // set x = x - 1
							chapter3.Call(chapter3.Var("even"), 888),            // call even
						), // end
					),
				), // end odd proc
			), // End ProcMap
			// Body of LetRec
			Begin(
				Assign("x", 13),                         // set x = 13
				chapter3.Call(chapter3.Var("odd"), 888), // call odd (initial call)
			),
		), // End LetRec
	) // End Let

	tc := chapter3.TestCase{Name: "impref_oddeven", Expected: 1, Expr: expr}
	RunImpRefTest(t, evaluator, &tc, nil)
}

func TestImpref_Counter(t *testing.T) {
	/* Python: (Same as before, uses 'set')
	   "counter": ("""
	           let g = let counter = 0
	                   in proc(dummy)
	                       begin
	                           set counter = -(counter, -1) ;
	                           counter
	                       end
	           in let a = (g 11)
	               in let b = (g 11)
	                   in -(a,b)
	       """, -1),
	*/
	evaluator := NewTestImpRefLangEval()
	expr := chapter3.Let(
		chapter3.ExprDict( // let g = ...
			"g", chapter3.Let( // Inner let for counter
				chapter3.ExprDict("counter", chapter3.Lit(0)), // counter = Ref(Lit(0))
				// Proc body is the value bound to 'g'
				chapter3.Proc([]string{"dummy"},
					Begin(
						// Equivalent to counter = counter + 1
						Assign("counter", chapter3.Op("-", chapter3.Var("counter"), -1)),
						chapter3.Var("counter"), // Return updated counter value
					),
				),
			), // End inner let
		), // End binding for g
		// Body of outer let
		chapter3.Let(chapter3.ExprDict("a", chapter3.Call(chapter3.Var("g"), 11)), // let a = (g 11)
			chapter3.Let(chapter3.ExprDict("b", chapter3.Call(chapter3.Var("g"), 11)), // let b = (g 11)
				chapter3.Op("-", chapter3.Var("a"), chapter3.Var("b")), // Body: -(a, b)
			),
		),
	) // End outer let

	tc := chapter3.TestCase{Name: "impref_counter", Expected: -1, Expr: expr}
	RunImpRefTest(t, evaluator, &tc, nil)
}

func TestImpref_RecProc(t *testing.T) {
	/* Python: (Same as before, uses 'set')
	   "recproc": ("""
	           let times4 = 0 in
	               begin
	                   set times4 = proc(x)
	                                   if isz(x)
	                                   then 0
	                               else -((times4 -(x,1)), -4) ;
	                   (times4 3)
	               end
	       """, 12),
	*/
	evaluator := NewTestImpRefLangEval()
	expr := chapter3.Let(
		chapter3.ExprDict("times4", chapter3.Lit(0)), // times4 = Ref(Lit(0))
		Begin(
			Assign("times4", // set times4 = proc...
				chapter3.Proc([]string{"x"},
					chapter3.If(chapter3.IsZero(chapter3.Var("x")),
						chapter3.Lit(0),
						// -( (times4 -(x,1)), -4) is equivalent to +((times4 -(x,1)), 4)
						chapter3.Op("-",
							chapter3.Call(chapter3.Var("times4"), chapter3.Op("-", chapter3.Var("x"), 1)),
							-4,
						),
					),
				),
			), // Assign returns the proc (discarded)
			chapter3.Call(chapter3.Var("times4"), 3), // Call the proc stored in times4
		),
	)

	tc := chapter3.TestCase{Name: "impref_recproc", Expected: 12, Expr: expr}
	RunImpRefTest(t, evaluator, &tc, nil)
}

func TestImpref_CallByRef(t *testing.T) {
	/* Reinterpretation for Go model using RefVar constructor
	   "callbyref": ("""
	           let a = 3
	           in let b = 4
	               in let swap = proc(x,y)
	                               let temp = deref(x)
	                               in begin
	                                   setref(x, deref(y));
	                                   setref(y,temp)
	                               end
	                   in begin ((swap ref a) ref b) ; -(a,b) end
	       """, 1),
	*/
	evaluator := NewTestImpRefLangEval()
	expr := chapter3.Let(chapter3.ExprDict("a", chapter3.Lit(3)), // a = Ref(Lit(3))
		chapter3.Let(chapter3.ExprDict("b", chapter3.Lit(4)), // b = Ref(Lit(4))
			chapter3.Let(chapter3.ExprDict("swap", // swap = Ref(BoundProc(...))
				chapter3.Proc([]string{"x", "y"}, // Params x, y will be bound to the Refs passed as args
					chapter3.Let(chapter3.ExprDict("temp", DeRef(chapter3.Var("x"))), // temp = Ref(deref(x)) -> temp = Ref(Lit(3))
						Begin(
							// Pass the variable directly; valueOfSetRef expects the 1st arg to eval to a Ref
							// Since x and y hold the Refs passed in from RefVar("a")/RefVar("b"), this works.
							SetRef(chapter3.Var("x"), DeRef(chapter3.Var("y"))),
							// Need the *value* from temp's ref for SetRef
							SetRef(chapter3.Var("y"), chapter3.Var("temp")),
						),
					),
				),
			), // End swap binding
				// Body of innermost let
				Begin(
					// Use RefVar constructor for RefExpr to pass the reference itself
					chapter3.Call(chapter3.Var("swap"), RefVar("a"), RefVar("b")), // Returns Lit(3) (discarded)
					chapter3.Op("-", chapter3.Var("a"), chapter3.Var("b")),        // -(a,b) -> -(4, 3) = 1
				),
			),
		),
	)

	tc := chapter3.TestCase{Name: "impref_callbyref_style", Expected: 1, Expr: expr}
	RunImpRefTest(t, evaluator, &tc, nil)
}

// Add Eq and Printable tests for AssignExpr
func TestAssignExprEqPrintable(t *testing.T) {
	// Eq
	assert.True(t, chapter3.ExprEq(Assign("x", 1), Assign("x", 1)))
	assert.False(t, chapter3.ExprEq(Assign("x", 1), Assign("y", 1)))
	assert.False(t, chapter3.ExprEq(Assign("x", 1), Assign("x", 2)))

	// Printable - rely on common_test.go for formatting checks
	// Structural check could be added if desired
}
