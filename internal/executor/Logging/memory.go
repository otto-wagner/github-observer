package Logging

import (
	"errors"
	"github-observer/internal/core"
	"sync"
)

type IMemory interface {
	StoreLastRepositoryWorkflow(workflow core.WorkflowRun) error
	GetLastWorkflowRun(core.WorkflowRun) (core.WorkflowRun, bool)
	StorePullRequests(repositoryName string, pr []core.GitPullRequest)
	GetPullRequest(repositoryName string, pr core.GitPullRequest) (core.GitPullRequest, bool)
}

type memory struct {
	workflow       map[string]map[int64]core.WorkflowRun
	muWorkflow     sync.RWMutex
	pullRequests   map[string][]core.GitPullRequest
	muPullRequests sync.RWMutex
}

func NewMemory() IMemory {
	return &memory{workflow: make(map[string]map[int64]core.WorkflowRun), pullRequests: make(map[string][]core.GitPullRequest)}
}

func (m *memory) StoreLastRepositoryWorkflow(w core.WorkflowRun) error {
	m.muWorkflow.Lock()
	defer m.muWorkflow.Unlock()

	if w.WorkflowId == 0 {
		return errors.New("workflow id is required")
	}

	if len(w.Repository.FullName) == 0 {
		return errors.New("repository name is required")
	}

	lastWorkflow, found := m.workflow[w.Repository.FullName][w.WorkflowId]
	if found {
		if lastWorkflow.RunNumber > w.RunNumber {
			return errors.New("workflow run number is lower than the last one")
		}
	}

	if _, ok := m.workflow[w.Repository.FullName]; !ok {
		m.workflow[w.Repository.FullName] = make(map[int64]core.WorkflowRun)
		m.workflow[w.Repository.FullName][w.WorkflowId] = w
		return nil
	}

	m.workflow[w.Repository.FullName][w.WorkflowId] = w
	return nil
}

func (m *memory) GetLastWorkflowRun(w core.WorkflowRun) (workflow core.WorkflowRun, found bool) {
	m.muWorkflow.Lock()
	defer m.muWorkflow.Unlock()

	workflow, found = m.workflow[w.Repository.FullName][w.WorkflowId]
	return
}

func (m *memory) StorePullRequests(repositoryName string, openPullRequests []core.GitPullRequest) {
	m.muPullRequests.Lock()
	defer m.muPullRequests.Unlock()

	if len(openPullRequests) == 0 {
		delete(m.pullRequests, repositoryName)
		return
	}
	m.pullRequests[repositoryName] = openPullRequests
	return
}

func (m *memory) GetPullRequest(repositoryName string, pr core.GitPullRequest) (core.GitPullRequest, bool) {
	m.muPullRequests.Lock()
	defer m.muPullRequests.Unlock()

	for _, p := range m.pullRequests[repositoryName] {
		if p == pr {
			return p, true
		}
	}

	return core.GitPullRequest{}, false
}
