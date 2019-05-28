// Code generated by MockGen. DO NOT EDIT.
// Source: questionnaire_create.go

// Package mock_application is a generated GoMock package.
package mock_application

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	application "github.com/opensaasstudio/meerkat/application"
	domain "github.com/opensaasstudio/meerkat/domain"
	reflect "reflect"
)

// MockQuestionnaireIDProvider is a mock of QuestionnaireIDProvider interface
type MockQuestionnaireIDProvider struct {
	ctrl     *gomock.Controller
	recorder *MockQuestionnaireIDProviderMockRecorder
}

// MockQuestionnaireIDProviderMockRecorder is the mock recorder for MockQuestionnaireIDProvider
type MockQuestionnaireIDProviderMockRecorder struct {
	mock *MockQuestionnaireIDProvider
}

// NewMockQuestionnaireIDProvider creates a new mock instance
func NewMockQuestionnaireIDProvider(ctrl *gomock.Controller) *MockQuestionnaireIDProvider {
	mock := &MockQuestionnaireIDProvider{ctrl: ctrl}
	mock.recorder = &MockQuestionnaireIDProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockQuestionnaireIDProvider) EXPECT() *MockQuestionnaireIDProviderMockRecorder {
	return m.recorder
}

// NewQuestionnaireID mocks base method
func (m *MockQuestionnaireIDProvider) NewQuestionnaireID() domain.QuestionnaireID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewQuestionnaireID")
	ret0, _ := ret[0].(domain.QuestionnaireID)
	return ret0
}

// NewQuestionnaireID indicates an expected call of NewQuestionnaireID
func (mr *MockQuestionnaireIDProviderMockRecorder) NewQuestionnaireID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewQuestionnaireID", reflect.TypeOf((*MockQuestionnaireIDProvider)(nil).NewQuestionnaireID))
}

// MockQuestionIDProvider is a mock of QuestionIDProvider interface
type MockQuestionIDProvider struct {
	ctrl     *gomock.Controller
	recorder *MockQuestionIDProviderMockRecorder
}

// MockQuestionIDProviderMockRecorder is the mock recorder for MockQuestionIDProvider
type MockQuestionIDProviderMockRecorder struct {
	mock *MockQuestionIDProvider
}

// NewMockQuestionIDProvider creates a new mock instance
func NewMockQuestionIDProvider(ctrl *gomock.Controller) *MockQuestionIDProvider {
	mock := &MockQuestionIDProvider{ctrl: ctrl}
	mock.recorder = &MockQuestionIDProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockQuestionIDProvider) EXPECT() *MockQuestionIDProviderMockRecorder {
	return m.recorder
}

// NewQuestionID mocks base method
func (m *MockQuestionIDProvider) NewQuestionID() domain.QuestionID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewQuestionID")
	ret0, _ := ret[0].(domain.QuestionID)
	return ret0
}

// NewQuestionID indicates an expected call of NewQuestionID
func (mr *MockQuestionIDProviderMockRecorder) NewQuestionID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewQuestionID", reflect.TypeOf((*MockQuestionIDProvider)(nil).NewQuestionID))
}

// MockCreatingQuestionnaireAuthorizationService is a mock of CreatingQuestionnaireAuthorizationService interface
type MockCreatingQuestionnaireAuthorizationService struct {
	ctrl     *gomock.Controller
	recorder *MockCreatingQuestionnaireAuthorizationServiceMockRecorder
}

// MockCreatingQuestionnaireAuthorizationServiceMockRecorder is the mock recorder for MockCreatingQuestionnaireAuthorizationService
type MockCreatingQuestionnaireAuthorizationServiceMockRecorder struct {
	mock *MockCreatingQuestionnaireAuthorizationService
}

// NewMockCreatingQuestionnaireAuthorizationService creates a new mock instance
func NewMockCreatingQuestionnaireAuthorizationService(ctrl *gomock.Controller) *MockCreatingQuestionnaireAuthorizationService {
	mock := &MockCreatingQuestionnaireAuthorizationService{ctrl: ctrl}
	mock.recorder = &MockCreatingQuestionnaireAuthorizationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCreatingQuestionnaireAuthorizationService) EXPECT() *MockCreatingQuestionnaireAuthorizationServiceMockRecorder {
	return m.recorder
}

// CanCreateQuestionnaire mocks base method
func (m *MockCreatingQuestionnaireAuthorizationService) CanCreateQuestionnaire(ctx context.Context, adminDescriptor application.AdminDescriptor, workspaceDescriptor application.WorkspaceDescriptor) (bool, domain.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CanCreateQuestionnaire", ctx, adminDescriptor, workspaceDescriptor)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(domain.Error)
	return ret0, ret1
}

// CanCreateQuestionnaire indicates an expected call of CanCreateQuestionnaire
func (mr *MockCreatingQuestionnaireAuthorizationServiceMockRecorder) CanCreateQuestionnaire(ctx, adminDescriptor, workspaceDescriptor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CanCreateQuestionnaire", reflect.TypeOf((*MockCreatingQuestionnaireAuthorizationService)(nil).CanCreateQuestionnaire), ctx, adminDescriptor, workspaceDescriptor)
}
