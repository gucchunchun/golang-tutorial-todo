// internal/validation/errors.go
package validation

import (
	"strings"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Errors is a collection of field errors and implements error.
type Errors []FieldError

func (e Errors) Error() string {
	if len(e) == 0 {
		return ""
	}
	b := strings.Builder{}
	for _, fe := range e {
		b.WriteString(fe.Field)
		b.WriteString(": ")
		b.WriteString(fe.Message)
		b.WriteByte('\n')
	}
	return b.String()
}

func (e *Errors) Add(field, message string) {
	*e = append(*e, FieldError{Field: field, Message: message})
}

func (e Errors) Has() bool { return len(e) > 0 }
