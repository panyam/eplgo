package chapter3

import (
	"fmt"
	"testing"

	epl "github.com/panyam/eplgo"
	"github.com/stretchr/testify/assert" // Use testify for better assertions
)

var ExprDict = epl.Dict[string, Expr]

// Only declare the interface we want at the caller side instead of at the provider side.
type Evaluator interface {
	Eval(expr Expr, env *epl.Env[any]) any
	SetOpFunc(string, OpFunc)
}

type TestCase struct {
	Name     string
	Expected any // Expected can be int, bool, etc.
	Expr     Expr
}

func RunTest(t *testing.T, e Evaluator, tc *TestCase, extraenv map[string]Expr) {
	env := epl.NewEnv[any](nil)
	for k, v := range extraenv {
		// Values in extraenv could be Go primitives or already Expr types (like LitExpr)
		// The env stores 'any', so we can put them in directly.
		// However, when ValueOfVar retrieves them, it gets 'any'.
		env.Set(k, v)
	}

	// log.Printf("======= Running TestCase: %s =======", tc.Name)
	// log.Println("Expr:", tc.Expr.Repr())
	// if extraenv != nil {
	//  log.Println("Initial Env:", extraenv)
	// }

	value := e.Eval(tc.Expr, env) // Eval returns any

	// Assert based on the *expected* type and value
	switch expected := tc.Expected.(type) {
	case int:
		finalValue, ok := value.(*LitExpr)
		if !ok {
			t.Fatalf("Test %s: Expected LitExpr(int) but got %T (%#v)", tc.Name, value, value)
		}
		assert.Equal(t, expected, finalValue.Value, "Test %s Failed", tc.Name)
	case bool:
		finalValue, ok := value.(*LitExpr)
		if !ok {
			t.Fatalf("Test %s: Expected LitExpr(bool) but got %T (%#v)", tc.Name, value, value)
		}
		assert.Equal(t, expected, finalValue.Value, "Test %s Failed", tc.Name)
	// Add cases for string, float64 if needed
	// Add case for expecting a BoundProc if some tests require it
	// case *BoundProc: ...
	default:
		// Fallback or error for unhandled expected types
		// Use assert.Equal for direct comparison if types might match (e.g., comparing two BoundProcs)
		assert.Equal(t, tc.Expected, value, "Test %s Failed (Default Comparison)", tc.Name)
		// Or fail if the expected type isn't handled yet:
		// t.Fatalf("Test %s: Unhandled expected type %T in RunTest assertion", tc.Name, tc.Expected)
	}
	// log.Printf("Test %s Passed. Found: %v (%T)\n", tc.Name, value, value)
}

func SetOpFuncs(e Evaluator) Evaluator {
	e.SetOpFunc("-", func(env *epl.Env[any], args []Expr) any {
		if len(args) != 2 {
			panic("'-' operator expects exactly 2 arguments")
		}
		v1Raw := e.Eval(args[0], env) // Returns any
		v2Raw := e.Eval(args[1], env) // Returns any
		// Assume ops work on literals and expect ints for '-'
		v1Lit, ok1 := v1Raw.(*LitExpr)
		v2Lit, ok2 := v2Raw.(*LitExpr)
		if !ok1 || !ok2 {
			panic(fmt.Sprintf("'-' operator requires LitExpr arguments, got %T and %T", v1Raw, v2Raw))
		}
		v1Int, ok1 := v1Lit.Value.(int)
		v2Int, ok2 := v2Lit.Value.(int)
		if !ok1 || !ok2 {
			panic(fmt.Sprintf("'-' operator requires integer values, got %T and %T", v1Lit.Value, v2Lit.Value))
		}
		return Lit(v1Int - v2Int) // Return LitExpr(int)
	})
	e.SetOpFunc("+", func(env *epl.Env[any], args []Expr) any {
		out := 0
		for _, arg := range args {
			vRaw := e.Eval(arg, env) // Returns any
			// Assume ops work on literals and expect ints for '+'
			vLit, ok := vRaw.(*LitExpr)
			if !ok {
				panic(fmt.Sprintf("'+' operator requires LitExpr arguments, got %T", vRaw))
			}
			vInt, ok := vLit.Value.(int)
			if !ok {
				panic(fmt.Sprintf("'+' operator requires integer values, got %T", vLit.Value))
			}
			out += vInt
		}
		return Lit(out) // Return LitExpr(int)
	})
	e.SetOpFunc("*", func(env *epl.Env[any], args []Expr) any {
		out := 1
		for _, arg := range args {
			vRaw := e.Eval(arg, env) // Returns any
			// Assume ops work on literals and expect ints for '*'
			vLit, ok := vRaw.(*LitExpr)
			if !ok {
				panic(fmt.Sprintf("'*' operator requires LitExpr arguments, got %T", vRaw))
			}
			vInt, ok := vLit.Value.(int)
			if !ok {
				panic(fmt.Sprintf("'*' operator requires integer values, got %T", vLit.Value))
			}
			out *= vInt
		}
		return Lit(out) // Return LitExpr(int)
	})
	return e
}
