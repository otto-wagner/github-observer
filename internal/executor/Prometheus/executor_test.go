//go:build all || unit

package Prometheus

import (
	"github-observer/internal/core"
	"github.com/google/go-github/v61/github"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

var exec *executor

func init() {
	exec = NewExecutor().(*executor)
}

func TestExecutorPrometheus(t *testing.T) {
	t.Run("Should return name", func(t *testing.T) {
		// given
		// when
		assert.Equal(t, "Prometheus", exec.Name())
	})

}

func TestExecutorPrometheusEvents(t *testing.T) {
	t.Run("Should count action", func(t *testing.T) {
		// given
		now := time.Now()
		event := github.CheckRunEvent{
			Action: github.String("completed"),
			CheckRun: &github.CheckRun{
				Name:        github.String("Analyze (go)"),
				HTMLURL:     github.String("https://github.com/otto-wagner/github-observer/actions/runs/8589035842/job/23534635896"),
				Status:      github.String("completed"),
				Conclusion:  github.String("success"),
				StartedAt:   &github.Timestamp{Time: now},
				CompletedAt: &github.Timestamp{Time: now},
			},
			Repo: &github.Repository{
				Name:    github.String("github-observer"),
				HTMLURL: github.String("https://github.com/otto-wagner/github-observer"),
			},
		}

		// when
		exec.RunEvent(event)

		// then
		assert.Equal(t, 1, testutil.CollectAndCount(exec.countAction))
		assert.Greater(t, testutil.ToFloat64(exec.countAction.WithLabelValues(
			event.GetAction(),
			event.GetCheckRun().GetName(),
			event.GetCheckRun().GetHTMLURL(),
			event.GetCheckRun().GetStatus(),
			event.GetCheckRun().GetConclusion(),
			event.GetCheckRun().GetStartedAt().String(),
			event.GetCheckRun().GetCompletedAt().String(),
			event.GetRepo().GetName(),
			event.GetRepo().GetHTMLURL(),
		)), float64(0))
	})

	t.Run("Should count pull request", func(t *testing.T) {
		// given
		now := time.Now()
		event := github.PullRequestEvent{
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
			Repo: &github.Repository{
				Name:    github.String("github-observer"),
				HTMLURL: github.String("https://github.com/otto-wagner/github-observer"),
			},
		}

		// when
		exec.PullRequestEvent(event)

		// then
		assert.Equal(t, 1, testutil.CollectAndCount(exec.countPR))
		assert.Greater(t, testutil.ToFloat64(exec.countPR.WithLabelValues(
			event.GetAction(),
			strconv.Itoa(event.GetPullRequest().GetNumber()),
			event.GetPullRequest().GetTitle(),
			event.GetPullRequest().GetBody(),
			event.GetPullRequest().GetUser().GetLogin(),
			event.GetPullRequest().GetHTMLURL(),
			event.GetPullRequest().GetState(),
			event.GetPullRequest().GetCreatedAt().String(),
			event.GetPullRequest().GetUpdatedAt().String(),
			event.GetPullRequest().GetClosedAt().String(),
			event.GetPullRequest().GetMergedAt().String(),
			event.GetRepo().GetName(),
			event.GetRepo().GetHTMLURL(),
		)), float64(0))
	})

	t.Run("Should count pull request review", func(t *testing.T) {
		// given
		event := github.PullRequestReviewEvent{
			Action: github.String("submitted"),
			PullRequest: &github.PullRequest{
				State:   github.String("open"),
				Title:   github.String("chore: test pullrequest_listener"),
				User:    &github.User{Login: github.String("otto-wagner")},
				HTMLURL: github.String("https://github.com/otto-wagner/github-observer/pull/2"),
			},
			Review: &github.PullRequestReview{
				State:   github.String("commented"),
				User:    &github.User{Login: github.String("otto-wagner")},
				Body:    github.String("LGTM"),
				HTMLURL: github.String("https://github.com/otto-wagner/github-observer/pull/2#pullrequestreview-1985121262"),
			},
			Repo: &github.Repository{
				Name:    github.String("github-observer"),
				HTMLURL: github.String("https://github.com/otto-wagner/github-observer"),
			},
		}

		// when
		exec.PullRequestReviewEvent(event)

		// then
		assert.Equal(t, 1, testutil.CollectAndCount(exec.countPRR))
		assert.Greater(t, testutil.ToFloat64(exec.countPRR.WithLabelValues(
			event.GetAction(),
			event.GetPullRequest().GetTitle(),
			event.GetPullRequest().GetUser().GetLogin(),
			event.GetPullRequest().GetHTMLURL(),
			event.GetPullRequest().GetState(),
			event.GetReview().GetBody(),
			event.GetReview().GetState(),
			event.GetReview().GetUser().GetLogin(),
			event.GetReview().GetHTMLURL(),
			event.GetRepo().GetName(),
			event.GetRepo().GetHTMLURL(),
		)), float64(0))
	})
}

func TestExecutorPrometheusWorkflowRuns(t *testing.T) {
	t.Run("Should set workflow runs", func(t *testing.T) {
		// given
		now := time.Now()
		workflowRun := github.WorkflowRun{
			WorkflowID:   github.Int64(1),
			RunNumber:    github.Int(1),
			Name:         github.String("test"),
			HeadBranch:   github.String("main"),
			Event:        github.String("push"),
			DisplayTitle: github.String("test"),
			Conclusion:   github.String("success"),
			HTMLURL:      github.String("https://github.com/otto-wagner/github-observer"),
			CreatedAt:    &github.Timestamp{Time: now},
			UpdatedAt:    &github.Timestamp{Time: now},
			Actor:        &github.User{Login: github.String("otto-wagner")},
			HeadCommit:   &github.HeadCommit{Message: github.String("test")},
			Repository: &github.Repository{
				Name:    github.String("github-observer"),
				HTMLURL: github.String("https://github.com/otto-wagner/github-observer"),
			},
		}
		workflow := core.ConvertToWorkflow(workflowRun)

		// when
		exec.WorkflowRuns([]*github.WorkflowRun{&workflowRun})

		// then
		assert.Equal(t, 1, testutil.CollectAndCount(exec.gaugeWorkflow))
		assert.Equal(t, float64(0), testutil.ToFloat64(exec.gaugeWorkflow.With(workflow.ToMap())))
	})

	//t.Run("Should set workflow runs each repository and last workflow", func(t *testing.T) {
	//	// given
	//	now := time.Now()
	//	workflowRun1 := github.WorkflowRun{
	//		WorkflowID:   github.Int64(1),
	//		RunNumber:    github.Int(1),
	//		Event:        github.String("push"),
	//		DisplayTitle: github.String("test"),
	//		Conclusion:   github.String("success"),
	//		HTMLURL:      github.String(""), // empty
	//
	//
	//}
}
