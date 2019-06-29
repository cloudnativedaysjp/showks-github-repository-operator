// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/gh/gh.go

// Package mock_gh is a generated GoMock package.
package mock_gh

import (
	gomock "github.com/golang/mock/gomock"
	github "github.com/google/go-github/github"
	reflect "reflect"
)

// MockGitHubClientInterface is a mock of GitHubClientInterface interface
type MockGitHubClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockGitHubClientInterfaceMockRecorder
}

// MockGitHubClientInterfaceMockRecorder is the mock recorder for MockGitHubClientInterface
type MockGitHubClientInterfaceMockRecorder struct {
	mock *MockGitHubClientInterface
}

// NewMockGitHubClientInterface creates a new mock instance
func NewMockGitHubClientInterface(ctrl *gomock.Controller) *MockGitHubClientInterface {
	mock := &MockGitHubClientInterface{ctrl: ctrl}
	mock.recorder = &MockGitHubClientInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGitHubClientInterface) EXPECT() *MockGitHubClientInterfaceMockRecorder {
	return m.recorder
}

// CreateRepository mocks base method
func (m *MockGitHubClientInterface) CreateRepository(org string, repo *github.Repository) (*github.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRepository", org, repo)
	ret0, _ := ret[0].(*github.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRepository indicates an expected call of CreateRepository
func (mr *MockGitHubClientInterfaceMockRecorder) CreateRepository(org, repo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRepository", reflect.TypeOf((*MockGitHubClientInterface)(nil).CreateRepository), org, repo)
}

// GetRepository mocks base method
func (m *MockGitHubClientInterface) GetRepository(org, repo string) (*github.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepository", org, repo)
	ret0, _ := ret[0].(*github.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepository indicates an expected call of GetRepository
func (mr *MockGitHubClientInterfaceMockRecorder) GetRepository(org, repo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepository", reflect.TypeOf((*MockGitHubClientInterface)(nil).GetRepository), org, repo)
}

// AddCollaborator mocks base method
func (m *MockGitHubClientInterface) AddCollaborator(owner, repo, user, permission string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCollaborator", owner, repo, user, permission)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCollaborator indicates an expected call of AddCollaborator
func (mr *MockGitHubClientInterfaceMockRecorder) AddCollaborator(owner, repo, user, permission interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCollaborator", reflect.TypeOf((*MockGitHubClientInterface)(nil).AddCollaborator), owner, repo, user, permission)
}

// GetPermissionLevel mocks base method
func (m *MockGitHubClientInterface) GetPermissionLevel(owner, repo, user string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPermissionLevel", owner, repo, user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPermissionLevel indicates an expected call of GetPermissionLevel
func (mr *MockGitHubClientInterfaceMockRecorder) GetPermissionLevel(owner, repo, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPermissionLevel", reflect.TypeOf((*MockGitHubClientInterface)(nil).GetPermissionLevel), owner, repo, user)
}

// UpdateBranchProtection mocks base method
func (m *MockGitHubClientInterface) UpdateBranchProtection(owner, repo, branch string, request *github.ProtectionRequest) (*github.Protection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBranchProtection", owner, repo, branch, request)
	ret0, _ := ret[0].(*github.Protection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateBranchProtection indicates an expected call of UpdateBranchProtection
func (mr *MockGitHubClientInterfaceMockRecorder) UpdateBranchProtection(owner, repo, branch, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBranchProtection", reflect.TypeOf((*MockGitHubClientInterface)(nil).UpdateBranchProtection), owner, repo, branch, request)
}

// CreateHook mocks base method
func (m *MockGitHubClientInterface) CreateHook(owner, repo string, hook *github.Hook) (*github.Hook, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateHook", owner, repo, hook)
	ret0, _ := ret[0].(*github.Hook)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateHook indicates an expected call of CreateHook
func (mr *MockGitHubClientInterfaceMockRecorder) CreateHook(owner, repo, hook interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateHook", reflect.TypeOf((*MockGitHubClientInterface)(nil).CreateHook), owner, repo, hook)
}
