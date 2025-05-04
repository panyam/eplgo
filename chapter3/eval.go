package chapter3

import (
	epl "github.com/panyam/eplgo"
)

func LetLangEval(expr Expr, env *epl.Env[any]) any {
	switch n := expr.(type) {
	case *LitExpr:
		return ValueOfLit(n, env)
	}
	return nil
}

func ValueOfLit(l *LitExpr, env *epl.Env[any]) any {
	return l.Value
}

/*

class Eval(CaseMatcher):
    # Tells which union type we are "case matching" on
    __caseon__ = Expr

    # These names MUST match the different cases in our "union_type"
    def valueOf(self, expr, env):
        # We expect the signature of "valueOf" and each selected
        # subexpression to have the same arity and return type
        return self(expr, env)

    @case("lit")
    def valueOfLit(self, lit, env = None):
        return lit

    @case("var")
    def valueOfVar(self, var, env):
        return env.get(var.name)

    @case("opexpr")
    def valueOfOpExpr(self, opexpr, env):
        # In this lang we make "op" expressions just hooks to external plugins
        # We can get even more generic once we have procedures (proclang onwards)
        opfunc = env.get(opexpr.op)
        assert opfunc is not None, "No plug in found for operator: %s" % opexpr.op
        return opfunc(self, env, opexpr.arguments)

    @case("tup")
    def valueOfTuple(self, tup, env):
        values = [self.valueOf(v,env) for v in tup.children]
        return TupleExpr(*values)

    @case("iszero")
    def valueOfIsZero(self, iszero, env):
        return Lit(self.valueOf(iszero.expr, env).value == 0)

    @case("ifexpr")
    def valueOfIf(self, ifexpr, env):
        result = self.valueOf(ifexpr.cond, env).value
        return self.valueOf(ifexpr.exp1 if result else ifexpr.exp2, env)

    @case("let")
    def valueOfLet(self, let, env):
        expvals = {var: self.valueOf(exp, env) for var,exp in let.mappings.items()}
        newenv = env.extend(**expvals)
        return self.valueOf(let.body, newenv)


class Expr(Union):
    # Convert this into a union metaclass
    lit = Variant(Lit)
    var = Variant(VarExpr)
    opexpr = Variant(OpExpr)
    tup = Variant(TupleExpr, checker = "is_tup", constructor = "as_tup")
    iszero = Variant(IsZeroExpr)
    ifexpr = Variant(IfExpr, checker = "is_if", constructor = "as_if")
    let = Variant(LetExpr, checker = "is_let", constructor = "as_let")

    @classmethod
    def as_diff(cls, e1, e2):
        return cls.as_opexpr("-", e1, e2)

    def printables(self):
        yield 0, self.variant_value.printables()
*/
