

import typing
from taggedunion import *
from epl.chapter3 import letlang
from epl.chapter4 import lazylang

class Eval(CaseMatcher):
    __caseon__ = lazylang.Expr

    def applyContinuation(self, cont, expr):
        pass

    def valueOf(self, expr, env, cont = None):
        if cont is None:
            # Create end continuation
            cont = EndCont()
        return self(expr, env, cont)

    @case("lit")
    def valueOfLit(self, lit, env, cont):
        return cont.apply(lit)

    @case("var")
    def valueOfVar(self, var, env, cont):
        value = env.get(var.name)
        return cont.apply(value)

    @case("procexpr")
    def valueOfProcExpr(self, procexpr, env, cont):
        return cont.apply(procexpr.bind(env))

    @case("iszero")
    def valueOfIsZero(self, iszero, env, cont):
        return self.valueOf(iszero.expr, env, IsZeroCont(cont))

    @case("ifexpr")
    def valueOfIf(self, ifexpr, env, cont):
        # Eval mappings with end continuations
        nextcont = IfCont(self, env, cont, ifexpr)
        return nextcont.start()

    @case("let")
    def valueOfLet(self, let, env, cont):
        # With let we want to chain the let bindings one after the other
        # Eval mappings with end continuations
        nextcont = LetCont(self, env, cont, let)
        return nextcont.start()

    @case("letrec")
    def valueOfLetRec(self, letrec, env, cont):
        newenv = env.push()
        for proc in letrec.procs.values():
            newenv.setone(proc.name, proc.bind(newenv))
        return self.valueOf(letrec.body, newenv, cont)

    @case("tup")
    def valueOfTupExpr(self, tup, env, cont):
        nextcont = ExprListCont(self, env, cont, tup.children, self.__caseon__.as_tup)
        return nextcont.start()

    @case("opexpr")
    def valueOfOpExpr(self, opexpr, env, cont):
        # With let we want to chain the let bindings one after the other
        # Eval mappings with end continuations
        opfunc = env.get(opexpr.op)
        assert opfunc is not None, "No plug in found for operator: %s" % opexpr.op
        nextcont = ExprListCont(self, env, cont, opexpr.arguments, opfunc)
        return nextcont.start()

    @case("ref")
    def valueOfRef(self, ref, env, cont):
        if not ref.is_var:
            ref = self.__caseon__.as_ref(self.valueOf(ref.expr, env)).ref
        return cont.apply(ref)

    @case("deref")
    def valueOfDeRef(self, deref, env, cont):
        return self.valueOf(deref.expr, env, DeRefCont(self, env, cont))

    @case("setref")
    def valueOfSetRef(self, setref, env, cont):
        return self.valueOf(setref.ref, env, SetRefCont(self, env, cont, setref))

    @case("block")
    def valueOfBlock(self, block, env, cont):
        nextcont = ExprListCont(self, env, cont, block.exprs,
                            lambda results: results[-1])
        return nextcont.start()

    @case("assign")
    def valueOfAssign(self, assign, env, cont):
        return self.valueOf(assign.expr, env,
                    AssignCont(self, env, cont, assign.varname))

    @case("lazy")
    def valueOfLazy(self, lazy_expr, env, cont):
        return cont.apply(lazy_expr.bind(env))

    @case("thunk")
    def valueOfThunk(self, thunk, env, cont):
        return self.valueOf(thunk.expr, env, cont)

    @case("callexpr")
    def valueOfCall(self, callexpr, env, cont):
        nextcont = CallCont(self, env, cont, callexpr)
        return nextcont.start()

class Cont(object):
    def __init__(self, Eval, env, nextcont = None):
        self.Eval = Eval
        self.nextcont = nextcont
        self.env = env

    def apply(self, expr) -> "Cont":
        assert False, "Implement this."

class EndCont(Cont):
    def __init__(self):
        Cont.__init__(self, None, None, None)

    def apply(self, value : int) -> Cont:
        assert type(value) is letlang.Lit
        return value

class IsZeroCont(Cont):
    def __init__(self, cont):
        Cont.__init__(self, None, None, cont)

    def apply(self, value : int):
        return self.nextcont.apply(value == 0)

class IfCont(Cont):
    def __init__(self, Eval, env, cont, ifexpr):
        Cont.__init__(self, Eval, env, cont)
        self.ifexpr = ifexpr

    def start(self):
        return self.Eval(self.ifexpr.cond, self.env, self)

    def apply(self, expr):
        if expr:
            return self.Eval(self.ifexpr.exp1, self.env, self.nextcont)
        else:
            return self.Eval(self.ifexpr.exp2, self.env, self.nextcont)

