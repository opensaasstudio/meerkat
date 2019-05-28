package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/opensaasstudio/meerkat/application"
	"github.com/opensaasstudio/meerkat/domain"
	"github.com/opensaasstudio/meerkat/domain/mock_domain"
	"github.com/stretchr/testify/assert"
)

func TestAskingUsecase(t *testing.T) {
	t.Run("AskAllIfNeeded", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		now := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)

		notificationService := mock_domain.NewMockNotificationService(ctrl)
		questionnaireSearcher := mock_domain.NewMockQuestionnaireSearcher(ctrl)
		notificationTargetSearcher := mock_domain.NewMockNotificationTargetSearcher(ctrl)
		lastExecutedRecorder := mock_domain.NewMockLastExecutedRecorder(ctrl)

		askingService := domain.NewAskingService(
			questionnaireSearcher,
			notificationTargetSearcher,
			notificationService,
			lastExecutedRecorder,
		)
		askingService.OverwriteNowFunc(func() time.Time { return now })

		u := application.NewAskingUsecase(askingService)

		// stub
		questionnaireSearcher.EXPECT().SearchExecutionNeeded(ctx).Return([]domain.Questionnaire{
			domain.NewQuestionnaire("id", "title", nil),
		}, nil)
		notificationTargetSearcher.EXPECT().SearchByQuestionnaireID(ctx, gomock.Any()).Return([]domain.NotificationTarget{
			domain.NewNotificationTargetBase("notifiationTargetID", "id", "answererID"),
		}, nil)

		notificationService.EXPECT().Notify(ctx, gomock.Any(), gomock.Any())
		lastExecutedRecorder.EXPECT().RecordLastExecuted(ctx, gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, questionnaire domain.Questionnaire, lastExecuted time.Time) domain.Error {
			assert.Equal(t, now, lastExecuted)
			return nil
		})
		err := u.AskAllIfNeeded(ctx)
		assert.NoError(t, err)
	})
}
