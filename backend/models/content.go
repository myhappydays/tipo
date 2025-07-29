package models

type ContentResponse struct {
	Title     string   `json:"title"`
	Source    string   `json:"source"`
	Category  string   `json:"category"`
	Image     string   `json:"image"`
	Sentences []string `json:"sentences"`
}
