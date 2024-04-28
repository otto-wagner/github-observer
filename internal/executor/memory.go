package executor

import "github-observer/internal/core"

type IMemory interface {
	Store(core.WorkflowRun)
	Get(core.WorkflowRun) (workflow core.WorkflowRun, found bool)
	StorePR(core.GitPullRequest)
	GetPR(core.GitPullRequest) (pr core.GitPullRequest, found bool)
}

type memory struct {
	mem   map[core.WorkflowRun]core.WorkflowRun
	memPr map[core.GitPullRequest]core.GitPullRequest
}

func NewMemory() IMemory {
	return &memory{make(map[core.WorkflowRun]core.WorkflowRun), make(map[core.GitPullRequest]core.GitPullRequest)}
}

func (m *memory) Store(workflow core.WorkflowRun) {
	m.mem[core.WorkflowRun{WorkflowId: workflow.WorkflowId, RunNumber: workflow.RunNumber, HeadBranch: workflow.HeadBranch}] = workflow
}

func (m *memory) Get(workflow core.WorkflowRun) (flow core.WorkflowRun, found bool) {
	flow, found = m.mem[core.WorkflowRun{WorkflowId: workflow.WorkflowId, RunNumber: workflow.RunNumber, HeadBranch: workflow.HeadBranch}]
	return
}

func (m *memory) StorePR(pr core.GitPullRequest) {
	m.memPr[core.GitPullRequest{PullRequest: core.PullRequest{Number: pr.PullRequest.Number}}] = pr
}

func (m *memory) GetPR(pr core.GitPullRequest) (pull core.GitPullRequest, found bool) {
	pull, found = m.memPr[core.GitPullRequest{PullRequest: core.PullRequest{Number: pr.PullRequest.Number}}]
	return
}
