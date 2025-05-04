package chapter4

import (
	epl "github.com/panyam/eplgo"
	"github.com/panyam/eplgo/chapter3"
)

// A few imports to not avoid having to prefix with chapter3 all over the place
type Expr = chapter3.Expr
type LitExpr = chapter3.LitExpr
type VarExpr = chapter3.VarExpr

var ExprDict = epl.Dict[string, Expr]

type TestCase = chapter3.TestCase
type Evaluator = chapter3.Evaluator

var SetOpFuncs = chapter3.SetOpFuncs
var Lit = chapter3.Lit
var Let = chapter3.Let
var LetRec = chapter3.LetRec
var Op = chapter3.Op
var If = chapter3.If
var IsZero = chapter3.IsZero
var Proc = chapter3.Proc
var Call = chapter3.Call
var ProcMap = chapter3.ProcMap
var Var = chapter3.Var
var AnyToExpr = chapter3.AnyToExpr
var ExprEq = chapter3.ExprEq
var ExprListEq = chapter3.ExprListEq
var ExprListRepr = chapter3.ExprListRepr
