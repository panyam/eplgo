package chapter5

import (
	"fmt"
	"reflect" // For detailed error message if needed

	epl "github.com/panyam/eplgo" // For potentially checking if value is Expr
	// For potentially checking if value is Expr
)

// RaisedError is a custom error type used to signal exceptions
// raised by the 'raise' expression within the interpreter.
// It wraps the actual value that was raised.
type RaisedError struct {
	Value any // The value passed to the 'raise' expression.
}

// Error implements the standard Go error interface.
func (e RaisedError) Error() string {
	// Provide a helpful error message including the raised value's representation.
	var valueRepr string
	if expr, ok := e.Value.(Expr); ok {
		// If the value is an Expr, use its Repr() method.
		valueRepr = expr.Repr()
	} else if ref, ok := e.Value.(*epl.Ref[any]); ok {
		// If it's a reference, show the value inside
		if exprVal, okVal := ref.Value.(Expr); okVal {
			valueRepr = fmt.Sprintf("Ref(%s)", exprVal.Repr())
		} else {
			valueRepr = fmt.Sprintf("Ref(%v:%T)", ref.Value, ref.Value)
		}
	} else {
		// Otherwise, use default formatting.
		valueRepr = fmt.Sprintf("%v:%T", e.Value, e.Value)
	}
	return fmt.Sprintf("raised value: %s", valueRepr)
}

// Is checks if the target error is a RaisedError.
// This allows using errors.Is(err, RaisedError{}) patterns if desired, although
// direct type assertion is more common for extracting the value.
func (e RaisedError) Is(target error) bool {
	_, ok := target.(RaisedError)
	// Check if target is the zero value RaisedError used for type checking
	if !ok {
		rv := reflect.ValueOf(target)
		if rv.Kind() == reflect.Struct && rv.Type() == reflect.TypeOf(RaisedError{}) && rv.IsZero() {
			ok = true
		}
	}
	return ok
}
