package Logging

import (
	"github-observer/internal/core"
	e "github-observer/internal/executor"
	"github.com/google/go-github/v61/github"
	"log/slog"
)

type executor struct {
	memory IMemory
	logger *slog.Logger
}

func NewExecutor(m IMemory, logger *slog.Logger) e.IExecutor {
	return &executor{m, logger}
}

func (e *executor) EventPullRequest(event github.PullRequestEvent) {
	e.logger.Info("Event", "PullRequest", core.ConvertPREToGitPullRequest(event))
	return
}

func (e *executor) EventPullRequestReview(event github.PullRequestReviewEvent) {
	e.logger.Info("Event", "PullRequestReview", core.ConvertToGitPullRequestReview(event))
	return
}

func (e *executor) EventWorkflowRun(event github.WorkflowRunEvent) {
	e.logger.Info("Event", "WorkflowRun", core.ConvertToWorkflowRun(event))
	return
}

func (e *executor) LastWorkflows(_ core.Repository, workflowRuns []*github.WorkflowRun) {
	for _, run := range workflowRuns {
		workflow := core.ConvertToWorkflow(*run)

		memWorkflow, exists := e.memory.GetLastWorkflowRun(workflow)
		if !exists {
			err := e.memory.StoreLastRepositoryWorkflow(workflow)
			if err != nil {
				e.logger.Error("WorkflowRun", "Action", workflow, "Error", err)
				continue
			}
			if workflow.Conclusion != "success" {
				e.logger.Info("WorkflowRun", "Action", workflow)
			}
			continue
		}
		if workflow.RunNumber > memWorkflow.RunNumber {
			err := e.memory.StoreLastRepositoryWorkflow(workflow)
			if err != nil {
				e.logger.Error("WorkflowRun", "Action", workflow, "Error", err)
				continue
			}
			if workflow.Conclusion != "success" {
				e.logger.Info("WorkflowRun", "Action", workflow)
			}
			continue
		}
	}
	return
}

func (e *executor) PullRequests(repository core.Repository, openPullRequests []*github.PullRequest) {
	var pullRequests []core.GitPullRequest

	for _, request := range openPullRequests {
		pullRequests = append(pullRequests, core.ConvertPRToGitPullRequest(*request))
	}

	for _, pr := range pullRequests {
		if pr.PullRequest.State == "closed" || pr.PullRequest.State == "merged" {
			continue
		}

		_, exists := e.memory.GetPullRequest(repository.Name, pr)
		if !exists {
			e.logger.Info("PullRequest", "PullRequest", pr)
		}
	}

	e.memory.StorePullRequests(repository.Name, pullRequests)
	return
}
