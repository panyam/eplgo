

from taggedunion import *
from epl.chapter3 import letlang

## Constructs for Procedures 

class ProcExpr(object):
    class BoundProc(object):
        """ A procedure bound to an environment. """
        def __init__(self, proc, env):
            self.procexpr = proc
            self.env = env

    def __init__(self, varnames, body):
        if type(varnames) is str: varnames = [varnames]
        self.name = None
        self.varnames = varnames
        self.body = body

    def bind(self, env):
        return ProcExpr.BoundProc(self, env)

    def printables(self):
        if self.name:
            yield 0, "Proc %s (%s) = " % (self.name, ", ".join(self.varnames))
        else:
            yield 0, "Proc (%s) = " % ", ".join(self.varnames)
        yield 1, self.body.printables()

    def __eq__(self, another):
        return self.varnames == another.varnames and \
                self.name == another.name and \
                self.body == another.body

    def __repr__(self):
        if self.name:
            return "<Proc %s (%s) { %s }" % (self.name, ", ".join(self.varnames), repr(self.body))
        else:
            return "<Proc(%s) { %s }" % (", ".join(self.varnames), repr(self.body))

class CallExpr(object):
    def __init__(self, operator, *arguments):
        self.operator = operator
        self.arguments = list(arguments)

    def printables(self):
        yield 0, "Call"
        yield 1, "Operator:"
        yield 2, self.operator.printables()
        yield 1, "Args:"
        for arg in self.arguments:
            yield 2, arg.printables()

    def __eq__(self, another):
        s1,s2 = set(self.varnames), set(another.varnames)
        return s1 == s2 and self.name == another.name and \
               self.body == another.body

    def __eq__(self, another):
        if self.operator != another.operator:
            return False
        if len(self.arguments) != len(another.arguments):
            return False
        for e1,e2 in zip(self.arguments, another.arguments):
            if e1 != e2:
                return False
        return True

    def __repr__(self):
        return "<Call (%s) in %s" % (self.operator, ", ".join(map(repr, self.arguments)))


class Expr(letlang.Expr):
    procexpr = Variant(ProcExpr, checker = "is_proc", constructor = "as_proc")
    callexpr = Variant(CallExpr, checker = "is_call", constructor = "as_call")

class Eval(letlang.Eval):
    __caseon__ = Expr

    @case("procexpr")
    def valueOfProc(self, procexpr, env):
        return procexpr.bind(env)

    @case("callexpr")
    def valueOfCall(self, callexpr, env):
        boundproc = self.valueOf(callexpr.operator, env)
        arguments = [self.valueOf(arg, env) for arg in callexpr.arguments]
        return self.apply_proc(boundproc, arguments)

    def apply_proc(self, boundproc, arguments):
        procexpr, saved_env = boundproc.procexpr, boundproc.env
        curr_procexpr = procexpr
        curr_env = saved_env
        curr_args = arguments

        while curr_args and curr_procexpr.varnames:
            nargs = len(curr_procexpr.varnames)
            arglen = len(curr_args)

            curr_args,rest_args = curr_args[:nargs], curr_args[nargs:]
            newargs = dict(zip(curr_procexpr.varnames, curr_args))
            newenv = curr_env.extend(**newargs)
            if nargs > arglen:  # Time to curry
                left_varnames = curr_procexpr.varnames[arglen:]
                newprocexpr = ProcExpr(left_varnames, curr_procexpr.body)
                return newprocexpr.bind(newenv)

            elif nargs == arglen:
                return self.valueOf(curr_procexpr.body, newenv)
            
            else: # nargs < arglen
                # Only take what we need and return rest as a call expr
                curr_procexpr = self.valueOf(curr_procexpr, newenv)
                # Should curr_env = newenv?
                # return self.apply_proc(newexpr, arguments[nargs:])

            # after all case
            curr_args = rest_args

        # Check atleast one application has happened
        assert curr_procexpr != procexpr, "Called entry is *not* a function"
