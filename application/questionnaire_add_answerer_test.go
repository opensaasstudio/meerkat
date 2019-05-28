package application_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/opensaasstudio/meerkat/application"
	"github.com/opensaasstudio/meerkat/application/mock_application"
	"github.com/opensaasstudio/meerkat/domain"
	"github.com/stretchr/testify/assert"
)

func TestAddingAnswererUsecase(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := mock_application.NewMockNotificationTargetRepository(ctrl)
	notificationTargetIDProvider := mock_application.NewMockNotificationTargetIDProvider(ctrl)
	u := application.NewAddingAnswererUsecase(repository, notificationTargetIDProvider)

	input := application.NewAddingAnswererUsecaseInput(
		domain.NewNotificationTargetBase(
			"",
			"questionnaireID",
			"answererID",
		),
	)

	notificationTargetIDProvider.EXPECT().NewNotificationTargetID().Return(domain.NotificationTargetID("notificationTargetID"))

	repository.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, notificationTarget domain.NotificationTarget) domain.Error {
			assert.Equal(t, domain.NotificationTargetID("notificationTargetID"), notificationTarget.ID())
			assert.Equal(t, domain.QuestionnaireID("questionnaireID"), notificationTarget.QuestionnaireID())
			return nil
		},
	)

	err := u.AddAnswerer(ctx, input)
	assert.NoError(t, err)
}
