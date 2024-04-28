package core

type Repository struct {
	Name    string `json:"name"`
	Owner   string `json:"owner"`
	HtmlUrl string `json:"html_url"`
}
