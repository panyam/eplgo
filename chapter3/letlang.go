package chapter3

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	epl "github.com/panyam/eplgo"
	gfn "github.com/panyam/goutils/fn"
)

type Expr interface {
	// Eq(another Expr) bool
	Printable() *epl.Printable
	Repr() string
}

func ExprEq(e1 Expr, e2 Expr) bool {
	if e1 == e2 {
		return true
	}
	if e1 == nil || e2 == nil {
		return false
	}
	t1 := reflect.TypeOf(e1)
	t2 := reflect.TypeOf(e2)
	if t1 != t2 {
		return false
	}

	// TODO - Call specific eq ops
	return true
}

func ExprListPrintable(level int, e []Expr, yield func(*epl.Printable) bool) bool {
	for _, child := range e {
		if !yield(child.Printable()) {
			return false
		}
	}
	return true
}

func ExprListRepr(e []Expr) string {
	return strings.Join(gfn.Map(e, func(e Expr) string { return e.Repr() }), ", ")
}

func ExprListEq(e1 []Expr, e2 []Expr) bool {
	if len(e1) != len(e2) {
		return false
	}

	for i, child := range e1 {
		if !ExprEq(child, e2[i]) {
			return false
		}
	}
	return true
}

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

func Op(op string, args ...Expr) *OpExpr {
	return &OpExpr{Op: op, Args: args}
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

func If(cond Expr, then Expr, els Expr) *IfExpr {
	return &IfExpr{cond, then, els}
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

func IsZero(e Expr) *IsZeroExpr {
	return &IsZeroExpr{e}
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
		for k, v := range v.Mappings {
			if !yield(epl.Printablef(2, "%s = ", k)) {
				return
			}
			if !yield(v.Printable()) {
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

func (l *LetLangEval) Eval(expr Expr, env *epl.Env[any]) any {
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
		log.Printf("Invalid type: %v", n)
	}
	panic("Invalid type")
}

// Evalute the value of a literal
func (l *LetLangEval) ValueOfLit(lit *LitExpr, env *epl.Env[any]) any {
	return lit
}

// Evaluate the value of a variable
func (l *LetLangEval) ValueOfVar(e *VarExpr, env *epl.Env[any]) any {
	// TODO - Error and type checking
	val, _ := env.Get(e.Name)
	return val
}

func (l *LetLangEval) ValueOfOpExpr(e *OpExpr, env *epl.Env[any]) any {
	// TODO - Error and type checking
	opfunc := l.GetOpFunc(e.Op)
	if opfunc == nil {
		log.Fatalf("opfunc not found: %s", e.Op)
		panic("opfunc not found")
	}
	return opfunc(env, e.Args)
}

func (l *LetLangEval) ValueOfIfExpr(e *IfExpr, env *epl.Env[any]) any {
	// TODO - Error and type checking
	cond := l.Eval(e.Cond, env)
	if cond == true {
		return l.Eval(e.Then, env)
	} else {
		return l.Eval(e.Else, env)
	}
}

func (l *LetLangEval) ValueOfTupleExpr(e *TupleExpr, env *epl.Env[any]) any {
	// TODO - Error and type checking
	vals := gfn.Map(e.Children, func(e Expr) any {
		return l.Eval(e, env)
	})
	return vals // Tuple(vals)
}

func (l *LetLangEval) ValueOfIsZeroExpr(e *IsZeroExpr, env *epl.Env[any]) any {
	// TODO - Error and type checking
	val := l.Eval(e.Expr, env).(*LitExpr)
	return Lit(val.Value == 0)
}

type ExprMap = map[string]Expr

func (l *LetLangEval) ValueOfLetExpr(e *LetExpr, env *epl.Env[any]) any {
	// TODO - Error and type checking
	bindings := map[string]any{}
	for k, v := range e.Mappings {
		bindings[k] = l.Eval(v, env)
	}
	newenv := env.Extend(bindings)
	return l.Eval(e.Body, newenv)
}

/*
   @case("let")
   def valueOfLet(self, let, env):
       expvals = {var: self.valueOf(exp, env) for var,exp in let.mappings.items()}
       newenv = env.extend(**expvals)
       return self.valueOf(let.body, newenv)
*/
