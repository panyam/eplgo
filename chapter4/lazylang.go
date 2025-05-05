package chapter4

import (
	"fmt"
	"log"

	epl "github.com/panyam/eplgo"
)

// --- Existing Structs (RefExpr, DeRefExpr, SetRefExpr, BlockExpr, AssignExpr) remain ---
// ...

// LazyExpr represents the 'lazy <expr>' construct.
// Its evaluation results in a Thunk value.
type LazyExpr struct {
	Expr Expr // The expression to be evaluated lazily.
}

func Lazy(expr any) *LazyExpr {
	return &LazyExpr{Expr: AnyToExpr(expr)}
}

func (e *LazyExpr) Printable() *epl.Printable {
	return epl.PrintableIter(func(yield func(v *epl.Printable) bool) {
		if !yield(epl.Printablef(0, "Lazy:")) {
			return
		}
		cP := e.Expr.Printable()
		cP.IndentLevel += 1
		if !yield(cP) {
			return
		}
	})
}

func (e *LazyExpr) Repr() string {
	return fmt.Sprintf("<Lazy(%s)>", e.Expr.Repr())
}

func (e *LazyExpr) Eq(another *LazyExpr) bool {
	return ExprEq(e.Expr, another.Expr)
}

// ThunkExpr represents the 'thunk <expr>' construct.
// It forces the evaluation of an expression that should yield a Thunk.
type ThunkExpr struct {
	Expr Expr // The expression expected to evaluate to a Thunk value.
}

func ForceThunk(expr any) *ThunkExpr { // Renamed constructor for clarity
	return &ThunkExpr{Expr: AnyToExpr(expr)}
}

func (e *ThunkExpr) Printable() *epl.Printable {
	return epl.PrintableIter(func(yield func(v *epl.Printable) bool) {
		if !yield(epl.Printablef(0, "Thunk:")) {
			return
		} // Represents the forcing operation
		cP := e.Expr.Printable()
		cP.IndentLevel += 1
		if !yield(cP) {
			return
		}
	})
}

func (e *ThunkExpr) Repr() string {
	return fmt.Sprintf("<Thunk(%s)>", e.Expr.Repr()) // Represents the forcing operation
}

func (e *ThunkExpr) Eq(another *ThunkExpr) bool {
	return ExprEq(e.Expr, another.Expr)
}

// Thunk is the *value* representing a delayed computation.
// It is NOT an AST node (Expr). It's returned by evaluating LazyExpr.
type Thunk struct {
	Expr Expr          // The unevaluated expression.
	Env  *epl.Env[any] // The environment captured when the LazyExpr was evaluated.
}

// Repr for Thunk value (for debugging evaluator results)
func (t *Thunk) Repr() string {
	return fmt.Sprintf("<ThunkValue Expr:%s Env:%p>", t.Expr.Repr(), t.Env)
}

// Note: We might need a way to compare Thunk values if they appear in test results,
// but strict equality based on Env pointer might be too restrictive. Comparing Expr might suffice.
// Add Eq if needed:
// func (t *Thunk) Eq(another *Thunk) bool { ... }

// --- Ensure ExprEq handles these via reflection ---
// No changes needed in chapter3/expr.go.:w

// LazyLangEval evaluates expressions including lazy evaluation constructs.
type LazyLangEval struct {
	ImpRefLangEval // Embed the previous evaluator
}

// NewLazyLangEval creates a new evaluator for the lazy language.
func NewLazyLangEval() *LazyLangEval {
	out := &LazyLangEval{}
	// CRITICAL: Set the Self pointer for the embedded BaseEval
	out.BaseEval.Self = out
	return out
}

// LocalEval handles expression types specific to LazyLang or delegates.
func (l *LazyLangEval) LocalEval(expr Expr, env *epl.Env[any]) (any, error) {
	// log.Printf("LazyLangEval evaluating: %s (%T)\n", expr.Repr(), expr)
	switch n := expr.(type) {
	case *LazyExpr: // Handle 'lazy <expr>'
		return l.valueOfLazyExpr(n, env)
	case *ThunkExpr: // Handle 'thunk <expr>' (force)
		return l.valueOfThunkExpr(n, env)
	default:
		// Delegate to the embedded ImpRefLangEval's LocalEval for other types
		return l.ImpRefLangEval.LocalEval(expr, env)
	}
}

// valueOfLazyExpr handles 'lazy <expr>'
func (l *LazyLangEval) valueOfLazyExpr(e *LazyExpr, env *epl.Env[any]) (any, error) {
	// Package the expression and the *current* environment into a Thunk value.
	// Do not evaluate e.Expr yet.
	thunkValue := &Thunk{
		Expr: e.Expr,
		Env:  env, // Capture the current environment
	}
	// log.Printf("lazy evaluated %s to Thunk %s\n", e.Expr.Repr(), thunkValue.Repr())
	return thunkValue, nil // Return the Thunk value itself
}

// valueOfThunkExpr handles 'thunk <expr>' (force evaluation)
func (l *LazyLangEval) valueOfThunkExpr(e *ThunkExpr, env *epl.Env[any]) (any, error) {
	// 1. Evaluate the expression that should yield a Thunk.
	value, err := l.Eval(e.Expr, env)
	if err != nil {
		return nil, err
	}

	// 2. Check if the result is actually a Thunk.
	thunkValue, ok := value.(*Thunk)
	if !ok {
		log.Panicf("thunk operator expected a thunk value, but got type %T for expr %s", value, e.Expr.Repr())
	}

	// log.Printf("thunk forcing evaluation of %s in Env %p\n", thunkValue.Expr.Repr(), thunkValue.Env)

	// 3. Force the evaluation by evaluating the Thunk's expression in its captured environment.
	//    Use l.Eval() for recursive dispatch.
	return l.Eval(thunkValue.Expr, thunkValue.Env)
}
