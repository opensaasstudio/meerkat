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

func TestCreatingQuestionnaireUsecase(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := mock_application.NewMockQuestionnaireRepository(ctrl)
	questionnaireIDProvider := mock_application.NewMockQuestionnaireIDProvider(ctrl)
	questionIDProvider := mock_application.NewMockQuestionIDProvider(ctrl)
	authorizationService := mock_application.NewMockCreatingQuestionnaireAuthorizationService(ctrl)

	u := application.NewCreatingQuestionnaireUsecase(
		repository,
		questionnaireIDProvider,
		questionIDProvider,
		authorizationService,
	)
	input := application.NewCreatingQuestionnaireUsecaseInput(
		"title",
		[]application.QuestionItem{
			application.NewQuestionItem(application.NewQuestion("question1"), true),
			application.NewQuestionItem(application.NewQuestion("question2"), false),
		},
	)
	input.SetSchedule(domain.NewWeekdayAndTimeSchedule(1, 2, 3, 0, true, true, true, true, true, true, true))

	authorizationService.EXPECT().CanCreateQuestionnaire(ctx, gomock.Any(), gomock.Any()).Return(true, nil)

	questionnaireIDProvider.EXPECT().NewQuestionnaireID().Return(domain.QuestionnaireID("id"))

	questionIDProvider.EXPECT().NewQuestionID().DoAndReturn(func() func() domain.QuestionID {
		i := 0
		return func() domain.QuestionID {
			i++
			return domain.QuestionID("question" + strconv.Itoa(i))
		}
	}()).AnyTimes()

	repository.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	output, err := u.Create(
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
