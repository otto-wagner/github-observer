//go:build all || unit

package Logger

import (
	"bytes"
	"github-observer/internal/core"
	"github.com/google/go-github/v61/github"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
	"time"
)

func TestExecutorLoggingEventWorkflowRun(t *testing.T) {
	t.Run("Should log workflow run event", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		logger := slog.New(slog.NewTextHandler(buf, nil))

		now := time.Now()
		event := github.WorkflowRunEvent{
			WorkflowRun: &github.WorkflowRun{
				WorkflowID:   github.Int64(1),
				RunNumber:    github.Int(1),
				Name:         github.String("WorkflowName"),
				HeadBranch:   github.String("main"),
				Event:        github.String("Run"),
				DisplayTitle: github.String("feat: add prometheus and grafana in docker-compose"),
				Status:       github.String("completed"),
				Conclusion:   github.String("failure"),
				HTMLURL:      github.String("anyUrl"),
				Actor:        &github.User{Login: github.String("otto-wagner")},
				CreatedAt:    &github.Timestamp{Time: now},
				UpdatedAt:    &github.Timestamp{Time: now},
				HeadCommit:   &github.HeadCommit{Message: github.String("chore: test")},
				Repository: &github.Repository{
					FullName: github.String("otto-wagner/github-observer"),
					HTMLURL:  github.String("repoUrl"),
				},
			},
		}

		// when
		NewExecutor(nil, logger).EventWorkflowRun(event)

		// then
		assert.Contains(t, buf.String(), "Event")
		assert.Contains(t, buf.String(), "Name")
		assert.Contains(t, buf.String(), "WorkflowName")
		assert.Contains(t, buf.String(), "DisplayTitle")
		assert.Contains(t, buf.String(), "feat: add prometheus and grafana in docker-compose")
	})
}

func TestExecutorLoggingEventPullRequest(t *testing.T) {

	t.Run("Should log pull request event", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		logger := slog.New(slog.NewTextHandler(buf, nil))

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
		NewExecutor(nil, logger).EventPullRequest(event)

		// then
		assert.Contains(t, buf.String(), "Event")
		assert.Contains(t, buf.String(), "Name")
		assert.Contains(t, buf.String(), "PullRequest")
		assert.Contains(t, buf.String(), "Title")
		assert.Contains(t, buf.String(), "chore: test pullrequest_listener")
	})

}

func TestExecutorLoggingEventPullRequestReview(t *testing.T) {

	t.Run("Should log pull request review event", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		logger := slog.New(slog.NewTextHandler(buf, nil))

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
		NewExecutor(nil, logger).EventPullRequestReview(event)

		// then
		assert.Contains(t, buf.String(), "Event")
		assert.Contains(t, buf.String(), "Name")
		assert.Contains(t, buf.String(), "PullRequestReview")
		assert.Contains(t, buf.String(), "Title")
		assert.Contains(t, buf.String(), "chore: test pullrequest_listener")
	})

}

