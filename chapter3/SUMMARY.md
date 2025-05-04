# Chapter 3 Summary: Core Interpreters

## Purpose

This module implements the core interpreters for the initial languages introduced in Chapter 3 of EPL:
1.  **LetLang:** Introduces basic expressions (literals, variables), arithmetic/boolean operations (`OpExpr`, `IsZeroExpr`), conditional logic (`IfExpr`), and lexical scope via `let` bindings (`LetExpr`).
2.  **ProcLang:** Extends LetLang with first-class procedures (`ProcExpr`), procedure calls (`CallExpr`), and currying.
3.  **LetRecLang:** Extends ProcLang with mutually recursive procedures using `letrec` (`LetRecExpr`).

## Python Files (Original)

*   `letlang.py`
*   `proclang.py`
*   `letreclang.py`
*   `nameless.py` (Not yet ported)

## Go Files (Converted)

*   `expr.go`: Defines the base `Expr` interface and related utilities (`ExprEq`, `AnyToExpr`). Shared across language variants.
*   `eval.go`: Defines the base `Evaluator` interface and `BaseEval` struct using embedding for inheritance. Shared across language variants.
*   `letlang.go`: Defines AST structs (`LitExpr`, `VarExpr`, `OpExpr`, `IfExpr`, `IsZeroExpr`, `LetExpr`, `TupleExpr`) and the `LetLangEval` evaluator.
*   `proclang.go`: Defines AST structs (`ProcExpr`, `CallExpr`, `BoundProc`) and the `ProcLangEval` evaluator, embedding `LetLangEval`.
*   `letreclang.go`: Defines the `LetRecExpr` AST struct and the `LetRecLangEval` evaluator, embedding `ProcLangEval`.
*   `testutils.go`: Helper functions (`RunTest`, `setOpFuncs`) for testing Chapter 3 evaluators.
*   `letlang_test.go`, `proclang_test.go`, `letreclang_test.go`, `expr_test.go`: Unit tests covering evaluation logic, expression equality (`ExprEq`), and pretty-printing (`Printable`) for Chapter 3 constructs.

## Status

*   The Go conversion for `letlang`, `proclang`, and `letreclang` (evaluation, AST, equality, printing) is considered **complete** and tested.
*   The `nameless.py` (likely de Bruijn index representation and translation) has **not** been ported.
*   AST representation uses Go structs implementing the `Expr` interface.
*   Evaluation uses embedded structs for an object-oriented feel.
*   Procedure application correctly handles lexical scope and currying.
*   `letrec` evaluation correctly sets up the environment for mutual recursion.
*   `ExprEq` uses reflection for dispatching equality checks.
*   `Printable` provides indented tree views for debugging.
