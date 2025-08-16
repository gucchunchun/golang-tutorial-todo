package apperr

type Code int

const (
	CodeUnknown Code = iota
	CodeNotFound
	CodeInvalid
	CodeConflict
	CodeUnauthorized
	CodeForbidden
)

type Error struct {
	Code    Code
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error() + ": " + e.Message
	}
	return e.Message
}
func (e *Error) Unwrap() error { return e.Err }

func E(code Code, msg string, err error) *Error { return &Error{Code: code, Message: msg, Err: err} }
