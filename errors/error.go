package errors

// Error is an interface that wraps the error interface.
type Error interface {
	error

	// If the last element of values is an error, it will be thrown.
	// The thrown error will be this Error interface wrapping the last element of values as its cause.
	Check(values ...any)
}

type stringError struct {
	msg string
}

var _ Error = new(stringError)

func (se *stringError) Error() string {
	return se.msg
}

func (se *stringError) Check(values ...any) {
	err := checkWithError(values)

	if err != nil {
		Throw(se, err)
	}
}

// New returns a new Error interface with the given message.
func New(msg string) Error {
	return &stringError{
		msg: msg,
	}
}
