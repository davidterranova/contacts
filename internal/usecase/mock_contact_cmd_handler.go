// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/davidterranova/contacts/internal/usecase (interfaces: ContactCmdHandler)

// Package usecase is a generated GoMock package.
package usecase

import (
        domain "github.com/davidterranova/contacts/internal/domain"
        eventsourcing "github.com/davidterranova/contacts/pkg/eventsourcing"
        "go.uber.org/mock/gomock"
        reflect "reflect"
)

// MockContactCmdHandler is a mock of ContactCmdHandler interface.
type MockContactCmdHandler struct {
        ctrl     *gomock.Controller
        recorder *MockContactCmdHandlerMockRecorder
}

// MockContactCmdHandlerMockRecorder is the mock recorder for MockContactCmdHandler.
type MockContactCmdHandlerMockRecorder struct {
        mock *MockContactCmdHandler
}

// NewMockContactCmdHandler creates a new mock instance.
func NewMockContactCmdHandler(ctrl *gomock.Controller) *MockContactCmdHandler {
        mock := &MockContactCmdHandler{ctrl: ctrl}
        mock.recorder = &MockContactCmdHandlerMockRecorder{mock}
        return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContactCmdHandler) EXPECT() *MockContactCmdHandlerMockRecorder {
        return m.recorder
}

// Handle mocks base method.
func (m *MockContactCmdHandler) Handle(arg0 eventsourcing.Command[*domain.Contact]) (*domain.Contact, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "Handle", arg0)
        ret0, _ := ret[0].(*domain.Contact)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// Handle indicates an expected call of Handle.
func (mr *MockContactCmdHandlerMockRecorder) Handle(arg0 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockContactCmdHandler)(nil).Handle), arg0)
}
