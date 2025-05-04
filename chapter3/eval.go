package chapter3

import (
	"fmt"

	epl "github.com/panyam/eplgo"
)

type OpFunc func(env *epl.Env[any], args []Expr) (any, error)

type evaluater interface {
	This() evaluater
	LocalEval(expr Expr, env *epl.Env[any]) (any, error)
}

type BaseEval struct {
	Self    evaluater
	OpFuncs map[string]OpFunc
}

func (b *BaseEval) This() evaluater {
	return b.Self
}

func (b *BaseEval) SetOpFunc(name string, fn OpFunc) {
	if b.OpFuncs == nil {
		b.OpFuncs = make(map[string]OpFunc)
	}
	b.OpFuncs[name] = fn
}

func (b *BaseEval) GetOpFunc(name string) OpFunc {
	if b.OpFuncs == nil {
		return nil
	}
	return b.OpFuncs[name]
}

func (b *BaseEval) Eval(expr Expr, env *epl.Env[any]) (any, error) {
	// Error check result before returning
	val, err := b.Self.LocalEval(expr, env)
	if err != nil {
		// Optionally log or wrap error here
		return nil, fmt.Errorf("evaluating %s: %w", expr.Repr(), err)
	}
	return val, nil
}

func (b *BaseEval) EvalExprList(exprs []Expr, env *epl.Env[any]) ([]any, error) {
	out := make([]any, len(exprs))
	for i, exp := range exprs {
		val, err := b.Self.LocalEval(exp, env) // Use Eval which returns (any, error)
		if err != nil {
			// If any expression fails, stop and return the error
			return nil, fmt.Errorf("evaluating argument %d (%s): %w", i, exp.Repr(), err)
		}
		out[i] = val
	}
	return out, nil
}
