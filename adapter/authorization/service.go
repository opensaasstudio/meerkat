package authorization

import (
	"context"

	"github.com/opensaasstudio/meerkat/application"
	"github.com/opensaasstudio/meerkat/domain"
)

//genconstructor
type Service struct {
}

// TODO implement
func (s Service) CanCreateQuestionnaire(ctx context.Context, adminDescriptor application.AdminDescriptor, workspaceDescriptor application.WorkspaceDescriptor) (bool, domain.Error) {
	return true, nil
}

// TODO implement
func (s Service) CanUpdateQuestionnaire(ctx context.Context, qestionnaire domain.Questionnaire, adminDescriptor application.AdminDescriptor, workspaceDescriptor application.WorkspaceDescriptor) (bool, domain.Error) {
	return true, nil
}
