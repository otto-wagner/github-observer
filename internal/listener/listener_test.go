//go:build all || unit

package listener

import (
	"encoding/json"
	"errors"
	"github-listener/internal/Executor"
	internalMocks "github-listener/internal/mocks"
	"github-listener/mocks"
	logger "github-listener/pkg/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v61/github"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestListenAction(t *testing.T) {

	t.Run("Should listen action", func(t *testing.T) {
		// given
		checkRunEvent := github.CheckRunEvent{}
		mockedExecutor := new(internalMocks.IExecutor)
		mockedExecutor.On("CheckRunEvent", checkRunEvent).Return(nil)

		event, _ := json.Marshal(checkRunEvent)
		context, recorder := mocks.MockContext("", string(event))

		// when
		NewListener([]Executor.IExecutor{mockedExecutor}).Action(context)

		// then
		var expectedResponse gin.H
		err := json.Unmarshal(recorder.Body.Bytes(), &expectedResponse)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, gin.H{"message": "Workflow received"}, expectedResponse)
		mockedExecutor.AssertExpectations(t)
	})

	t.Run("Should listen action return error", func(t *testing.T) {
		// given
		logs := logger.MockedLogger()

		checkRunEvent := github.CheckRunEvent{}
		mockedExecutor := new(internalMocks.IExecutor)
		mockedExecutor.On("CheckRunEvent", checkRunEvent).Return(errors.New("logging error occurred"))
		mockedExecutor.On("Name").Return("Logging")
		mockedExecutor2 := new(internalMocks.IExecutor)
		mockedExecutor2.On("CheckRunEvent", checkRunEvent).Return(errors.New("prometheus error occurred"))
		mockedExecutor2.On("Name").Return("Prometheus")

		event, _ := json.Marshal(checkRunEvent)
		context, recorder := mocks.MockContext("", string(event))

		// when
		NewListener([]Executor.IExecutor{mockedExecutor, mockedExecutor2}).Action(context)

		// then
		var expectedResponse gin.H
		err := json.Unmarshal(recorder.Body.Bytes(), &expectedResponse)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, gin.H{"message": "failed to execute workflow"}, expectedResponse)
		assert.Contains(t, logs.All()[0].Message, "Failed to execute workflow")
		assert.Contains(t, logs.All()[0].Context[0].Key, "executor")
		assert.Contains(t, logs.All()[0].Context[0].String, "Logging")
		assert.Contains(t, logs.All()[0].Context[1].Key, "error")
		assert.Contains(t, logs.All()[0].Context[1].String, "logging error occurred")
		assert.Contains(t, logs.All()[1].Message, "Failed to execute workflow")
		assert.Contains(t, logs.All()[1].Context[0].Key, "executor")
		assert.Contains(t, logs.All()[1].Context[0].String, "Prometheus")
		assert.Contains(t, logs.All()[1].Context[1].Key, "error")
		assert.Contains(t, logs.All()[1].Context[1].String, "prometheus error occurred")
		mockedExecutor.AssertExpectations(t)
		mockedExecutor2.AssertExpectations(t)
	})

}

