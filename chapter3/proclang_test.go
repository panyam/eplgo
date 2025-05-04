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

/*
def test_proc2():
    runtest(*(cases.proclang["proc2"]))


def test_proc_multiargs():
    runtest(*(cases.proclang["proc_multiargs"]))

def test_proc_currying():
    runtest(*(cases.proclang["proc_currying"]))

def test_proc_currying2():
    runtest(*(cases.proclang["proc_currying2"]))
proclang = {
    "proc2": ("""
        let x = 200 in
            let f = proc(z) -(z,x) in
                let x = 100 in
                    let g = proc(z) -(z,x) in
                        -((f 1), (g 1))
        """, -100),
    "proc_multiargs": ("""
        let f = proc(x,y) -(x,y) in
            -((f 1 10), (f 10 5))
    """, -14),
    "proc_currying": (""" let f = proc(x,y) -(x,y) in ((f 5) 3) """, 2),
    "proc_currying2": ("""
        let f = proc(x,y)
                if (isz y)
                then x
                else proc(a,b) (if isz b then +(a,x,y) else +(a,b,x,y))
        in
        (f 1 2 2 0)
    """, 5),
}
*/
