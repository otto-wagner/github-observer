//go:build all || integration

package router

import (
	"encoding/json"
	"github-observer/internal/Executor"
	"github-observer/internal/Executor/Logging"
	"github-observer/internal/listener"
	"github-observer/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v61/github"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRouterIntegration(t *testing.T) {
	executor := Logging.NewExecutor()
	l := listener.NewListener([]Executor.IExecutor{executor})
	engine := gin.New()
	InitializeRoutes(engine, l)

	t.Run("Should return ok", func(t *testing.T) {
		// given
		// when
		response := mocks.PerformRequest(engine, http.MethodGet, "/health", "")

		// then
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "{\"health\":\"ok\"}", response.Body.String())
	})

	t.Run("Should listen action", func(t *testing.T) {
		// given
		event, _ := json.Marshal(github.CheckRunEvent{
			CheckRun: &github.CheckRun{
				Name: github.String("Analyze (go)"),
			},
			Repo: &github.Repository{
				Name: github.String("github-observer"),
			},
		})

		// when
		response := mocks.PerformRequest(engine, http.MethodPost, "/listen/action", string(event))

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
		assert.Equal(t, "{\"message\":\"Workflow received\"}", response.Body.String())
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
		assert.Equal(t, "{\"message\":\"Workflow received\"}", response.Body.String())
	})
}
