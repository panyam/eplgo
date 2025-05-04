package chapter3

import "testing"

// No additional specific imports needed for tests beyond testing and chapter3 itself

// NewTestLetRecLangEval creates an evaluator and sets up standard operators.
func NewTestLetRecLangEval() Evaluator {
	return setOpFuncs(NewLetRecLangEval())
}

/*
letreclang = {
    "oddeven": ("""
            letrec
                even(x) = if isz(x) then 1 else (odd -(x,1))
                odd(x) = if isz(x) then 0 else (even -(x,1))
            in (odd 13)
        """, 1),
    "currying": ("""
            letrec f(x,y) = if (isz y)
                            then x
                            else (f +(x,y))
            in
            (f 1 2 3 4 5 0)
        """, 15)
}
*/

func TestLetRecDouble(t *testing.T) {
	/* Python:
	   "double": ("""
	       letrec double(x) = if isz(x) then 0 else - ((double -(x,1)), -2)
	       in (double 6)
	   """, 12),
	*/
	expr := LetRec(
		ProcMap( // Use helper ProcMap
			"double", Proc([]string{"x"},
				If(IsZero(Var("x")),
					Lit(0),
					// -( (double -(x,1)), -2) is equivalent to +((double -(x,1)), 2)
					// Assuming '-' takes 2 args, '*' takes >= 1, '+' takes >= 0
					// Let's stick to the original structure:
					Op("-",
						Call(Var("double"), Op("-", Var("x"), Lit(1))),
						Lit(-2), // Use negative literal
					),
				),
			),
		), // End ProcMap
		Call(Var("double"), Lit(6)), // Body
	)
	tc := TestCase{Name: "letrec_double", Expected: 12, Expr: expr}
	RunTest(t, NewTestLetRecLangEval(), &tc, nil)
}

func TestLetRecOddEven(t *testing.T) {
	/* Python:
	   "oddeven": ("""
	           letrec
	               even(x) = if isz(x) then 1 else (odd -(x,1))
	               odd(x) = if isz(x) then 0 else (even -(x,1))
	           in (odd 13)
	       """, 1),
	*/
	expr := LetRec(
		ProcMap( // Use helper ProcMap
			"even", Proc([]string{"x"},
				If(IsZero(Var("x")),
					Lit(1),
					Call(Var("odd"), Op("-", Var("x"), Lit(1))),
				),
			),
			"odd", Proc([]string{"x"},
				If(IsZero(Var("x")),
					Lit(0),
					Call(Var("even"), Op("-", Var("x"), Lit(1))),
				),
			),
		), // End ProcMap
		Call(Var("odd"), Lit(13)), // Body
	)
	tc := TestCase{Name: "letrec_oddeven", Expected: 1, Expr: expr}
	RunTest(t, NewTestLetRecLangEval(), &tc, nil)
}

func TestLetRecCurrying(t *testing.T) {
	/* Python:
	   "currying": ("""
	           letrec f(x,y) = if (isz y)
	                           then x
	                           else (f +(x,y)) // NOTE: Python code calls f with only one arg here!
	           in
	           (f 1 2 3 4 5 0)
	       """, 15)
	   Let's assume the Python meant `(f (+ x y) 0)` or similar for termination,
	   or perhaps relies on currying in a way not directly obvious.
	   The recursive call `(f +(x,y))` only provides *one* argument to `f` which expects two (`x`, `y`).
	   This implies the result of `f (+ x y)` should be a curried function expecting the `y` argument.
	   However, the structure seems intended for iterative summation until `y` is 0.
	   Let's try to implement the summation logic directly assuming that was the intent,
	   as the provided Python code seems broken or relies on subtle currying behavior
	   that might differ in Go's strict applyProc.

	   Revised interpretation for summation:
	   f(current_sum, next_val) = if next_val == 0 then current_sum else f(current_sum + next_val, ???)
	   This structure doesn't quite match the call `(f 1 2 3 4 5 0)`.

	   Let's try to implement the python *literally* and see if our `applyProc` handles it.
	   `f(x,y) = if isz(y) then x else (f +(x,y))`
	   Call: `(f 1 2 3 4 5 0)`
	   - f(1, 2) -> y != 0 -> call (f +(1,2)) -> call (f 3)
	   - (f 3) is a call to a function expecting two args with only one. `applyProc` should curry.
	   - It returns `proc(y') = if isz(y') then 3 else (f +(3, y'))` bound to the letrec env. Let's call this `f_curried_3`.
	   - Now we call `(f_curried_3 3)`.
	   - f_curried_3(3) -> y' != 0 -> call (f +(3, 3)) -> call (f 6)
	   - (f 6) -> returns `proc(y') = if isz(y') then 6 else (f +(6, y'))` (let's call this `f_curried_6`)
	   - Now call `(f_curried_6 4)`
	   - f_curried_6(4) -> y' != 0 -> call (f +(6, 4)) -> call (f 10)
	   - (f 10) -> returns `f_curried_10`
	   - Call `(f_curried_10 5)`
	   - f_curried_10(5) -> y' != 0 -> call (f +(10, 5)) -> call (f 15)
	   - (f 15) -> returns `f_curried_15`
	   - Call `(f_curried_15 0)`
	   - f_curried_15(0) -> y' == 0 -> returns x (which was 15) -> 15.

	   Okay, the literal interpretation seems to work with currying. Let's implement that.
	*/
	expr := LetRec(
		ProcMap(
			"f", Proc([]string{"x", "y"},
				If(IsZero(Var("y")),
					Var("x"),
					// Recursive call providing only the first argument 'x'
					Call(Var("f"), Op("+", Var("x"), Var("y"))),
				),
			),
		), // End ProcMap
		// Body: Call f with all arguments
		Call(Var("f"), Lit(1), Lit(2), Lit(3), Lit(4), Lit(5), Lit(0)),
	)
	tc := TestCase{Name: "letrec_currying_sum", Expected: 15, Expr: expr}
	RunTest(t, NewTestLetRecLangEval(), &tc, nil)
}
