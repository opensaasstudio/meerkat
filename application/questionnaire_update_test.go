package application_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/opensaasstudio/meerkat/application"
	"github.com/opensaasstudio/meerkat/application/mock_application"
	"github.com/opensaasstudio/meerkat/domain"
	"github.com/stretchr/testify/assert"
)

func TestUpdatingQuestionnaireUsecase(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := mock_application.NewMockQuestionnaireRepository(ctrl)
	questionIDProvider := mock_application.NewMockQuestionIDProvider(ctrl)
	authorizationService := mock_application.NewMockUpdatingQuestionnaireAuthorizationService(ctrl)

	u := application.NewUpdatingQuestionnaireUsecase(
		repository,
		questionIDProvider,
		authorizationService,
	)

	repository.EXPECT().FindByID(ctx, domain.QuestionnaireID("id")).Return(
		domain.NewQuestionnaire(
			"id",
			"title",
			nil,
		), 1, nil,
	)

	creatingInput := application.NewCreatingQuestionnaireUsecaseInput(
		"title",
		[]application.QuestionItem{
			application.NewQuestionItem(application.NewQuestion("question1"), true),
			application.NewQuestionItem(application.NewQuestion("question2"), false),
		},
	)
	creatingInput.SetSchedule(domain.NewWeekdayAndTimeSchedule(1, 2, 3, 0, true, true, true, true, true, true, true))
	input := application.NewUpdatingQuestionnaireUsecaseInput(
		domain.QuestionnaireID("id"),
		creatingInput,
	)

	authorizationService.EXPECT().CanUpdateQuestionnaire(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)

	questionIDProvider.EXPECT().NewQuestionID().DoAndReturn(func() func() domain.QuestionID {
		i := 0
		return func() domain.QuestionID {
			i++
			return domain.QuestionID("question" + strconv.Itoa(i))
		}
	}()).AnyTimes()

	repository.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).Return(nil)

	output, err := u.Update(
		ctx,
		application.AdminDescriptor{},     // TODO
		application.WorkspaceDescriptor{}, // TODO
		input,
	)
	assert.NoError(t, err)

	assert.NotEmpty(t, output.ID())
	assert.Equal(t, "title", output.Title())
	assert.Equal(t, "question1", output.QuestionItems()[0].Question().Text())
	assert.NotEmpty(t, output.Schedule())
}
