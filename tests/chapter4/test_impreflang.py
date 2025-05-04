
from epl.chapter4 import impreflang
from tests import settings
from tests.chapter4 import cases

Expr = impreflang.Expr
Eval = impreflang.Eval

def runtest(input, exp, **extra_env):
    from tests.utils import runevaltest
    with settings.push(Expr = Expr, Eval = Eval):
        return runevaltest(input, exp, **extra_env)

def test_oddeven():
    runtest(*(cases.imprefs["oddeven"]))

def test_counter():
    runtest(*(cases.imprefs["counter"]))

def test_recproc():
    runtest(*(cases.imprefs["recproc"]))

def test_callbyref():
    runtest(*(cases.imprefs["callbyref"]))

