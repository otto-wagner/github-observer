package listener

import (
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v61/github"
	"go.uber.org/zap"
	"net/http"
)

type IListener interface {
	Action(*gin.Context)
	PullRequest(*gin.Context)
}

type listener struct {
}

func NewListener() IListener {
	return &listener{}
}

func (l *listener) Action(c *gin.Context) {
	var runEvent github.CheckRunEvent
	if err := c.BindJSON(&runEvent); err != nil {
		zap.S().Errorw("Failed to bind CheckRunEvent", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
	}

	zap.S().Infow("Workflow received",
		"repo", runEvent.Repo.GetName(),
		"repo_html_url", runEvent.Repo.GetHTMLURL(),
		"name", runEvent.CheckRun.GetName(),
		"html_url", runEvent.CheckRun.GetHTMLURL(),
		"action", runEvent.GetAction(),
		"status", runEvent.CheckRun.GetStatus(),
		"conclusion", runEvent.CheckRun.GetConclusion(),
	)
	c.JSON(http.StatusOK, gin.H{"message": "Workflow received"})
}

func (l *listener) PullRequest(c *gin.Context) {
	var prEvent github.PullRequestEvent
	if err := c.BindJSON(&prEvent); err != nil {
		zap.S().Errorw("Failed to bind PullRequestEvent", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
	}

	zap.S().Infow("Workflow received",
		"repo", prEvent.Repo.GetName(),
		"repo_html_url", prEvent.Repo.GetHTMLURL(),
		"title", prEvent.PullRequest.GetTitle(),
		"user", prEvent.PullRequest.GetUser().GetLogin(),
		"html_url", prEvent.PullRequest.GetHTMLURL(),
		"action", prEvent.GetAction(),
		"status", prEvent.PullRequest.GetState(),
	)
	c.JSON(http.StatusOK, gin.H{"message": "Workflow received"})
}
