package watcher

import (
	"github-observer/internal/Executor"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v61/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

type IWatcher interface {
	Actions(*gin.Context)
	PullRequests(*gin.Context)
}

type watcher struct {
	executors []Executor.IExecutor
	token     string
}

func NewWatcher(executors []Executor.IExecutor) IWatcher {
	// todo: wo anders einlesen
	token := os.Getenv("GITHUB_TOKEN")
	return &watcher{executors, token}
}

func (w *watcher) Actions(c *gin.Context) {
	client := github.NewClient(oauth2.NewClient(c, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: w.token})))
	runs, r, err := client.Actions.ListRepositoryWorkflowRuns(c, "otto-wagner", "github-observer", nil)
	if err != nil {
		zap.S().Errorw("Failed to list workflow runs", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to list workflow runs"})
		return
	}
	if r.StatusCode > 299 {
		zap.S().Errorw("Failed to list workflow runs", "status_code", r.StatusCode)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to list workflow runs"})
		return
	}

	c.JSON(http.StatusOK, runs.WorkflowRuns)
}

func (w *watcher) PullRequests(c *gin.Context) {
	client := github.NewClient(oauth2.NewClient(c, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: w.token})))
	pullRequests, r, err := client.PullRequests.List(c, "otto-wagner", "github-observer", nil)
	if err != nil {
		zap.S().Errorw("Failed to list pull requests", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to list pull requests"})
		return
	}
	if r.StatusCode > 299 {
		zap.S().Errorw("Failed to list pull requests", "status_code", r.StatusCode)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to list pull requests"})
		return
	}

	c.JSON(http.StatusOK, pullRequests)
}
