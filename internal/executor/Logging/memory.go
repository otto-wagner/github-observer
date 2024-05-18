package Logging

import (
	"github-observer/internal/core"
	"sync"
)

type IMemory interface {
	StoreLastRepositoryWorkflow(repositoryName string, workflow core.WorkflowRun)
	GetLastRepositoryWorkflow(repositoryName string) (core.WorkflowRun, bool)
	StorePullRequests(repositoryName string, pr []core.GitPullRequest)
	GetPullRequest(repositoryName string, pr core.GitPullRequest) (core.GitPullRequest, bool)
}

type memory struct {
	lastWorkflow   map[string]core.WorkflowRun
	muWorkflow     sync.RWMutex
	pullRequests   map[string][]core.GitPullRequest
	muPullRequests sync.RWMutex
}

func NewMemory() IMemory {
	return &memory{lastWorkflow: make(map[string]core.WorkflowRun), pullRequests: make(map[string][]core.GitPullRequest)}
}

func (m *memory) StoreLastRepositoryWorkflow(repositoryName string, workflow core.WorkflowRun) {
	m.muWorkflow.Lock()
	defer m.muWorkflow.Unlock()

	m.lastWorkflow[repositoryName] = workflow
	return
}

func (m *memory) GetLastRepositoryWorkflow(repositoryName string) (workflow core.WorkflowRun, found bool) {
	m.muWorkflow.Lock()
	defer m.muWorkflow.Unlock()

	workflow, found = m.lastWorkflow[repositoryName]
	return
}

func (m *memory) StorePullRequests(repositoryName string, openPullRequests []core.GitPullRequest) {
	m.muPullRequests.Lock()
	defer m.muPullRequests.Unlock()

	m.pullRequests[repositoryName] = openPullRequests
	return
}

func (m *memory) GetPullRequest(repositoryName string, pr core.GitPullRequest) (pull core.GitPullRequest, found bool) {
	m.muPullRequests.Lock()
	defer m.muPullRequests.Unlock()

	for _, p := range m.pullRequests[repositoryName] {
		if p == pr {
			return p, true
		}
	}
	return
}
