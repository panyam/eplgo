
from epl import bp
import pytest
from epl.chapter5 import trylang
from epl.chapter7 import typed
from tests import settings
from tests.chapter7 import cases
from tests.chapter7 import utils

@pytest.mark.parametrize("input, expected", cases.checked.values())
def test_checked(input, expected):
    with settings.push(Expr = trylang.Expr,
                       TypeOf = typed.TypeOf,
                       Type = typed.Type, 
                       print_tree = False):
        utils.runtest(input, expected)
