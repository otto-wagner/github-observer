//go:build all || unit

package listener

import (
	"encoding/json"
	internalMocks "github-listener/internal/mocks"
	"github-listener/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v61/github"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestListen(t *testing.T) {

	t.Run("Should listen action", func(t *testing.T) {
		// given
		checkRunEvent := github.CheckRunEvent{}
		mockedExecutor := new(internalMocks.IExecutor)
		mockedExecutor.On("CheckRunEvent", checkRunEvent)

		event, _ := json.Marshal(checkRunEvent)
		context, recorder := mocks.MockContext("", string(event))

		// when
		NewListener(mockedExecutor).Action(context)

		// then
		var expectedResponse gin.H
		err := json.Unmarshal(recorder.Body.Bytes(), &expectedResponse)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, gin.H{"message": "Workflow received"}, expectedResponse)
		mockedExecutor.AssertExpectations(t)
	})

	t.Run("Should listen pull request", func(t *testing.T) {
		// given
		checkRunEvent := github.PullRequestEvent{}
		mockedExecutor := new(internalMocks.IExecutor)
		mockedExecutor.On("PullRequestEvent", checkRunEvent)

		event, _ := json.Marshal(checkRunEvent)
		context, recorder := mocks.MockContext("", string(event))

		// when
		NewListener(mockedExecutor).PullRequest(context)

		// then
		var expectedResponse gin.H
		err := json.Unmarshal(recorder.Body.Bytes(), &expectedResponse)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, gin.H{"message": "Workflow received"}, expectedResponse)
		mockedExecutor.AssertExpectations(t)
	})

	t.Run("Should listen pull request review", func(t *testing.T) {
		// given
		pullRequestReviewEvent := github.PullRequestReviewEvent{}
		mockedExecutor := new(internalMocks.IExecutor)
		mockedExecutor.On("PullRequestReviewEvent", pullRequestReviewEvent)

		event, _ := json.Marshal(pullRequestReviewEvent)
		context, recorder := mocks.MockContext("", string(event))

		// when
		NewListener(mockedExecutor).PullRequestReview(context)

		// then
		var expectedResponse gin.H
		err := json.Unmarshal(recorder.Body.Bytes(), &expectedResponse)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, gin.H{"message": "Workflow received"}, expectedResponse)
		mockedExecutor.AssertExpectations(t)
	})

}
