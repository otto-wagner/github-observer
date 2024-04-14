package Executor

import "github.com/google/go-github/v61/github"

type IExecutor interface {
	Name() string
	CheckRunEvent(github.CheckRunEvent) error
	PullRequestEvent(github.PullRequestEvent) error
	PullRequestReviewEvent(github.PullRequestReviewEvent) error
}
