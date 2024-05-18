package listener

import (
	"github-observer/internal/executor"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v61/github"
	"go.uber.org/zap"
	"net/http"
)

type IListener interface {
	Action(*gin.Context)
	PullRequest(*gin.Context)
	PullRequestReview(*gin.Context)
}

type listener struct {
	executors []executor.IExecutor
}

func NewListener(executors []executor.IExecutor) IListener {
	return &listener{executors}
}

func (l *listener) Action(c *gin.Context) {
	var event github.CheckRunEvent
	if err := c.BindJSON(&event); err != nil {
		zap.S().Errorw("Failed to bind EventRun", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
	}

	for _, e := range l.executors {
		e.EventRun(event)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Workflow received"})
}

func (l *listener) PullRequest(c *gin.Context) {
	var event github.PullRequestEvent
	if err := c.BindJSON(&event); err != nil {
		zap.S().Errorw("Failed to bind EventPullRequest", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
	}

	for _, e := range l.executors {
		e.EventPullRequest(event)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Workflow received"})
}

func (l *listener) PullRequestReview(c *gin.Context) {
	var event github.PullRequestReviewEvent
	if err := c.BindJSON(&event); err != nil {
		zap.S().Errorw("Failed to bind EventPullRequestReview", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
	}

	for _, e := range l.executors {
		e.EventPullRequestReview(event)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Workflow received"})
}
