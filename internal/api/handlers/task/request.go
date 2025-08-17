package task

type createRequest struct {
	Name  string  `json:"name"`
	DueAt *string `json:"due_at"`
}

type updateRequest struct {
	Name     *string `json:"name"`
	DueAt    *string `json:"due_at"`
	Status   *string `json:"status"`
	ClearDue *bool   `json:"clear_due"`
}
