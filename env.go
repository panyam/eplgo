package epl

import (
	"fmt"
)

// Env[T] holds the runtime values for identifiers (variables, functions, components).
// Supports basic scoping via the 'outer' environment.
type Env[T any] struct {
	store map[string]*Ref[T]
	outer *Env[T]
}

// NewEnv[T] creates a new environment nested within an outer one.
// If outer is nil then returns a fresh top-level environment.
// Useful for function calls or block scopes.
func NewEnv[T any](outer *Env[T]) *Env[T] {
	s := make(map[string]*Ref[T])
	return &Env[T]{store: s, outer: outer}
}

// Get retrieves a value by name. It checks the current environment first,
// then recursively checks outer environments.
func (e *Env[T]) GetRef(name string) *Ref[T] {
	ref, ok := e.store[name]
	if (!ok || ref == nil) && e.outer != nil {
		// Not found here, try the outer scope
		ref = e.outer.GetRef(name)
	}
	return ref
}

func (e *Env[T]) Get(name string) (out T, found bool) {
	ref := e.GetRef(name)
	if ref != nil {
		out = ref.Value
		found = true
	}
	return
}

func (e *Env[T]) Set(key string, value T) {
	// Create or update the VarState
	e.store[key] = &Ref[T]{Value: value}
}

func (e *Env[T]) SetMany(kvpairs map[string]T) {
	for k, v := range kvpairs {
		e.Set(k, v)
	}
}

// String representation for debugging
func (e *Env[T]) String() string {
	keys := make([]string, 0, len(e.store))
	for k := range e.store {
		keys = append(keys, k)
	}
	return fmt.Sprintf("Env[T]{store: %v, outer: %v}", keys, e.outer != nil)
}

// References to values
type Ref[T any] struct {
	Value T
}
