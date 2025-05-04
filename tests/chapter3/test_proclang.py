
from epl import bp
from epl.bp import eprint
from epl.chapter3 import proclang
from tests import settings
from tests.utils import runevaltest
from tests.chapter3 import cases

Expr = proclang.Expr
Eval = proclang.Eval

def runtest(input, exp, **extra_env):
    with settings.push(Expr = Expr, Eval = Eval):
        return runevaltest(input, exp, **extra_env)

def test_proc1():
    runtest(*(cases.proclang["proc1"]))

def test_proc2():
    runtest(*(cases.proclang["proc2"]))


def test_proc_multiargs():
    runtest(*(cases.proclang["proc_multiargs"]))

def test_proc_currying():
    runtest(*(cases.proclang["proc_currying"]))

def test_proc_currying2():
    runtest(*(cases.proclang["proc_currying2"]))

