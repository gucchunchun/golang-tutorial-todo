package apperr

import (
	"net/http"
)

type Code int

const (
	CodeUnknown Code = iota
	CodeNotFound
	CodeInvalid
)

func (c Code) String() string {
	switch c {
	case CodeNotFound:
		return "NOT_FOUND"
	case CodeInvalid:
		return "INVALID"
	default:
		return "UNKNOWN"
	}
}

func (c Code) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

func (c Code) HTTPStatus() int {
	switch c {
	case CodeNotFound:
		return http.StatusNotFound
	case CodeInvalid:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

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
