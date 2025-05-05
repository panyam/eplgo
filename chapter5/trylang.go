package chapter5

import (
	"errors"
	"fmt"

	// "sort" // Not needed here yet
	// "strings" // Not needed here yet

	epl "github.com/panyam/eplgo"
	"github.com/panyam/eplgo/chapter3"
	"github.com/panyam/eplgo/chapter4"
	// Import chapter3 for the base Expr
	// gfn "github.com/panyam/goutils/fn"  // Not needed here yet
)

// TryExpr represents the 'try E catch (x) H' construct.
type TryExpr struct {
	TryBody     Expr   // The expression E to try.
	VarName     string // The variable name 'x' to bind the exception value to.
	HandlerExpr Expr   // The handler expression H.
}

func Try(tryBody any, varName string, handlerExpr any) *TryExpr {
	return &TryExpr{
		TryBody:     AnyToExpr(tryBody),
		VarName:     varName,
		HandlerExpr: AnyToExpr(handlerExpr),
	}
}

func (e *TryExpr) Printable() *epl.Printable {
	return epl.PrintableIter(func(yield func(v *epl.Printable) bool) {
		if !yield(epl.Printablef(0, "Try:")) {
			return
		}
		// Try Body
		tP := e.TryBody.Printable()
		tP.IndentLevel += 1
		if !yield(tP) {
			return
		}
		// Catch clause
		if !yield(epl.Printablef(1, "Catch (%s):", e.VarName)) {
			return
		}
		// Handler Body
		hP := e.HandlerExpr.Printable()
		hP.IndentLevel += 2 // Indent handler under catch
		if !yield(hP) {
			return
		}
	})
}

func (e *TryExpr) Repr() string {
	return fmt.Sprintf("<Try(%s) Catch(%s -> %s)>",
		e.TryBody.Repr(),
		e.VarName,
		e.HandlerExpr.Repr())
}

func (e *TryExpr) Eq(another *TryExpr) bool {
	return e.VarName == another.VarName &&
		ExprEq(e.TryBody, another.TryBody) &&
		ExprEq(e.HandlerExpr, another.HandlerExpr)
}

// RaiseExpr represents the 'raise E' construct.
type RaiseExpr struct {
	RaiseValueExpr Expr // The expression E whose value is raised.
}

func Raise(valueExpr any) *RaiseExpr {
	return &RaiseExpr{RaiseValueExpr: AnyToExpr(valueExpr)}
}

func (e *RaiseExpr) Printable() *epl.Printable {
	return epl.PrintableIter(func(yield func(v *epl.Printable) bool) {
		if !yield(epl.Printablef(0, "Raise:")) {
			return
		}
		vP := e.RaiseValueExpr.Printable()
		vP.IndentLevel += 1
		if !yield(vP) {
			return
		}
	})
}

func (e *RaiseExpr) Repr() string {
	return fmt.Sprintf("<Raise(%s)>", e.RaiseValueExpr.Repr())
}

func (e *RaiseExpr) Eq(another *RaiseExpr) bool {
	return ExprEq(e.RaiseValueExpr, another.RaiseValueExpr)
}

// TryLangEval evaluates expressions including try/catch and raise.
type TryLangEval struct {
	chapter4.LazyLangEval // Embed the previous evaluator
}

// NewTryLangEval creates a new evaluator for the Chapter 5 language.
func NewTryLangEval() *TryLangEval {
	out := &TryLangEval{}
	// CRITICAL: Set the Self pointer for the embedded BaseEval
	out.BaseEval.Self = out
	return out
}

// LocalEval handles expression types specific to Chapter 5 or delegates.
func (l *TryLangEval) LocalEval(expr chapter3.Expr, env *epl.Env[any]) (any, error) {
	// log.Printf("TryLangEval evaluating: %s (%T)\n", expr.Repr(), expr)
	switch n := expr.(type) {
	case *TryExpr:
		return l.valueOfTry(n, env)
	case *RaiseExpr:
		return l.valueOfRaise(n, env)
	default:
		// Delegate to the embedded LazyLangEval's LocalEval for other types
		return l.LazyLangEval.LocalEval(expr, env)
	}
}

// valueOfRaise handles 'raise E'.
func (l *TryLangEval) valueOfRaise(e *RaiseExpr, env *epl.Env[any]) (any, error) {
	// 1. Evaluate the expression E whose value will be raised.
	raisedValue, err := l.Eval(e.RaiseValueExpr, env)
	if err != nil {
		// If evaluating the value itself causes an error, propagate that error.
		return nil, fmt.Errorf("evaluating expression for raise: %w", err)
	}

	// 2. Wrap the evaluated value in our custom RaisedError type.
	//    Return nil value and the RaisedError.
	// log.Printf("Raising value %v (%T)\n", raisedValue, raisedValue)
	return nil, RaisedError{Value: raisedValue}
}

// valueOfTry handles 'try E catch (x) H'.
func (l *TryLangEval) valueOfTry(e *TryExpr, env *epl.Env[any]) (any, error) {
	// 1. Evaluate the 'try' body (E).
	// log.Printf("Entering try block for: %s\n", e.TryBody.Repr())
	tryResult, tryErr := l.Eval(e.TryBody, env)

	// 2. Check the error returned from evaluating the body.
	if tryErr == nil {
		// No error occurred in the try block. Return its result directly.
		// log.Printf("Try block completed normally with result: %v (%T)\n", tryResult, tryResult)
		return tryResult, nil
	} else {
		// An error occurred. Check if it's a RaisedError we should catch.
		var raisedErr RaisedError
		// Use errors.As for robust checking, even if error is wrapped.
		if errors.As(tryErr, &raisedErr) {
			// It's a raised exception that we might handle.
			// log.Printf("Caught raised error in try block: %v\n", raisedErr)
			// log.Printf("Binding '%s' to value %v (%T)\n", e.VarName, raisedErr.Value, raisedErr.Value)

			// 3. Create a new environment for the handler, extending the current one.
			//    Bind the caught value (raisedErr.Value) to the specified variable name (e.VarName).
			handlerEnv := env.Extend(epl.Dict[string, any](e.VarName, raisedErr.Value))

			// 4. Evaluate the handler expression (H) in the new environment.
			// log.Printf("Evaluating handler: %s\n", e.HandlerExpr.Repr())
			// The result of the entire 'try' expression is the result of the handler.
			// Propagate any result or error from the handler evaluation.
			return l.Eval(e.HandlerExpr, handlerEnv)
		} else {
			// It's some other kind of error (e.g., variable not found, op type error).
			// This 'try/catch' does not handle it. Propagate the original error.
			// log.Printf("Propagating non-raised error through try block: %v\n", tryErr)
			return nil, tryErr
		}
	}
}
