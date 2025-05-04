
import typing
from taggedunion import *
from epl.chapter3 import letlang
from epl.chapter3 import proclang
from epl.chapter3 import letreclang

class NVar(object):
    """ A nameless variable. """
    def __init__(self, depth, index, lr_bound = False):
        self.depth = depth
        self.index = index
        self.lr_bound = lr_bound

    def printables(self):
        yield 0, "NVar (%d,%d,LR=%s)" % (self.depth, self.index, str(self.lr_bound))

    def __repr__(self):
        return "<NVar(%d:%d)>" % (self.depth, self.index)

class NProcExpr(object):
    def __init__(self, body):
        self.body = body

    def printables(self):
        yield 0, "NProc:"
        yield 1, self.body.printables()

    def __repr__(self):
        return "<NProc{ %s }" % repr(self.body)

class NLetExpr(object):
    """ Nameless let expressions.  No more mappings from var -> expr.  Only the "expr" is left and will be referred by indexes. """
    def __init__(self, values, body):
        self.values = values
        self.body = body

    def printables(self):
        yield 0, "NLet:"
        for v in self.values:
            yield 3, v.printables()
        yield 1, "in:"
        yield 2, self.body.printables()

class NLetRecExpr(object):
    """ To enable recursive procedures of the form:

        letrec F(x,y,z) = expression in expression
    """
    def __init__(self, proc_map, body):
        """ proc_map ::  : typing.Dict[str, (typing.List[str], "Expr")] """
        self.procs = {k: ProcExpr(k, v[0], v[1]) for k,v in proc_map.items()}
        self.body = body

    def printables(self):
        yield 0, "NLetRec:"
        for proc in self.procs.values():
            yield 2, "%s (%s) = " % (proc.name, ",".join(proc.varnames))
            yield 3, proc.body.printables()
        yield 1, "in:"
        yield 2, self.body.printables()

    def printables(self):
        yield 0, "NLet:"
        for v in self.values:
            yield 3, v.printables()
        yield 1, "in:"
        yield 2, self.body.printables()

class NExpr(Union):
    lit = Variant(letlang.Lit)
    iszero = Variant(letlang.IsZeroExpr)
    opexpr = Variant(letlang.OpExpr)
    tup = Variant(letlang.TupleExpr, checker = "is_tup", constructor = "as_tup")
    ifexpr = Variant(letlang.IfExpr, checker = "is_if", constructor = "as_if")
    callexpr = Variant(proclang.CallExpr, checker = "is_call", constructor = "as_call")

    # Nameless versions
    nvar = Variant(NVar)
    nlet = Variant(NLetExpr, checker = "is_nlet", constructor = "as_nlet")
    nletrec = Variant(NLetRecExpr)
    nprocexpr = Variant(NProcExpr, checker = "is_nproc", constructor = "as_nproc")

