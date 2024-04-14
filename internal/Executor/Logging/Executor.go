package Logging

import (
	"github-observer/internal/Executor"
	"github.com/google/go-github/v61/github"
	"go.uber.org/zap"
)

type executor struct {
}

func NewExecutor() Executor.IExecutor {
	return &executor{}
}

func (e executor) Name() string {
	return "Logging"
}

func (e executor) CheckRunEvent(runEvent github.CheckRunEvent) (err error) {
	zap.S().Infow("Workflow received",
		"repo", runEvent.GetRepo().GetName(),
		"repo_html_url", runEvent.GetRepo().GetHTMLURL(),
		"name", runEvent.GetCheckRun().GetName(),
		"html_url", runEvent.GetCheckRun().GetHTMLURL(),
		"action", runEvent.GetAction(),
		"status", runEvent.GetCheckRun().GetStatus(),
		"conclusion", runEvent.GetCheckRun().GetConclusion(),
	)
	return
}

func (e executor) PullRequestEvent(event github.PullRequestEvent) (err error) {
	zap.S().Infow("Workflow received",
		"repo", event.GetRepo().GetName(),
		"repo_html_url", event.GetRepo().GetHTMLURL(),
		"title", event.GetPullRequest().GetTitle(),
		"user", event.GetPullRequest().GetUser().GetLogin(),
		"html_url", event.GetPullRequest().GetHTMLURL(),
		"action", event.GetAction(),
		"status", event.GetPullRequest().GetState(),
	)
	return
}

func (e executor) PullRequestReviewEvent(event github.PullRequestReviewEvent) (err error) {
	zap.S().Infow("Workflow received",
		"repo", event.GetRepo().GetName(),
		"repo_html_url", event.GetRepo().GetHTMLURL(),
		"title", event.GetPullRequest().GetTitle(),
		"user", event.GetPullRequest().GetUser().GetLogin(),
		"html_url", event.GetPullRequest().GetHTMLURL(),
		"action", event.GetAction(),
		"status", event.GetPullRequest().GetState(),
		"review", event.GetReview().GetBody(),
		"state", event.GetReview().GetState(),
		"reviewer", event.GetReview().GetUser().GetLogin(),
	)
	return
}
