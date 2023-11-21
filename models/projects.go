package models

type Projects struct {
	ID         int      `json:"id,omitempty"`
	UserID     int      `json:"userID,omitempty"`
	Title      string   `json:"title,omitempty"`
	Body       string   `json:"body,omitempty"`
	Image      string   `json:"image,omitempty"`
	Techstacks []string `json:"techstacks,omitempty"`
}
