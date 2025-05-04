
# Project Summary: EPLGo

## Goal

This project aims to convert the Python implementations of interpreters and type checkers for the small languages defined in the book "Elements of Programming Languages" (EPL, Friedman and Wand, 3rd Ed) into Go.

## Original Python Source

The original Python code (found in `chapter*/*.py` and `tests/**/*.py`) utilizes:
*   `taggedunion` library for representing Abstract Syntax Trees (ASTs) and Types.
*   `CaseMatcher` for dispatching logic in evaluators and type checkers.
*   Progressive language features introduced chapter by chapter (Let, Proc, LetRec, Refs, Assignment, Lazy Eval, Continuations, Exceptions, Types).
*   `pytest` based testing infrastructure.

## Go Approach

The Go conversion employs:
*   Go interfaces (`chapter3/Expr`) and structs (`*LitExpr`, `*VarExpr`, etc.) to represent the AST.
*   Go generics for the scoped runtime environment (`epl.Env[T]` in `env.go`).
*   Embedded structs and method overriding (`BaseEval`, `LetLangEval`, `ProcLangEval`, `LetRecLangEval`) for the evaluator hierarchy. Type switches or further interface embedding might be used for more complex dispatch later.
*   Reflection (`reflect` package) is used in `ExprEq` for dynamic equality checks based on concrete types' `Eq` methods.
*   Standard Go testing (`testing` package, `testify/assert`) for unit tests (`*_test.go`).
*   Utility functions for common tasks and pretty-printing (`common.go` and `printable.go`).

## Current Status (As of Chapter 3 Completion)

*   **Core Infrastructure:**
    *   Generic Scoped Environment (`env.go`): Complete and functional.
    *   Common Utilities (`common.go`): Basic utilities for printing, dict creation, list equality are present.
    *   Go Module (`go.mod`, `go.sum`): Set up.
*   **Chapter 3 (Let, Proc, LetRec):**
    *   AST (`expr.go`, `letlang.go`, `proclang.go`, `letreclang.go`): Structs for all Ch3 expressions are defined.
    *   Evaluation (`eval.go`, `letlang.go`, `proclang.go`, `letreclang.go`): Evaluators for Let, Proc, and LetRec languages are implemented, including handling lexical scope and currying.
    *   Equality (`expr.go`): `ExprEq` implemented using reflection to call specific `Eq` methods. Concrete `Eq` methods implemented for Ch3 types.
    *   Printing (`common.go`, expr structs): `Printable` interface and implementations allow for indented tree printing of expressions.
    *   Testing (`chapter3/*_test.go`): Unit tests covering evaluation, equality, and printing for Chapter 3 constructs are implemented and passing.
*   **Chapters 4, 5, 7:** Implementations (AST, Eval, Type Checking) are **not yet ported** from Python.
*   **Parser:** A Go parser to convert source code strings into the Go AST is **not yet implemented**. Tests currently construct the AST directly.
*   **Testing Infrastructure:** Python test utilities (`tests/settings.py`, `tests/utils.py`, `tests/externs.py`) are **not yet ported**. Go tests currently use basic test runners and direct AST construction.

## Key Go Components

*   `env.go`: Generic, scoped environment implementation.
*   `common.go`: Utility functions, `Printable` struct for debugging output.
*   `chapter3/expr.go`: Defines the base `Expr` interface and core equality/utility functions (`ExprEq`, `ExprListEq`, `AnyToExpr`).
*   `chapter3/eval.go`: Defines the base `Evaluator` interface structure using embedding for extension.
*   `chapter3/letlang.go`: AST node and evaluator logic for the basic Let language.
*   `chapter3/proclang.go`: Extends letlang with procedures and calls (incl. currying).
*   `chapter3/letreclang.go`: Extends proclang with mutual recursion via `letrec`.
*   `chapter3/*_test.go`: Go unit tests for Chapter 3 functionality.
