
import typing
from taggedunion import *
from epl.chapter3 import letreclang

class RefExpr(object):
    """ This expression both captures a reference to a cell as well as a reference to a variable. """
    def __init__(self, expr_or_var):
        self.expr = expr_or_var

    @property
    def is_var(self):
        return type(self.expr) is str

    def printables(self):
        yield 0, "Ref:"
        yield 1, self.expr.printables()

    def __eq__(self, another):
        return self.expr == another.expr

class DeRefExpr(object):
    def __init__(self, expr):
        self.expr = expr

    def printables(self):
        yield 0, "DeRef:"
        yield 1, self.expr.printables()

    def __eq__(self, another):
        return self.expr == another.expr

class SetRefExpr(object):
    def __init__(self, ref, value):
        self.ref = ref
        self.value = value

    def printables(self):
        yield 0, "SetRef:"
        yield 1, "Ref:"
        yield 2, self.ref.printables()
        yield 1, "Value:"
        yield 2, self.value.printables()

    def __eq__(self, another):
        return self.ref == another.ref and \
               self.value == another.value

class BlockExpr(object):
    def __init__(self, exprs):
        self.exprs = exprs

    def printables(self):
        yield 0, "Begin:"
        for expr in self.exprs:
            yield 1, expr.printables()
        yield 0, "End"

    def __eq__(self, another):
        return self.exprs == another.exprs

class Expr(letreclang.Expr):
    ref = Variant(RefExpr)
    deref = Variant(DeRefExpr)
    setref = Variant(SetRefExpr)
    block = Variant(BlockExpr)

class Eval(letreclang.Eval):
    __caseon__ = Expr

    def valueOf(self, expr, env):
        return self(expr, env)

    @case("ref")
    def valueOfRef(self, ref, env):
        # evaluate ref value if it is an Expr only
        # otherwise we could have a value or a variable (in which case it is a ref to a var)
        if not ref.is_var:
            return self.__caseon__.as_ref(self.valueOf(ref.expr, env)).ref
        else:
            # Return the ref cell as is - upto caller to use 
            # this reference and the value in it as it sees fit
            return ref

    @case("deref")
    def valueOfDeRef(self, deref, env):
        ref = self(deref.expr, env)
        assert type(ref) is RefExpr
        if ref.is_var:
            # Then get the value of the named ref
            return env.get(ref.expr)
        else:
            return ref.expr

    @case("setref")
    def valueOfSetRef(self, setref, env):
        val1 = self(setref.ref, env)
        assert type(val1) is RefExpr
        val2 = self(setref.value, env)
        if val1.is_var:
            # since a named ref get the ref by name first - we have an extra level of indirection here
            env.replace(val1.expr, val2)
        else:
            # Set ref as is
            val1.expr = val2
        return val2

    @case("block")
    def valueOfBlock(self, block, env):
        value = self.__caseon__.as_lit(0)
        for expr in block.exprs:
            value = self(expr, env)
        return value
