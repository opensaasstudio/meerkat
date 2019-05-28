package application_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/opensaasstudio/meerkat/application"
	"github.com/opensaasstudio/meerkat/application/mock_application"
	"github.com/opensaasstudio/meerkat/domain"
	"github.com/opensaasstudio/meerkat/domain/mock_domain"
	"github.com/stretchr/testify/assert"
)

func TestRemovingAnswererUsecase(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	searcher := mock_domain.NewMockNotificationTargetSearcher(ctrl)
	repository := mock_application.NewMockNotificationTargetRepository(ctrl)
	u := application.NewRemovingAnswererUsecase(searcher, repository)

	input := application.NewRemovingAnswererUsecaseInput(
		"questionnaireID",
		"answererID",
	)

	searcher.EXPECT().SearchByQuestionnaireIDAndAnswererID(
		ctx,
		domain.QuestionnaireID("questionnaireID"),
		domain.AnswererID("answererID"),
	).Return([]domain.NotificationTarget{
		domain.NewNotificationTargetBase("notificationTarget1", "questionnaireID", "answererID"),
		domain.NewNotificationTargetBase("notificationTarget2", "questionnaireID", "answererID"),
	}, nil)

	repository.EXPECT().Delete(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, notificationTarget domain.NotificationTarget) domain.Error {
			assert.Equal(t, domain.QuestionnaireID("questionnaireID"), notificationTarget.QuestionnaireID())
			assert.Equal(t, domain.AnswererID("answererID"), notificationTarget.AnswererID())
			return nil
		},
	).Times(2)

	err := u.RemoveAnswerer(ctx, input)
	assert.NoError(t, err)
}
