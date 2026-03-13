package executor

import (
	"github.com/otto-wagner/github-observer/internal/core"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v73/github"
)

type IExecutor interface {
	Handler() gin.HandlerFunc
	EventPullRequest(github.PullRequestEvent)
	EventPullRequestReview(github.PullRequestReviewEvent)
	EventWorkflowRun(github.WorkflowRunEvent)
	LastWorkflows(core.Repository, []*github.WorkflowRun)
	PullRequests(core.Repository, []*github.PullRequest)
}
