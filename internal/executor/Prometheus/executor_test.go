//go:build all || unit

package Prometheus

import (
	"github-observer/internal/core"
	"github.com/google/go-github/v61/github"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

var exec *executor

func init() {
	exec = NewExecutor().(*executor)
}

func TestExecutorPrometheusEventRun(t *testing.T) {
	t.Run("Should count action", func(t *testing.T) {
		// given
		event := github.CheckRunEvent{
			Action: github.String("completed"), CheckRun: &github.CheckRun{
				ID:         github.Int64(1),
				Name:       github.String("Analyze (go)"),
				Status:     github.String("completed"),
				Conclusion: github.String("success"),
			},
			Repo: &github.Repository{
				FullName: github.String("otto-wagner/github-observer"),
			},
		}

		// when
		exec.EventRun(event)

		// then
		assert.Equal(t, 1, testutil.CollectAndCount(exec.eventRun))
		assert.Greater(t, testutil.ToFloat64(exec.eventRun.WithLabelValues(
			event.GetAction(),
			strconv.FormatInt(event.GetCheckRun().GetID(), 10),
			event.GetCheckRun().GetName(),
			event.GetCheckRun().GetStatus(),
			event.GetCheckRun().GetConclusion(),
			event.GetRepo().GetFullName(),
		)), float64(0))
	})
}

func TestExecutorPrometheusEventPullRequest(t *testing.T) {

	t.Run("Should count pull request", func(t *testing.T) {
		// given
		event := github.PullRequestEvent{
			Action: github.String("opened"),
			PullRequest: &github.PullRequest{
				Title: github.String("chore: test pullrequest_listener"),
				State: github.String("open"),
			},
			Repo: &github.Repository{
				FullName: github.String("otto-wagner/github-observer"),
			},
		}

		// when
		exec.EventPullRequest(event)

		// then
		assert.Equal(t, 1, testutil.CollectAndCount(exec.eventPullRequest))
		assert.Greater(t, testutil.ToFloat64(exec.eventPullRequest.WithLabelValues(
			event.GetAction(),
			event.GetPullRequest().GetTitle(),
			event.GetPullRequest().GetState(),
			event.GetRepo().GetFullName(),
		)), float64(0))

	})

}

func TestExecutorPrometheusEventPullRequestReview(t *testing.T) {
	// ignored
}

func TestExecutorPrometheusWorkflowRuns(t *testing.T) {
	t.Run("Should count failed workflow run", func(t *testing.T) {
		// given
		repository := core.Repository{Name: "github-observer", Owner: "otto-wagner"}
		runs := []*github.WorkflowRun{{
			Conclusion: github.String("failure"),
			Repository: &github.Repository{FullName: github.String("otto-wagner/github-observer")},
		}, {
			Conclusion: github.String("success"),
			Repository: &github.Repository{FullName: github.String("otto-wagner/github-observer")},
		}}

		// when
		exec.LastWorkflows(repository, runs)

		// then
		assert.Equal(t, 1, testutil.CollectAndCount(exec.workflowRun))
	})
}

func TestExecutorPrometheusPullRequests(t *testing.T) {
	t.Run("Should count pull request", func(t *testing.T) {
		// given
		repository := core.Repository{Name: "github-observer", Owner: "otto-wagner"}

		pullRequests := []*github.PullRequest{{
			State: github.String("open"),
			Base: &github.PullRequestBranch{
				Repo: &github.Repository{FullName: github.String("otto-wagner/github-observer")},
			}},
			{State: github.String("closed"),
				Base: &github.PullRequestBranch{
					Repo: &github.Repository{FullName: github.String("otto-wagner/github-observer")},
				}},
			{State: github.String("merged"),
				Base: &github.PullRequestBranch{
					Repo: &github.Repository{FullName: github.String("otto-wagner/github-observer")},
				}},
		}

		// when
		exec.PullRequests(repository, pullRequests)

		// then
		assert.Equal(t, 1, testutil.CollectAndCount(exec.pullRequest))
	})

	t.Run("Should not count pull requests when no pull requests given", func(t *testing.T) {
		// given
		repository := core.Repository{Name: "github-observer", Owner: "otto-wagner"}
		var pullRequests []*github.PullRequest

		// when
		exec.PullRequests(repository, pullRequests)

		// then
		assert.Equal(t, 0, testutil.CollectAndCount(exec.pullRequest))
	})

}

func TestExecutorPrometheusPullRequestReview(t *testing.T) {
	// ignored
	t.Run("Ignored", func(t *testing.T) {
		// given
		// when
		exec.EventPullRequestReview(github.PullRequestReviewEvent{})
		// then
	})
}
