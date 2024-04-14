//go:build all || unit

package Prometheus

import (
	"github.com/google/go-github/v61/github"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewExecutor(t *testing.T) {

	t.Run("Should return name", func(t *testing.T) {
		// given
		// when
		name := NewExecutor().Name()

		// then
		assert.Equal(t, "Prometheus", name)
	})

	t.Run("Should log listener action", func(t *testing.T) {
		// given
		event := github.CheckRunEvent{
			Repo: &github.Repository{
				Name:    github.String("github-listener"),
				HTMLURL: github.String("https://github.com/otto-wagner/github-listener"),
			},
			Action: github.String("completed"),
			CheckRun: &github.CheckRun{
				ID:         github.Int64(1),
				Name:       github.String("Analyze (go)"),
				HTMLURL:    github.String("https://github.com/otto-wagner/github-listener/actions/runs/8589035842/job/23534635896"),
				Status:     github.String("completed"),
				Conclusion: github.String("success"),
			},
		}

		executor := NewExecutor().(*executor)

		// when
		err := executor.CheckRunEvent(event)

		// then
		assert.NoError(t, err)
		assert.Equal(t, 1, testutil.CollectAndCount(executor.lastRequestReceivedTimeAction))
		assert.Greater(t, testutil.ToFloat64(executor.lastRequestReceivedTimeAction.WithLabelValues(
			event.GetRepo().GetName(),
			event.GetAction(),
			event.GetCheckRun().GetStatus(),
			event.GetCheckRun().GetConclusion(),
			"",
		)), float64(0))
		assert.Equal(t, 1, testutil.CollectAndCount(executor.countAction))
		assert.Greater(t, testutil.ToFloat64(executor.countAction.WithLabelValues(
			event.GetRepo().GetName(),
			event.GetAction(),
			event.GetCheckRun().GetStatus(),
			event.GetCheckRun().GetConclusion(),
		)), float64(0))
	})

	t.Run("Should log listener pull request", func(t *testing.T) {
		// given
		event := github.PullRequestEvent{
			Repo: &github.Repository{
				Name:    github.String("github-listener"),
				HTMLURL: github.String("https://github.com/otto-wagner/github-listener"),
			},
			Action: github.String("opened"),
			PullRequest: &github.PullRequest{
				Title:   github.String("chore: test pullrequest_listener"),
				User:    &github.User{Login: github.String("otto-wagner")},
				HTMLURL: github.String("https://github.com/otto-wagner/github-listener/pull/2"),
				State:   github.String("open"),
			},
		}

		// when
		_ = NewExecutor().PullRequestEvent(event)

		// then

	})

	t.Run("Should log listener pull request review", func(t *testing.T) {
		// given
		event := github.PullRequestReviewEvent{
			Repo: &github.Repository{
				Name:    github.String("github-listener"),
				HTMLURL: github.String("https://github.com/otto-wagner/github-listener"),
			},
			Action: github.String("submitted"),
			PullRequest: &github.PullRequest{
				Title:   github.String("chore: test pullrequest_listener"),
				User:    &github.User{Login: github.String("otto-wagner")},
				HTMLURL: github.String("https://github.com/otto-wagner/github-listener/pull/2"),
				State:   github.String("open"),
			},
			Review: &github.PullRequestReview{
				Body:  github.String("LGTM"),
				State: github.String("commented"),
				User:  &github.User{Login: github.String("otto-wagner")},
			},
		}

		// when
		_ = NewExecutor().PullRequestReviewEvent(event)

		// then

	})
}
