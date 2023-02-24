package errors

// Throw joins all errs to one and panics with it.
//
// The first element of errs is the key error.
// If the key error is nil, Throw does nothing.
func Throw(key error, errs ...error) {
	if key == nil {
		return
	}

	if he, ok := key.(*handlerError); ok {
		he.errors = append(he.errors, errs...)
		panic(he)
	}

	if len(errs) == 0 {
		panic(Join(key))
	}

	allErrors := make([]error, 0, len(errs)+1)
	allErrors = append(allErrors, key)
	allErrors = append(allErrors, errs...)
	panic(Join(allErrors...))
}

func Rethrow(err error) {
	e := recover()

	if e == nil {
		return
	}

	if he, ok := e.(*handlerError); ok {
		Throw(err, he.errors...)
		return
	}

	if e, ok := e.(error); ok {
		Throw(err, e)
		return
	}

	panic(e)
}

// Thrower wraps original function call results with error handling.
type Thrower interface {
	Throw(err error)
}

type thrower struct {
	err error
}

var _ Thrower = thrower{}

func (t thrower) Throw(err error) {
	if err != nil && t.err != nil {
		Throw(err, t.err)
	}
}

// If checks the error value returned by a function call and returns a thrower.
// If the function call returns an error,
// calling thrower's Throw(err) will join err with the error.
//
// If is useful when we want to return a business error, e.g. error code, with a cause.
// If it's not the case, use Check instead.
func If(values ...any) Thrower {
	err := checkWithError(values)
	return thrower{
		err: err,
	}
}

// Thrower1 wraps original function call results with error handling.
type Thrower1[T1 any] interface {
	Throw(err error) T1
}

type thrower1[T1 any] struct {
	err error
	In1 T1
}

var _ Thrower1[int] = &thrower1[int]{}

func (t *thrower1[T1]) Throw(err error) T1 {
	if t.err != nil && err != nil {
		Throw(err, t.err)
	}

	return t.In1
}

// If1 stores the return value of a function call and returns a thrower.
// If the function call returns an error,
// calling thrower's Throw(err) will join err with the error.
//
// If1 is useful when we want to return a business error, e.g. error code, with a cause.
// If it's not the case, use Check1 instead.
func If1[T1 any](in1 T1, err error) Thrower1[T1] {
	return &thrower1[T1]{
		err: err,
		In1: in1,
	}
}

// Thrower2 wraps original function call results with error handling.
type Thrower2[T1, T2 any] interface {
	Throw(err error) (T1, T2)
}

type thrower2[T1, T2 any] struct {
	err error
	In1 T1
	In2 T2
}

var _ Thrower2[int, int] = &thrower2[int, int]{}

func (t *thrower2[T1, T2]) Throw(err error) (T1, T2) {
	if t.err != nil && err != nil {
		Throw(err, t.err)
	}

	return t.In1, t.In2
}

// If2 stores the return value of a function call and returns a thrower.
// If the function call returns an error,
// calling thrower's Throw(err) will join err with the error.
//
// If2 is useful when we want to return a business error, e.g. error code, with a cause.
// If it's not the case, use Check2 instead.
func If2[T1, T2 any](in1 T1, in2 T2, err error) Thrower2[T1, T2] {
	return &thrower2[T1, T2]{
		err: err,
		In1: in1,
		In2: in2,
	}
}

// Thrower3 wraps original function call results with error handling.
type Thrower3[T1, T2, T3 any] interface {
	Throw(err error) (T1, T2, T3)
}

type thrower3[T1, T2, T3 any] struct {
	err error
	In1 T1
	In2 T2
	In3 T3
}

var _ Thrower3[int, int, int] = &thrower3[int, int, int]{}

func (t *thrower3[T1, T2, T3]) Throw(err error) (T1, T2, T3) {
	if t.err != nil && err != nil {
		Throw(err, t.err)
	}

	return t.In1, t.In2, t.In3
}

