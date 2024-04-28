package Logging

import (
	"github-observer/internal/core"
	e "github-observer/internal/executor"
	"github.com/google/go-github/v61/github"
	"go.uber.org/zap"
)

type executor struct {
	memory e.IMemory
}

func NewExecutor(m e.IMemory) e.IExecutor {
	return &executor{m}
}

func (e *executor) Name() string {
	return "Logging"
}

func (e *executor) RunEvent(runEvent github.CheckRunEvent) {
	zap.S().Infow("Event", "Run", core.ConvertToGitAction(runEvent))
	return
}

func (e *executor) PullRequestEvent(event github.PullRequestEvent) {
	zap.S().Infow("Event", "PullRequest", core.ConvertPREToGitPullRequest(event))
	return
}

func (e *executor) PullRequestReviewEvent(event github.PullRequestReviewEvent) {
	zap.S().Infow("Event", "PullRequestReview", core.ConvertToGitPullRequestReview(event))
	return
}

func (e *executor) WorkflowRuns(runs []*github.WorkflowRun) {
	for _, run := range runs {
		workflow := core.ConvertToWorkflow(*run)

		memWorkflow, found := e.memory.Get(workflow)
		if !found || (workflow.WorkflowId == memWorkflow.WorkflowId && workflow.RunNumber > memWorkflow.RunNumber) {
			e.memory.Store(workflow)
			if workflow.Conclusion != "success" {
				zap.S().Infow("WorkflowRun", "Action", core.ConvertToWorkflow(*run))
			}
		}
	}
	return
}

func (e *executor) PullRequests(repository core.Repository, pullRequests []*github.PullRequest) {
	for _, pr := range pullRequests {
		pullrequest := core.ConvertPRToGitPullRequest(repository, *pr)
		if pullrequest.PullRequest.State == "closed" || pullrequest.PullRequest.State == "merged" {
			continue
		}

		_, found := e.memory.GetPR(pullrequest)
		if !found {
			e.memory.StorePR(pullrequest)
			zap.S().Infow("PullRequest", "PullRequest", pullrequest)
		}
	}
	return
}
