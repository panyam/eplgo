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

	// 3. Use reflection to find and call the Eq method
	v1 := reflect.ValueOf(e1)
	v2 := reflect.ValueOf(e2)

	// Look for a method named "Eq" on the first value (v1).
	// Assumes the method signature is `Eq(another T) bool` where T is the type of v1.
	eqMethod := v1.MethodByName("Eq")

	// Check if the method exists
	if !eqMethod.IsValid() {
		// If no Eq method, what's the desired behavior?
		// Option 1: Panic - indicates an incomplete implementation for an Expr type.
		log.Panicf("ExprEq: Type %T does not have an Eq method", e1)
		// Option 2: Fallback to deep equality (might be slow or incorrect for unexported fields).
		// return reflect.DeepEqual(e1, e2)
		// Option 3: Return false (safest default if Eq is optional).
		// return false
		return false // Let's default to false if Eq is missing
	}

	// Check the method signature (optional but recommended for robustness)
	methodType := eqMethod.Type()
	// Eq should take 1 argument (the other value) and return 1 bool result.
	if methodType.NumIn() != 1 || methodType.NumOut() != 1 || methodType.Out(0).Kind() != reflect.Bool {
		log.Panicf("ExprEq: Type %T has an Eq method with incorrect signature: %s (expected func(T) bool)", e1, methodType)
		return false
	}
	// Check if the input argument type matches the type of v1/v2
	if methodType.In(0) != t1 {
		log.Panicf("ExprEq: Type %T has an Eq method with incorrect argument type: %s (expected %s)", e1, methodType.In(0), t1)
		return false
	}

	// Prepare arguments for the call: just the second value (v2)
	args := []reflect.Value{v2}

	// Call the Eq method
	results := eqMethod.Call(args)

	// Extract the boolean result
	if len(results) == 1 && results[0].Kind() == reflect.Bool {
		return results[0].Bool()
	} else {
		// This shouldn't happen if the signature check passed, but handle defensively.
		log.Panicf("ExprEq: Call to Eq method on type %T returned unexpected results: %v", e1, results)
		return false
	}
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
	if (e1 == nil && e2 != nil) || (e1 != nil && e2 == nil) {
		return false
	}

	for i, child1 := range e1 {
		child2 := e2[i]
		// Use the top-level ExprEq for recursive comparison
		if !ExprEq(child1, child2) { // This now uses the reflection version
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
