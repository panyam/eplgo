package chapter4

import (
	"fmt"
	"log"

	// For deterministic printing if needed
	epl "github.com/panyam/eplgo" // Import chapter3 for the base Expr
	"github.com/panyam/eplgo/chapter3"
	gfn "github.com/panyam/goutils/fn" // For Map
)

// RefExpr represents the 'newref' operation.
// Note: The Python 'RefExpr' covered both 'newref' and variable references ('ref x').
// In Go, 'newref x' will likely be an AST node, while 'ref x' for Call-by-Reference
// might be handled differently if implemented (perhaps via evaluator logic or a different node).
// Let's focus on 'newref' as defined by its evaluation semantics first.
// This AST node corresponds to `newref(expr)`.
type RefExpr struct {
	// If IsVarRef is true, ExprOrVar contains the variable name (string).
	// If IsVarRef is false, ExprOrVar contains the expression for newref (Expr).
	ExprOrVar any
	IsVarRef  bool
}

func NewRef(expr any) *RefExpr {
	return &RefExpr{
		ExprOrVar: AnyToExpr(expr),
		IsVarRef:  false,
	}
}

// Constructor for 'ref varname'
func RefVar(varname string) *RefExpr {
	// Store the varname directly
	return &RefExpr{
		ExprOrVar: varname,
		IsVarRef:  true,
	}
}

func (e *RefExpr) Printable() *epl.Printable {
	if e.IsVarRef {
		return epl.Printablef(0, "RefVar: %s", e.ExprOrVar.(string))
	} else {
		// Reuse PrintableIter logic from previous NewRef
		return epl.PrintableIter(func(yield func(v *epl.Printable) bool) {
			if !yield(epl.Printablef(0, "NewRef:")) {
				return
			}
			cP := e.ExprOrVar.(Expr).Printable()
			cP.IndentLevel += 1
			if !yield(cP) {
				return
			}
		})
	}
}

func (e *RefExpr) Repr() string {
	if e.IsVarRef {
		return fmt.Sprintf("<RefVar(%s)>", e.ExprOrVar.(string))
	} else {
		return fmt.Sprintf("<NewRef(%s)>", e.ExprOrVar.(Expr).Repr())
	}
}

func (e *RefExpr) Eq(another *RefExpr) bool {
	if e.IsVarRef != another.IsVarRef {
		return false
	}
	if e.IsVarRef {
		// Both are RefVar, compare varnames
		return e.ExprOrVar.(string) == another.ExprOrVar.(string)
	} else {
		// Both are NewRef, compare expressions
		return ExprEq(e.ExprOrVar.(Expr), another.ExprOrVar.(Expr))
	}
}

// DeRefExpr represents the 'deref' operation.
type DeRefExpr struct {
	RefExpr Expr // The expression that should evaluate to a reference (*epl.Ref[any]).
}

func DeRef(refExpr any) *DeRefExpr {
	return &DeRefExpr{RefExpr: AnyToExpr(refExpr)}
}

func (e *DeRefExpr) Printable() *epl.Printable {
	return epl.PrintableIter(func(yield func(v *epl.Printable) bool) {
		if !yield(epl.Printablef(0, "DeRef:")) {
			return
		}
		cP := e.RefExpr.Printable()
		cP.IndentLevel += 1
		if !yield(cP) {
			return
		}
	})
}

func (e *DeRefExpr) Repr() string {
	return fmt.Sprintf("<DeRef(%s)>", e.RefExpr.Repr())
}

func (e *DeRefExpr) Eq(another *DeRefExpr) bool {
	return ExprEq(e.RefExpr, another.RefExpr)
}

// SetRefExpr represents the 'setref' operation.
type SetRefExpr struct {
	RefExpr   Expr // The expression that should evaluate to a reference (*epl.Ref[any]).
	ValueExpr Expr // The expression providing the new value.
}

func SetRef(refExpr, valueExpr any) *SetRefExpr {
	return &SetRefExpr{
		RefExpr:   AnyToExpr(refExpr),
		ValueExpr: AnyToExpr(valueExpr),
	}
}

func (e *SetRefExpr) Printable() *epl.Printable {
	return epl.PrintableIter(func(yield func(v *epl.Printable) bool) {
		if !yield(epl.Printablef(0, "SetRef:")) {
			return
		}
		if !yield(epl.Printablef(1, "Ref:")) {
			return
		}
		rP := e.RefExpr.Printable()
		rP.IndentLevel += 2
		if !yield(rP) {
			return
		}
		if !yield(epl.Printablef(1, "Value:")) {
			return
		}
		vP := e.ValueExpr.Printable()
		vP.IndentLevel += 2
		if !yield(vP) {
			return
		}
	})
}

func (e *SetRefExpr) Repr() string {
	return fmt.Sprintf("<SetRef(%s, %s)>", e.RefExpr.Repr(), e.ValueExpr.Repr())
}

func (e *SetRefExpr) Eq(another *SetRefExpr) bool {
	return ExprEq(e.RefExpr, another.RefExpr) &&
		ExprEq(e.ValueExpr, another.ValueExpr)
}

// BlockExpr represents the 'begin ... end' sequence.
type BlockExpr struct {
	Exprs []Expr
}

func Begin(exprs ...any) *BlockExpr {
	// Convert anys to Exprs
	goExprs := gfn.Map(exprs, AnyToExpr)
	return &BlockExpr{Exprs: goExprs}
}

