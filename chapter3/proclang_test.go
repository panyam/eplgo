package chapter3

import "testing"

func NewTestProcLangEval() Evaluator {
	return setOpFuncs(NewProcLangEval())
}

func TestProc1(t *testing.T) {
	// """ let f = proc (x) -(x,11) in (f (f 77)) """, 55),
	tc := TestCase{"proc1", 55,
		Let(ExprDict("f", Proc([]string{"x"}, Op("-", Var("x"), Lit(11)))),
			Call(Var("f"), Call(Var("f"), Lit(77))))}
	RunTest(t, NewTestProcLangEval(), &tc, nil)
}

func TestProc2(t *testing.T) {
	/*
	   "proc2": ("""
	       let x = 200 in
	           let f = proc(z) -(z,x) in
	               let x = 100 in
	                   let g = proc(z) -(z,x) in
	                       -((f 1), (g 1))
	       """, -100),
	*/
	tc := TestCase{"proc2", -100,
		Let(ExprDict("x", Lit(200)),
			Let(ExprDict("f", Proc([]string{"z"}, Op("-", Var("z"), Var("x")))),
				Let(ExprDict("x", Lit(100)),
					Let(ExprDict("g", Proc([]string{"z"}, Op("-", Var("z"), Var("x")))),
						Op("-", Call(Var("f"), Lit(1)), Call(Var("g"), Lit(1)))))))}
	RunTest(t, NewTestProcLangEval(), &tc, nil)
}

func TestProcMultiArgs(t *testing.T) {
	/*
	   let f = proc(x,y) -(x,y) in
	       -((f 1 10), (f 10 5))
	*/
	tc := TestCase{"proc_multiargs", -14,
		Let(ExprDict("f", Proc([]string{"x", "y"}, Op("-", Var("x"), Var("y")))),
			Op("-", Call(Var("f"), Lit(1), Lit(10)), Call(Var("f"), Lit(10), Lit(5))))}
	RunTest(t, NewTestProcLangEval(), &tc, nil)
}

func TestProcCurrying(t *testing.T) {
	/*
	   let f = proc(x,y) -(x,y) in ((f 5) 3)
	*/
	tc := TestCase{"proc_currying", 2,
		Let(ExprDict("f", Proc([]string{"x", "y"}, Op("-", Var("x"), Var("y")))),
			Call(Call(Var("f"), Lit(5)), Lit(3)))}
	RunTest(t, NewTestProcLangEval(), &tc, nil)
}

func TestProcCurrying2(t *testing.T) {
	/*
	   let f = proc(x,y)
	           if (isz y)
	           then x
	           else proc(a,b) (if isz b then +(a,x,y) else +(a,b,x,y))
	   in
	   (f 1 2 2 0)
	*/
	expected := 5
	expr := Let(ExprDict("f",
		Proc([]string{"x", "y"},
			If(IsZero("y"),
				"x",
				Proc([]string{"a", "b"},
					If(IsZero("b"),
						Op("+", "a", "x", "y"),
						Op("+", "a", "b", "x", "y")))))),
		Call("f", 1, 2, 2, 0),
	)
	tc := TestCase{"proc_currying2", expected, expr}
	RunTest(t, NewTestProcLangEval(), &tc, nil)
}
