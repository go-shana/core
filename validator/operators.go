package validator

import "github.com/go-shana/core/errors"

type invertRange[T any] struct {
	ranger Ranger[T]
}

var _ Ranger[int] = invertRange[int]{}

func (invert invertRange[T]) InRange(v T) bool {
	return !invert.ranger.InRange(v)
}

// Not returns a new Ranger which inverts the result of the ranger.
func Not[T any](ranger Ranger[T]) Ranger[T] {
	return invertRange[T]{
		ranger: ranger,
	}
}

type andRange[T any] struct {
	rangers []Ranger[T]
}

var _ Ranger[int] = andRange[int]{}

func (and andRange[T]) InRange(v T) bool {
	for _, r := range and.rangers {
		if !r.InRange(v) {
			return false
		}
	}

	return true
}

// And returns a new Ranger which returns true if all of the rangers return true.
func And[T any](rangers ...Ranger[T]) Ranger[T] {
	return andRange[T]{
		rangers: rangers,
	}
}

type orRange[T any] struct {
	rangers []Ranger[T]
}

var _ Ranger[int] = orRange[int]{}

func (or orRange[T]) InRange(v T) bool {
	for _, r := range or.rangers {
		if r.InRange(v) {
			return true
		}
	}

	return false
}

// Or returns a new Ranger which returns true if any one of the rangers returns true.
func Or[T any](rangers ...Ranger[T]) Ranger[T] {
	return orRange[T]{
		rangers: rangers,
	}
}

// Equal asserts that v must equal to target.
// Otherwise, it throws an error.
func Equal[T comparable](v, target T) {
	if v == target {
		return
	}

	errors.Throw(errValidator)
}

// NotEqual asserts that v must not equal to target.
// Otherwise, it throws an error.
func NotEqual[T comparable](v, target T) {
	if v != target {
		return
	}

	errors.Throw(errValidator)
}

// In asserts that v must be in all of the rangers.
// Otherwise, it throws an error.
func In[T any](v T, rangers ...Ranger[T]) {
	for _, ranger := range rangers {
		if !ranger.InRange(v) {
			errors.Throw(errValidator)
			return
		}
	}
}

// InAny asserts that v must be in one of the rangers.
// Otherwise, it throws an error.
func InAny[T any](v T, rangers ...Ranger[T]) {
	for _, ranger := range rangers {
		if ranger.InRange(v) {
			return
		}
	}

	errors.Throw(errValidator)
}
