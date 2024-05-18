package Logging

import (
	"github-observer/internal/core"
	e "github-observer/internal/executor"
	"github.com/google/go-github/v61/github"
	"go.uber.org/zap"
)

type executor struct {
	memory IMemory
}

func NewExecutor(m IMemory) e.IExecutor {
	return &executor{m}
}

func (e *executor) EventRun(runEvent github.CheckRunEvent) {
	zap.S().Infow("Event", "Run", core.ConvertToGitAction(runEvent))
	return
}

func (e *executor) EventPullRequest(event github.PullRequestEvent) {
	zap.S().Infow("Event", "PullRequest", core.ConvertPREToGitPullRequest(event))
	return
}

func (e *executor) EventPullRequestReview(event github.PullRequestReviewEvent) {
	zap.S().Infow("Event", "PullRequestReview", core.ConvertToGitPullRequestReview(event))
	return
}

func (e *executor) LastWorkflows(repository core.Repository, runs []*github.WorkflowRun) {
	for _, run := range runs {
		workflow := core.ConvertToWorkflow(*run)

		memWorkflow, exists := e.memory.GetLastRepositoryWorkflow(repository.FullName)
		if !exists {
			e.memory.StoreLastRepositoryWorkflow(repository.FullName, workflow)
			if workflow.Conclusion != "success" {
				zap.S().Infow("WorkflowRun", "Action", workflow)
			}
			continue
		}
		if workflow.WorkflowId != memWorkflow.WorkflowId && workflow.RunNumber > memWorkflow.RunNumber {
			e.memory.StoreLastRepositoryWorkflow(repository.FullName, workflow)
			if workflow.Conclusion != "success" {
				zap.S().Infow("WorkflowRun", "Action", workflow)
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
			zap.S().Infow("PullRequest", "PullRequest", pr)
		}
	}

	e.memory.StorePullRequests(repository.Name, pullRequests)
	return
}
