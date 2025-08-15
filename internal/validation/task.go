package validation

type CreateTaskInput struct {
	Name    string  `json:"name" validate:"required"`
	DueDate *string `json:"due_date" validate:"omitempty,datetime=2006-01-02"`
}

func ValidateCreateTaskInput(input CreateTaskInput) Errors {
	return Validate(input)
}
