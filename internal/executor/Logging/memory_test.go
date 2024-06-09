//go:build all || unit

package Logging

import (
	"github-observer/internal/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWorkflow(t *testing.T) {

	t.Run("Should store and return workflow", func(t *testing.T) {
		// given
		run := core.WorkflowRun{WorkflowId: 1, Repository: core.Repository{FullName: "github-observer"}}
		m := NewMemory()

		// when
		_ = m.StoreLastRepositoryWorkflow(run)
		workflow, _ := m.GetLastWorkflowRun(run)

		// then
		assert.Equal(t, run, workflow)
	})

	t.Run("Should validate workflow", func(t *testing.T) {
		// given
		m := NewMemory()

		// when
		err := m.StoreLastRepositoryWorkflow(core.WorkflowRun{WorkflowId: 1})

		// then
		assert.EqualError(t, err, "repository name is required")

		// when
		err = m.StoreLastRepositoryWorkflow(core.WorkflowRun{Repository: core.Repository{FullName: "github-observer"}})

		// then
		assert.EqualError(t, err, "workflow id is required")
	})

	t.Run("Should store and return last workflow", func(t *testing.T) {
		// given
		workflowRun := core.WorkflowRun{WorkflowId: 1, Repository: core.Repository{FullName: "github-observer"}}
		first := workflowRun
		first.RunNumber = 1

		second := workflowRun
		second.RunNumber = 2

		third := workflowRun
		third.RunNumber = 3

		m := NewMemory()

		// when
		err := m.StoreLastRepositoryWorkflow(second)
		lastRun, _ := m.GetLastWorkflowRun(workflowRun)

		// then
		assert.NoError(t, err)
		assert.Equal(t, second, lastRun)

		// when
		err = m.StoreLastRepositoryWorkflow(first)

		// then
		assert.EqualError(t, err, "workflow run number is lower than the last one")

		// when
		err = m.StoreLastRepositoryWorkflow(third)
		lastRun, _ = m.GetLastWorkflowRun(workflowRun)

		// then
		assert.NoError(t, err)
		assert.Equal(t, third, lastRun)
	})

}

func TestPullRequest(t *testing.T) {

	t.Run("Should store and return pullrequests", func(t *testing.T) {
		// given
		pr := core.GitPullRequest{Repository: core.Repository{FullName: "github-observer"}}
		m := NewMemory()

		// when
		m.StorePullRequests(pr.Repository.FullName, []core.GitPullRequest{pr})
		pullRequest, _ := m.GetPullRequest(pr.Repository.FullName, pr)

		// then
		assert.Equal(t, pr, pullRequest)
	})

	t.Run("Should delete pullrequests", func(t *testing.T) {
		// given
		pr := core.GitPullRequest{Repository: core.Repository{FullName: "github-observer"}}
		m := NewMemory()

		// when
		m.StorePullRequests(pr.Repository.FullName, []core.GitPullRequest{pr})
		m.StorePullRequests(pr.Repository.FullName, []core.GitPullRequest{})
		_, exists := m.GetPullRequest(pr.Repository.FullName, pr)

		// then
		assert.False(t, exists)
	})

}
