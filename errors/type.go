package errors

// Error is an error struct used for most errors encountered.
type Error struct {
	Status  int
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

// HTTPStatus returns the status that the error represents.
func (e *Error) HTTPStatus() int {
	return e.Status
}
