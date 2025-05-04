package chapter3

import (
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

func AnyToExpr(x any) Expr {
	if x == nil {
		return x.(Expr)
	}
	switch x := x.(type) {
	case Expr:
		return x
	case string:
		return Var(x)
	case int:
		return Lit(x)
	case bool:
		return Lit(x)
	default:
		log.Fatalf("Cannot convert %v (type: %v) to Expr", x, reflect.TypeOf(x))
	}
	panic("Cannot convert to Expr")
}
