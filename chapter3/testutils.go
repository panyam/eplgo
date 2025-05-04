package chapter3

import (
	"testing"

	epl "github.com/panyam/eplgo"
	"gotest.tools/assert"
)

var ExprDict = epl.Dict[string, Expr]

type Evaluator interface {
	Eval(expr Expr, env *epl.Env[any]) any
	SetOpFunc(string, OpFunc)
}

type TestCase struct {
	Name     string
	Expected any
	Expr     Expr
}

func RunTest(t *testing.T, e Evaluator, tc *TestCase, extraenv map[string]Expr) {
	env := epl.NewEnv[any](nil)
	for k, v := range extraenv {
		env.Set(k, v)
	}
	value := e.Eval(tc.Expr, env)
	found := value.(*LitExpr)
	assert.Equal(t, found.Value, tc.Expected)
	/*log.Println("======= TestCase: ", tc.Name, "=======")
	log.Println("Found: ", found.Value)
	log.Println("Expected: ", tc.Expected)
	*/
	// assert.Equal(tc.Expected, value)
}

func setOpFuncs(e Evaluator) Evaluator {
	e.SetOpFunc("-", func(env *epl.Env[any], args []Expr) any {
		v1 := e.Eval(args[0], env).(*LitExpr)
		v2 := e.Eval(args[1], env).(*LitExpr)
		return Lit(v1.Value.(int) - v2.Value.(int))
	})
	return e
}
