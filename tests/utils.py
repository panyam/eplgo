
from epl import parser
from epl.env import DefaultEnv as Env
from epl.bp import eprint
from tests import externs
from tests import settings

def parse(input, expr_class, type_class, optable = None):
    mixins = [  parser.BasicMixin, parser.ProcMixin, parser.LetMixin,
                parser.LetRecMixin, parser.RefMixin, parser.TryMixin,
                parser.TypingMixin ]
    theparser = parser.make_parser(expr_class, type_class, optable, *mixins)
    return theparser.parse(input)

def prepareenv(starting_env = None, default_env = None, **extra_env):
    starting_env = starting_env or default_env or Env().set(**externs.env())
    newenv = starting_env.push().set(**extra_env)
    return newenv

def runevaltest(input, exp, starting_env = None, **extra_env):
    Expr = settings.get("Expr")
    Eval = settings.get("Eval")
    expr,tree = parse(input, Expr, None)
    value = Eval().valueOf(expr, prepareenv(starting_env, **extra_env))
    assert value == exp

