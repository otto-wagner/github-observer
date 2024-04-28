package executor

import (
	"github-observer/internal/core"
	"github.com/google/go-github/v61/github"
)

type IExecutor interface {
	Name() string
	RunEvent(github.CheckRunEvent)
	PullRequestEvent(github.PullRequestEvent)
	PullRequestReviewEvent(github.PullRequestReviewEvent)
	WorkflowRuns([]*github.WorkflowRun)
	PullRequests(core.Repository, []*github.PullRequest)
}
