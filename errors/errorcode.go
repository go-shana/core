package errors

// ErrorCode is an error with error code.
type ErrorCode[T any] interface {
	Error

	Code() T
}

type codeError[T any] struct {
	*stringError
	code T
}

var _ ErrorCode[int] = &codeError[int]{}

// NewErrorCode returns a new ErrorCode interface with the given code and message.
func NewErrorCode[T any](code T, msg string) ErrorCode[T] {
	return &codeError[T]{
		code:        code,
		stringError: &stringError{msg: msg},
	}
}

func (ce *codeError[T]) Error() string {
	return ce.stringError.Error()
}

func (ce *codeError[T]) Code() T {
	return ce.code
}
