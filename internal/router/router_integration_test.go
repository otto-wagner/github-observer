//go:build all || integration

package router

import (
	"bytes"
	"encoding/json"
	"github-observer/conf"
	"github-observer/internal/core"
	"github-observer/internal/executor"
	eLogging "github-observer/internal/executor/Logging"
	ePrometheus "github-observer/internal/executor/Prometheus"
	l "github-observer/internal/listener"
	"github-observer/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v61/github"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"testing"
)

func TestRouterIntegration(t *testing.T) {
	engine := gin.New()
	buf := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(buf, nil))
	repositories := []core.Repository{{Owner: "otto-wagner", Name: "github-observer", Branch: "main"}}
	executors := []executor.IExecutor{eLogging.NewExecutor(eLogging.NewMemory(), logger), ePrometheus.NewExecutor()}
	listener := l.NewListener(repositories, executors, logger)

	InitializeRoutes(engine, listener, conf.Config{Secret: "your-secret"})

	t.Run("Should return ok", func(t *testing.T) {
		// given
		// when
		response := mocks.PerformRequest(engine, http.MethodGet, "/health", "")

		// then
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "{\"health\":\"ok\"}", response.Body.String())
	})

	t.Run("Should return 404", func(t *testing.T) {
		// given
		// when
		response := mocks.PerformRequest(engine, http.MethodGet, "/not-found", "")

		// then
		assert.Equal(t, http.StatusNotFound, response.Code)
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
		response := mocks.PerformRequestAuthorisation(engine, http.MethodPost, "/listen/workflow", string(event), "your-secret")

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
		response := mocks.PerformRequestAuthorisation(engine, http.MethodPost, "/listen/pullrequest", string(event), "your-secret")

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
		response := mocks.PerformRequestAuthorisation(engine, http.MethodPost, "/listen/pullrequest/review", string(event), "your-secret")

		// then
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "{\"message\":\"Pullrequest review received\"}", response.Body.String())
	})

	t.Run("Should return 400 when invalid json", func(t *testing.T) {
		// given
		// when
		response := mocks.PerformRequestAuthorisation(engine, http.MethodPost, "/listen/workflow", "{", "your-secret")

		// then
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Should return 400 when invalid request", func(t *testing.T) {
		// given
		// when
		response := mocks.PerformRequestAuthorisation(engine, http.MethodPost, "/listen/workflow", "", "your-secret")

		// then
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Should return 401 when invalid secret", func(t *testing.T) {
		// given
		// when
		response := mocks.PerformRequestAuthorisation(engine, http.MethodPost, "/listen/workflow", "{},", "invalid")

		// then
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	})

	t.Run("Should return 400 when missing X-Hub-Signature-256", func(t *testing.T) {
		// given
		// when
		response := mocks.PerformRequest(engine, http.MethodPost, "/listen/workflow", "{}")

		// then
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Should return 400 when body is empty", func(t *testing.T) {
		// given
		// when
		response := mocks.PerformRequestAuthorisation(engine, http.MethodPost, "/listen/workflow", "", "your-secret")

		// then
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

}
