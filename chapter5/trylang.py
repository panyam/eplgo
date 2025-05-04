
import typing
from taggedunion import *
from epl.chapter4 import lazylang
from epl.chapter5 import continuations

class TryCatch(object):
    """ A try catch expression. """
    def __init__(self, expr, varname, handlerexpr):
        self.expr = expr
        self.varname = varname
        self.handler = handler

    def printables(self):
        yield 0, "Try:"
        yield 1, self.expr.printables()
        yield 1, "Catch (%s)" % self.varname
        yield 2, self.exception.printables()

    def __eq__(self, another):
        return self.expr == another.expr and        \
                self.varname == another.varname and \
                self.exception == another.exception

class Raise(object):
    def __init__(self, expr):
        self.expr = expr

    def printables(self):
        yield 0, "Raise:"
        yield 1, self.expr.printables()

    def __eq__(self, another):
        return self.expr == another.expr

class Expr(lazylang.Expr):
    tryexpr = Variant(TryCatch)
    raiseexpr = Variant(Raise)

class Eval(continuations.Eval):
    __caseon__ = Expr

    @case("tryexpr")
    def valueOfTry(self, tryexpr, env, cont):
        return self(tryexpr.expr, env, TryCont(cont, tryexpr))

    @case("raiseexpr")
    def valueOfRaise(self, raiseexpr, env, cont):
        return self(raiseexpr, env, RaiseCont(self, env, cont))

class TryCont(continuations.Cont):
    def __init__(self, cont, tryexpr):
        continuations.Cont.__init__(self, None, None, cont)
        self.tryexpr = tryexpr

    def apply(self, result):
        return self.nextcont.apply(result)

class RaiseCont(continuations.Cont):
    def __init__(self, Eval, env, cont):
        continuations.Cont(Eval, env, cont)

    def apply(self, exc):
        # We have an exception here so we need to find the closes handler 
        # that can handle this exception
        cont = self
        while cont and type(cont) is not TryCont:
            cont = cont.nextcont
        assert cont is not None
        tryexpr = cont.tryexpr
        newenv = cont.env.push().setone(tryexpr.varname, exc)
        return self.Eval(tryexpr.handler, newenv, cont)
