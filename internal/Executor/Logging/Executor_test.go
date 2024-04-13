//go:build all || unit

package Logging

import (
	logger "github-listener/pkg/mocks"
	"github.com/google/go-github/v61/github"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewExecutor(t *testing.T) {

	t.Run("Should return name", func(t *testing.T) {
		// given
		// when
		name := NewExecutor().Name()

		// then
		assert.Equal(t, "Logging", name)
	})

	t.Run("Should log listener action", func(t *testing.T) {
		// given
		logs := logger.MockedLogger()

		event := github.CheckRunEvent{
			Repo: &github.Repository{
				Name:    github.String("github-listener"),
				HTMLURL: github.String("https://github.com/otto-wagner/github-listener"),
			},
			Action: github.String("completed"),
			CheckRun: &github.CheckRun{
				Name:       github.String("Analyze (go)"),
				HTMLURL:    github.String("https://github.com/otto-wagner/github-listener/actions/runs/8589035842/job/23534635896"),
				Status:     github.String("completed"),
				Conclusion: github.String("success"),
			},
		}

		// when
		_ = NewExecutor().CheckRunEvent(event)

		// then
		assert.Contains(t, logs.All()[0].Message, "Workflow received")
		assert.Contains(t, logs.All()[0].Context[0].Key, "repo")
		assert.Contains(t, logs.All()[0].Context[0].String, "github-listener")
		assert.Contains(t, logs.All()[0].Context[1].Key, "repo_html_url")
		assert.Contains(t, logs.All()[0].Context[1].String, "https://github.com/otto-wagner/github-listener")
		assert.Contains(t, logs.All()[0].Context[2].Key, "name")
		assert.Contains(t, logs.All()[0].Context[2].String, "Analyze (go)")
		assert.Contains(t, logs.All()[0].Context[3].Key, "html_url")
		assert.Contains(t, logs.All()[0].Context[3].String, "https://github.com/otto-wagner/github-listener/actions/runs/8589035842/job/23534635896")
		assert.Contains(t, logs.All()[0].Context[4].Key, "action")
		assert.Contains(t, logs.All()[0].Context[4].String, "completed")
		assert.Contains(t, logs.All()[0].Context[5].Key, "status")
		assert.Contains(t, logs.All()[0].Context[5].String, "completed")
		assert.Contains(t, logs.All()[0].Context[6].Key, "conclusion")
		assert.Contains(t, logs.All()[0].Context[6].String, "success")
	})

	t.Run("Should log listener pull request", func(t *testing.T) {
		// given
		logs := logger.MockedLogger()
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
		assert.Contains(t, logs.All()[0].Message, "Workflow received")
		assert.Contains(t, logs.All()[0].Context[0].Key, "repo")
		assert.Contains(t, logs.All()[0].Context[0].String, "github-listener")
		assert.Contains(t, logs.All()[0].Context[1].Key, "repo_html_url")
		assert.Contains(t, logs.All()[0].Context[1].String, "https://github.com/otto-wagner/github-listener")
		assert.Contains(t, logs.All()[0].Context[2].Key, "title")
		assert.Contains(t, logs.All()[0].Context[2].String, "chore: test pullrequest_listener")
		assert.Contains(t, logs.All()[0].Context[3].Key, "user")
		assert.Contains(t, logs.All()[0].Context[3].String, "otto-wagner")
		assert.Contains(t, logs.All()[0].Context[4].Key, "html_url")
		assert.Contains(t, logs.All()[0].Context[4].String, "https://github.com/otto-wagner/github-listener/pull/2")
		assert.Contains(t, logs.All()[0].Context[5].Key, "action")
		assert.Contains(t, logs.All()[0].Context[5].String, "opened")
		assert.Contains(t, logs.All()[0].Context[6].Key, "status")
		assert.Contains(t, logs.All()[0].Context[6].String, "open")
	})

	t.Run("Should log listener pull request review", func(t *testing.T) {
		// given
		logs := logger.MockedLogger()
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
		assert.Contains(t, logs.All()[0].Message, "Workflow received")
		assert.Contains(t, logs.All()[0].Context[0].Key, "repo")
		assert.Contains(t, logs.All()[0].Context[0].String, "github-listener")
		assert.Contains(t, logs.All()[0].Context[1].Key, "repo_html_url")
		assert.Contains(t, logs.All()[0].Context[1].String, "https://github.com/otto-wagner/github-listener")
		assert.Contains(t, logs.All()[0].Context[2].Key, "title")
		assert.Contains(t, logs.All()[0].Context[2].String, "chore: test pullrequest_listener")
		assert.Contains(t, logs.All()[0].Context[3].Key, "user")
		assert.Contains(t, logs.All()[0].Context[3].String, "otto-wagner")
		assert.Contains(t, logs.All()[0].Context[4].Key, "html_url")
		assert.Contains(t, logs.All()[0].Context[4].String, "https://github.com/otto-wagner/github-listener/pull/2")
		assert.Contains(t, logs.All()[0].Context[5].Key, "action")
		assert.Contains(t, logs.All()[0].Context[5].String, "submitted")
		assert.Contains(t, logs.All()[0].Context[6].Key, "status")
		assert.Contains(t, logs.All()[0].Context[6].String, "open")
		assert.Contains(t, logs.All()[0].Context[7].Key, "review")
		assert.Contains(t, logs.All()[0].Context[7].String, "LGTM")
		assert.Contains(t, logs.All()[0].Context[8].Key, "state")
		assert.Contains(t, logs.All()[0].Context[8].String, "commented")
		assert.Contains(t, logs.All()[0].Context[9].Key, "reviewer")
		assert.Contains(t, logs.All()[0].Context[9].String, "otto-wagner")
	})
}
