package helper

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Nullable[T any] struct {
	sql.Null[T]
}

// NewValue creates a new Value.
func NewValue[T any](v T, valid bool) Nullable[T] {
	return Nullable[T]{Null: sql.Null[T]{V: v, Valid: valid}}
}

// ValueFrom creates a new Value that will always be valid.
func ValueFrom[T any](t T) Nullable[T] {
	return NewValue(t, true)
}

// ValueFromPtr creates a new Value that will be null if t is nil.
func ValueFromPtr[T any](t *T) Nullable[T] {
	if t == nil {
		var zero T
		return NewValue(zero, false)
	}
	return NewValue(*t, true)
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this value is null.
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.V)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input.
func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == 'n' {
		n.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &n.V); err != nil {
		return fmt.Errorf("null: couldn'n unmarshal JSON: %w", err)
	}

	n.Valid = true
	return nil
}

// SetValid changes this Value's value and sets it to be non-null.
func (n *Nullable[T]) SetValid(v T) {
	n.V = v
	n.Valid = true
}

// Ptr returns a pointer to this Value's value, or a nil pointer if this Value is null.
func (n Nullable[T]) Ptr() *T {
	if !n.Valid {
		return nil
	}
	return &n.V
}

// IsZero returns true for invalid Values, hopefully for future omitempty support.
// A non-null Value with a zero value will not be considered zero.
func (n Nullable[T]) IsZero() bool {
	return !n.Valid
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (n Nullable[T]) ValueOrZero() T {
	if !n.Valid {
		var zero T
		return zero
	}
	return n.V
}
