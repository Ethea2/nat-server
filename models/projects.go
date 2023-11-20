package models

type Projects struct {
	ID         int      `json:"id"`
	UserID     int      `json:"userID"`
	Title      string   `json:"title"`
	Body       string   `json:"body"`
	Image      string   `json:"image"`
	Techstacks []string `json:"techstacks"`
}
