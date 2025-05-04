package chapter3

import (
	"fmt"
	"log"
	"reflect"

	epl "github.com/panyam/eplgo"
	gfn "github.com/panyam/goutils/fn"
)

type LitExpr struct {
	// can only be string, int, float or bool or one of the other lit types
	Value any
}

func Lit(val any) *LitExpr {
	// while type(value) is Lit: value = value.value
	// assert type(value) in (str, int, float, bool)
	return &LitExpr{Value: val}
}

func (l *LitExpr) Eq(another *LitExpr) bool {
	// if type(another) == type(self.value): return self.value == another
	// elif type(self) != type(another): return False
	return l.Value == another.Value
}

func (l *LitExpr) Repr() string {
	return fmt.Sprintf("Val(%v:%v)", l.Value, reflect.TypeOf(l.Value).Name())
}

func (l *LitExpr) Printable() *epl.Printable {
	return &epl.Printable{0, l.Repr(), nil}
}

type VarExpr struct {
	Name string
}

func Var(n string) *VarExpr {
	return &VarExpr{Name: n}
}

func (v *VarExpr) Printable() *epl.Printable {
	return epl.Printablef(0, "var %s", v.Name)
}

func (v *VarExpr) Eq(another *VarExpr) bool {
	return v.Name == another.Name
}

func (v *VarExpr) Repr() string {
	return fmt.Sprintf("<Var(%s)>", v.Name)
}

type TupleExpr struct {
	Children []Expr
}

func Tuple(children ...Expr) *TupleExpr {
	return &TupleExpr{Children: children}
}

func (e *TupleExpr) Printable() *epl.Printable {
	return epl.PrintableIter(func(yield func(v *epl.Printable) bool) {
		// yield 0, "var %s" % self.name
		if !yield(epl.Printablef(0, "Tuple")) {
			return
		}
		ExprListPrintable(1, e.Children, yield)
	})
}

func (e *TupleExpr) Eq(another *TupleExpr) bool {
	return ExprListEq(e.Children, another.Children)
}

func (e *TupleExpr) Repr() string {
	return fmt.Sprintf("<Tuple(%s)>", ExprListRepr(e.Children))
}

type OpExpr struct {
	Op   string
	Args []Expr
}

func Op(op string, args ...any) *OpExpr {
	return &OpExpr{Op: op, Args: gfn.Map(args, AnyToExpr)}
}

func (v *OpExpr) Printable() *epl.Printable {
	return epl.PrintableIter(func(yield func(v *epl.Printable) bool) {
		// yield 0, "var %s" % self.name
		if !yield(epl.Printablef(0, "Op<%s>", v.Op)) {
			return
		}
		ExprListPrintable(1, v.Args, yield)
	})
}

func (v *OpExpr) Eq(another *OpExpr) bool {
	return v.Op == another.Op && ExprListEq(v.Args, another.Args)
}

func (v *OpExpr) Repr() string {
	return fmt.Sprintf("<Op(%s, [%s])>", v.Op, ExprListRepr(v.Args))
}

type IfExpr struct {
	Cond Expr
	Then Expr
	Else Expr
}

func If(cond any, then any, els any) *IfExpr {
	return &IfExpr{AnyToExpr(cond), AnyToExpr(then), AnyToExpr(els)}
}

func (v *IfExpr) Printable() *epl.Printable {
	return epl.PrintableIter(func(yield func(v *epl.Printable) bool) {
		if !yield(epl.Printablef(0, "If:")) {
			return
		}
		if !yield(epl.Printablef(1, "Cond")) {
			return
		}
		if !yield(v.Cond.Printable()) {
			return
		}
		if !yield(epl.Printablef(1, "Then")) {
			return
		}
		if !yield(v.Then.Printable()) {
			return
		}
		if !yield(epl.Printablef(1, "Else")) {
			return
		}
		if !yield(v.Else.Printable()) {
			return
		}
	})
}

func (v *IfExpr) Eq(another *IfExpr) bool {
	return ExprEq(v.Cond, another.Cond) && ExprEq(v.Then, another.Then) && ExprEq(v.Else, another.Else)
}

func (v *IfExpr) Repr() string {
	return fmt.Sprintf("<If(%s) { %s } else { %s }>", v.Cond.Repr(), v.Then.Repr(), v.Else.Repr())
}

type IsZeroExpr struct {
	Expr Expr
}

func IsZero(e any) *IsZeroExpr {
	return &IsZeroExpr{AnyToExpr(e)}
}

func (v *IsZeroExpr) Printable() *epl.Printable {
	return epl.PrintableIter(func(yield func(v *epl.Printable) bool) {
		if !yield(epl.Printablef(0, "IsZero:")) {
			return
		}
		if !yield(v.Expr.Printable()) {
			return
		}
	})
}

func (v *IsZeroExpr) Eq(another *IsZeroExpr) bool {
	return ExprEq(v.Expr, another.Expr)
}

func (v *IsZeroExpr) Repr() string {
	return fmt.Sprintf("<IsZero(%s)>", v.Expr.Repr())
}

type LetExpr struct {
	Mappings map[string]Expr
	Body     Expr
}

func Let(mappings map[string]Expr, body Expr) *LetExpr {
	return &LetExpr{Body: body, Mappings: mappings}
}

func (v *LetExpr) Printable() *epl.Printable {
	return epl.PrintableIter(func(yield func(v *epl.Printable) bool) {
		if !yield(epl.Printablef(0, "LetExpr:")) {
			return
		}
		for _, k := range epl.SortedKeys(v.Mappings) {
			val := v.Mappings[k]
			if !yield(epl.Printablef(2, "%s = ", k)) {
				return
			}
			vp := val.Printable()
			vp.IndentLevel += 2
			if !yield(vp) {
				return
			}
		}
		if !yield(epl.Printablef(1, "in:")) {
			return
		}
		if !yield(v.Body.Printable()) {
			return
		}
	})
}

