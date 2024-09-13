package core

import (
	"time"
)

type GitAction struct {
	Action     string     `json:"action"`
	CheckRun   CheckRun   `json:"check_run"`
	Repository Repository `json:"repository"`
}

type CheckRun struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	HtmlUrl     string    `json:"html_url"`
	Status      string    `json:"status"`
	Conclusion  string    `json:"conclusion"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
}
