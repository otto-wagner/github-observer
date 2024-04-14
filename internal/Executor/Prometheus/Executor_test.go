//go:build all || unit

package Prometheus

import (
	"github.com/google/go-github/v61/github"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewExecutor(t *testing.T) {
	e := NewExecutor().(*executor)

	t.Run("Should return name", func(t *testing.T) {
		// given
		// when
		assert.Equal(t, "Prometheus", e.Name())
	})

	t.Run("Should log listener action", func(t *testing.T) {
		// given
		event := github.CheckRunEvent{
			Repo: &github.Repository{
				Name: github.String("github-observer"),
			},
			Action: github.String("completed"),
			CheckRun: &github.CheckRun{
				Status:     github.String("completed"),
				Conclusion: github.String("success"),
			},
		}

		// when
		err := e.CheckRunEvent(event)

		// then
		assert.NoError(t, err)
		assert.Equal(t, 1, testutil.CollectAndCount(e.lastRequestReceivedTimeAction))
		assert.Greater(t, testutil.ToFloat64(e.lastRequestReceivedTimeAction.WithLabelValues(
			event.GetRepo().GetName(),
			event.GetAction(),
			event.GetCheckRun().GetStatus(),
			event.GetCheckRun().GetConclusion(),
			"",
		)), float64(0))
		assert.Equal(t, 1, testutil.CollectAndCount(e.countAction))
		assert.Greater(t, testutil.ToFloat64(e.countAction.WithLabelValues(
			event.GetRepo().GetName(),
			event.GetAction(),
			event.GetCheckRun().GetStatus(),
			event.GetCheckRun().GetConclusion(),
			"",
		)), float64(0))
	})

	t.Run("Should log listener pull request", func(t *testing.T) {
		// given
		event := github.PullRequestEvent{
			Repo: &github.Repository{
				Name: github.String("github-observer"),
			},
			Action: github.String("opened"),
			PullRequest: &github.PullRequest{
				State: github.String("open"),
			},
		}

		// when
		err := e.PullRequestEvent(event)

		// then
		assert.NoError(t, err)
		assert.Equal(t, 1, testutil.CollectAndCount(e.lastRequestReceivedTimePR))
		assert.Greater(t, testutil.ToFloat64(e.lastRequestReceivedTimePR.WithLabelValues(
			event.GetRepo().GetName(),
			event.GetAction(),
			event.GetPullRequest().GetState(),
			"",
			"",
		)), float64(0))
		assert.Equal(t, 1, testutil.CollectAndCount(e.countPR))
		assert.Greater(t, testutil.ToFloat64(e.countPR.WithLabelValues(
			event.GetRepo().GetName(),
			event.GetAction(),
			event.GetPullRequest().GetState(),
			"",
			"",
		)), float64(0))
	})

	t.Run("Should log listener pull request review", func(t *testing.T) {
		// given
		event := github.PullRequestReviewEvent{
			Repo: &github.Repository{
				Name: github.String("github-observer"),
			},
			Action: github.String("submitted"),
			PullRequest: &github.PullRequest{
				State: github.String("open"),
			},
			Review: &github.PullRequestReview{
				State: github.String("commented"),
			},
		}

		// when
		err := e.PullRequestReviewEvent(event)

		// then
		assert.NoError(t, err)
		assert.Equal(t, 1, testutil.CollectAndCount(e.lastRequestReceivedTimePRR))
		assert.Greater(t, testutil.ToFloat64(e.lastRequestReceivedTimePRR.WithLabelValues(
			event.GetRepo().GetName(),
			event.GetAction(),
			event.GetPullRequest().GetState(),
			"",
			event.GetReview().GetState(),
		)), float64(0))
		assert.Equal(t, 1, testutil.CollectAndCount(e.countPRR))
		assert.Greater(t, testutil.ToFloat64(e.countPRR.WithLabelValues(
			event.GetRepo().GetName(),
			event.GetAction(),
			event.GetPullRequest().GetState(),
			"",
			event.GetReview().GetState(),
		)), float64(0))
	})
}
