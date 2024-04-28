//go:build all || unit

package listener

import (
	"encoding/json"
	"github-observer/internal/executor"
	internalMocks "github-observer/internal/mocks"
	"github-observer/mocks"
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
		mockedExecutor.On("RunEvent", checkRunEvent)

		event, _ := json.Marshal(checkRunEvent)
		context, recorder := mocks.MockContext("", string(event))

		// when
		NewListener([]executor.IExecutor{mockedExecutor}).Action(context)

		// then
		var expectedResponse gin.H
		err := json.Unmarshal(recorder.Body.Bytes(), &expectedResponse)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, gin.H{"message": "Workflow received"}, expectedResponse)
		mockedExecutor.AssertExpectations(t)
	})

}

func TestListenPullRequest(t *testing.T) {

	t.Run("Should listen pull request", func(t *testing.T) {
		// given
		checkRunEvent := github.PullRequestEvent{}
		mockedExecutor := new(internalMocks.IExecutor)
		mockedExecutor.On("PullRequestEvent", checkRunEvent)

		event, _ := json.Marshal(checkRunEvent)
		context, recorder := mocks.MockContext("", string(event))

		// when
		NewListener([]executor.IExecutor{mockedExecutor}).PullRequest(context)

		// then
		var expectedResponse gin.H
		err := json.Unmarshal(recorder.Body.Bytes(), &expectedResponse)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, gin.H{"message": "Workflow received"}, expectedResponse)
		mockedExecutor.AssertExpectations(t)
	})

}

func TestListenPullRequestReview(t *testing.T) {

	t.Run("Should listen pull request review", func(t *testing.T) {
		// given
		pullRequestReviewEvent := github.PullRequestReviewEvent{}
		mockedExecutor := new(internalMocks.IExecutor)
		mockedExecutor.On("PullRequestReviewEvent", pullRequestReviewEvent)

		event, _ := json.Marshal(pullRequestReviewEvent)
		context, recorder := mocks.MockContext("", string(event))

		// when
		NewListener([]executor.IExecutor{mockedExecutor}).PullRequestReview(context)

		// then
		var expectedResponse gin.H
		err := json.Unmarshal(recorder.Body.Bytes(), &expectedResponse)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, gin.H{"message": "Workflow received"}, expectedResponse)
		mockedExecutor.AssertExpectations(t)
	})

}
