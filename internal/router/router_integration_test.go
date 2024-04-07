//go:build all || integration

package router

import (
	"encoding/json"
	"github-listener/internal/listener"
	"github-listener/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v61/github"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRouterIntegration(t *testing.T) {
	l := listener.NewListener()
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

	t.Run("Should listen", func(t *testing.T) {
		// given
		event, _ := json.Marshal(github.CheckRunEvent{
			CheckRun: &github.CheckRun{
				Name: github.String("Analyze (go)"),
			},
			Repo: &github.Repository{
				Name: github.String("github-listener"),
			},
		})

		// when
		response := mocks.PerformRequest(engine, http.MethodPost, "/listener", string(event))

		// then
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "{\"message\":\"Workflow received\"}", response.Body.String())
	})
}
