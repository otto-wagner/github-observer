package watcher

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/go-github/v73/github"
	"github.com/otto-wagner/github-observer/internal/core"
	"github.com/otto-wagner/github-observer/internal/executor"
)

type IWatcher interface {
	Start()
	PullRequests(core.Repository)
	WorkflowRuns(core.Repository)
}

type watcher struct {
	ctx          context.Context
	client       *github.Client
	repositories []core.Repository
	workflows    map[core.Repository][]*github.Workflow
	executors    []executor.IExecutor
}

func NewWatcher(client *github.Client, repositories []core.Repository, executors []executor.IExecutor) IWatcher {
	return &watcher{
		ctx:          context.Background(),
		client:       client,
		repositories: repositories,
		executors:    executors,
		workflows:    make(map[core.Repository][]*github.Workflow),
	}
}

func (w *watcher) Start() {
	scheduler := gocron.NewScheduler(time.UTC)

	// Ratelimit: 5000 requests pro Stunde
	_, err := scheduler.Every(15).Minute().Do(func() {
		w.checkRateLimit()
	})
	if err != nil {
		slog.Error("check rate limit cron scheduler failed", "error", err)
	}

	for _, repository := range w.repositories {
		r := repository
		_, err = scheduler.Every(1).Hours().Do(func() {
			w.PullRequests(r)
		})
		if err != nil {
			slog.Error("pull requests cron scheduler failed", "error", err)
		}

		// Jeder Workflow Check verbraucht im Schnitt 20 requests.
		// Hängt davon ab, wie viele Workflows in einem Repository sind
		_, err = scheduler.Every(15).Minute().Do(func() {
			w.updateExistingWorkflows(r)
			w.WorkflowRuns(r)
		})
		if err != nil {
			slog.Error("workflow runs cron scheduler failed", "error", err)
		}
	}
	scheduler.StartAsync()
}

func (w *watcher) PullRequests(repository core.Repository) {
	pullRequests, _, err := w.client.PullRequests.List(w.ctx, repository.Owner, repository.Name, &github.PullRequestListOptions{})
	if err != nil {
		slog.Error("Failed to list pull requests", "error", err)
		return
	}
	for _, e := range w.executors {
		e.PullRequests(repository, pullRequests)
	}
}

func (w *watcher) WorkflowRuns(repository core.Repository) {
	var latestWorkflowRuns []*github.WorkflowRun
	for _, workflow := range w.workflows[repository] {
		lastRun, _, err := w.client.Actions.ListWorkflowRunsByID(w.ctx, repository.Owner, repository.Name, workflow.GetID(),
			&github.ListWorkflowRunsOptions{ListOptions: github.ListOptions{PerPage: 1}, Branch: repository.Branch, Event: "push"},
		)
		if err != nil {
			slog.Error("Failed to list workflow runs", "error", err)
			return
		}
		latestWorkflowRuns = append(latestWorkflowRuns, lastRun.WorkflowRuns...)
	}

	for _, e := range w.executors {
		e.LastWorkflows(repository, latestWorkflowRuns)
	}
}

func (w *watcher) checkRateLimit() {
	rateLimit, _, err := w.client.RateLimit.Get(w.ctx)
	if err != nil {
		slog.Error("Failed to get rate limit", "error", err)
		return
	}
	slog.Info("RateLimit", "Rate", rateLimit)
}

func (w *watcher) updateExistingWorkflows(repository core.Repository) {
	listedWorkFlows, _, err := w.client.Actions.ListWorkflows(context.Background(), repository.Owner, repository.Name, &github.ListOptions{})
	if err != nil {
		slog.Error("Failed to list workflows", "error", err)
		return
	}

	workflowNames := make([]string, 0, len(listedWorkFlows.Workflows))
	for _, wf := range listedWorkFlows.Workflows {
		workflowNames = append(workflowNames, wf.GetName())
	}
	slog.Info("Workflows",
		"repository", repository.Owner+"/"+repository.Name,
		"branch", repository.Branch,
		"count", listedWorkFlows.GetTotalCount(),
		"workflows", workflowNames,
	)
	w.workflows[repository] = listedWorkFlows.Workflows
}
