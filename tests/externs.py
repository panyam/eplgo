
from ipdb import set_trace
import functools
from epl.common import Lit

# "Custom" functions are just "calls" as macros!!!
def contenv():
    def isz_expr_eval(exprs):
        return Lit(exprs[0].value == 0)

    def minus_expr_eval(exprs):
        val1 = exprs[0]
        val2 = exprs[1]
        return val1 - val2

    def div_expr_eval(exprs):
        val1 = exprs[0]
        val2 = exprs[1]
        return Lit(val1.value / val2.value)

    def plus_expr_eval(exprs):
        vals = [e.value for e in exprs]
        return Lit(sum(vals))

    def mult_expr_eval(exprs):
        vals = [e.value for e in exprs]
        return Lit(functools.reduce(lambda x,y: x*y, vals, 1))

    return {
        "isz": isz_expr_eval,
        "+": plus_expr_eval,
        "-": minus_expr_eval,
        "/": div_expr_eval,
        "*": mult_expr_eval,
    }

def env():
    def isz_expr_eval(evalfunc, env, exprs):
        return Lit(evalfunc(exprs[0], env).value == 0)

    def minus_expr_eval(evalfunc, env, exprs):
        val1 = evalfunc(exprs[0], env)
        val2 = evalfunc(exprs[1], env)
        return val1 - val2

    def div_expr_eval(evalfunc, env, exprs):
        val1 = evalfunc(exprs[0], env).value
        val2 = evalfunc(exprs[1], env).value
        return Lit(val1 / val2)

    def plus_expr_eval(evalfunc, env, exprs):
        vals = [evalfunc(exp, env).value for exp in exprs]
        return Lit(sum(vals))

    def mult_expr_eval(evalfunc, env, exprs):
        vals = [evalfunc(exp, env).value for exp in exprs]
        return Lit(functools.reduce(lambda x,y: x*y, vals, 1))

    return {
        "isz": isz_expr_eval,
        "+": plus_expr_eval,
        "-": minus_expr_eval,
        "/": div_expr_eval,
        "*": mult_expr_eval,
    }

