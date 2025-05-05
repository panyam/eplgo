// ./chapter4/eval.go
package chapter4

import (
	"fmt"
	"log"

	epl "github.com/panyam/eplgo"
	"github.com/panyam/eplgo/chapter3"
)

// AssignExpr represents the 'set var = expr' operation.
type AssignExpr struct {
	Varname string        // Name of the variable to assign to.
	Expr    chapter3.Expr // The expression providing the new value.
}

func Assign(varname string, expr any) *AssignExpr {
	return &AssignExpr{
		Varname: varname,
		Expr:    chapter3.AnyToExpr(expr),
	}
}

func (e *AssignExpr) Printable() *epl.Printable {
	return epl.PrintableIter(func(yield func(v *epl.Printable) bool) {
		// Example Format: "Assign: varname = <Expr>"
		if !yield(epl.Printablef(0, "Assign: %s =", e.Varname)) {
			return
		}
		// Value expr at level 1 relative to Assign node
		vP := e.Expr.Printable()
		vP.IndentLevel += 1 // Indent value by 1 level
		if !yield(vP) {
			return
		}
	})
}

func (e *AssignExpr) Repr() string {
	return fmt.Sprintf("<Assign(%s = %s)>", e.Varname, e.Expr.Repr())
}

func (e *AssignExpr) Eq(another *AssignExpr) bool {
	// Compare variable names and the expression structure
	return e.Varname == another.Varname &&
		chapter3.ExprEq(e.Expr, another.Expr)
}

// --- Make sure chapter3.ExprEq handles this via reflection ---
// No changes needed in chapter3/expr.go if using reflection approach.

// ImpRefLangEval evaluates expressions including implicit variable assignment.
type ImpRefLangEval struct {
	ExpRefLangEval // Embed the previous evaluator
}

// NewImpRefLangEval creates a new evaluator for the impref language.
func NewImpRefLangEval() *ImpRefLangEval {
	out := &ImpRefLangEval{}
	// CRITICAL: Set the Self pointer for the embedded BaseEval
	out.BaseEval.Self = out
	return out
}

// LocalEval handles expression types specific to ImpRefLang or delegates.
func (l *ImpRefLangEval) LocalEval(expr chapter3.Expr, env *epl.Env[any]) (any, error) {
	// log.Printf("ImpRefLangEval evaluating: %s (%T)\n", expr.Repr(), expr)
	switch n := expr.(type) {
	case *AssignExpr: // Handle the new type
		return l.valueOfAssign(n, env)
	default:
		// Delegate to the embedded ExpRefLangEval's LocalEval for other types
		return l.ExpRefLangEval.LocalEval(expr, env)
	}
}

// valueOfAssign handles 'set var = expr'
func (l *ImpRefLangEval) valueOfAssign(e *AssignExpr, env *epl.Env[any]) (any, error) {
	// Evaluate the expression for the new value
	newValue, err := l.Eval(e.Expr, env)
	if err != nil {
		return nil, err
	}

	// Get the *reference* associated with the variable name from the environment.
	// env.Get looks up the value *inside* the ref. We need the ref itself.
	varRef := env.GetRef(e.Varname) // Use GetRef which finds the *epl.Ref[any]

	// Check if the variable exists (i.e., if GetRef found it)
	if varRef == nil {
		log.Panicf("set: variable '%s' not found in environment", e.Varname)
	}

	// log.Printf("assign evaluated %s to %v. Found Ref %p for var %s. Updating ref.\n", e.Expr.Repr(), newValue, varRef, e.Varname)

	// Update the value *inside* the existing reference cell for the variable
	varRef.Value = newValue

	// 'set' returns the new value
	return newValue, nil
}
