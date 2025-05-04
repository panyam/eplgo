
import pytest
from epl.bp import eprint
from epl.chapter4 import lazylang
from epl.chapter5 import continuations
from tests.chapter4 import cases
from tests import settings

Expr = lazylang.Expr
Eval = continuations.Eval

def runtest(input, exp, **extra_env):
    from tests import externs
    from tests.utils import runevaltest
    with settings.push(Expr = Expr, Eval = Eval):
        from epl.env import DefaultEnv as Env
        starting_env = Env().set(**externs.contenv())
        return runevaltest(input, exp, starting_env, **extra_env)

@pytest.mark.parametrize("input, expected", cases.lazy.values())
def test_lazy(input, expected):
    runtest(input, expected)

@pytest.mark.parametrize("input, expected", cases.misc.values())
def test_misc(input, expected):
    runtest(input, expected)

@pytest.mark.parametrize("input, expected", cases.exprefs.values())
def test_exprefs(input, expected):
    runtest(input, expected)

@pytest.mark.parametrize("input, expected", cases.imprefs.values())
def test_imprefs(input, expected):
    runtest(input, expected)

