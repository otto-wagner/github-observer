//go:build all || unit

package watcher

import (
	"github-observer/internal/core"
	"github-observer/internal/executor"
	internalMocks "github-observer/internal/mocks"
	"github.com/google/go-github/v61/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"testing"
)

func TestWatcherPullRequests(t *testing.T) {

	t.Run("Should watch pull requests and send them to executors", func(t *testing.T) {
		// given
		repository := core.Repository{Owner: "otto-wagner", Name: "github-observer"}
		pullRequests := []*github.PullRequest{{Title: github.String(repository.Owner)}}

		mockedHTTPClient := mock.NewMockedHTTPClient(mock.WithRequestMatch(mock.GetReposPullsByOwnerByRepo, pullRequests))

		mockedGithubClient := github.NewClient(mockedHTTPClient)
		mockedExecutor := new(internalMocks.IExecutor)
		mockedExecutor.On("PullRequests", pullRequests)

		// when
		w := NewWatcher(mockedGithubClient, nil, []executor.IExecutor{mockedExecutor})
		w.PullRequests(repository)

		// then
		mockedExecutor.AssertExpectations(t)
	})
}

func TestWatcherWorkflowRuns(t *testing.T) {

	t.Run("Should watch workflow runs and send them to executors", func(t *testing.T) {
		// given
		repository := core.Repository{Owner: "otto-wagner", Name: "github-observer"}
		workflows := []*github.Workflow{{ID: github.Int64(1), Name: github.String("CodeQL")}}
		workflowRuns := []*github.WorkflowRun{{ID: github.Int64(1)}}

		mockedHTTPClient := mock.NewMockedHTTPClient(
			mock.WithRequestMatch(mock.GetReposActionsWorkflowsByOwnerByRepo,
				github.Workflows{TotalCount: github.Int(1), Workflows: workflows}),
			mock.WithRequestMatch(mock.GetReposActionsWorkflowsRunsByOwnerByRepoByWorkflowId,
				github.WorkflowRuns{TotalCount: github.Int(1), WorkflowRuns: workflowRuns}),
		)

		mockedExecutor := new(internalMocks.IExecutor)
		mockedExecutor.On("LastWorkflows", workflowRuns)

		// when
		w := NewWatcher(github.NewClient(mockedHTTPClient), nil, []executor.IExecutor{mockedExecutor})
		w.WorkflowRuns(repository)

		// then
		mockedExecutor.AssertExpectations(t)
	})
}
