
from epl import bp
import pytest
from epl.chapter5 import trylang
from epl.chapter7 import typed
from tests.chapter7 import cases
from tests.chapter7 import utils
from tests import settings

def runtest(input, exp, starting_env = None, **extra_env):
    with settings.push(Expr = trylang.Expr,
                       TypeOf = typed.TypeOf,
                       Type = typed.Type):
        utils.runtest(input, exp, starting_env, **extra_env)

@pytest.mark.parametrize("input, expected", cases.basic_checks.values())
def test_basic_checks(input, expected):
    runtest(input, expected)

