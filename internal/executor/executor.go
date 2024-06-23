package executor

import (
	"github-observer/internal/core"
	"github.com/google/go-github/v61/github"
)

type IExecutor interface {
	EventPullRequest(github.PullRequestEvent)
	EventPullRequestReview(github.PullRequestReviewEvent)
	EventWorkflowRun(github.WorkflowRunEvent)
	LastWorkflows(core.Repository, []*github.WorkflowRun)
	PullRequests(core.Repository, []*github.PullRequest)
}