func TestExecutorLoggingWorkflowRuns(t *testing.T) {

	t.Run("Should log failed workflow run", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		logger := slog.New(slog.NewTextHandler(buf, nil))
		now := time.Now()

		repository := core.Repository{Name: "github-observer", Owner: "otto-wagner"}
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
		NewExecutor(NewMemory(), logger).LastWorkflows(repository, []*github.WorkflowRun{&workflowRun})

		// then
		assert.Contains(t, buf.String(), "Event")
		assert.Contains(t, buf.String(), "Name")
		assert.Contains(t, buf.String(), "Action")
		assert.Contains(t, buf.String(), "Title")
		assert.Contains(t, buf.String(), "feat: add prometheus and grafana in docker-compose")
		assert.Contains(t, buf.String(), "Conclusion")
		assert.Contains(t, buf.String(), "failure")
	})

	t.Run("Should not log success workflow run", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		logger := slog.New(slog.NewTextHandler(buf, nil))

		repository := core.Repository{Name: "github-observer", Owner: "otto-wagner"}
		runs := []*github.WorkflowRun{
			{WorkflowID: github.Int64(1), Repository: &github.Repository{FullName: github.String("otto-wagner/github-observer")}, RunNumber: github.Int(1), HeadBranch: github.String("main"), Conclusion: github.String("success")},
		}

		// when
		NewExecutor(NewMemory(), logger).LastWorkflows(repository, runs)

		// then
		assert.Empty(t, buf.String())
	})

	t.Run("Should not log dual times the same log", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		logger := slog.New(slog.NewTextHandler(buf, nil))

		repository := core.Repository{Name: "github-observer", Owner: "otto-wagner"}
		runs := []*github.WorkflowRun{
			{WorkflowID: github.Int64(1), RunNumber: github.Int(1), Name: github.String("first"), HeadBranch: github.String("main"), Conclusion: github.String("failed"), Repository: &github.Repository{FullName: github.String("otto-wagner/github-observer")}},
			{WorkflowID: github.Int64(1), RunNumber: github.Int(1), Name: github.String("duplicate"), HeadBranch: github.String("main"), Conclusion: github.String("failed"), Repository: &github.Repository{FullName: github.String("otto-wagner/github-observer")}},
			{WorkflowID: github.Int64(1), RunNumber: github.Int(2), Name: github.String("second"), HeadBranch: github.String("main"), Conclusion: github.String("failed"), Repository: &github.Repository{FullName: github.String("otto-wagner/github-observer")}},
		}

		// when
		NewExecutor(NewMemory(), logger).LastWorkflows(repository, runs)

		// then
		assert.Contains(t, buf.String(), "Event")
		assert.Contains(t, buf.String(), "Name")
		assert.Contains(t, buf.String(), "first", "Expected to contain first, got %s", buf.String())
		assert.NotContainsf(t, buf.String(), "duplicate", "Expected not to contain duplicate, got %s", buf.String())
		assert.Contains(t, buf.String(), "second", "Expected to contain second, got %s", buf.String())
	})

}

func TestExecutorLoggingPullRequests(t *testing.T) {

	t.Run("Should log pull request", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		logger := slog.New(slog.NewTextHandler(buf, nil))
		now := time.Now()

		repository := core.Repository{Name: "github-observer", Owner: "otto-wagner"}
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
		NewExecutor(NewMemory(), logger).PullRequests(repository, pullRequests)

		// then
		assert.Contains(t, buf.String(), "PullRequest")
		assert.Contains(t, buf.String(), "State")
		assert.Contains(t, buf.String(), "open")
		assert.Contains(t, buf.String(), "Title")
		assert.Contains(t, buf.String(), "chore: test pullrequest_listener")
	})

	t.Run("Should not log closed and not merged pull request", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		logger := slog.New(slog.NewTextHandler(buf, nil))
		repository := core.Repository{Name: "github-observer", Owner: "otto-wagner"}
		pullRequests := []*github.PullRequest{{State: github.String("closed")}, {State: github.String("merged")}}

		// when
		NewExecutor(NewMemory(), logger).PullRequests(repository, pullRequests)

		// then
		assert.Empty(t, buf.String())
	})

	t.Run("Should only log the pull request one time", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		logger := slog.New(slog.NewTextHandler(buf, nil))

		repository := core.Repository{Name: "github-observer", Owner: "otto-wagner"}
		pullRequests := []*github.PullRequest{
			{Number: github.Int(1), State: github.String("open"), Title: github.String("Ensure unique logging")},
		}

		// when
		newMemory := NewMemory()
		NewExecutor(newMemory, logger).PullRequests(repository, pullRequests)
		NewExecutor(newMemory, logger).PullRequests(repository, pullRequests)

		// then
		assert.Equal(t, 1, bytes.Count([]byte(buf.String()), []byte("Ensure unique logging")), "The pull request should be logged only once")
	})

}
