//go:build all || unit

package Logging

import (
	"github-observer/internal/core"
	e "github-observer/internal/executor"
	logger "github-observer/pkg/mocks"
	"github.com/google/go-github/v61/github"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExecutorLogging(t *testing.T) {

	t.Run("Should return name", func(t *testing.T) {
		// given
		// when
		name := NewExecutor(nil).Name()

		// then
		assert.Equal(t, "Logging", name)
	})

}

func TestExecutorLoggingEvents(t *testing.T) {

	t.Run("Should return name", func(t *testing.T) {
		// given
		// when
		name := NewExecutor(nil).Name()

		// then
		assert.Equal(t, "Logging", name)
	})

	t.Run("Should log run event", func(t *testing.T) {
		// given
		logs := logger.MockedLogger()

		now := time.Now()
		event := github.CheckRunEvent{
			Repo: &github.Repository{
				Name:    github.String("github-observer"),
				HTMLURL: github.String("https://github.com/otto-wagner/github-observer"),
			},
			Action: github.String("completed"),
			CheckRun: &github.CheckRun{
				Name:        github.String("Analyze (go)"),
				HTMLURL:     github.String("https://github.com/otto-wagner/github-observer/actions/runs/8589035842/job/23534635896"),
				Status:      github.String("completed"),
				Conclusion:  github.String("success"),
				StartedAt:   &github.Timestamp{Time: now},
				CompletedAt: &github.Timestamp{Time: now},
			},
		}

		// when
		NewExecutor(nil).RunEvent(event)

		// then
		assert.Contains(t, logs.All()[0].Message, "Event")
		assert.Contains(t, logs.All()[0].Context[0].Key, "Run")
		assert.Equal(t, logs.All()[0].Context[0].Interface, core.ConvertToGitAction(event))
	})

	t.Run("Should log pull request event", func(t *testing.T) {
		// given
		logs := logger.MockedLogger()
		now := time.Now()
		event := github.PullRequestEvent{
			Repo: &github.Repository{
				Name:    github.String("github-observer"),
				HTMLURL: github.String("https://github.com/otto-wagner/github-observer"),
			},
			Action: github.String("opened"),
			PullRequest: &github.PullRequest{
				Title:     github.String("chore: test pullrequest_listener"),
				User:      &github.User{Login: github.String("otto-wagner")},
				HTMLURL:   github.String("https://github.com/otto-wagner/github-observer/pull/2"),
				State:     github.String("open"),
				CreatedAt: &github.Timestamp{Time: now},
				UpdatedAt: &github.Timestamp{Time: now},
				ClosedAt:  &github.Timestamp{Time: now},
				MergedAt:  &github.Timestamp{Time: now},
			},
		}

		// when
		NewExecutor(nil).PullRequestEvent(event)

		// then
		assert.Contains(t, logs.All()[0].Message, "Event")
		assert.Contains(t, logs.All()[0].Context[0].Key, "PullRequest")
		assert.Equal(t, logs.All()[0].Context[0].Interface, core.ConvertPREToGitPullRequest(event))
	})

	t.Run("Should log pull request review event", func(t *testing.T) {
		// given
		logs := logger.MockedLogger()
		now := time.Now()
		event := github.PullRequestReviewEvent{
			Repo: &github.Repository{
				Name:    github.String("github-observer"),
				HTMLURL: github.String("https://github.com/otto-wagner/github-observer"),
			},
			Action: github.String("submitted"),
			PullRequest: &github.PullRequest{
				Title:     github.String("chore: test pullrequest_listener"),
				User:      &github.User{Login: github.String("otto-wagner")},
				HTMLURL:   github.String("https://github.com/otto-wagner/github-observer/pull/2"),
				State:     github.String("open"),
				CreatedAt: &github.Timestamp{Time: now},
				UpdatedAt: &github.Timestamp{Time: now},
				ClosedAt:  &github.Timestamp{Time: now},
				MergedAt:  &github.Timestamp{Time: now},
			},
			Review: &github.PullRequestReview{
				Body:    github.String("LGTM"),
				State:   github.String("commented"),
				User:    &github.User{Login: github.String("otto-wagner")},
				HTMLURL: github.String("https://github.com/otto-wagner/github-observer/pull/2#pullrequestreview-1985121262"),
			},
		}

		// when
		NewExecutor(nil).PullRequestReviewEvent(event)

		// then
		assert.Contains(t, logs.All()[0].Message, "Event")
		assert.Contains(t, logs.All()[0].Context[0].Key, "PullRequestReview")
		assert.Equal(t, logs.All()[0].Context[0].Interface, core.ConvertToGitPullRequestReview(event))
	})

}

func TestExecutorLoggingWorkflowRuns(t *testing.T) {

	t.Run("Should log failed workflow run", func(t *testing.T) {
		// given
		logs := logger.MockedLogger()
		memory := e.NewMemory()
		now := time.Now()

		workflowRun := github.WorkflowRun{
			WorkflowID:   github.Int64(1),
			Name:         github.String("main"),
			HeadBranch:   github.String("main"),
			Event:        github.String("push"),
			DisplayTitle: github.String("feat: add prometheus and grafana in docker-compose"),
			Status:       github.String("completed"),
			Conclusion:   github.String("failure"),
			HTMLURL:      github.String("https://github.com/otto-wagner/github-observer/actions/runs/8712037001"),
			Actor:        &github.User{Login: github.String("otto-wagner")},
			CreatedAt:    &github.Timestamp{Time: now},
			UpdatedAt:    &github.Timestamp{Time: now},
			HeadCommit:   &github.HeadCommit{Message: github.String("chore: test")},
			Repository:   &github.Repository{Name: github.String("github-observer")},
		}

		// when
		NewExecutor(memory).WorkflowRuns([]*github.WorkflowRun{&workflowRun})

		// then
		assert.Contains(t, logs.All()[0].Message, "WorkflowRun")
		assert.Contains(t, logs.All()[0].Context[0].Key, "Action")
		assert.Equal(t, logs.All()[0].Context[0].Interface, core.ConvertToWorkflow(workflowRun))
	})

	t.Run("Should not log success workflow run", func(t *testing.T) {
		// given
		logs := logger.MockedLogger()
		memory := e.NewMemory()

		runs := []*github.WorkflowRun{
			{RunNumber: github.Int(1), HeadBranch: github.String("main"), Conclusion: github.String("success")},
		}

		// when
		NewExecutor(memory).WorkflowRuns(runs)

		// then
		assert.Lenf(t, logs.All(), 0, "Expected 0 logs, got %d", len(logs.All()))
	})

	t.Run("Should log failed workflow run one time", func(t *testing.T) {
		// given
		logs := logger.MockedLogger()
		memory := e.NewMemory()

		runs := []*github.WorkflowRun{
			{RunNumber: github.Int(1), HeadBranch: github.String("main"), Conclusion: github.String("failed")},
			{RunNumber: github.Int(1), HeadBranch: github.String("main"), Conclusion: github.String("failed")},
		}

		// when
		NewExecutor(memory).WorkflowRuns(runs)

		// then
		assert.Lenf(t, logs.All(), 1, "Expected 1 log, got %d", len(logs.All()))
	})

	t.Run("Should only log newest workflow runs", func(t *testing.T) {
		// given
		logs := logger.MockedLogger()
		memory := e.NewMemory()

		runs := []*github.WorkflowRun{
			{WorkflowID: github.Int64(1), RunNumber: github.Int(1), Name: github.String("w1 r1 - should be called"), HeadBranch: github.String("main"), Conclusion: github.String("failed")},
			{WorkflowID: github.Int64(1), RunNumber: github.Int(1), Name: github.String("w1 r1 - should be ignored"), HeadBranch: github.String("main"), Conclusion: github.String("failed")},
			{WorkflowID: github.Int64(2), RunNumber: github.Int(1), Name: github.String("w2 r1 - should be called"), HeadBranch: github.String("main"), Conclusion: github.String("failed")},
		}

		// when
		NewExecutor(memory).WorkflowRuns(runs)

		// then
		assert.Lenf(t, logs.All(), 2, "Expected 2 logs, got %d", len(logs.All()))
		assert.Equal(t, "w1 r1 - should be called", logs.All()[0].Context[0].Interface.(core.WorkflowRun).Name)
		assert.Equal(t, "w2 r1 - should be called", logs.All()[1].Context[0].Interface.(core.WorkflowRun).Name)
	})

	t.Run("Should also log newest workflow run", func(t *testing.T) {
		// given
		logs := logger.MockedLogger()
		memory := e.NewMemory()

		runs := []*github.WorkflowRun{
			{RunNumber: github.Int(1), Name: github.String("old"), HeadBranch: github.String("main"), Conclusion: github.String("failed")},
			{RunNumber: github.Int(1), Name: github.String("ignored"), HeadBranch: github.String("main"), Conclusion: github.String("failed")},
			{RunNumber: github.Int(2), Name: github.String("new"), HeadBranch: github.String("main"), Conclusion: github.String("failed")},
		}

		// when
		NewExecutor(memory).WorkflowRuns(runs)

		// then
		assert.Lenf(t, logs.All(), 2, "Expected 2 logs, got %d", len(logs.All()))
		assert.Equal(t, "old", logs.All()[0].Context[0].Interface.(core.WorkflowRun).Name)
		assert.Equal(t, "new", logs.All()[1].Context[0].Interface.(core.WorkflowRun).Name)
	})

	t.Run("Should log workflow run in each branch", func(t *testing.T) {
		// given
		logs := logger.MockedLogger()
		memory := e.NewMemory()

		runs := []*github.WorkflowRun{
			{RunNumber: github.Int(1), HeadBranch: github.String("main"), Conclusion: github.String("failed")},
			{RunNumber: github.Int(1), HeadBranch: github.String("another"), Conclusion: github.String("failed")},
		}

		// when
		NewExecutor(memory).WorkflowRuns(runs)

		// then
		assert.Lenf(t, logs.All(), 2, "Expected 2 logs, got %d", len(logs.All()))
		assert.Equal(t, "main", logs.All()[0].Context[0].Interface.(core.WorkflowRun).HeadBranch)
		assert.Equal(t, "another", logs.All()[1].Context[0].Interface.(core.WorkflowRun).HeadBranch)
	})

}

func TestExecutorLoggingPullRequests(t *testing.T) {

	t.Run("Should log pull request", func(t *testing.T) {
		// given
		logs := logger.MockedLogger()
		memory := e.NewMemory()
		now := time.Now()

		repository := core.Repository{
			Name:    "github-observer",
			Owner:   "otto-wagner",
			HtmlUrl: "https://github.com/otto-wagner/github-observer",
		}

		pullRequests := []*github.PullRequest{{
			Number:    github.Int(1),
			State:     github.String("open"),
			Title:     github.String("chore: test pullrequest_listener"),
			Body:      github.String("LGTM"),
			User:      &github.User{Login: github.String("otto-wagner")},
			HTMLURL:   github.String("https://github.com/otto-wagner/github-observer/actions/runs/8712037001"),
			CreatedAt: &github.Timestamp{Time: now},
			UpdatedAt: &github.Timestamp{Time: now},
			ClosedAt:  &github.Timestamp{Time: now},
			MergedAt:  &github.Timestamp{Time: now},
		}}

		// when
		NewExecutor(memory).PullRequests(repository, pullRequests)

		// then
		assert.Lenf(t, logs.All(), 1, "Expected 1 log, got %d", len(logs.All()))
		assert.Contains(t, logs.All()[0].Message, "PullRequest")
		assert.Contains(t, logs.All()[0].Context[0].Key, "PullRequest")
		assert.Equal(t, logs.All()[0].Context[0].Interface, core.ConvertPRToGitPullRequest(repository, *pullRequests[0]))
	})

	t.Run("Should not log closed and not merged pull request", func(t *testing.T) {
		// given
		logs := logger.MockedLogger()
		pullRequests := []*github.PullRequest{{State: github.String("closed")}, {State: github.String("merged")}}

		// when
		NewExecutor(nil).PullRequests(core.Repository{}, pullRequests)

		// then
		assert.Lenf(t, logs.All(), 0, "Expected 0 logs, got %d", len(logs.All()))
	})

	t.Run("Should only log the pull request one time", func(t *testing.T) {
		// given
		logs := logger.MockedLogger()
		memory := e.NewMemory()
		pullRequests := []*github.PullRequest{
			{Number: github.Int(1), State: github.String("open")},
			{Number: github.Int(1), State: github.String("open")},
		}

		// when
		NewExecutor(memory).PullRequests(core.Repository{}, pullRequests)

		// then
		assert.Lenf(t, logs.All(), 1, "Expected 1 log, got %d", len(logs.All()))
	})

}
