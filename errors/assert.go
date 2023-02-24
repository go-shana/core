package errors

var errAssert = New("errors: assert failed")

// Assert panics if any of the exprs is false.
func Assert(exprs ...bool) {
	for _, expr := range exprs {
		if !expr {
			Throw(errAssert)
		}
	}
}
