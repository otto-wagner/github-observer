package watcher

import (
	"context"
	"github-observer/internal/core"
	"github-observer/internal/executor"
	"github.com/go-co-op/gocron"
	"github.com/google/go-github/v61/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"time"
)

type IWatcher interface {
	Watch()
	CheckRateLimit()
	PullRequests(core.Repository)
	WorkflowRuns(core.Repository)
}

type watcher struct {
	token        string
	client       *github.Client
	repositories []core.Repository
	workflows    map[core.Repository][]*github.Workflow
	executors    []executor.IExecutor
}

func NewWatcher(token string, client *github.Client, repositories []core.Repository, executors []executor.IExecutor) IWatcher {
	return &watcher{token: token, client: client, repositories: repositories, executors: executors, workflows: make(map[core.Repository][]*github.Workflow)}
}

func (w *watcher) Watch() {
	scheduler := gocron.NewScheduler(time.UTC)

	_, err := scheduler.Every(30).Minute().Do(func() {
		w.client = github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: w.token})))
	})
	if err != nil {
		zap.S().Errorw("renew client cron scheduler failed", "error", err)
	}

	// Ratelimit:  5000 requests pro Stunde
	_, err = scheduler.Every(15).Minute().Do(func() {
		w.CheckRateLimit()
	})
	if err != nil {
		zap.S().Errorw("check rate limit cron scheduler failed", "error", err)
	}

	for _, repository := range w.repositories {
		r := repository
		_, err = scheduler.Every(1).Hours().Do(func() {
			w.PullRequests(r)
		})
		if err != nil {
			zap.S().Errorw("pull requests cron scheduler failed", "error", err)
		}

		// Jeder Workflow Check verbraucht im Schnitt 20 requests.
		// Hängt davon ab, wie viele Workflows in einem Repository sind
		_, err = scheduler.Every(15).Minute().Do(func() {
			w.updateExistingWorkflows(r)
			w.WorkflowRuns(r)
		})
		if err != nil {
			zap.S().Errorw("workflow runs cron scheduler failed", "error", err)
		}
	}
	scheduler.StartAsync()
}

func (w *watcher) CheckRateLimit() {
	rateLimit, _, err := w.client.RateLimit.Get(context.Background())
	if err != nil {
		zap.S().Errorw("Failed to get rate limit", "error", err)
		return
	}
	zap.S().Infow("RateLimit", "Rate", rateLimit)
}

func (w *watcher) PullRequests(repository core.Repository) {
	pullRequests, _, err := w.client.PullRequests.List(context.Background(), repository.Owner, repository.Name, &github.PullRequestListOptions{})
	if err != nil {
		zap.S().Errorw("Failed to list pull requests", "error", err)
		return
	}
	for _, e := range w.executors {
		e.PullRequests(repository, pullRequests)
	}
}

func (w *watcher) updateExistingWorkflows(repository core.Repository) {
	workflows, err := w.listWorkflows(repository)
	if err != nil {
		zap.S().Errorw("Failed to list workflows", "error", err)
		return
	}
	zap.S().Infow("Workflows", "Repository", repository, "Workflows", workflows)
	w.workflows[repository] = workflows
}

func (w *watcher) listWorkflows(repository core.Repository) ([]*github.Workflow, error) {
	workflows, _, err := w.client.Actions.ListWorkflows(context.Background(), repository.Owner, repository.Name, &github.ListOptions{})
	if err != nil {
		return nil, err
	}
	return workflows.Workflows, nil
}

func (w *watcher) WorkflowRuns(repository core.Repository) {
	latestWorkflowRuns, err := w.getLatestWorkflowRuns(repository)
	if err != nil {
		zap.S().Errorw("Failed to get latest workflow runs", "error", err)
		return
	}

	for _, e := range w.executors {
		e.LastWorkflows(repository, latestWorkflowRuns)
	}
}

func (w *watcher) getLatestWorkflowRuns(repository core.Repository) (latestWorkflowRuns []*github.WorkflowRun, err error) {
	for _, workflow := range w.workflows[repository] {
		lastRun, _, err := w.client.Actions.ListWorkflowRunsByID(context.Background(), repository.Owner, repository.Name, workflow.GetID(),
			&github.ListWorkflowRunsOptions{ListOptions: github.ListOptions{PerPage: 1}, Branch: repository.Branch, Event: "push"},
		)
		if err != nil {
			return nil, err
		}
		latestWorkflowRuns = append(latestWorkflowRuns, lastRun.WorkflowRuns...)
	}
	return
}
