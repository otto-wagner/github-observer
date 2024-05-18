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
			FullName: event.GetRepo().GetFullName(),
			HtmlUrl:  event.GetRepo().GetHTMLURL(),
		},
	}
}
