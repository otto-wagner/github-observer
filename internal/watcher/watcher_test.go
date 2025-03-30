//go:build all || unit

package watcher

import (
	"bytes"
	"github-observer/internal/core"
	"github-observer/internal/executor"
	"github-observer/internal/executor/mocks"
	"github.com/google/go-github/v61/github"
	gitMock "github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"testing"
)

func TestWatchPullRequests(t *testing.T) {
	t.Run("Should list pull requests and send them to executors", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		logger := slog.New(slog.NewTextHandler(buf, nil))

		repository := core.Repository{Owner: "otto-wagner", Name: "github-observer"}

		mockedExecutor := new(mocks.IExecutor)
		mockedExecutor.On("PullRequests", repository, mock.Anything)
		mockedSecondExecutor := new(mocks.IExecutor)
		mockedSecondExecutor.On("PullRequests", repository, mock.Anything)

		mockedGithubClient := gitMock.NewMockedHTTPClient(
			gitMock.WithRequestMatch(gitMock.GetReposPullsByOwnerByRepo,
				[]github.PullRequest{{ID: github.Int64(1)}, {ID: github.Int64(2)}}),
		)

		// when
		w := NewWatcher("token", github.NewClient(mockedGithubClient), []core.Repository{repository}, []executor.IExecutor{mockedExecutor, mockedSecondExecutor}, logger)
		w.PullRequests(repository)

		// then
		mockedExecutor.AssertExpectations(t)
		mockedSecondExecutor.AssertExpectations(t)
	})
}

func TestWatchWorkflows(t *testing.T) {
	t.Run("Should get latest workflow runs and send them to executors", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		logger := slog.New(slog.NewTextHandler(buf, nil))
		repository := core.Repository{Owner: "otto-wagner", Name: "github-observer"}

		workflows := []*github.Workflow{{ID: github.Int64(1), Name: github.String("CodeQL")}}
		workflowRuns := []*github.WorkflowRun{{ID: github.Int64(1)}, {ID: github.Int64(2)}}

		mockedExecutor := new(mocks.IExecutor)
		mockedExecutor.On("LastWorkflows", repository, mock.Anything)
		mockedSecondExecutor := new(mocks.IExecutor)
		mockedSecondExecutor.On("LastWorkflows", repository, mock.Anything)

		mockedGithubClient := gitMock.NewMockedHTTPClient(
			gitMock.WithRequestMatch(gitMock.GetReposActionsWorkflowsByOwnerByRepo,
				github.Workflows{Workflows: workflows}),
			gitMock.WithRequestMatch(gitMock.GetReposActionsWorkflowsRunsByOwnerByRepoByWorkflowId,
				github.WorkflowRuns{WorkflowRuns: workflowRuns}),
		)

		// when
		w := NewWatcher("token", github.NewClient(mockedGithubClient), []core.Repository{repository}, []executor.IExecutor{mockedExecutor, mockedSecondExecutor}, logger)
		w.WorkflowRuns(repository)

		// then
		mockedExecutor.AssertExpectations(t)
		mockedSecondExecutor.AssertExpectations(t)
	})
}
