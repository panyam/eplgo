
import typing
from taggedunion import *
from epl.chapter4 import expreflang

class AssignExpr(object):
    def __init__(self, varname, expr):
        self.varname = varname
        self.expr = expr

    def printables(self):
        yield 0, "AssignExpr:"
        yield 1, "%s = " % self.varname
        yield 1, self.expr.printables()

    def __eq__(self, another):
        return self.varname == another.varname and \
               self.expr == another.expr


class Expr(expreflang.Expr):
    assign = Variant(AssignExpr)

class Eval(expreflang.Eval):
    __caseon__ = Expr

    def valueOf(self, expr, env):
        return self(expr, env)

    @case("assign")
    def valueOfAssign(self, assign, env):
        val = self(assign.expr, env)
        env.replace(assign.varname, val)
        return val
