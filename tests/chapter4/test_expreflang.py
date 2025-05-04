
from epl.chapter4 import expreflang
from tests import settings
from tests.chapter4 import cases

Expr = expreflang.Expr
Eval = expreflang.Eval

def runtest(input, exp, **extra_env):
    from tests.utils import runevaltest
    with settings.push(Expr = Expr, Eval = Eval):
        return runevaltest(input, exp, **extra_env)

def test_oddeven():
    runtest(*(cases.exprefs["oddeven"]))

def test_counter():
    runtest(*(cases.exprefs["counter"]))

