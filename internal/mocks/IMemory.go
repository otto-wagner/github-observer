// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	core "github-observer/internal/core"

	mock "github.com/stretchr/testify/mock"
)

// IMemory is an autogenerated mock type for the IMemory type
type IMemory struct {
	mock.Mock
}

// GetLastWorkflowRun provides a mock function with given fields: _a0
func (_m *IMemory) GetLastWorkflowRun(_a0 core.WorkflowRun) (core.WorkflowRun, bool) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetLastWorkflowRun")
	}

	var r0 core.WorkflowRun
	var r1 bool
	if rf, ok := ret.Get(0).(func(core.WorkflowRun) (core.WorkflowRun, bool)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(core.WorkflowRun) core.WorkflowRun); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(core.WorkflowRun)
	}

	if rf, ok := ret.Get(1).(func(core.WorkflowRun) bool); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetPullRequest provides a mock function with given fields: repositoryName, pr
func (_m *IMemory) GetPullRequest(repositoryName string, pr core.GitPullRequest) (core.GitPullRequest, bool) {
	ret := _m.Called(repositoryName, pr)

	if len(ret) == 0 {
		panic("no return value specified for GetPullRequest")
	}

	var r0 core.GitPullRequest
	var r1 bool
	if rf, ok := ret.Get(0).(func(string, core.GitPullRequest) (core.GitPullRequest, bool)); ok {
		return rf(repositoryName, pr)
	}
	if rf, ok := ret.Get(0).(func(string, core.GitPullRequest) core.GitPullRequest); ok {
		r0 = rf(repositoryName, pr)
	} else {
		r0 = ret.Get(0).(core.GitPullRequest)
	}

	if rf, ok := ret.Get(1).(func(string, core.GitPullRequest) bool); ok {
		r1 = rf(repositoryName, pr)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// StoreLastRepositoryWorkflow provides a mock function with given fields: workflow
func (_m *IMemory) StoreLastRepositoryWorkflow(workflow core.WorkflowRun) error {
	ret := _m.Called(workflow)

	if len(ret) == 0 {
		panic("no return value specified for StoreLastRepositoryWorkflow")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(core.WorkflowRun) error); ok {
		r0 = rf(workflow)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StorePullRequests provides a mock function with given fields: repositoryName, pr
func (_m *IMemory) StorePullRequests(repositoryName string, pr []core.GitPullRequest) {
	_m.Called(repositoryName, pr)
}

// NewIMemory creates a new instance of IMemory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIMemory(t interface {
	mock.TestingT
	Cleanup(func())
}) *IMemory {
	mock := &IMemory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
