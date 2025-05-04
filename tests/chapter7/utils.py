
import taggedunion
from epl import bp
from epl.chapter7 import typed
from tests import externs
from tests import settings
from tests.utils import parse, runevaltest, prepareenv

def runtest(input, exp, starting_env = None, **extra_env):
    Expr = settings.get("Expr")
    Type = settings.get("Type")
    TypeOf = settings.get("TypeOf")
    expr,tree = parse(input, Expr, Type)
    if settings.get("print_tree", no_throw = True):
        bp.eprint(expr)
    env = prepareenv(starting_env, **extra_env)
    try:
        foundtype = TypeOf()(expr, env)
        stringifier = typed.Stringifier()
        if settings.get("print_type", no_throw = True):
            print("Found Type: ", stringifier(foundtype))
        if type(exp) is not bool and foundtype != exp:
            assert False, "Type Found: %s, Expected: %s" % (stringifier(foundtype), stringifier(exp))
    except taggedunion.InvalidVariantError as ive:
        assert exp == False
    except typed.TypeError as te:
        assert exp == False
