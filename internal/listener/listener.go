package listener

import (
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v61/github"
	"go.uber.org/zap"
	"net/http"
)

type IListener interface {
	Listen(*gin.Context)
}

type listener struct {
}

func NewListener() IListener {
	return &listener{}
}

func (l *listener) Listen(c *gin.Context) {
	var runEvent github.CheckRunEvent

	if err := c.BindJSON(&runEvent); err != nil {
		zap.S().Errorw("Failed to bind CheckRunEvent", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind check run JSON"})
	}

	zap.S().Infow("Workflow received",
		"name", runEvent.CheckRun.Name,
		"html_url", runEvent.CheckRun.HTMLURL,
		"status", runEvent.CheckRun.Status,
		"conclusion", runEvent.CheckRun.Conclusion,
		"repo", runEvent.Repo.Name,
		"repo_html_url", runEvent.Repo.HTMLURL)
	c.JSON(http.StatusOK, gin.H{"message": "Workflow received"})
}
