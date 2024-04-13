package listener

import (
	"github-listener/internal/Executor"
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
	executor Executor.IExecutor
}

func NewListener(executor Executor.IExecutor) IListener {
	return &listener{executor}
}

func (l *listener) Action(c *gin.Context) {
	var runEvent github.CheckRunEvent
	if err := c.BindJSON(&runEvent); err != nil {
		zap.S().Errorw("Failed to bind CheckRunEvent", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
	}

	l.executor.CheckRunEvent(runEvent)
	c.JSON(http.StatusOK, gin.H{"message": "Workflow received"})
}

func (l *listener) PullRequest(c *gin.Context) {
	var event github.PullRequestEvent
	if err := c.BindJSON(&event); err != nil {
		zap.S().Errorw("Failed to bind PullRequestEvent", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
	}

	l.executor.PullRequestEvent(event)
	c.JSON(http.StatusOK, gin.H{"message": "Workflow received"})
}

func (l *listener) PullRequestReview(c *gin.Context) {
	var event github.PullRequestReviewEvent
	if err := c.BindJSON(&event); err != nil {
		zap.S().Errorw("Failed to bind PullRequestReviewEvent", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
	}

	l.executor.PullRequestReviewEvent(event)
	c.JSON(http.StatusOK, gin.H{"message": "Workflow received"})
}
