package Executor

import "github.com/google/go-github/v61/github"

type IExecutor interface {
	CheckRunEvent(github.CheckRunEvent)
	PullRequestEvent(github.PullRequestEvent)
	PullRequestReviewEvent(github.PullRequestReviewEvent)
}