// If3 stores the return value of a function call and returns a thrower.
// If the function call returns an error,
// calling thrower's Throw(err) will join err with the error.
//
// If3 is useful when we want to return a business error, e.g. error code, with a cause.
// If it's not the case, use Check3 instead.
func If3[T1, T2, T3 any](in1 T1, in2 T2, in3 T3, err error) Thrower3[T1, T2, T3] {
	return &thrower3[T1, T2, T3]{
		err: err,
		In1: in1,
		In2: in2,
		In3: in3,
	}
}

// Thrower4 wraps original function call results with error handling.
type Thrower4[T1, T2, T3, T4 any] interface {
	Throw(err error) (T1, T2, T3, T4)
}

type thrower4[T1, T2, T3, T4 any] struct {
	err error
	In1 T1
	In2 T2
	In3 T3
	In4 T4
}

var _ Thrower4[int, int, int, int] = &thrower4[int, int, int, int]{}

func (t *thrower4[T1, T2, T3, T4]) Throw(err error) (T1, T2, T3, T4) {
	if t.err != nil && err != nil {
		Throw(err, t.err)
	}

	return t.In1, t.In2, t.In3, t.In4
}

// If4 stores the return value of a function call and returns a thrower.
// If the function call returns an error,
// calling thrower's Throw(err) will join err with the error.
//
// If4 is useful when we want to return a business error, e.g. error code, with a cause.
// If it's not the case, use Check4 instead.
func If4[T1, T2, T3, T4 any](in1 T1, in2 T2, in3 T3, in4 T4, err error) Thrower4[T1, T2, T3, T4] {
	return &thrower4[T1, T2, T3, T4]{
		err: err,
		In1: in1,
		In2: in2,
		In3: in3,
		In4: in4,
	}
}

// Thrower5 wraps original function call results with error handling.
type Thrower5[T1, T2, T3, T4, T5 any] interface {
	Throw(err error) (T1, T2, T3, T4, T5)
}

type thrower5[T1, T2, T3, T4, T5 any] struct {
	err error
	In1 T1
	In2 T2
	In3 T3
	In4 T4
	In5 T5
}

var _ Thrower5[int, int, int, int, int] = &thrower5[int, int, int, int, int]{}

func (t *thrower5[T1, T2, T3, T4, T5]) Throw(err error) (T1, T2, T3, T4, T5) {
	if t.err != nil && err != nil {
		Throw(err, t.err)
	}

	return t.In1, t.In2, t.In3, t.In4, t.In5
}

// If5 stores the return value of a function call and returns a thrower.
// If the function call returns an error,
// calling thrower's Throw(err) will join err with the error.
//
// If5 is useful when we want to return a business error, e.g. error code, with a cause.
// If it's not the case, use Check5 instead.
func If5[T1, T2, T3, T4, T5 any](in1 T1, in2 T2, in3 T3, in4 T4, in5 T5, err error) Thrower5[T1, T2, T3, T4, T5] {
	return &thrower5[T1, T2, T3, T4, T5]{
		err: err,
		In1: in1,
		In2: in2,
		In3: in3,
		In4: in4,
		In5: in5,
	}
}

// Thrower6 wraps original function call results with error handling.
type Thrower6[T1, T2, T3, T4, T5, T6 any] interface {
	Throw(err error) (T1, T2, T3, T4, T5, T6)
}

type thrower6[T1, T2, T3, T4, T5, T6 any] struct {
	err error
	In1 T1
	In2 T2
	In3 T3
	In4 T4
	In5 T5
	In6 T6
}

var _ Thrower6[int, int, int, int, int, int] = &thrower6[int, int, int, int, int, int]{}

func (t *thrower6[T1, T2, T3, T4, T5, T6]) Throw(err error) (T1, T2, T3, T4, T5, T6) {
	if t.err != nil && err != nil {
		Throw(err, t.err)
	}

	return t.In1, t.In2, t.In3, t.In4, t.In5, t.In6
}

