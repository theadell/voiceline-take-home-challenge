// Code generated by MockGen. DO NOT EDIT.
// Source: internal/db/querier.go
//
// Generated by this command:
//
//	mockgen -source=internal/db/querier.go -destination=internal/db/mock/querier_mock.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	db "adelhub.com/voiceline/internal/db"
	gomock "go.uber.org/mock/gomock"
)

// MockQuerier is a mock of Querier interface.
type MockQuerier struct {
	ctrl     *gomock.Controller
	recorder *MockQuerierMockRecorder
}

// MockQuerierMockRecorder is the mock recorder for MockQuerier.
type MockQuerierMockRecorder struct {
	mock *MockQuerier
}

// NewMockQuerier creates a new mock instance.
func NewMockQuerier(ctrl *gomock.Controller) *MockQuerier {
	mock := &MockQuerier{ctrl: ctrl}
	mock.recorder = &MockQuerierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuerier) EXPECT() *MockQuerierMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockQuerier) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, arg)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockQuerierMockRecorder) CreateUser(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockQuerier)(nil).CreateUser), ctx, arg)
}

// CreateUserProvider mocks base method.
func (m *MockQuerier) CreateUserProvider(ctx context.Context, arg db.CreateUserProviderParams) (db.UserProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserProvider", ctx, arg)
	ret0, _ := ret[0].(db.UserProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUserProvider indicates an expected call of CreateUserProvider.
func (mr *MockQuerierMockRecorder) CreateUserProvider(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserProvider", reflect.TypeOf((*MockQuerier)(nil).CreateUserProvider), ctx, arg)
}

// GetUserAndProviderInfo mocks base method.
func (m *MockQuerier) GetUserAndProviderInfo(ctx context.Context, arg db.GetUserAndProviderInfoParams) (db.GetUserAndProviderInfoRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserAndProviderInfo", ctx, arg)
	ret0, _ := ret[0].(db.GetUserAndProviderInfoRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserAndProviderInfo indicates an expected call of GetUserAndProviderInfo.
func (mr *MockQuerierMockRecorder) GetUserAndProviderInfo(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserAndProviderInfo", reflect.TypeOf((*MockQuerier)(nil).GetUserAndProviderInfo), ctx, arg)
}

// GetUserByEmail mocks base method.
func (m *MockQuerier) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, email)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockQuerierMockRecorder) GetUserByEmail(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockQuerier)(nil).GetUserByEmail), ctx, email)
}

// UpdateUserPassword mocks base method.
func (m *MockQuerier) UpdateUserPassword(ctx context.Context, arg db.UpdateUserPasswordParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPassword", ctx, arg)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserPassword indicates an expected call of UpdateUserPassword.
func (mr *MockQuerierMockRecorder) UpdateUserPassword(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPassword", reflect.TypeOf((*MockQuerier)(nil).UpdateUserPassword), ctx, arg)
}
