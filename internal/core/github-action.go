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
			Name:        event.GetCheckRun().GetName(),
			HtmlUrl:     event.GetCheckRun().GetHTMLURL(),
			Status:      event.GetCheckRun().GetStatus(),
			Conclusion:  event.GetCheckRun().GetConclusion(),
			StartedAt:   event.GetCheckRun().GetStartedAt().Time,
			CompletedAt: event.GetCheckRun().GetCompletedAt().Time,
		},
		Repository: Repository{
			Name:    event.GetRepo().GetName(),
			HtmlUrl: event.GetRepo().GetHTMLURL(),
		},
	}
}

func (g *GitAction) ToMap() map[string]string {
	return map[string]string{
		"action":               g.Action,
		"check_run_name":       g.CheckRun.Name,
		"check_run_url":        g.CheckRun.HtmlUrl,
		"check_run_status":     g.CheckRun.Status,
		"check_run_conclusion": g.CheckRun.Conclusion,
		"check_run_started_at": g.CheckRun.StartedAt.String(),
		"check_run_completed":  g.CheckRun.CompletedAt.String(),
		"repository_name":      g.Repository.Name,
		"repository_url":       g.Repository.HtmlUrl,
	}
}
