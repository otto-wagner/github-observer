package core

import (
	"github.com/google/go-github/v61/github"
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
			User:      User{Login: event.GetPullRequest().GetUser().GetLogin()},
			CreatedAt: event.GetPullRequest().GetCreatedAt().String(),
			UpdatedAt: event.GetPullRequest().GetUpdatedAt().String(),
			ClosedAt:  event.GetPullRequest().GetClosedAt().String(),
			MergedAt:  event.GetPullRequest().GetMergedAt().String(),
		},
		Repository: Repository{
			FullName: event.GetRepo().GetFullName(),
			HtmlUrl:  event.GetRepo().GetHTMLURL(),
		},
	}
}

func ConvertPRToGitPullRequest(p github.PullRequest) GitPullRequest {
	return GitPullRequest{
		PullRequest: PullRequest{
			Number:    p.GetNumber(),
			HtmlUrl:   p.GetHTMLURL(),
			State:     p.GetState(),
			Title:     p.GetTitle(),
			User:      User{Login: p.GetUser().GetLogin()},
			CreatedAt: p.GetCreatedAt().String(),
			UpdatedAt: p.GetUpdatedAt().String(),
			ClosedAt:  p.GetClosedAt().String(),
			MergedAt:  p.GetMergedAt().String(),
		},
		Repository: Repository{
			FullName: p.GetBase().GetRepo().GetFullName(),
			HtmlUrl:  p.GetBase().GetRepo().GetHTMLURL(),
		},
	}
}
