
import typing
from taggedunion import *
from epl.chapter4 import impreflang

class LazyExpr(object):
    """ A lazy expression is like a bound proc that captures the environment it is being evaluated in. """

    def __init__(self, expr):
        self.expr = expr

    def bind(self, env):
        return Thunk(self, env)

    def printables(self):
        yield 0, "LazyExpr:"
        yield 1, self.expr.printables()

    def __eq__(self, another):
        return self.expr == another.expr

class Thunk(object):
    def __init__(self, expr, env):
        self.expr = expr
        self.saved_env = env

class Expr(impreflang.Expr):
    lazy = Variant(LazyExpr)
    thunk = Variant(Thunk)

class Eval(impreflang.Eval):
    __caseon__ = Expr

    def valueOf(self, expr, env):
        return self(expr, env)

    @case("lazy")
    def valueOfLazyExpr(self, lazy_expr, env):
        return lazy_expr.bind(env)

    @case("thunk")
    def valueOfThunk(self, thunk, env):
        return self.valueOf(thunk.expr, thunk.saved_env)
