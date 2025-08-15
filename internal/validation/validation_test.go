package validation_test

import (
	"testing"

	"golang/tutorial/todo/internal/validation"

	"github.com/stretchr/testify/assert"
)

type TestValidationRequest struct {
	Name     string `validate:"required"`
	Birthday string `validate:"required,datetime=2006-01-02"`
}

func TestValidationOK(t *testing.T) {
	errs := validation.Validate(TestValidationRequest{
		Name:     "John Doe",
		Birthday: "1990-01-01",
	})

	assert.Nil(t, errs)
}
func TestValidationOKPointer(t *testing.T) {
	errs := validation.Validate(&TestValidationRequest{
		Name:     "John Doe",
		Birthday: "1990-01-01",
	})

	assert.Nil(t, errs)
}

func TestValidationNil(t *testing.T) {
	var x any = nil
	errs := validation.Validate(x)

	assert.NotNil(t, errs)
}

func TestValidationNilPointer(t *testing.T) {
	var x any = nil
	errs := validation.Validate(&x)

	assert.NotNil(t, errs)
}

func TestValidationRequiredError(t *testing.T) {
	errs := validation.Validate(TestValidationRequest{
		Name:     "",
		Birthday: "1990-01-01",
	})

	assert.NotNil(t, errs)
}
func TestValidationRequiredErrorPointer(t *testing.T) {
	errs := validation.Validate(&TestValidationRequest{
		Name:     "",
		Birthday: "1990-01-01",
	})

	assert.NotNil(t, errs)
}

func TestValidationDateFormatError(t *testing.T) {
	errs := validation.Validate(TestValidationRequest{
		Name:     "",
		Birthday: "1990-01",
	})

	assert.NotNil(t, errs)
}
func TestValidationDateFormatErrorPointer(t *testing.T) {
	errs := validation.Validate(&TestValidationRequest{
		Name:     "",
		Birthday: "1990-01",
	})

	assert.NotNil(t, errs)
}
