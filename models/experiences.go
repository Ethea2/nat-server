package models

type Experiences struct {
	ID         int      `json:"id"`
	UserID     int      `json:"userID"`
	Title      string   `json:"title"`
	Body       string   `json:"body"`
	Position   string   `json:"position"`
	Techstacks []string `json:"teckstacks"`
}
