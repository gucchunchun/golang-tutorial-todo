package cmd

type updateInput struct {
	Name     *string `json:"name"`
	DueDate  *string `json:"due_date"`
	Status   *string `json:"status"`
	ClearDue *bool   `json:"clear_due_date"`
}
