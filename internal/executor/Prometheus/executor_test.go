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

func TestExecutorPrometheusEventWorkflowRun(t *testing.T) {
	t.Run("Should count workflow run and delete success runs", func(t *testing.T) {
		// given
		run := github.WorkflowRunEvent{
			WorkflowRun: &github.WorkflowRun{
				Name:       github.String("main"),
				ID:         github.Int64(1),
				RunNumber:  github.Int(1),
				Conclusion: github.String("failure"),
				Status:     github.String("completed"),
				Repository: &github.Repository{FullName: github.String("otto-wagner/github-observer")},
			},
		}

		successRun := github.WorkflowRunEvent{
			WorkflowRun: &github.WorkflowRun{
				ID:         github.Int64(1),
				RunNumber:  github.Int(1),
				Conclusion: github.String("success"),
				Repository: &github.Repository{FullName: github.String("otto-wagner/github-observer")},
			},
		}

		// when
		exec.EventWorkflowRun(run)

		// then
		assert.Equal(t, 1, testutil.CollectAndCount(exec.workflowRun))
		// []string{"repository_full_name", "workflow_name", "workflow_run_id", "run_number","state", "conclusion"})
		assert.Greater(t, testutil.ToFloat64(exec.workflowRun.WithLabelValues(
			run.GetWorkflowRun().GetRepository().GetFullName(),
			run.GetWorkflowRun().GetName(),
			strconv.FormatInt(run.GetWorkflowRun().GetID(), 10),
			strconv.Itoa(run.GetWorkflowRun().GetRunNumber()),
			run.GetWorkflowRun().GetStatus(),
			run.GetWorkflowRun().GetConclusion(),
		)), float64(0))

		// when
		exec.EventWorkflowRun(successRun)

		// then
		assert.Equal(t, 0, testutil.CollectAndCount(exec.workflowRun))
	})

}

func TestExecutorPrometheusLastWorkflowRuns(t *testing.T) {
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
