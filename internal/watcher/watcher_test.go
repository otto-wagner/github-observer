//go:build all || unit

package watcher

import (
	"testing"

	"github.com/google/go-github/v73/github"
	gitMock "github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/otto-wagner/github-observer/internal/core"
	"github.com/otto-wagner/github-observer/internal/executor"
	"github.com/otto-wagner/github-observer/mocks"
	"github.com/stretchr/testify/mock"
)

func TestWatchPullRequests(t *testing.T) {
	t.Run("Should list pull requests and send them to executors", func(t *testing.T) {
		// given
		repository := core.Repository{Owner: "otto-wagner", Name: "github-observer"}

		mockedExecutor := new(mocks.MockIExecutor)
		mockedExecutor.On("PullRequests", repository, mock.Anything)
		mockedSecondExecutor := new(mocks.MockIExecutor)
		mockedSecondExecutor.On("PullRequests", repository, mock.Anything)

		mockedHTTPClient := gitMock.NewMockedHTTPClient(
			gitMock.WithRequestMatch(gitMock.GetReposPullsByOwnerByRepo,
				[]github.PullRequest{{ID: github.Ptr(int64(1))}, {ID: github.Ptr(int64(2))}}),
		)

		// when
		w := NewWatcher(github.NewClient(mockedHTTPClient), []core.Repository{repository}, []executor.IExecutor{mockedExecutor, mockedSecondExecutor})
		w.PullRequests(repository)

		// then
		mockedExecutor.AssertExpectations(t)
		mockedSecondExecutor.AssertExpectations(t)
	})
}

func TestWatchWorkflows(t *testing.T) {
	t.Run("Should get latest workflow runs and send them to executors", func(t *testing.T) {
		// given
		repository := core.Repository{Owner: "otto-wagner", Name: "github-observer"}

		workflows := []*github.Workflow{{ID: github.Ptr(int64(1)), Name: github.Ptr("CodeQL")}}
		workflowRuns := []*github.WorkflowRun{{ID: github.Ptr(int64(1))}, {ID: github.Ptr(int64(2))}}

		mockedExecutor := new(mocks.MockIExecutor)
		mockedExecutor.On("LastWorkflows", repository, mock.Anything)
		mockedSecondExecutor := new(mocks.MockIExecutor)
		mockedSecondExecutor.On("LastWorkflows", repository, mock.Anything)

		mockedGithubClient := gitMock.NewMockedHTTPClient(
			gitMock.WithRequestMatch(gitMock.GetReposActionsWorkflowsByOwnerByRepo,
				github.Workflows{Workflows: workflows}),
			gitMock.WithRequestMatch(gitMock.GetReposActionsWorkflowsRunsByOwnerByRepoByWorkflowId,
				github.WorkflowRuns{WorkflowRuns: workflowRuns}),
		)

		// when
		w := NewWatcher(github.NewClient(mockedGithubClient), []core.Repository{repository}, []executor.IExecutor{mockedExecutor, mockedSecondExecutor})
		w.WorkflowRuns(repository)

		// then
		mockedExecutor.AssertExpectations(t)
		mockedSecondExecutor.AssertExpectations(t)
	})
}
