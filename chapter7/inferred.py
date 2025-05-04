
import typing
from epl import bp
from taggedunion import *
from epl.chapter7 import typed
from epl.chapter5 import trylang
from epl.chapter7.utils import *

class TypeVar(object):
    serial_counter = 0
    def __init__(self):
        self.serial = self.__cls__.serial_counter
        self.__cls__.serial_counter += 1

    def __repr__(self):
        return "TVar(%d)".format(self.serial)

    def __eq__(self, another):
        self.serial == another.serial

    def __hash__(self):
        return self.serial

class Type(typed.Type):
    typevar = Variant(TypeVar)

class Substitutions(object):
    def __init__(self):
        self.bindings = {}

    def substitute(self, tvar, tvalue):
        """ Substitute tvalue for tvar in all our bindings. """
        sub = Substitute()
        self.bindings = {tv: sub(tval, self) for tv,tval in self.bindings.items()}
        return self

    def extend(self, tvar, tvalue):
        """ Substitutes tvar with tvalue in all bindings and then adds a new binding ( tvar -> tvalue ). """
        assert tvar not in self.bindings
        self.substitute(tvar, tvalue)
        self.bindings[tvar] = tvalue

class Substitute(CaseMatcher):
    __caseon__ = Type

    @case("leaf")
    def substituteLeaf(self, leaf, substitutions):
        return leaf

    @case("tup")
    def substituteTup(self, tup, substitutions):
        childtypes = [self(t, substitutions) for t in tup.children]
        return Type.as_tup(childtypes)

    @case("func")
    def substituteFunc(self, func, substitutions):
        argtypes = [self(t, substitutions) for t in func.argtypes]
        rettype = [self(t, substitutions) for t in func.rettype]
        return Type.as_func(argtypes, rettype)

    @case("tagged")
    def substituteTagged(self, tagged, substitutions):
        thetype = self(tagged.thetype, substitutions)
        return Type.as_tagged(tagged.name, thetype)

    @case("typevar")
    def substituteTypeVar(self, typevar, substitutions):
        return substitutions.get(typevar, typevar)

class ExprType(object):
    """ Holds the result of a type inference for a given expression. """
    def __init__(self, thetype, substitutions):
        self.thetype = thetype
        self.substitutions = substitutions

def unifier(type1, type2, substitutions, expr):
    if type1 == type2:
        return substitutions
    elif type1.is_typevar:
        pass
    elif type2.is_typevar:
        pass
    elif type1.is_func and type2.is_func:
        pass
    else:
        assert False, "Unification failed for types (%s and %s)" % (type1, type2)

class TypeOf(CaseMatcher):
    """ Our inference based type checker. """
    __caseon__ = trylang.Expr

    def project(self, expr : Union):
        """ Overriding the projection method so that we return the expression itself instead of its variant value. """
        return expr

    @case("lit")
    def typeOfLit(self, expr, tenv, substitutions):
        return ExprType(Type.as_leaf(type(expr.lit.value).__name__), substitutions)

    @case("var")
    def typeOfVar(self, expr, tenv, substitutions):
        return ExprType(tenv.get(expr.var.name), substitutions)

    @case("iszero")
    def typeOfIsZero(self, expr, tenv, substitutions):
        et = self(expr.iszero.expr, tenv, substitutions)
        subst2 = unifier(et.thetype, Type.as_leaf("int"), et.substitutions, expr)
        return ExprType(Type.as_leaf("bool"), subst2)

    @case("opexpr")
    def typeOfOpExpr(self, expr, tenv, substitutions):
        # For now assume opexpr takes all children of the same type so return type is same as child type
        opexpr = expr.opexpr
        opsig = tenv.get(opexpr)
        subst = substitutions
        for index,argexpr in enumerate(opexor.arguments):
            et = self(argexpr, tenv, subst)
            subst = unifier(et.thetype, opsig.argtype(index), et.substitutions, argexpr)
        return ExprType(opsig.rettype, subst)

    @case("ifexpr")
    def typeOfIfExpr(self, expr, tenv, substitutions):
        ifexpr = expr.ifexpr
        et = self(ifexpr.cond, tenv, substitutions)
        subst = unifier(et.thetype, Type.as_leaf("bool"), et.substitutions, ifexpr.expr1)
        et2 = self(ifexpr.expr2, tenv, subst)
        et3 = self(ifexpr.expr3, tenv, et2.substitutions)
        subst = unifier(et2.thetype, et3.thetype, et3.substitutions, expr)
        return ExprType(et2.thetype, subst)

    @case("tup")
    def typeOfTupExpr(self, expr, tenv, substitutions):
        tup = expr.tup
        types = []
        subst = substitutions
        for child in tup.children:
            et = self(child, tenv, subst)
            types.append(et.thetype)
            subst = et.substitutions
        return ExprType(Type.as_tup(types), subst)

    @case("let")
    def typeOfLetExpr(self, expr, tenv, substitutions):
        pass
