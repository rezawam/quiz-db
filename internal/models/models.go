package models

type Question struct {
	ID         int      `json:"id"`
	Question   string   `json:"question"`
	Type       string   `json:"type"` // multiple, date
	Options    []string `json:"options,omitempty"` // not empty only if type is "multiple"
	Answer     string   `json:"corrext_answer"`
	Category   string   `json:"category,omitempty"`
	Difficulty string   `json:"difficulty,omitempty"`
}
