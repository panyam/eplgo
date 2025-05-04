package chapter3

import (
	"bytes"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- Tests for ExprEq ---
func makeProcXBodyX() *ProcExpr { return Proc([]string{"x"}, Var("x")) }
func makeCallFX() *CallExpr     { return Call(Var("f"), Var("x")) }
func makeLetRecF() *LetRecExpr {
	return LetRec(ProcMap("f", makeProcXBodyX()), makeCallFX())
}

func TestExprEq(t *testing.T) {
	// Reusable expressions
	lit1 := Lit(1)
	lit1b := Lit(1)
	lit2 := Lit(2)
	litTrue := Lit(true)
	litFalse := Lit(false)
	varX := Var("x")
	varXB := Var("x")
	varY := Var("y")

	tests := []struct {
		name     string
		e1       Expr
		e2       Expr
		expected bool
	}{
		// --- Use factories/direct construction inside tests ---
		{"nil vs nil", nil, nil, true},
		{"nil vs lit", nil, lit1, false},
		{"lit vs nil", lit1, nil, false},
		{"lit identity", lit1, lit1, true},
		{"lit equal int", lit1, lit1b, true},
		{"lit different int", lit1, lit2, false},
		{"lit equal bool", litTrue, Lit(true), true}, // Recreate
		{"lit different bool", litTrue, litFalse, false},
		{"lit int vs bool", lit1, litTrue, false},
		{"var identity", varX, varX, true},
		{"var equal", varX, varXB, true},
		{"var different", varX, varY, false},
		{"lit vs var", lit1, varX, false},
		{"op equal", Op("-", varX, varY), Op("-", varXB, varY), true},
		{"op different op", Op("-", varX, varY), Op("+", varX, varY), false},
		{"op different args", Op("-", varX, varY), Op("-", varX, lit1), false},
		{"iszero equal", IsZero(varX), IsZero(varXB), true},
		{"iszero different", IsZero(varX), IsZero(varY), false},
		{"if equal", If(IsZero(varX), lit1, lit2), If(IsZero(varXB), lit1b, lit2), true},
		{"if different cond", If(IsZero(varX), lit1, lit2), If(IsZero(varY), lit1, lit2), false},
		{"if different then", If(IsZero(varX), lit1, lit2), If(IsZero(varX), Lit(9), lit2), false},
		{"if different else", If(IsZero(varX), lit1, lit2), If(IsZero(varX), lit1, Lit(9)), false},
		// Proc tests - Ensure fresh Proc instances are compared
		{"proc equal", Proc([]string{"x"}, Var("x")), Proc([]string{"x"}, Var("x")), true}, // Fresh instances
		{"proc different param", Proc([]string{"x"}, Var("x")), Proc([]string{"y"}, Var("x")), false},
		{"proc different body", Proc([]string{"x"}, Var("x")), Proc([]string{"x"}, Var("y")), false},
		{"call equal", Call(Var("f"), Var("x")), Call(Var("f"), Var("x")), true},
		{"call different op", Call(Var("f"), Var("x")), Call(Var("g"), Var("x")), false},
		{"call different arg", Call(Var("f"), Var("x")), Call(Var("f"), Var("y")), false},
		{"let equal", Let(ExprDict("x", Lit(1)), Var("x")), Let(ExprDict("x", Lit(1)), Var("x")), true},
		{"let different key", Let(ExprDict("x", Lit(1)), Var("x")), Let(ExprDict("y", Lit(1)), Var("x")), false},
		{"let different val", Let(ExprDict("x", Lit(1)), Var("x")), Let(ExprDict("x", Lit(2)), Var("x")), false},
		{"let different body", Let(ExprDict("x", Lit(1)), Var("x")), Let(ExprDict("x", Lit(1)), Var("y")), false},
		// LetRec tests - Create ALL components freshly for each definition to avoid mutation side effects
		{"letrec equal",
			LetRec(ProcMap("f", Proc([]string{"x"}, Var("x"))), Call(Var("f"), Var("x"))),
			LetRec(ProcMap("f", Proc([]string{"x"}, Var("x"))), Call(Var("f"), Var("x"))), // All fresh
			true},
		{"letrec different proc name",
			LetRec(ProcMap("f", Proc([]string{"x"}, Var("x"))), Call(Var("f"), Var("x"))),
			LetRec(ProcMap("g", Proc([]string{"x"}, Var("x"))), Call(Var("f"), Var("x"))), // Fresh, different name
			false},
		{"letrec different proc def",
			LetRec(ProcMap("f", Proc([]string{"x"}, Var("x"))), Call(Var("f"), Var("x"))),
			LetRec(ProcMap("f", Proc([]string{"y"}, Var("x"))), Call(Var("f"), Var("x"))), // Fresh, different proc
			false},
		{"letrec different body",
			LetRec(ProcMap("f", Proc([]string{"x"}, Var("x"))), Call(Var("f"), Var("x"))),
			LetRec(ProcMap("f", Proc([]string{"x"}, Var("x"))), Call(Var("g"), Var("x"))), // Fresh, different body
			false},

		{"tuple equal", Tuple(lit1, varX), Tuple(lit1b, varXB), true},
		{"tuple different len", Tuple(lit1, varX), Tuple(lit1), false},
		{"tuple different elem", Tuple(lit1, varX), Tuple(lit1, varY), false},
		{"proc equal", makeProcXBodyX(), makeProcXBodyX(), true}, // Call factory twice
		{"proc different param", makeProcXBodyX(), Proc([]string{"y"}, Var("x")), false},
		{"proc different body", makeProcXBodyX(), Proc([]string{"x"}, Var("y")), false},
		{"letrec equal", makeLetRecF(), makeLetRecF(), true}, // Call factory twice
		{"letrec different proc name",
			makeLetRecF(),
			LetRec(ProcMap("g", makeProcXBodyX()), makeCallFX()), // Use factory for parts
			false},
		{"letrec different proc def",
			makeLetRecF(),
			LetRec(ProcMap("f", Proc([]string{"y"}, Var("x"))), makeCallFX()), // Use factory for parts
			false},
		{"letrec different body",
			makeLetRecF(),
			LetRec(ProcMap("f", makeProcXBodyX()), Call(Var("g"), Var("x"))), // Use factory for parts
			false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, ExprEq(tc.e1, tc.e2))
		})
	}
}

// --- Helper to Capture Log Output ---

var originalLogOutput io.Writer

func captureOutput() *bytes.Buffer {
	originalLogOutput = log.Writer() // Using log.Writer() is better than relying on os.Stderr directly
	var buf bytes.Buffer
	log.SetOutput(&buf)
	// Disable standard log prefixes (date, time) for cleaner comparison
	log.SetFlags(0)
	return &buf
}

func restoreOutput(buf *bytes.Buffer) string {
	log.SetOutput(originalLogOutput)
	log.SetFlags(log.LstdFlags) // Restore default flags
	return buf.String()
}
