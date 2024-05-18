package executor

import (
	"github-observer/internal/core"
	"github.com/google/go-github/v61/github"
)

type IExecutor interface {
	EventRun(github.CheckRunEvent)
	EventPullRequest(github.PullRequestEvent)
	EventPullRequestReview(github.PullRequestReviewEvent)
	LastWorkflows(core.Repository, []*github.WorkflowRun)
	PullRequests(core.Repository, []*github.PullRequest)
}