class Translator(CaseMatcher):
    """ Translates an expression with variables to a nameless one with depths (and indexes).  """
    __caseon__ = letreclang.Expr

    class StaticEnv(object):
        """ A static environment used in the translation of a named program to its nameless counterpart. """
        def __init__(self):
            self.nvars_stack = [[]]

        def push(self, *varnames):
            self.nvars_stack.insert(0, [])
            for varname in varnames: self.register(varname)

        def pop(self):
            self.nvars_stack.pop(0)

        def nvar_for(self, name):
            for depth,nvars in enumerate(self.nvars_stack):
                for index,vname in enumerate(nvars):
                    if vname == name:
                        return depth,index
            return None

        def register(self, varname):
            if varname in self.nvars_stack[0]:
                assert False, "Var already exists in this scope.  Push first?"
            self.nvars_stack[0].append(varname)

    def translate(self, expr, senv = None):
        senv = senv or Translator.StaticEnv()
        return self(expr, senv)

    @case("lit")
    def translateLit(self, lit, senv):
        return NExpr.as_lit(lit)

    @case("opexpr")
    def translateDiff(self, opexpr, senv):
        exprs = [self(exp, senv) for exp in opexpr.exprs]
        return NExpr.as_op(opexpr.op, exprs)

    @case("tup")
    def translateTupExpr(self, tup, senv):
        values = [self(v,env) for v in tup.children]
        return NExpr.as_tup(*values)

    @case("iszero")
    def translateIsZero(self, iszero, senv):
        return NExpr.as_iszero(self(iszero.expr, senv))

    @case("ifexpr")
    def translateIf(self, ifexpr, senv):
        return NExpr.as_if(self(ifexpr.exp1, senv), self(ifexpr.exp2, senv), self(ifexpr.exp3, senv))

    @case("var")
    def translateVar(self, var, senv):
        depth,index = senv.nvar_for(var.name)
        return NExpr.as_nvar(depth,index)

    @case("let")
    def translateLet(self, let, senv):
        keys = let.mappings.keys()
        newvalues = [self(let.mappings[k], senv) for k in keys]

        senv.push(*keys)
        newbody = self(let.body, senv)
        out = NExpr.as_nlet(newvalues, newbody)
        senv.pop()
        return out

    @case("letrec")
    def translateLetRec(self, letrec, senv):
        set_trace()
        newvalues = [self(v,senv) for v in let.values]
        newargs = let.mappings.keys()
        newenv = senv.push()
        map(newenv.register, let.mappings.keys())
        newbody = self(let.body, newenv)
        return NExpr.as_nlet(newvalues, newbody)

    @case("procexpr")
    def translateProc(self, procexpr, senv):
        newenv = env.push()
        map(newenv.register, procexpr.varnames)
        newbody = self(body, newenv)
        return Expr.as_nproc(newbody)

    @case("callexpr")
    def translateCall(self, callexpr, senv):
        newbody = self(body, senv)
        newargs = [self(arg, senv) for arg in callexpr.args]
        return Expr.as_call(newbody, *newargs)

class Eval(CaseMatcher):
    """ An evaluator on the nameless expressions. """
    __caseon__ = NExpr

    class NEnv(object):
        """ A nameless environment used for evaluations. """
        def __init__(self, parent = None):
            self.parent = parent
            self.values = []

        def get(self, nvar):
            curr = self
            i = 0
            while i < nvar.depth and curr:
                curr = curr.parent

            if not curr: return None
            return curr.values[nvar.index]

        def push(self):
            """ Create a new environment. """
            return self.__class__(self)

        def extend(self, *values):
            """ Create a new nenv by extending this with new vals. """
            out = self.push()
            out.values.extend(values)
            return out

    def valueOf(self, nexpr, env = None):
        env = env or Eval.NEnv()
        return self(nexpr, env)

    @case("lit")
    def valueOfLit(self, lit, env):
        return lit

    @case("opexpr")
    def valueOfOpExpr(self, opexpr, env):
        set_trace()

    @case("tup")
    def valueOfTupExpr(self, tup, env):
        values = [self(v,env) for v in tup.children]
        return TupleExpr(*values)

    @case("iszero")
    def valueOfIsZero(self, iszero, env):
        return self(iszero.expr, env) == 0

    @case("ifexpr")
    def valueOfIf(self, ifexpr, env):
        if self(ifexpr.cond, env):
            return self(ifexpr.exp1, env)
        else:
            return self(ifexpr.exp2, env)

    @case("nprocexpr")
    def valueOfProc(self, procexpr, env):
        return proclang.BoundProc(procexpr, env)

    @case("nvar")
    def valueOfNVar(self, nvar, env):
        return env.get(nvar)

    @case("callexpr")
    def valueOfCall(self, callexpr, env):
        boundproc = self(callexpr.operator, env)
        newargs = [self(arg, env) for arg in callexpr.args]
        return self.apply_proc(boundproc, args)

    def apply_proc(self, boundproc, args):
        procexpr, saved_env = boundproc.procexpr, boundproc.env
        newenv = saved_env.push.extend(args)
        return self(procexpr.body, newenv)

    @case("nlet")
    def valueOfLet(self, let, env):
        newvalues = [self(v,env) for v in let.values]
        newenv = env.extend(*newvalues)
        return self(let.body, newenv)

    @case("nletrec")
    def valueOfLetRec(self, letrec, senv):
        newvalues = [self(v,senv) for v in let.values]
        newargs = let.mappings.keys()
        newenv = senv.push()
        map(newenv.register, let.mappings.keys())
        newbody = self(let.body, newenv)
        return Expr.as_nletrec(newvalues, newbody)
