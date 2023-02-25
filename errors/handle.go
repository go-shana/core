package errors

import "fmt"

// Handle recovers from panic and returns error to the target if panic is an error.
func Handle(target *error) {
	r := recover()

	if r == nil {
		return
	}

	if err, ok := r.(error); ok {
		*target = err
		return
	}

	*target = New(fmt.Sprintf("errors: uncaught panic: %v", r))
}
