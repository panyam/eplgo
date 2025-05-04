# Next Steps for EPLGo Conversion

This document outlines the remaining tasks and suggested next steps for completing the Python-to-Go conversion of the EPL implementations.

## Outstanding Work

1.  **Port Chapter 4 (State & Laziness):**
    *   **Concepts:** Mutable state via explicit references (`newref`, `deref`, `setref`), implicit references/assignment (`set`), sequencing (`begin`/`block`), lazy evaluation (`lazy`, `thunk`).
    *   **Python Files:** `expreflang.py`, `impreflang.py`, `lazylang.py`.
    *   **Go Tasks:**
        *   Define corresponding AST structs (e.g., `RefExpr`, `DeRefExpr`, `SetRefExpr`, `BlockExpr`, `AssignExpr`, `LazyExpr`, `Thunk`). Extend the `Expr` interface ecosystem.
        *   Implement evaluators (`ExpRefLangEval`, `ImpRefLangEval`, `LazyLangEval`) extending the Chapter 3 hierarchy. This will likely involve modifying how the `Env` stores values (perhaps distinguishing between direct values and references/locations) or introducing a separate `Store`.
        *   Write unit tests for evaluation, equality, and printing.

2.  **Implement Go Parser:**
    *   **Need:** Currently, tests construct ASTs directly. A parser is needed to process source code strings (like those in Python test cases) into the Go `Expr` AST structure.
    *   **Approach:**
        *   Choose a parsing strategy/library (e.g., standard library `text/scanner`, `go/scanner`+`go/parser` if Go-like syntax is desired, or a parser generator like `goyacc`, `antlr`, or a recursive descent parser).
        *   Define the grammar for the EPL languages.
        *   Implement parsing functions to build the Go AST nodes defined in `chapter3`, `chapter4`, etc.
    *   **Integration:** Update tests to parse input strings instead of using direct AST constructors.

3.  **Port Chapter 5 (Continuations & Exceptions):**
    *   **Concepts:** Continuation-Passing Style (CPS) evaluation, exception handling (`try`, `raise`). Trampolining might be needed if implementing direct CPS for deep recursion (`trampoline.py`).
    *   **Python Files:** `continuations.py`, `trylang.py`, `trampoline.py`.
    *   **Go Tasks:**
        *   Define AST structs (`TryExpr`, `RaiseExpr`).
        *   Define `Continuation` types/interfaces.
        *   Implement a CPS-based evaluator. This is a significant shift from the current direct-style evaluation. Consider if direct style with explicit error handling (Go's `error` interface) is a viable alternative or if CPS is strictly required by the EPL book's approach for this chapter.
        *   If implementing CPS, address potential stack overflow (trampolining or other techniques).
        *   Write unit tests.

4.  **Port Chapter 7 (Types):**
    *   **Concepts:** Type representation (base types, function types, tuple types, tagged types, type variables), type checking (explicitly typed language), type inference.
    *   **Python Files:** `typed.py`, `inferred.py`, `utils.py`.
    *   **Go Tasks:**
        *   Define Go structs/interfaces for representing `Type` variants (similar to `Expr`).
        *   Implement type environments (similar to `epl.Env[T]`, but storing types).
        *   Implement the `TypeOf` logic (likely as a separate struct/method set) for both explicit and inferred type checking, potentially using type switches or visitors on the `Expr` AST.
        *   Implement type unification for inference.
        *   Define and implement `TypeError`.
        *   Write unit tests for type checking/inference.

5.  **Port/Refine Testing Infrastructure:**
    *   **Python Files:** `tests/externs.py`, `tests/settings.py`, `tests/utils.py`, `tests/**/cases.py`.
    *   **Go Tasks:**
        *   Implement Go equivalents for defining external functions/operators used in tests (currently partially done in `testutils.go`).
        *   Port test cases from `cases.py` for Chapters 4, 5, 7.
        *   Decide if the `settings.py` context management is needed in Go tests.
        *   Refine test runners (`RunTest`, etc.) to handle parsing and different evaluation/type-checking modes.

6.  **Port Chapter 3 `nameless.py` (Optional):**
    *   **Concepts:** Nameless representation (e.g., de Bruijn indices), translation from named AST.
    *   **Tasks:** Define nameless AST, implement translator, implement nameless evaluator. Assess if this is a priority.

7.  **Code Cleanup & Refinements:**
    *   Review error handling (use Go idioms like `error` interface where appropriate, instead of panics in non-fatal situations).
    *   Improve comments and documentation within the code.
    *   Standardize formatting (`gofmt`/`goimports`).
    *   Optimize performance if necessary (e.g., reconsider reflection in `ExprEq` if it becomes a bottleneck).

## Suggested Order

1.  **Chapter 4:** Builds directly on Chapter 3's evaluation structure, introducing state. Implementing the `Store` concept early is beneficial.
2.  **Parser:** Having a parser makes testing Chapters 5 and 7 much easier, as you can use the original source strings.
3.  **Chapter 5:** Introduces significant changes to evaluation (CPS or error handling). Tackle this after state is handled.
4.  **Chapter 7:** Type checking builds on the AST and requires its own environment and logic. Can be done mostly independently after the AST for relevant chapters is stable.
5.  **Testing Infrastructure:** Port test cases and utilities incrementally as each chapter is ported.
6.  **Nameless / Cleanup:** Lower priority, can be done last or as needed.
