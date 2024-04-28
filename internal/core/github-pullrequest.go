package core

import (
	"github.com/google/go-github/v61/github"
	"strconv"
)

type GitPullRequest struct {
	Action      string      `json:"action"`
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
}

type PullRequest struct {
	Number    int    `json:"number"`
	HtmlUrl   string `json:"html_url"`
	IssueUrl  string `json:"issue_url"`
	State     string `json:"state"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	User      User   `json:"user"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	ClosedAt  string `json:"closed_at"`
	MergedAt  string `json:"merged_at"`
}

func ConvertPREToGitPullRequest(event github.PullRequestEvent) GitPullRequest {
	return GitPullRequest{
		Action: event.GetAction(),
		PullRequest: PullRequest{
			Number:    event.GetPullRequest().GetNumber(),
			HtmlUrl:   event.GetPullRequest().GetHTMLURL(),
			State:     event.GetPullRequest().GetState(),
			Title:     event.GetPullRequest().GetTitle(),
			Body:      event.GetPullRequest().GetBody(),
			User:      User{Login: event.GetPullRequest().GetUser().GetLogin()},
			CreatedAt: event.GetPullRequest().GetCreatedAt().String(),
			UpdatedAt: event.GetPullRequest().GetUpdatedAt().String(),
			ClosedAt:  event.GetPullRequest().GetClosedAt().String(),
			MergedAt:  event.GetPullRequest().GetMergedAt().String(),
		},
		Repository: Repository{
			Name:    event.GetRepo().GetName(),
			HtmlUrl: event.GetRepo().GetHTMLURL(),
		},
	}
}

func ConvertPRToGitPullRequest(r Repository, p github.PullRequest) GitPullRequest {
	return GitPullRequest{
		PullRequest: PullRequest{
			Number:    p.GetNumber(),
			HtmlUrl:   p.GetHTMLURL(),
			State:     p.GetState(),
			Title:     p.GetTitle(),
			Body:      p.GetBody(),
			User:      User{Login: p.GetUser().GetLogin()},
			CreatedAt: p.GetCreatedAt().String(),
			UpdatedAt: p.GetUpdatedAt().String(),
			ClosedAt:  p.GetClosedAt().String(),
			MergedAt:  p.GetMergedAt().String(),
		},
		Repository: Repository{
			Name:    r.Name,
			Owner:   r.Owner,
			HtmlUrl: r.HtmlUrl,
		},
	}

}

func (c *GitPullRequest) ToMap() map[string]string {
	return map[string]string{
		"action":                  c.Action,
		"pull_request_number":     strconv.Itoa(c.PullRequest.Number),
		"pull_request_title":      c.PullRequest.Title,
		"pull_request_body":       c.PullRequest.Body,
		"pull_request_user":       c.PullRequest.User.Login,
		"pull_request_url":        c.PullRequest.HtmlUrl,
		"pull_request_state":      c.PullRequest.State,
		"pull_request_created_at": c.PullRequest.CreatedAt,
		"pull_request_updated_at": c.PullRequest.UpdatedAt,
		"pull_request_closed_at":  c.PullRequest.ClosedAt,
		"pull_request_merged_at":  c.PullRequest.MergedAt,
		"repository_name":         c.Repository.Name,
		"repository_url":          c.Repository.HtmlUrl,
	}
}
