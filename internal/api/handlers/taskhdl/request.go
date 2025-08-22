package taskhdl

type createRequest struct {
	Name  string  `json:"name"`
	DueAt *string `json:"due_at"`
	/*
		Reference: O'REILLY「実用GO言語」8.1 p.171
		非公開のフィールドはエンコード, デコードの対象外。
	*/
	// Secret *string `json:"secret"`
}

type updateRequest struct {
	Name     *string `json:"name"`
	DueAt    *string `json:"due_at"`
	Status   *string `json:"status"`
	ClearDue *bool   `json:"clear_due"`
}
