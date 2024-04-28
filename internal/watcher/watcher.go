package watcher

import (
	"context"
	"github-observer/internal/core"
	"github-observer/internal/executor"
	"github.com/go-co-op/gocron"
	"github.com/google/go-github/v61/github"
	"go.uber.org/zap"
	"time"
)

type watcher struct {
	client    *github.Client
	executors []executor.IExecutor
}

func Watch(client *github.Client, repositories []core.Repository, executors []executor.IExecutor) {
	w := &watcher{client, executors}

	for _, repository := range repositories {
		scheduler := gocron.NewScheduler(time.UTC)
		scheduler.SetMaxConcurrentJobs(1, gocron.WaitMode)
		_, err := scheduler.Every(1).Minute().Do(func() {
			go w.Actions(repository)
			go w.PullRequests(repository)
		})
		if err != nil {
			zap.S().Fatalw("cron scheduler failed", "error", err)
		}
		scheduler.StartAsync()
	}

	return
}

func (w *watcher) Actions(repository core.Repository) {
	actions, r, err := w.client.Actions.ListRepositoryWorkflowRuns(context.Background(), repository.Owner, repository.Name, nil)
	if err != nil {
		zap.S().Errorw("Failed to list workflow runs", "error", err)
	}
	if r.StatusCode > 299 {
		zap.S().Errorw("Failed to list workflow runs", "status_code", r.StatusCode)
	}
	for _, e := range w.executors {
		go e.WorkflowRuns(actions.WorkflowRuns)
	}
}

func (w *watcher) PullRequests(repository core.Repository) {
	pullRequests, r, err := w.client.PullRequests.List(context.Background(), repository.Owner, repository.Name, nil)
	if err != nil {
		zap.S().Errorw("Failed to list pull requests", "error", err)
	}
	if r.StatusCode > 299 {
		zap.S().Errorw("Failed to list pull requests", "status_code", r.StatusCode)
	}
	for _, e := range w.executors {
		go e.PullRequests(repository, pullRequests)
	}
}
