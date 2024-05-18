package watcher

import (
	"context"
	"github-observer/internal/core"
	"github-observer/internal/executor"
	"github.com/go-co-op/gocron"
	"github.com/google/go-github/v61/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"strconv"
	"time"
)

type IWatcher interface {
	Watch()
	PullRequests(core.Repository)
	WorkflowRuns(core.Repository)
}

type watcher struct {
	token        string
	client       *github.Client
	repositories []core.Repository
	executors    []executor.IExecutor
}

func NewWatcher(token string, repositories []core.Repository, executors []executor.IExecutor) IWatcher {
	client := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})))

	return &watcher{token, client, repositories, executors}
}

func (w *watcher) Watch() {
	scheduler := gocron.NewScheduler(time.UTC)

	_, err := scheduler.Every(3).Minute().Do(func() {
		w.client = github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: w.token})))
	})
	if err != nil {
		zap.S().Errorw("renew client cron scheduler failed", "error", err)
	}

	for _, repository := range w.repositories {
		r := repository
		_, err = scheduler.Every(1).Minute().Do(func() {
			w.PullRequests(r)
		})
		if err != nil {
			zap.S().Errorw("pull requests cron scheduler failed", "error", err)
		}

		_, err = scheduler.Every(1).Minute().Do(func() {
			w.WorkflowRuns(r)
		})
		if err != nil {
			zap.S().Errorw("workflow runs cron scheduler failed", "error", err)
		}
	}
	scheduler.StartAsync()
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

func (w *watcher) WorkflowRuns(repository core.Repository) {
	workflows, err := w.listWorkflows(repository)
	if err != nil {
		zap.S().Errorw("Failed to list workflows", "error", err)
		return
	}

	latestWorkflowRuns, err := w.getLatestWorkflowRuns(repository, workflows)
	if err != nil {
		zap.S().Errorw("Failed to get latest workflow runs", "error", err)
		return
	}

	for _, e := range w.executors {
		e.LastWorkflows(repository, latestWorkflowRuns)
	}
}

func (w *watcher) listWorkflows(repository core.Repository) ([]*github.Workflow, error) {
	workflows, _, err := w.client.Actions.ListWorkflows(context.Background(), repository.Owner, repository.Name, &github.ListOptions{})
	if err != nil {
		return nil, err
	}
	return workflows.Workflows, nil
}

func (w *watcher) getLatestWorkflowRuns(repository core.Repository, workflows []*github.Workflow) (latestWorkflowRuns []*github.WorkflowRun, err error) {
	for _, workflow := range workflows {
		lastRun, _, err := w.client.Actions.ListWorkflowRunsByFileName(context.Background(), repository.Owner, repository.Name, strconv.FormatInt(workflow.GetID(), 10), &github.ListWorkflowRunsOptions{
			ListOptions: github.ListOptions{PerPage: 1}, Branch: repository.Branch,
		})
		if err != nil {
			return nil, err
		}
		latestWorkflowRuns = append(latestWorkflowRuns, lastRun.WorkflowRuns...)
	}
	return
}
