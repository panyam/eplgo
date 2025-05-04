package chapter3

import (
	"fmt" // Needed for Eq comparison
	"strings"

	epl "github.com/panyam/eplgo"
)

// LetRecExpr represents the 'letrec' construct for mutual recursion.
// Example: letrec f(x) = ..., g(y) = ... in body
type LetRecExpr struct {
	// Procs maps procedure names to their ProcExpr definitions.
	Procs map[string]*ProcExpr
	// Body is the expression evaluated in the environment extended with the recursive procedures.
	Body Expr
}

// LetRec is a constructor for LetRecExpr.
// It also ensures the Name field within each ProcExpr is set.
func LetRec(procs map[string]*ProcExpr, body Expr) *LetRecExpr {
	// Ensure the Name field is set correctly in the provided ProcExprs
	for name, proc := range procs {
		if proc.Name != "" && proc.Name != name {
			// Optionally panic or log a warning if name is pre-set inconsistently
			panic(fmt.Sprintf("Inconsistent name in LetRec proc map: key '%s', proc.Name '%s'", name, proc.Name))
		}
		proc.Name = name // Set or overwrite name based on map key
	}
	return &LetRecExpr{Procs: procs, Body: body}
}

// Printable generates a printable representation for debugging.
func (v *LetRecExpr) Printable() *epl.Printable {
	return epl.PrintableIter(func(yield func(v *epl.Printable) bool) {
		if !yield(epl.Printablef(0, "LetRec:")) {
			return
		}
		for _, name := range epl.SortedKeys(v.Procs) {
			proc := v.Procs[name]
			// Reuse ProcExpr's Printable, adjusting indentation maybe?
			// Or construct manually here:
			if !yield(epl.Printablef(1, "%s (%s) =", name, strings.Join(proc.Varnames, ", "))) {
				return
			}
			// Indent the body of the proc under its declaration
			procBodyPrintable := proc.Body.Printable()
			procBodyPrintable.IndentLevel += 2 // Adjust relative indent
			if !yield(procBodyPrintable) {
				return
			}
		}
		if !yield(epl.Printablef(1, "in:")) {
			return
		}
		// Indent the main body
		bodyPrintable := v.Body.Printable()
		bodyPrintable.IndentLevel += 2 // Adjust relative indent
		if !yield(bodyPrintable) {
			return
		}
	})
}

// Eq checks for equality with another LetRecExpr.
func (v *LetRecExpr) Eq(another *LetRecExpr) bool {
	if len(v.Procs) != len(another.Procs) {
		return false
	}
	// Compare Procs map
	for name, proc1 := range v.Procs {
		proc2, ok := another.Procs[name]
		if !ok || !proc1.Eq(proc2) { // Use ProcExpr.Eq
			return false
		}
	}
	// Compare Body
	return ExprEq(v.Body, another.Body)
}

// Repr generates a string representation for debugging.
func (v *LetRecExpr) Repr() string {
	var procStrs []string
	for name, proc := range v.Procs {
		// Simplified representation for brevity
		procStrs = append(procStrs, fmt.Sprintf("%s(%s)=%s", name, strings.Join(proc.Varnames, ","), proc.Body.Repr()))
	}
	return fmt.Sprintf("<LetRec {%s} in %s>", strings.Join(procStrs, "; "), v.Body.Repr())
}

// LetRecLangEval extends ProcLangEval to handle 'letrec'.
type LetRecLangEval struct {
	ProcLangEval // Embed ProcLangEval to inherit its methods
}

// NewLetRecLangEval creates a new evaluator for the letrec language.
func NewLetRecLangEval() *LetRecLangEval {
	out := &LetRecLangEval{}
	// CRITICAL: Set the Self pointer for the embedded BaseEval
	out.BaseEval.Self = out
	return out
}

// LocalEval handles expression types specific to LetRecLang or delegates to ProcLangEval.
func (l *LetRecLangEval) LocalEval(expr Expr, env *epl.Env[any]) any {
	// log.Printf("LetRecLangEval evaluating: %s (%T)", expr.Repr(), expr)
	switch n := expr.(type) {
	case *LetRecExpr:
		return l.ValueOfLetRec(n, env)
	default:
		// Delegate to the embedded ProcLangEval's LocalEval for other types
		// (Lit, Var, Op, If, Let, IsZero, Proc, Call)
		return l.ProcLangEval.LocalEval(expr, env)
	}
}

// ValueOfLetRec evaluates a 'letrec' expression.
func (l *LetRecLangEval) ValueOfLetRec(e *LetRecExpr, env *epl.Env[any]) any {
	// 1. Create a new environment nested within the current one.
	//    This new environment will hold the mutually recursive bindings.
	newenv := env.Push() // newenv.outer points to env

	// 2. Bind each procedure in the letrec block to this *new* environment.
	//    This is the key step that enables recursion.
	boundProcs := map[string]*BoundProc{}
	for name, procExpr := range e.Procs {
		// log.Printf("Binding proc '%s' in letrec to newenv %p (outer: %p)", name, newenv, env)
		boundProcs[name] = procExpr.Bind(newenv) // Bind uses newenv
	}

	// 3. Populate the new environment with these bound procedures.
	// TODO can steps 2 and 3 be merged in a single loop?
	for name, boundProc := range boundProcs {
		newenv.Set(name, boundProc)
	}
	// log.Printf("Newenv after binding letrec procs: %s", newenv)

	// 4. Evaluate the body expression within this new environment.
	//    Calls within the body can now resolve the recursive procedure names.
	// log.Printf("Evaluating letrec body %s in newenv %p", e.Body.Repr(), newenv)
	return l.Eval(e.Body, newenv)
}

// --- Helper to create Proc Maps easily in tests ---

// ProcMap creates a map suitable for LetRec's Procs field.
// Args are alternating string (name) and *ProcExpr (proc).
func ProcMap(args ...any) map[string]*ProcExpr {
	out := make(map[string]*ProcExpr)
	if len(args)%2 != 0 {
		panic("ProcMap requires an even number of arguments (name, proc, name, proc, ...)")
	}
	for i := 0; i < len(args); i += 2 {
		name, okName := args[i].(string)
		proc, okProc := args[i+1].(*ProcExpr)
		if !okName || !okProc {
			panic(fmt.Sprintf("Invalid arguments at index %d,%d for ProcMap: expected string, *ProcExpr, got %T, %T", i, i+1, args[i], args[i+1]))
		}
		out[name] = proc
	}
	return out
}