class ExprListCont(Cont):
    """ A general continuation that needs to evaluate and "collect" N expressions in before another operation that depends on these results can be performed.  Also it is required that we chain results from one to another."""
    def __init__(self, Eval, env, cont, exprs, onresults = None):
        Cont.__init__(self, Eval, env, cont)
        self.curr = 0
        self.exprs = exprs
        self.onresults = onresults
        self.results = []

    def start(self):
        return self.Eval(self.exprs[0], self.env, self)

    def apply(self, expr):
        self.curr += 1
        self.results.append(expr)
        if self.curr < len(self.exprs):
            # we have more
            nextexpr = self.exprs[self.curr]
            return self.Eval(nextexpr, self.env, self)
        else:
            result = self.results
            if self.onresults:
                result = self.onresults(result)
            return self.nextcont.apply(result)

class DeRefCont(Cont):
    def __init__(self, Eval, env, cont):
        Cont.__init__(self, Eval, env, cont)

    def apply(self, ref):
        result = ref.expr
        if ref.is_var:
            # Then get the value of the named ref
            result = self.env.get(ref.expr)
        return self.nextcont.apply(result)

class SetRefCont(Cont):
    def __init__(self, Eval, env, cont, setref):
        Cont.__init__(self, Eval, env, cont)
        self.state = 0      # 0 = processing ref
                            # 1 == processing value
        self.setref = setref
        # results of setref child evaluations
        self.val1 = None

    def apply(self, ref):
        if self.state == 0:
            self.state += 1
            self.val1 = ref
            return self.Eval(self.setref.value, self.env, self)
        else:
            val2 = ref
            if self.val1.is_var:
                # since a named ref get the ref by name first - 
                # we have an extra level of indirection here
                self.env.replace(self.val1.expr, val2)
            else:
                # Set ref as is
                self.val1.expr = val2
            return self.nextcont.apply(val2)

class AssignCont(Cont):
    def __init__(self, Eval, env, cont, varname):
        Cont.__init__(self, Eval, env, cont)
        self.varname = varname
    
    def apply(self, result):
        self.env.replace(self.varname, result)
        return self.nextcont.apply(result)

class LetCont(Cont):
    def __init__(self, Eval, env, cont, letexpr):
        Cont.__init__(self, Eval, env, cont)
        self.curr = 0
        self.letexpr = letexpr
        self.varnames = list(letexpr.mappings.keys())
        self.newenv = env.push()

    def start(self):
        self.currvar = self.varnames[0]
        expr1 = self.letexpr.mappings[self.currvar]
        return self.Eval(expr1, self.env, self)

    def apply(self, expr):
        # Here we are called with the "expr" of the ith var being processed
        lastvar = self.currvar
        self.newenv.setone(lastvar, expr)
        self.curr += 1
        if self.curr < len(self.varnames):
            # Now evaluate the next var as if it is the body
            self.currvar = self.varnames[self.curr]
            nextexpr = self.letexpr.mappings[self.currvar]
            return self.Eval(nextpexr, self.newenv, self)
        else:
            # Here all bindings have been evaluated
            return self.Eval(self.letexpr.body, self.newenv, self.nextcont)

class CallCont(Cont):
    def __init__(self, Eval, env, cont, callexpr):
        Cont.__init__(self, Eval, env, cont)
        self.callexpr = callexpr

    def start(self):
        return self.Eval(self.callexpr.operator, self.env, self)

    def apply(self, boundproc):
        # We just received operator result
        # so kick off arg 
        proc_cont = ApplyProcCont(self.Eval, self.env, self.nextcont, boundproc)
        nextcont = ExprListCont(self.Eval, self.env, proc_cont, self.callexpr.arguments)
        return nextcont.start()

class ApplyProcCont(Cont):
    def __init__(self, Eval, env, cont, boundproc):
        Cont.__init__(self, Eval, env, cont)
        self.boundproc = boundproc

    def apply(self, args):
        # At this point operator and operands have been evaluated
        # So we need to do the "call" continuation
        # so each result is one "application"
        oldargs = currargs = args
        boundproc = self.boundproc
        procexpr, saved_env = boundproc.procexpr, boundproc.env
        if not currargs or not procexpr.varnames:
            assert False, "Called entry is *not* a function"

        nargs = len(procexpr.varnames)
        arglen = len(currargs)

        currargs,rest_args = currargs[:nargs], currargs[nargs:]
        newargs = dict(zip(procexpr.varnames, currargs))
        newenv = saved_env.extend(**newargs)

        if nargs > arglen:  # Time to curry
            left_varnames = procexpr.varnames[arglen:]
            newprocexpr = self.Eval.__caseon__.as_proc(left_varnames, procexpr.body).procexpr
            return self.nextcont.apply(newprocexpr.bind(newenv))

        elif nargs == arglen:
            return self.Eval(procexpr.body, newenv, self.nextcont)
        
        else: # nargs < arglen
            # Only take what we need and return rest as a call expr
            return self.Eval(procexpr, newenv, self.nextcont)
