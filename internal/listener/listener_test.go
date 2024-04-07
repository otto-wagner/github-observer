//go:build all || unit

package listener

import (
	"encoding/json"
	"github-listener/mocks"
	logger "github-listener/pkg/mocks"
	"github.com/google/go-github/v61/github"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestListen(t *testing.T) {
	event, _ := json.Marshal(github.CheckRunEvent{
		CheckRun: &github.CheckRun{
			Name:       github.String("Analyze (go)"),
			HTMLURL:    github.String("https://github.com/otto-wagner/github-listener/actions/runs/8589035842/job/23534635896"),
			Status:     github.String("completed"),
			Conclusion: github.String("success"),
		},
		Repo: &github.Repository{
			Name:    github.String("github-listener"),
			HTMLURL: github.String("https://github.com/otto-wagner/github-listener"),
		},
	})

	t.Run("Should start listening", func(t *testing.T) {
		// given
		logs := logger.MockedLogger()
		context, recorder := mocks.MockContext("", string(event))

		// when
		NewListener().Listen(context)

		// then
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "{\"message\":\"Workflow received\"}", recorder.Body.String())
		assert.Contains(t, logs.All()[0].Message, "Workflow received")
		assert.Contains(t, logs.All()[0].Context[0].Key, "name")
		assert.Contains(t, logs.All()[0].Context[0].String, "Analyze (go)")
		assert.Contains(t, logs.All()[0].Context[1].Key, "html_url")
		assert.Contains(t, logs.All()[0].Context[1].String, "https://github.com/otto-wagner/github-listener/actions/runs/8589035842/job/23534635896")
		assert.Contains(t, logs.All()[0].Context[2].Key, "status")
		assert.Contains(t, logs.All()[0].Context[2].String, "completed")
		assert.Contains(t, logs.All()[0].Context[3].Key, "conclusion")
		assert.Contains(t, logs.All()[0].Context[3].String, "success")
		assert.Contains(t, logs.All()[0].Context[4].Key, "repo")
		assert.Contains(t, logs.All()[0].Context[4].String, "github-listener")
		assert.Contains(t, logs.All()[0].Context[5].Key, "repo_html_url")
		assert.Contains(t, logs.All()[0].Context[5].String, "https://github.com/otto-wagner/github-listener")
	})
}