// If6 stores the return value of a function call and returns a thrower.
// If the function call returns an error,
// calling thrower's Throw(err) will join err with the error.
//
// If6 is useful when we want to return a business error, e.g. error code, with a cause.
// If it's not the case, use Check6 instead.
func If6[T1, T2, T3, T4, T5, T6 any](in1 T1, in2 T2, in3 T3, in4 T4, in5 T5, in6 T6, err error) Thrower6[T1, T2, T3, T4, T5, T6] {
	return &thrower6[T1, T2, T3, T4, T5, T6]{
		err: err,
		In1: in1,
		In2: in2,
		In3: in3,
		In4: in4,
		In5: in5,
		In6: in6,
	}
}

// Thrower7 wraps original function call results with error handling.
type Thrower7[T1, T2, T3, T4, T5, T6, T7 any] interface {
	Throw(err error) (T1, T2, T3, T4, T5, T6, T7)
}

type thrower7[T1, T2, T3, T4, T5, T6, T7 any] struct {
	err error
	In1 T1
	In2 T2
	In3 T3
	In4 T4
	In5 T5
	In6 T6
	In7 T7
}

var _ Thrower7[int, int, int, int, int, int, int] = &thrower7[int, int, int, int, int, int, int]{}

func (t *thrower7[T1, T2, T3, T4, T5, T6, T7]) Throw(err error) (T1, T2, T3, T4, T5, T6, T7) {
	if t.err != nil && err != nil {
		Throw(err, t.err)
	}

	return t.In1, t.In2, t.In3, t.In4, t.In5, t.In6, t.In7
}

// If7 stores the return value of a function call and returns a thrower.
// If the function call returns an error,
// calling thrower's Throw(err) will join err with the error.
//
// If7 is useful when we want to return a business error, e.g. error code, with a cause.
// If it's not the case, use Check7 instead.
func If7[T1, T2, T3, T4, T5, T6, T7 any](in1 T1, in2 T2, in3 T3, in4 T4, in5 T5, in6 T6, in7 T7, err error) Thrower7[T1, T2, T3, T4, T5, T6, T7] {
	return &thrower7[T1, T2, T3, T4, T5, T6, T7]{
		err: err,
		In1: in1,
		In2: in2,
		In3: in3,
		In4: in4,
		In5: in5,
		In6: in6,
		In7: in7,
	}
}

// Thrower8 wraps original function call results with error handling.
type Thrower8[T1, T2, T3, T4, T5, T6, T7, T8 any] interface {
	Throw(err error) (T1, T2, T3, T4, T5, T6, T7, T8)
}

type thrower8[T1, T2, T3, T4, T5, T6, T7, T8 any] struct {
	err error
	In1 T1
	In2 T2
	In3 T3
	In4 T4
	In5 T5
	In6 T6
	In7 T7
	In8 T8
}

var _ Thrower8[int, int, int, int, int, int, int, int] = &thrower8[int, int, int, int, int, int, int, int]{}

func (t *thrower8[T1, T2, T3, T4, T5, T6, T7, T8]) Throw(err error) (T1, T2, T3, T4, T5, T6, T7, T8) {
	if t.err != nil && err != nil {
		Throw(err, t.err)
	}

	return t.In1, t.In2, t.In3, t.In4, t.In5, t.In6, t.In7, t.In8
}

// If8 stores the return value of a function call and returns a thrower.
// If the function call returns an error,
// calling thrower's Throw(err) will join err with the error.
//
// If8 is useful when we want to return a business error, e.g. error code, with a cause.
// If it's not the case, use Check8 instead.
func If8[T1, T2, T3, T4, T5, T6, T7, T8 any](in1 T1, in2 T2, in3 T3, in4 T4, in5 T5, in6 T6, in7 T7, in8 T8, err error) Thrower8[T1, T2, T3, T4, T5, T6, T7, T8] {
	return &thrower8[T1, T2, T3, T4, T5, T6, T7, T8]{
		err: err,
		In1: in1,
		In2: in2,
		In3: in3,
		In4: in4,
		In5: in5,
		In6: in6,
		In7: in7,
		In8: in8,
	}
}
