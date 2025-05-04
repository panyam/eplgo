package chapter3

import (
	"testing"

	epl "github.com/panyam/eplgo"
)

var ExprDict = epl.Dict[string, Expr]

func NewTestLetLangEval() *LetLangEval {
	out := NewLetLangEval()
	out.SetOpFunc("-", func(env *epl.Env[any], args []Expr) any {
		v1 := out.Eval(args[0], env).(*LitExpr)
		v2 := out.Eval(args[1], env).(*LitExpr)
		return Lit(v1.Value.(int) - v2.Value.(int))
	})
	return out
}

func TestNum(t *testing.T) {
	tc := TestCase{"num", 3, Lit(3)}
	RunTest(t, NewTestLetLangEval(), &tc, nil)
}

func TestVar(t *testing.T) {
	tc := TestCase{"var", 5, Var("x")}
	RunTest(t, NewTestLetLangEval(), &tc, map[string]Expr{
		"x": Lit(5),
	})
}

func TestZeroTrue(t *testing.T) {
	tc := TestCase{"isz_true", true, IsZero(Lit(0))}
	RunTest(t, NewTestLetLangEval(), &tc, nil)
}

func TestZeroFalse(t *testing.T) {
	tc := TestCase{"isz_false", false, IsZero(Lit(1))}
	RunTest(t, NewTestLetLangEval(), &tc, nil)
}

func TestDiff(t *testing.T) {
	tc := TestCase{"diff", 3,
		Op("-",
			Op("-",
				Var("x"),
				Lit(3),
			),
			Op("-",
				Var("v"),
				Var("i"),
			))}
	RunTest(t, NewTestLetLangEval(), &tc, map[string]Expr{
		"i": Lit(1),
		"v": Lit(5),
		"x": Lit(10),
	})
}

func TestIf(t *testing.T) {
	tc := TestCase{"if", 18,
		If(
			IsZero(Op("-", Var("x"), Lit(11))),
			Op("-", Var("y"), Lit(2)),
			Op("-", Var("y"), Lit(4)))}
	RunTest(t, NewTestLetLangEval(), &tc, map[string]Expr{
		"x": Lit(33),
		"y": Lit(22),
	})
}

func TestLet(t *testing.T) {
	tc := TestCase{"let", 2,
		Let(ExprDict("x", Lit(5)),
			Op("-", Var("x"), Lit(3)),
		),
	}
	RunTest(t, NewTestLetLangEval(), &tc, nil)
}

func TestLetNested(t *testing.T) {
	tc := TestCase{"letnested", 3,
		// ("let z = 5 in let x = 3 in let y = -(x, 1) in let x = 4 in -(z, -(x,y))", 3),
		Let(ExprDict("z", Lit(5)),
			Let(ExprDict("x", Lit(3)),
				Let(ExprDict("y", Op("-", Var("x"), Lit(1))),
					Let(ExprDict("x", Lit(4)),
						Op("-", Var("z"), Op("-", Var("x"), Var("y"))),
					),
				),
			),
		)}
	RunTest(t, NewTestLetLangEval(), &tc, nil)
}

func TestLet3(t *testing.T) {
	/*
	   "let3": ("""
	       let x = 7 in
	           let y = 2 in
	               let y = let x = -(x,1) in -(x,y)
	               in -(-(x,8), y)
	       """, -5),
	*/
	tc := TestCase{"let3", -5,
		Let(ExprDict("x", Lit(7)),
			Let(ExprDict("y", Lit(2)),
				Let(ExprDict("y", Let(ExprDict("x", Op("-", Var("x"), Lit(1))), Op("-", Var("x"), Var("y")))),
					Op("-", Op("-", Var("x"), Lit(8)), Var("y")),
				),
			),
		)}
	RunTest(t, NewTestLetLangEval(), &tc, nil)
}

func TestLetMultiArgs(t *testing.T) {
	/*
	   "letmultiargs": ("""
	       let x = 7 y = 2 in
	           let y = let x = -(x,1) in -(x,y)
	           in -(-(x, 8),y)
	       """, -5),
	*/
	tc := TestCase{"letmultiargs", -5,
		Let(ExprDict("x", Lit(7), "y", Lit(2)),
			Let(ExprDict("y",
				Let(ExprDict("x", Op("-", Var("x"), Lit(1))),
					Op("-", Var("x"), Var("y")))),
				Op("-", Op("-", Var("x"), Lit(8)), Var("y")))),
	}
	RunTest(t, NewTestLetLangEval(), &tc, nil)
}
