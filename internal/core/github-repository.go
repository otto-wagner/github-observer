package core

import (
	"strings"
)

type Repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Branch   string `json:"branch"`
	Owner    string `json:"owner"`
	HtmlUrl  string `json:"html_url"`
}

func ToRepository(repo string) Repository {
	first := strings.Split(repo, "/")
	if len(first) != 2 {
		return Repository{}
	}
	second := strings.Split(first[1], "@")
	if len(second) != 2 {
		return Repository{}
	}
	return Repository{
		Owner:  first[0],
		Name:   second[0],
		Branch: second[1],
	}
}
