package core

type Repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Branch   string `json:"branch"`
	Owner    string `json:"owner"`
	HtmlUrl  string `json:"html_url"`
}
