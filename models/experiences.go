package models

type Experiences struct {
	ID         int      `json:"id,omitempty"`
	UserID     int      `json:"userID,omitempty"`
	Title      string   `json:"title,omitempty"`
	Body       string   `json:"body,omitempty"`
	Position   string   `json:"position,omitempty"`
	Techstacks []string `json:"techstacks,omitempty"`
}
