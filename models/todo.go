package models

type Todo struct {
	ID              int64  `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	CreatorID       int64  `json:"creator_id"`
	CreatorUsername string `json:"creator_username"`
	CompletedAnyel  bool   `json:"completed_anyel"`
	CompletedAlexis bool   `json:"completed_alexis"`
	CreatedAt       string `json:"created_at"`
}
