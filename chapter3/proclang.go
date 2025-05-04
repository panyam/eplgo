package chapter3

import (
	"fmt"
	"strings"

	epl "github.com/panyam/eplgo"
)

// Constructs for Procedures
type BoundProc struct {
	ProcExpr *ProcExpr
	Env      *epl.Env[any]
}

type ProcExpr struct {
	Name     string
	Varnames []string
	Body     Expr
}

func Proc(varnames []string, body Expr) *ProcExpr {
	return &ProcExpr{Varnames: varnames, Body: body}
}

func (v *ProcExpr) Printable() *epl.Printable {
	return epl.PrintableIter(func(yield func(v *epl.Printable) bool) {
		if v.Name != "" {
			if !yield(epl.Printablef(0, "Proc %s (%s) = ", v.Name, strings.Join(v.Varnames, ", "))) {
				return
			}
		} else {
			if !yield(epl.Printablef(0, "Proc (%s) = ", strings.Join(v.Varnames, ", "))) {
				return
			}
		}

		yield(v.Body.Printable())
	})
}

func (v *ProcExpr) Eq(another *ProcExpr) bool {
	// TODO - Also check varnames
	if v.Name != another.Name || !epl.StringListEq(v.Varnames, another.Varnames) {
		return false
	}
	return ExprEq(v.Body, another.Body)
}

func (v *ProcExpr) Repr() string {
	if v.Name != "" {
		return fmt.Sprintf("<Proc %s (%s) { %s }", v.Name, strings.Join(v.Varnames, ", "), v.Body.Repr())
	} else {
		return fmt.Sprintf("<Proc (%s) { %s }", strings.Join(v.Varnames, ", "), v.Body.Repr())
	}
}

func (v *ProcExpr) Bind(env *epl.Env[any]) *BoundProc {
	return &BoundProc{v, env}
}

type CallExpr struct {
	Operator Expr
	Args     []Expr
}

func Call(operator Expr, args ...Expr) *CallExpr {
	return &CallExpr{Operator: operator, Args: args}
}

func (v *CallExpr) Printable() *epl.Printable {
	return epl.PrintableIter(func(yield func(v *epl.Printable) bool) {
		if !yield(&epl.Printable{0, "Call", nil}) {
			return
		}
		if !yield(&epl.Printable{1, "Operator", nil}) {
			return
		}
		if !yield(v.Operator.Printable()) {
			return
		}
		if !yield(&epl.Printable{1, "Args", nil}) {
			return
		}
		ExprListPrintable(2, v.Args, yield)
	})
}

func (v *CallExpr) Eq(another *CallExpr) bool {
	// TODO - Also check varnames
	return ExprEq(v.Operator, another.Operator) && ExprListEq(v.Args, another.Args)
}

func (v *CallExpr) Repr() string {
	return fmt.Sprintf("<Call (%s) in %s", v.Operator.Repr(), ExprListRepr(v.Args))
}

type ProcLangEval struct {
	LetLangEval
}

func NewProcLangEval() *ProcLangEval {
	out := &ProcLangEval{}
	out.BaseEval.Self = out
	return out
}

func (l *ProcLangEval) LocalEval(expr Expr, env *epl.Env[any]) any {
	// log.Println("ProcLangEval for: ", reflect.TypeOf(expr), expr.Repr())
	switch n := expr.(type) {
	case *ProcExpr:
		return l.ValueOfProc(n, env)
	case *CallExpr:
		return l.ValueOfCall(n, env)
	default:
		// Call super method
		// log.Printf("Calling super for type: %v - %v", n, reflect.TypeOf(n))
		return l.LetLangEval.LocalEval(expr, env)
	}
}

func (l *ProcLangEval) ValueOfProc(e *ProcExpr, env *epl.Env[any]) any {
	return e.Bind(env)
}

func (l *ProcLangEval) ValueOfCall(e *CallExpr, env *epl.Env[any]) any {
	boundproc := l.Eval(e.Operator, env).(*BoundProc)
	args := l.EvalExprList(e.Args, env)
	return l.applyProc(boundproc, args)
}

func (l *ProcLangEval) applyProc(boundproc *BoundProc, args []any) any {
	procexpr, savedEnv := boundproc.ProcExpr, boundproc.Env
	currProcexpr := procexpr
	currEnv := savedEnv
	currArgs := args
	var restArgs []any

	for len(currArgs) > 0 && len(currProcexpr.Varnames) > 0 {
		nargs := len(currProcexpr.Varnames)
		arglen := len(currArgs)

		currArgs, restArgs = currArgs[:nargs], currArgs[nargs:]
		newargs := epl.DictZip(currProcexpr.Varnames, currArgs)
		newenv := currEnv.Extend(newargs)
		if nargs > arglen { // Time to curry
			leftVarnames := currProcexpr.Varnames[arglen:]
			newprocexpr := Proc(leftVarnames, currProcexpr.Body)
			return newprocexpr.Bind(newenv)
		} else if nargs == arglen {
			return l.Eval(currProcexpr.Body, newenv)
		} else { // nargs < arglen
			// Only take what we need and return rest as a call expr
			// TODO - check types
			currProcexpr = l.Eval(currProcexpr, newenv).(*ProcExpr)
		}

		// after all case
		currArgs = restArgs
	}

	// Check atleast one application has happened
	if currProcexpr == procexpr {
		panic("Called entry is *not* a function")
	}
	return nil
}
