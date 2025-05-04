package chapter3

import epl "github.com/panyam/eplgo"

type OpFunc func(env *epl.Env[any], args []Expr) any

type Evaluater interface {
	This() Evaluater
	Eval(expr Expr, env *epl.Env[any]) any
}

type BaseEval struct {
	Self    Evaluater
	OpFuncs map[string]OpFunc
}

func (b *BaseEval) This() Evaluater {
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
