package quotes

type Quote struct {
	Text   string `json:"q"`
	Author string `json:"a,omitempty"`
}
