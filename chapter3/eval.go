package chapter3

import (
	epl "github.com/panyam/eplgo"
)

type OpFunc func(env *epl.Env[any], args []Expr) any

type evaluater interface {
	This() evaluater
	LocalEval(expr Expr, env *epl.Env[any]) any
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

func (b *BaseEval) Eval(expr Expr, env *epl.Env[any]) any {
	// log.Println("SelfType: ", reflect.TypeOf(b.Self))
	// log.Println("ExprType: ", reflect.TypeOf(expr), expr.Repr())
	return b.Self.LocalEval(expr, env)
}

func (b *BaseEval) EvalExprList(exprs []Expr, env *epl.Env[any]) []any {
	var out []any
	for _, exp := range exprs {
		out = append(out, b.Self.LocalEval(exp, env))
	}
	return out
}