func TestListenPullRequest(t *testing.T) {

	t.Run("Should listen pull request", func(t *testing.T) {
		// given
		checkRunEvent := github.PullRequestEvent{}
		mockedExecutor := new(internalMocks.IExecutor)
		mockedExecutor.On("PullRequestEvent", checkRunEvent).Return(nil)

		event, _ := json.Marshal(checkRunEvent)
		context, recorder := mocks.MockContext("", string(event))

		// when
		NewListener([]Executor.IExecutor{mockedExecutor}).PullRequest(context)

		// then
		var expectedResponse gin.H
		err := json.Unmarshal(recorder.Body.Bytes(), &expectedResponse)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, gin.H{"message": "Workflow received"}, expectedResponse)
		mockedExecutor.AssertExpectations(t)
	})

	t.Run("Should listen pull request return error", func(t *testing.T) {
		// given
		logs := logger.MockedLogger()

		checkRunEvent := github.PullRequestEvent{}
		mockedExecutor := new(internalMocks.IExecutor)
		mockedExecutor.On("PullRequestEvent", checkRunEvent).Return(errors.New("logging error occurred"))
		mockedExecutor.On("Name").Return("Logging")
		mockedExecutor2 := new(internalMocks.IExecutor)
		mockedExecutor2.On("PullRequestEvent", checkRunEvent).Return(errors.New("prometheus error occurred"))
		mockedExecutor2.On("Name").Return("Prometheus")

		event, _ := json.Marshal(checkRunEvent)
		context, recorder := mocks.MockContext("", string(event))

		// when
		NewListener([]Executor.IExecutor{mockedExecutor, mockedExecutor2}).PullRequest(context)

		// then
		var expectedResponse gin.H
		err := json.Unmarshal(recorder.Body.Bytes(), &expectedResponse)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, gin.H{"message": "failed to execute workflow"}, expectedResponse)
		assert.Contains(t, logs.All()[0].Message, "Failed to execute workflow")
		assert.Contains(t, logs.All()[0].Context[0].Key, "executor")
		assert.Contains(t, logs.All()[0].Context[0].String, "Logging")
		assert.Contains(t, logs.All()[0].Context[1].Key, "error")
		assert.Contains(t, logs.All()[0].Context[1].String, "logging error occurred")
		assert.Contains(t, logs.All()[1].Message, "Failed to execute workflow")
		assert.Contains(t, logs.All()[1].Context[0].Key, "executor")
		assert.Contains(t, logs.All()[1].Context[0].String, "Prometheus")
		assert.Contains(t, logs.All()[1].Context[1].Key, "error")
		assert.Contains(t, logs.All()[1].Context[1].String, "prometheus error occurred")
		mockedExecutor.AssertExpectations(t)
		mockedExecutor2.AssertExpectations(t)
	})

}

func TestListenPullRequestReview(t *testing.T) {

	t.Run("Should listen pull request review", func(t *testing.T) {
		// given
		pullRequestReviewEvent := github.PullRequestReviewEvent{}
		mockedExecutor := new(internalMocks.IExecutor)
		mockedExecutor.On("PullRequestReviewEvent", pullRequestReviewEvent).Return(nil)

		event, _ := json.Marshal(pullRequestReviewEvent)
		context, recorder := mocks.MockContext("", string(event))

		// when
		NewListener([]Executor.IExecutor{mockedExecutor}).PullRequestReview(context)

		// then
		var expectedResponse gin.H
		err := json.Unmarshal(recorder.Body.Bytes(), &expectedResponse)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, gin.H{"message": "Workflow received"}, expectedResponse)
		mockedExecutor.AssertExpectations(t)
	})

	t.Run("Should listen pull request review return error", func(t *testing.T) {
		// given
		logs := logger.MockedLogger()

		checkRunEvent := github.PullRequestReviewEvent{}
		mockedExecutor := new(internalMocks.IExecutor)
		mockedExecutor.On("PullRequestReviewEvent", checkRunEvent).Return(errors.New("logging error occurred"))
		mockedExecutor.On("Name").Return("Logging")
		mockedExecutor2 := new(internalMocks.IExecutor)
		mockedExecutor2.On("PullRequestReviewEvent", checkRunEvent).Return(errors.New("prometheus error occurred"))
		mockedExecutor2.On("Name").Return("Prometheus")

		event, _ := json.Marshal(checkRunEvent)
		context, recorder := mocks.MockContext("", string(event))

		// when
		NewListener([]Executor.IExecutor{mockedExecutor, mockedExecutor2}).PullRequestReview(context)

		// then
		var expectedResponse gin.H
		err := json.Unmarshal(recorder.Body.Bytes(), &expectedResponse)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, gin.H{"message": "failed to execute workflow"}, expectedResponse)
		assert.Contains(t, logs.All()[0].Message, "Failed to execute workflow")
		assert.Contains(t, logs.All()[0].Context[0].Key, "executor")
		assert.Contains(t, logs.All()[0].Context[0].String, "Logging")
		assert.Contains(t, logs.All()[0].Context[1].Key, "error")
		assert.Contains(t, logs.All()[0].Context[1].String, "logging error occurred")
		assert.Contains(t, logs.All()[1].Message, "Failed to execute workflow")
		assert.Contains(t, logs.All()[1].Context[0].Key, "executor")
		assert.Contains(t, logs.All()[1].Context[0].String, "Prometheus")
		assert.Contains(t, logs.All()[1].Context[1].Key, "error")
		assert.Contains(t, logs.All()[1].Context[1].String, "prometheus error occurred")
		mockedExecutor.AssertExpectations(t)
		mockedExecutor2.AssertExpectations(t)
	})

}
