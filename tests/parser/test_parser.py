
from ipdb import set_trace
from epl import bp
from epl import parser
from epl.common import Lit
from epl.chapter3.letreclang import Expr

def parse(input, exp):
    mixins = [ parser.BasicMixin, parser.LetMixin, parser.ProcMixin, parser.LetRecMixin ]
    theparser = parser.make_parser(Expr, None, None, *mixins)
    return theparser.parse(input)

def runtest(input, exp):
    e,t = parse(input, Expr)
    assert e == exp

def test_parse_num():
    runtest("3", Expr.as_lit(3))

def test_parse_varname():
    runtest("x", Expr.as_var("x"))

def test_parse_paren():
    runtest("( ( ( 666 )) )", Expr.as_lit(666))

def test_operators():
    runtest("/(x,y)", Expr.as_opexpr("/", Expr.as_var("x"), Expr.as_var("y")))
    runtest(">>(x,y)", Expr.as_opexpr(">>", Expr.as_var("x"), Expr.as_var("y")))

def test_parse_iszero():
    e1 = Expr.as_opexpr("?", Expr.as_lit(0))
    runtest("? ( 0 )", e1)

def test_parse_iszero_cust():
    e2 = Expr.as_iszero(Expr.as_lit(33))
    runtest("isz ( ( ( 33 ) ) )", e2)

def test_parse_diff():
    e2 = Expr.as_diff(Expr.as_lit(33), Expr.as_lit(44))
    runtest("- ( 33, 44)", e2)

def test_parse_tuple():
    e2 = Expr.as_opexpr("$", Expr.as_lit(3), Expr.as_lit(4))
    runtest("$ (3, 4)", e2)

def test_parse_letrec_double():
    e2 = Expr.as_letrec({
        "double": Expr.as_proc(["x"],
                    Expr.as_if(Expr.as_iszero(Expr.as_var("x")),
                               Expr.as_lit(0),
                               Expr.as_diff(
                                   Expr.as_call(Expr.as_var("double"),
                                                Expr.as_diff(
                                                    Expr.as_var("x"),
                                                    Expr.as_lit(1))),
                                   Expr.as_lit(-2)))).procexpr
        }, Expr.as_call(Expr.as_var("double"), Expr.as_lit(6)))
    runtest("""
        letrec
            double(x) = (if (isz x) then 0 else -((double -(x,1)), -2))
        in (double 6)
    """, e2)


def test_parse_letrec_oddeven():
    even = Expr.as_proc(["x"],
                Expr.as_if(
                    Expr.as_iszero(Expr.as_var("x")),
                    Expr.as_lit(1),
                    Expr.as_call(
                        Expr.as_var("odd"),
                            Expr.as_diff(
                                Expr.as_var("x"),
                                Expr.as_lit(1))))).procexpr
    odd = Expr.as_proc(["x"],
                Expr.as_if(
                    Expr.as_iszero(Expr.as_var("x")),
                    Expr.as_lit(0),
                    Expr.as_call(
                        Expr.as_var("even"),
                            Expr.as_diff(
                                Expr.as_var("x"),
                                Expr.as_lit(1))))).procexpr
    e2 = Expr.as_letrec({"even": even, "odd": odd},
            Expr.as_call(Expr.as_var("odd"), Expr.as_lit(13)))
    runtest("""
        letrec
            even(x) = if (isz x) then 1 else (odd -(x,1))
            odd(x) = if (isz x) then 0 else (even -(x,1))
        in (odd 13)
    """, e2)
