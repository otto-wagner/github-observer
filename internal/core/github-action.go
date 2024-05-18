package core

import (
	"github.com/google/go-github/v61/github"
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

func ConvertToGitAction(event github.CheckRunEvent) GitAction {
	return GitAction{
		Action: event.GetAction(),
		CheckRun: CheckRun{
			Id:          event.GetCheckRun().GetID(),
			Name:        event.GetCheckRun().GetName(),
			HtmlUrl:     event.GetCheckRun().GetHTMLURL(),
			Status:      event.GetCheckRun().GetStatus(),
			Conclusion:  event.GetCheckRun().GetConclusion(),
			StartedAt:   event.GetCheckRun().GetStartedAt().Time,
			CompletedAt: event.GetCheckRun().GetCompletedAt().Time,
		},
		Repository: Repository{
			FullName: event.GetRepo().GetFullName(),
			HtmlUrl:  event.GetRepo().GetHTMLURL(),
		},
	}
}
