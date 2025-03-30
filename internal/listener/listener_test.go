//go:build all || unit

package listener

import (
	"bytes"
	"encoding/json"
	"github-observer/internal/core"
	"github-observer/internal/executor"
	internalMocks "github-observer/internal/executor/mocks"
	"github-observer/server/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v61/github"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"testing"
)

func TestListenWorkflow(t *testing.T) {
	repository := core.Repository{Owner: "otto-wagner", Name: "github-observer", Branch: "main"}

	t.Run("Should listen workflow", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		logger := slog.New(slog.NewTextHandler(buf, nil))

		workflowRunEvent := github.WorkflowRunEvent{
			WorkflowRun: &github.WorkflowRun{
				HeadBranch: github.String("main"),
			},
			Repo: &github.Repository{
				Name: github.String("github-observer"),
			},
		}

		mockedExecutor := new(internalMocks.IExecutor)
		mockedExecutor.On("EventWorkflowRun", workflowRunEvent)

		event, _ := json.Marshal(workflowRunEvent)
		context, recorder := mocks.MockContext("", string(event))

		// when
		NewListener([]core.Repository{repository}, []executor.IExecutor{mockedExecutor}, logger).Workflow(context)

		// then
		var expectedResponse gin.H
		err := json.Unmarshal(recorder.Body.Bytes(), &expectedResponse)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, gin.H{"message": "Workflow received"}, expectedResponse)
		mockedExecutor.AssertExpectations(t)
	})

	t.Run("Should not listen workflows in other branch", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		logger := slog.New(slog.NewTextHandler(buf, nil))

		workflowRunEvent := github.WorkflowRunEvent{
			WorkflowRun: &github.WorkflowRun{
				HeadBranch: github.String("another"),
			},
			Repo: &github.Repository{
				Name: github.String("github-observer"),
			},
		}

		event, _ := json.Marshal(workflowRunEvent)
		context, recorder := mocks.MockContext("", string(event))

		// when
		NewListener([]core.Repository{repository}, []executor.IExecutor{nil}, logger).Workflow(context)

		// then
		var expectedResponse gin.H
		err := json.Unmarshal(recorder.Body.Bytes(), &expectedResponse)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusAccepted, recorder.Code)
		assert.Equal(t, gin.H{"message": "Branch ignored"}, expectedResponse)
	})

}

func TestListenPullRequest(t *testing.T) {
	repository := core.Repository{Owner: "otto-wagner", Name: "github-observer", Branch: "main"}

	t.Run("Should listen pull request", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		logger := slog.New(slog.NewTextHandler(buf, nil))

		checkRunEvent := github.PullRequestEvent{}
		mockedExecutor := new(internalMocks.IExecutor)
		mockedExecutor.On("EventPullRequest", checkRunEvent)

		event, _ := json.Marshal(checkRunEvent)
		context, recorder := mocks.MockContext("", string(event))

		// when
		NewListener([]core.Repository{repository}, []executor.IExecutor{mockedExecutor}, logger).PullRequest(context)

		// then
		var expectedResponse gin.H
		err := json.Unmarshal(recorder.Body.Bytes(), &expectedResponse)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, gin.H{"message": "Pullrequest received"}, expectedResponse)
		mockedExecutor.AssertExpectations(t)
	})

}

func TestListenPullRequestReview(t *testing.T) {
	repository := core.Repository{Owner: "otto-wagner", Name: "github-observer", Branch: "main"}

	t.Run("Should listen pull request review", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		logger := slog.New(slog.NewTextHandler(buf, nil))

		pullRequestReviewEvent := github.PullRequestReviewEvent{}
		mockedExecutor := new(internalMocks.IExecutor)
		mockedExecutor.On("EventPullRequestReview", pullRequestReviewEvent)

		event, _ := json.Marshal(pullRequestReviewEvent)
		context, recorder := mocks.MockContext("", string(event))

		// when
		NewListener([]core.Repository{repository}, []executor.IExecutor{mockedExecutor}, logger).PullRequestReview(context)

		// then
		var expectedResponse gin.H
		err := json.Unmarshal(recorder.Body.Bytes(), &expectedResponse)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, gin.H{"message": "Pullrequest review received"}, expectedResponse)
		mockedExecutor.AssertExpectations(t)
	})

}
