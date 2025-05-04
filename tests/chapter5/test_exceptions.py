

import pytest
from epl.bp import eprint
from epl.chapter5 import trylang
from epl.chapter5 import continuations
from tests.chapter5 import cases
from tests import settings

Expr = trylang.Expr
Eval = continuations.Eval

def runtest(input, exp, **extra_env):
    from tests.utils import runevaltest
    with settings.push(Expr = Expr, Eval = Eval):
        from tests import externs
        from epl.env import DefaultEnv as Env
        starting_env = Env().set(**externs.contenv())
        return runevaltest(input, exp, starting_env, **extra_env)

@pytest.mark.parametrize("input, expected", cases.exceptions.values())
def _test_lazy(input, expected):
    runtest(input, expected)