func (e *BlockExpr) Printable() *epl.Printable {
	return epl.PrintableIter(func(yield func(v *epl.Printable) bool) {
		if !yield(epl.Printablef(0, "Begin:")) {
			return
		}
		for _, expr := range e.Exprs {
			cP := expr.Printable()
			cP.IndentLevel += 1
			if !yield(cP) {
				return
			}
		}
		// Optional: Add an "End:" indicator? The Python didn't explicitly, but might be clearer.
		// if !yield(epl.Printablef(0, "End:")) { return }
	})
}

func (e *BlockExpr) Repr() string {
	return fmt.Sprintf("<Begin(%s)>", ExprListRepr(e.Exprs))
}

func (e *BlockExpr) Eq(another *BlockExpr) bool {
	return ExprListEq(e.Exprs, another.Exprs)
}

// Ensure ExprEq knows about these types
// We need to modify chapter3/expr.go later for this.

// ExpRefLangEval evaluates expressions including explicit references.
type ExpRefLangEval struct {
	chapter3.LetRecLangEval // Embed the previous evaluator
}

// NewExpRefLangEval creates a new evaluator for the expref language.
func NewExpRefLangEval() *ExpRefLangEval {
	out := &ExpRefLangEval{}
	// CRITICAL: Set the Self pointer for the embedded BaseEval
	out.BaseEval.Self = out
	return out
}

// LocalEval handles expression types specific to ExpRefLang or delegates.
func (l *ExpRefLangEval) LocalEval(expr Expr, env *epl.Env[any]) any {
	// log.Printf("ExpRefLangEval evaluating: %s (%T)\n", expr.Repr(), expr)
	switch n := expr.(type) {
	case *RefExpr:
		return l.valueOfRef(n, env)
	case *DeRefExpr:
		return l.valueOfDeRef(n, env)
	case *SetRefExpr:
		return l.valueOfSetRef(n, env)
	case *BlockExpr:
		return l.valueOfBlock(n, env)
	default:
		// Delegate to the embedded LetRecLangEval's LocalEval for other types
		return l.LetRecLangEval.LocalEval(expr, env)
	}
}

// --- Implement valueOf... methods ---

func (l *ExpRefLangEval) valueOfBlock(e *BlockExpr, env *epl.Env[any]) any {
	var result any
	// Default result if block is empty (or consider panic/error?)
	// Python returned Lit(0). Let's return nil for now, maybe change to Lit(0) or Void later.
	result = nil // Or perhaps Lit(0)? Check Python behavior/tests.
	// Python: value = self.__caseon__.as_lit(0); result = value
	// Let's mimic Python for now.
	result = Lit(0)

	for _, expr := range e.Exprs {
		result = l.Eval(expr, env) // Use Eval to allow dispatch back to Self
	}
	return result // Return the value of the last expression
}

func (l *ExpRefLangEval) valueOfRef(e *RefExpr, env *epl.Env[any]) any {
	if e.IsVarRef {
		// Mode: 'ref varname'
		varname := e.ExprOrVar.(string)
		varRef := env.GetRef(varname) // Get the *epl.Ref[any] itself
		if varRef == nil {
			log.Panicf("ref: variable '%s' not found in environment", varname)
		}
		// log.Printf("ref var evaluated %s to Ref %p\n", varname, varRef)
		return varRef // Return the existing reference
	} else {
		// Mode: 'newref(expr)'
		initialValueExpr := e.ExprOrVar.(Expr)
		initialValue := l.Eval(initialValueExpr, env)
		// Create a *new* reference cell containing this value
		newRef := &epl.Ref[any]{Value: initialValue}
		// log.Printf("newref evaluated %s to %v, created Ref %p\n", initialValueExpr.Repr(), initialValue, newRef)
		return newRef // Return the pointer to the new Ref struct
	}
}

func (l *ExpRefLangEval) valueOfDeRef(e *DeRefExpr, env *epl.Env[any]) any {
	// Evaluate the expression which *should* yield a reference
	refVal := l.Eval(e.RefExpr, env)

	// Check if the result is actually a reference (*epl.Ref[any])
	theRef, ok := refVal.(*epl.Ref[any])
	if !ok {
		log.Panicf("deref expected a reference argument, but got type %T for expr %s", refVal, e.RefExpr.Repr())
	}
	// log.Printf("deref evaluated %s to Ref %p, returning Value %v\n", e.RefExpr.Repr(), theRef, theRef.Value)
	// Return the value *inside* the reference cell
	return theRef.Value
}

func (l *ExpRefLangEval) valueOfSetRef(e *SetRefExpr, env *epl.Env[any]) any {
	// Evaluate the expression which *should* yield a reference
	refVal := l.Eval(e.RefExpr, env)
	theRef, ok := refVal.(*epl.Ref[any])
	if !ok {
		log.Panicf("setref expected a reference argument for the first expression, but got type %T for expr %s", refVal, e.RefExpr.Repr())
	}

	// Evaluate the expression for the new value
	newValue := l.Eval(e.ValueExpr, env)

	// log.Printf("setref evaluated %s to Ref %p, evaluated %s to %v. Updating ref.\n", e.RefExpr.Repr(), theRef, e.ValueExpr.Repr(), newValue)

	// Update the value inside the reference cell
	theRef.Value = newValue

	// setref returns the new value
	return newValue
}
