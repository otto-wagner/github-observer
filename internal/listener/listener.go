package listener

import (
	"github-observer/internal/core"
	"github-observer/internal/executor"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v61/github"
	"log/slog"
	"net/http"
)

type IListener interface {
	Workflow(*gin.Context)
	PullRequest(*gin.Context)
	PullRequestReview(*gin.Context)
}

type listener struct {
	repositories []core.Repository
	executors    []executor.IExecutor
	logger       *slog.Logger
}

func NewListener(repositories []core.Repository, executors []executor.IExecutor, logger *slog.Logger) IListener {
	return &listener{repositories, executors, logger}
}

func (l *listener) Workflow(c *gin.Context) {
	var event github.WorkflowRunEvent
	if err := c.BindJSON(&event); err != nil {
		l.logger.Error("Failed to bind EventWorkflow", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
		return
	}

	var repository core.Repository
	for _, r := range l.repositories {
		if event.GetRepo().GetName() == r.Name {
			repository = r
		}
	}
	l.logger.Info("Workflow received",
		"repository", repository.Name,
		"workflow", event.GetWorkflowRun().Name,
		"url", event.GetWorkflowRun().GetHTMLURL(),
		"status", event.GetWorkflowRun().GetStatus(),
		"conclusion", event.GetWorkflowRun().GetConclusion(),
	)

	if event.GetWorkflowRun().GetHeadBranch() == repository.Branch {
		for _, e := range l.executors {
			e.EventWorkflowRun(event)
		}
		c.JSON(http.StatusOK, gin.H{"message": "Workflow received"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Branch ignored"})
}

func (l *listener) PullRequest(c *gin.Context) {
	var event github.PullRequestEvent
	if err := c.BindJSON(&event); err != nil {
		l.logger.Error("Failed to bind EventPullRequest", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
		return
	}
	l.logger.Info("Pullrequest received", "repository", event.GetRepo().GetName(), "pullrequest", event.GetPullRequest().GetTitle(), "url", event.GetPullRequest().GetHTMLURL())

	for _, e := range l.executors {
		e.EventPullRequest(event)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pullrequest received"})
}

func (l *listener) PullRequestReview(c *gin.Context) {
	var event github.PullRequestReviewEvent
	if err := c.BindJSON(&event); err != nil {
		l.logger.Error("Failed to bind EventPullRequestReview", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
		return
	}
	l.logger.Info("Pullrequest review received", "repository", event.GetRepo().GetName(), "pullrequest", event.GetPullRequest().GetTitle(), "url", event.GetPullRequest().GetHTMLURL())

	for _, e := range l.executors {
		e.EventPullRequestReview(event)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pullrequest review received"})
}
