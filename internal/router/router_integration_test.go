//go:build all || integration

package router

import (
	"encoding/json"
	"github-observer/internal/core"
	"github-observer/internal/executor"
	eLogging "github-observer/internal/executor/Logging"
	ePrometheus "github-observer/internal/executor/Prometheus"
	l "github-observer/internal/listener"
	"github-observer/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v61/github"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRouterIntegration(t *testing.T) {
	engine := gin.New()
	repositories := []core.Repository{{Owner: "otto-wagner", Name: "github-observer", Branch: "main"}}
	executors := []executor.IExecutor{eLogging.NewExecutor(eLogging.NewMemory()), ePrometheus.NewExecutor()}
	listener := l.NewListener(repositories, executors)

	InitializeRoutes(engine, listener)

	t.Run("Should return ok", func(t *testing.T) {
		// given
		// when
		response := mocks.PerformRequest(engine, http.MethodGet, "/health", "")

		// then
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "{\"health\":\"ok\"}", response.Body.String())
	})

	t.Run("Should listen workflows", func(t *testing.T) {
		// given
		event, _ := json.Marshal(github.WorkflowRunEvent{
			WorkflowRun: &github.WorkflowRun{
				HeadBranch: github.String("main"),
			},
			Repo: &github.Repository{
				Name: github.String("github-observer"),
			},
		})

		// when
		response := mocks.PerformRequest(engine, http.MethodPost, "/listen/workflow", string(event))

		// then
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "{\"message\":\"Workflow received\"}", response.Body.String())
	})

	t.Run("Should listen pull requests", func(t *testing.T) {
		// given
		event, _ := json.Marshal(github.PullRequestEvent{
			PullRequest: &github.PullRequest{
				Title: github.String("chore: test pullrequest_listener"),
			},
			Repo: &github.Repository{
				Name: github.String("github-observer"),
			},
		})

		// when
		response := mocks.PerformRequest(engine, http.MethodPost, "/listen/pullrequest", string(event))

		// then
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "{\"message\":\"Pullrequest received\"}", response.Body.String())
	})

	t.Run("Should listen pull request review", func(t *testing.T) {
		// given
		event, _ := json.Marshal(github.PullRequestReviewEvent{
			PullRequest: &github.PullRequest{
				Title: github.String("chore: test pullrequest_listener"),
			},
			Repo: &github.Repository{
				Name: github.String("github-observer"),
			},
			Review: &github.PullRequestReview{
				Body: github.String("LGTM"),
			},
		})

		// when
		response := mocks.PerformRequest(engine, http.MethodPost, "/listen/pullrequest/review", string(event))

		// then
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "{\"message\":\"Pullrequest review received\"}", response.Body.String())
	})
}