func (v *LetExpr) Eq(another *LetExpr) bool {
	if len(v.Mappings) != len(another.Mappings) {
		return false
	}

	for k, v := range v.Mappings {
		if !ExprEq(v, another.Mappings[k]) {
			return false
		}
	}
	return ExprEq(v.Body, another.Body)
}

func (v *LetExpr) Repr() string {
	out := "<Let "
	first := true
	for k, v := range v.Mappings {
		if first {
			out += "("
		} else {
			out += ", "
		}
		out += k
		out += " = "
		out += v.Repr()
		first = false
	}
	if !first {
		out += ")"
	}
	out += " in "
	out += v.Body.Repr()
	return out
}

// Evaluator for the LetLang

type LetLangEval struct {
	BaseEval
}

func NewLetLangEval() *LetLangEval {
	out := &LetLangEval{}
	out.BaseEval.Self = out
	return out
}

func (l *LetLangEval) LocalEval(expr Expr, env *epl.Env[any]) (any, error) {
	// log.Println("LetLangEval for: ", reflect.TypeOf(expr), expr.Repr())
	switch n := expr.(type) {
	case *LitExpr:
		return l.ValueOfLit(n, env)
	case *VarExpr:
		return l.ValueOfVar(n, env)
	case *OpExpr:
		return l.ValueOfOpExpr(n, env)
	case *IsZeroExpr:
		return l.ValueOfIsZeroExpr(n, env)
	case *IfExpr:
		return l.ValueOfIfExpr(n, env)
	case *LetExpr:
		return l.ValueOfLetExpr(n, env)
	case *TupleExpr:
		return l.ValueOfTupleExpr(n, env)
	default:
		// log.Printf("Invalid type: %v - %v", n, reflect.TypeOf(n))
	}
	panic("Invalid type")
}

// Evalute the value of a literal
func (l *LetLangEval) ValueOfLit(lit *LitExpr, env *epl.Env[any]) (any, error) {
	return lit, nil
}

// Evaluate the value of a variable
func (l *LetLangEval) ValueOfVar(e *VarExpr, env *epl.Env[any]) (any, error) {
	// TODO - Error and type checking
	val, found := env.Get(e.Name)
	if !found {
		return nil, fmt.Errorf("variable '%s' not found in environment", e.Name)
	}
	return val, nil
}

func (l *LetLangEval) ValueOfOpExpr(e *OpExpr, env *epl.Env[any]) (any, error) {
	// TODO - Error and type checking
	opfunc := l.GetOpFunc(e.Op)
	if opfunc == nil {
		log.Fatalf("opfunc not found: %s", e.Op)
		panic("opfunc not found")
	}
	return opfunc(env, e.Args)
}

func (l *LetLangEval) ValueOfIfExpr(e *IfExpr, env *epl.Env[any]) (any, error) {
	condVal, err := l.Eval(e.Cond, env) // Returns (any, error)
	if err != nil {
		return nil, err // Propagate error
	}

	condBool := false

	// Define truthiness: only Lit(true) is true, others (incl. Lit(false)) are false
	if litCond, ok := condVal.(*LitExpr); ok {
		if boolVal, ok2 := litCond.Value.(bool); ok2 {
			condBool = boolVal // Use the actual boolean value
		}
	} // Any non-literal result (like a BoundProc) also counts as false

	// log.Printf("If condition %s evaluated to %v (%T), bool result: %v\n", e.Cond.Repr(), condVal, condVal, condBool)

	if condBool {
		return l.Eval(e.Then, env)
	} else {
		return l.Eval(e.Else, env)
	}
}

func (l *LetLangEval) ValueOfTupleExpr(e *TupleExpr, env *epl.Env[any]) (any, error) {
	vals, err := l.EvalExprList(e.Children, env)
	if err != nil {
		return nil, err // Propagate error from list eval
	}
	// Return the slice directly
	return vals, nil
}

func (l *LetLangEval) ValueOfIsZeroExpr(e *IsZeroExpr, env *epl.Env[any]) (any, error) {
	val, err := l.Eval(e.Expr, env) // Returns any
	if err != nil {
		return nil, err
	}
	litVal, ok := val.(*LitExpr)
	if !ok {
		return nil, fmt.Errorf("iszero expected a LitExpr argument, got %T (%v) for expr %s", val, val, e.Expr.Repr())
	}
	// For now, assume IsZero only works on ints
	intVal, ok := litVal.Value.(int)
	if !ok {
		panic(fmt.Sprintf("iszero expected an integer value, got %T (%v)", litVal.Value, litVal.Value))
	}
	return Lit(intVal == 0), nil // Return *LitExpr(bool)
}

type ExprMap = map[string]Expr

func (l *LetLangEval) ValueOfLetExpr(e *LetExpr, env *epl.Env[any]) (any, error) {
	bindings := map[string]any{}
	for k, v := range e.Mappings {
		// log.Printf("Evaluating let binding: %s = %s\n", k, v.Repr())
		val, err := l.Eval(v, env) // Eval returns (any, error)
		if err != nil {
			return nil, fmt.Errorf("evaluating binding '%s': %w", k, err) // Propagate error
		}
		bindings[k] = val
		// log.Printf("Binding %s evaluated to: %v (%T)\n", k, bindings[k], bindings[k])
	}
	newenv := env.Extend(bindings)
	// log.Printf("Evaluating let body %s in new env: %s\n", e.Body.Repr(), newenv)
	return l.Eval(e.Body, newenv) // Propagate result/error from body
}
