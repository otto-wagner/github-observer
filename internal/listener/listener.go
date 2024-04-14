package listener

import (
	"github-observer/internal/Executor"
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
	executors []Executor.IExecutor
}

func NewListener(executors []Executor.IExecutor) IListener {
	return &listener{executors}
}

func (l *listener) Action(c *gin.Context) {
	var event github.CheckRunEvent
	if err := c.BindJSON(&event); err != nil {
		zap.S().Errorw("Failed to bind CheckRunEvent", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
	}

	var errOccurred bool
	for _, e := range l.executors {
		err := e.CheckRunEvent(event)
		if err != nil {
			zap.S().Errorw("Failed to execute workflow", "executor", e.Name(), "error", err.Error())
			errOccurred = true
		}
	}
	if errOccurred {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to execute workflow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workflow received"})
}

func (l *listener) PullRequest(c *gin.Context) {
	var event github.PullRequestEvent
	if err := c.BindJSON(&event); err != nil {
		zap.S().Errorw("Failed to bind PullRequestEvent", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
	}

	var errOccurred bool
	for _, e := range l.executors {
		err := e.PullRequestEvent(event)
		if err != nil {
			zap.S().Errorw("Failed to execute workflow", "executor", e.Name(), "error", err.Error())
			errOccurred = true
		}
	}
	if errOccurred {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to execute workflow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workflow received"})
}

func (l *listener) PullRequestReview(c *gin.Context) {
	var event github.PullRequestReviewEvent
	if err := c.BindJSON(&event); err != nil {
		zap.S().Errorw("Failed to bind PullRequestReviewEvent", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
	}

	var errOccurred bool
	for _, e := range l.executors {
		err := e.PullRequestReviewEvent(event)
		if err != nil {
			zap.S().Errorw("Failed to execute workflow", "executor", e.Name(), "error", err.Error())
			errOccurred = true
		}
	}
	if errOccurred {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to execute workflow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workflow received"})
}
