package core

import "github.com/google/go-github/v61/github"

type GitPullRequestReview struct {
	Action      string      `json:"action"`
	PullRequest PullRequest `json:"pull_request"`
	Review      Review      `json:"review"`
	Repository  Repository  `json:"repository"`
}

type Review struct {
	User    User   `json:"user"`
	Body    string `json:"body"`
	State   string `json:"state"`
	HtmlUrl string `json:"html_url"`
}

func ConvertToGitPullRequestReview(event github.PullRequestReviewEvent) GitPullRequestReview {
	return GitPullRequestReview{
		Action: event.GetAction(),
		PullRequest: PullRequest{
			HtmlUrl: event.GetPullRequest().GetHTMLURL(),
			Title:   event.GetPullRequest().GetTitle(),
			User:    User{Login: event.GetPullRequest().GetUser().GetLogin()},
			State:   event.GetPullRequest().GetState(),
		},
		Review: Review{
			User:    User{Login: event.GetReview().GetUser().GetLogin()},
			Body:    event.GetReview().GetBody(),
			State:   event.GetReview().GetState(),
			HtmlUrl: event.GetReview().GetHTMLURL(),
		},
		Repository: Repository{
			Name:    event.GetRepo().GetName(),
			HtmlUrl: event.GetRepo().GetHTMLURL(),
		},
	}
}

func (c *GitPullRequestReview) ToMap() map[string]string {
	return map[string]string{
		"action":             c.Action,
		"pull_request_title": c.PullRequest.Title,
		"pull_request_user":  c.PullRequest.User.Login,
		"pull_request_url":   c.PullRequest.HtmlUrl,
		"pull_request_state": c.PullRequest.State,
		"review_body":        c.Review.Body,
		"review_state":       c.Review.State,
		"review_user":        c.Review.User.Login,
		"review_url":         c.Review.HtmlUrl,
		"repository_name":    c.Repository.Name,
		"repository_url":     c.Repository.HtmlUrl,
	}
}
