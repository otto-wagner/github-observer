// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import (
	gin "github.com/gin-gonic/gin"

	mock "github.com/stretchr/testify/mock"
)

// IListener is an autogenerated mock type for the IListener type
type IListener struct {
	mock.Mock
}

// Action provides a mock function with given fields: _a0
func (_m *IListener) Action(_a0 *gin.Context) {
	_m.Called(_a0)
}

// PullRequest provides a mock function with given fields: _a0
func (_m *IListener) PullRequest(_a0 *gin.Context) {
	_m.Called(_a0)
}

// NewIListener creates a new instance of IListener. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIListener(t interface {
	mock.TestingT
	Cleanup(func())
}) *IListener {
	mock := &IListener{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
