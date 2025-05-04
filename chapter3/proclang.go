package chapter3

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	epl "github.com/panyam/eplgo"
	gfn "github.com/panyam/goutils/fn"
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

func Call(operator any, args ...any) *CallExpr {
	return &CallExpr{Operator: AnyToExpr(operator), Args: gfn.Map(args, AnyToExpr)}
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
	currProcexpr, currEnv := boundproc.ProcExpr, boundproc.Env
	currArgs := args
	var result any
	initialCall := true

	for { // Loop until explicitly returned or error
		numParams := len(currProcexpr.Varnames)
		numArgVals := len(currArgs)

		// log.Printf("applyProc Loop Start: Proc(%v), Args: %v, Env: %s\n", currProcexpr.Varnames, currArgs, currEnv)

		if numParams == 0 {
			// If proc takes 0 params, evaluate its body.
			// It *must not* be called with arguments.
			if numArgVals > 0 {
				panic(fmt.Sprintf("Procedure %s takes 0 arguments, but called with %d arguments: %v", currProcexpr.Repr(), numArgVals, currArgs))
			}
			// log.Println("Proc takes 0 params, evaluating body")
			result = l.Eval(currProcexpr.Body, currEnv)
			// If the body returned *another* 0-arg proc, we need to evaluate that too.
			// This handles chains like `proc() proc() 5`
			if bp, ok := result.(*BoundProc); ok && len(bp.ProcExpr.Varnames) == 0 {
				// log.Println("Body returned another 0-arg proc, continuing")
				currProcexpr = bp.ProcExpr
				currEnv = bp.Env
				// currArgs remains []
				initialCall = false
				continue // Re-evaluate the new 0-arg proc
			} else {
				// log.Println("Returning result from 0-arg proc body")
				return result // Final value or a proc requiring args
			}
		}

		// If we have a proc expecting params, but no args left, it means we have a partial application.
		if numArgVals == 0 {
			if initialCall {
				// This should not happen if Call() is constructed correctly.
				panic(fmt.Sprintf("Initial call to Proc(%v) with no arguments.", currProcexpr.Varnames))
			} else {
				// We consumed args in previous iterations, now none left. Return the current proc bound to its env.
				// log.Printf("No more args, returning partially applied Proc(%v)\n", currProcexpr.Varnames)
				return currProcexpr.Bind(currEnv) // Return the *current* bound proc
			}
		}

		// --- Consume arguments ---
		initialCall = false // An application is happening
		maxArgs := min(numArgVals, numParams)
		consumedArgs, restArgs := currArgs[:maxArgs], currArgs[maxArgs:]
		// Only map params that are being consumed in this step
		newArgsMap := epl.DictZip(currProcexpr.Varnames[:maxArgs], consumedArgs)
		newenv := currEnv.Extend(newArgsMap)
		// log.Printf("Consumed %d args (%v), %d remaining (%v). New Env: %s\n", maxArgs, consumedArgs, len(restArgs), restArgs, newenv)

		if numParams > numArgVals { // Curry: Not enough args provided in this call
			// The procedure expects more arguments than were supplied *in this chunk*.
			leftVarnames := currProcexpr.Varnames[numArgVals:] // Params not covered by current args
			newprocexpr := Proc(leftVarnames, currProcexpr.Body)
			// log.Printf("Currying: Returning Proc(%v) bound to env %s\n", leftVarnames, newenv)
			// The environment *must* include the args just consumed.
			return newprocexpr.Bind(newenv) // Return the new curried proc

		} else { // Exact match (numParams == numArgVals) OR More args than params (numParams < numArgVals)
			// We have enough (or more) arguments to satisfy the current procedure's parameters.
			// log.Printf("Evaluating body of Proc(%v) with env %s\n", currProcexpr.Varnames, newenv)
			result = l.Eval(currProcexpr.Body, newenv) // Evaluate body with the consumed args
			// log.Printf("Body evaluation returned: %v (%T)\n", result, result)

			if bp, ok := result.(*BoundProc); ok {
				// Body returned another procedure. Continue the loop with this new proc and remaining args.
				// log.Println("Body returned another proc, continuing loop")
				currProcexpr = bp.ProcExpr
				currEnv = bp.Env
				currArgs = restArgs // Use remaining args for the *new* proc
				// Loop continues without returning here
			} else {
				// Body returned a non-procedure value.
				if len(restArgs) == 0 {
					// No arguments left, this value is the final result.
					// log.Println("Body returned value and no args left. Returning final result.")
					return result
				} else {
					// Body returned a value, but we still have args left. This is an error.
					panic(fmt.Sprintf("Procedure %s returned non-procedure value %v (%T), but %d arguments remain: %v", currProcexpr.Repr(), result, result, len(restArgs), restArgs))
				}
			}
		}
		// If we didn't return, the loop continues with updated currProcexpr, currEnv, currArgs
	} // End for loop

}

func (l *ProcLangEval) applyProcOld(boundproc *BoundProc, args []any) any {
	currProcexpr, currEnv := boundproc.ProcExpr, boundproc.Env
	currArgs := args
	var restArgs []any

	numParams := len(currProcexpr.Varnames)
	numArgVals := len(currArgs)
	for numArgVals > 0 && numParams > 0 {
		maxargs := min(numArgVals, numParams)
		currArgs, restArgs = currArgs[:maxargs], currArgs[maxargs:]
		newargs := epl.DictZip(currProcexpr.Varnames, currArgs)
		newenv := currEnv.Extend(newargs)
		if numParams > numArgVals { // Time to curry
			leftVarnames := currProcexpr.Varnames[numArgVals:]
			newprocexpr := Proc(leftVarnames, currProcexpr.Body)
			return newprocexpr.Bind(newenv)
		} else if numParams == numArgVals {
			return l.Eval(currProcexpr.Body, newenv)
		} else { // numParams < numArgVals
			// Only take what we need and return rest as a call expr
			// TODO - check types
			log.Println("Before : ", currProcexpr, newenv.String())
			res := l.Eval(currProcexpr, newenv)

			log.Println("Result: ", res, reflect.TypeOf(res))
			currProcexpr = res.(*BoundProc).ProcExpr
		}

		// after all case
		currArgs = restArgs
		numParams = len(currProcexpr.Varnames)
		numArgVals = len(currArgs)
	}

	// Check atleast one application has happened
	if currProcexpr == boundproc.ProcExpr {
		panic("Called entry is *not* a function")
	}
	return nil
}
