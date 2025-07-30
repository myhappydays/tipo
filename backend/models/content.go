package models

type ContentResponse struct {
	Title     string   `json:"title"`
	Source    string   `json:"source"`
	Keyword   string   `json:"keyword"`
	PubDate   string   `json:"pubDate"`
	Sentences []string `json:"sentences"`
}
