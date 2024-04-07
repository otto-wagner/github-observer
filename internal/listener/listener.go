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
	PullRequestReview(*gin.Context)
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
		"repo", runEvent.GetRepo().GetName(),
		"repo_html_url", runEvent.GetRepo().GetHTMLURL(),
		"name", runEvent.GetCheckRun().GetName(),
		"html_url", runEvent.GetCheckRun().GetHTMLURL(),
		"action", runEvent.GetAction(),
		"status", runEvent.GetCheckRun().GetStatus(),
		"conclusion", runEvent.GetCheckRun().GetConclusion(),
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
		"repo", prEvent.GetRepo().GetName(),
		"repo_html_url", prEvent.GetRepo().GetHTMLURL(),
		"title", prEvent.GetPullRequest().GetTitle(),
		"user", prEvent.GetPullRequest().GetUser().GetLogin(),
		"html_url", prEvent.GetPullRequest().GetHTMLURL(),
		"action", prEvent.GetAction(),
		"status", prEvent.GetPullRequest().GetState(),
	)
	c.JSON(http.StatusOK, gin.H{"message": "Workflow received"})
}

func (l *listener) PullRequestReview(c *gin.Context) {
	var prEvent github.PullRequestReviewEvent
	if err := c.BindJSON(&prEvent); err != nil {
		zap.S().Errorw("Failed to bind PullRequestReviewEvent", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
	}

	zap.S().Infow("Workflow received",
		"repo", prEvent.GetRepo().GetName(),
		"repo_html_url", prEvent.GetRepo().GetHTMLURL(),
		"title", prEvent.GetPullRequest().GetTitle(),
		"user", prEvent.GetPullRequest().GetUser().GetLogin(),
		"html_url", prEvent.GetPullRequest().GetHTMLURL(),
		"action", prEvent.GetAction(),
		"status", prEvent.GetPullRequest().GetState(),
		"review", prEvent.GetReview().GetBody(),
		"state", prEvent.GetReview().GetState(),
		"reviewer", prEvent.GetReview().GetUser().GetLogin(),
	)
	c.JSON(http.StatusOK, gin.H{"message": "Workflow received"})
}
