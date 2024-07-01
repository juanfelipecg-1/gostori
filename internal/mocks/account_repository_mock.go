// Code generated by MockGen. DO NOT EDIT.
// Source: internal/ports/account_repository.go
//
// Generated by this command:
//
//	mockgen -source=internal/ports/account_repository.go -destination=internal/mocks/account_repository_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	domain "github.com/juanfcgarcia/gostori/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockAccountRepository is a mock of AccountRepository interface.
type MockAccountRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAccountRepositoryMockRecorder
}

// MockAccountRepositoryMockRecorder is the mock recorder for MockAccountRepository.
type MockAccountRepositoryMockRecorder struct {
	mock *MockAccountRepository
}

// NewMockAccountRepository creates a new mock instance.
func NewMockAccountRepository(ctrl *gomock.Controller) *MockAccountRepository {
	mock := &MockAccountRepository{ctrl: ctrl}
	mock.recorder = &MockAccountRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountRepository) EXPECT() *MockAccountRepositoryMockRecorder {
	return m.recorder
}

// CreateAccount mocks base method.
func (m *MockAccountRepository) CreateAccount(ctx context.Context, email string) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", ctx, email)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockAccountRepositoryMockRecorder) CreateAccount(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockAccountRepository)(nil).CreateAccount), ctx, email)
}

// GetAccountByID mocks base method.
func (m *MockAccountRepository) GetAccountByID(ctx context.Context, id int32) (*domain.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountByID", ctx, id)
	ret0, _ := ret[0].(*domain.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountByID indicates an expected call of GetAccountByID.
func (mr *MockAccountRepositoryMockRecorder) GetAccountByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountByID", reflect.TypeOf((*MockAccountRepository)(nil).GetAccountByID), ctx, id)
}
