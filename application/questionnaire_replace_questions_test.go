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

func TestReplacingQuestionsUsecase(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := mock_application.NewMockQuestionnaireRepository(ctrl)

	repository.EXPECT().FindByID(ctx, domain.QuestionnaireID("id")).Return(
		domain.NewQuestionnaire(
			"id",
			"title",
			[]domain.QuestionItem{
				domain.NewQuestionItem(domain.NewQuestion("question1ID", "question1"), false), // overwrite
				domain.NewQuestionItem(domain.NewQuestion("question3ID", "question3"), false), // delete
			},
		), 1, nil,
	)

	newQuestionItemParam := func(questionID domain.QuestionID, text string, required bool) application.QuestionItem {
		q := application.NewQuestion(text)
		q.SetID(questionID)
		return application.NewQuestionItem(q, required)
	}

	input := application.NewReplacingQuestionsUsecaseInput(
		"id",
		[]application.QuestionItem{
			newQuestionItemParam("question1ID", "new question1", true),              // overwrite
			application.NewQuestionItem(application.NewQuestion("question1"), true), // new
			newQuestionItemParam("question2ID", "new question2", false),             // append and overwrite
		},
	)

	questionIDProvider := mock_application.NewMockQuestionIDProvider(ctrl)

	questionIDProvider.EXPECT().NewQuestionID().DoAndReturn(func() func() domain.QuestionID {
		i := 4
		return func() domain.QuestionID {
			i++
			return domain.QuestionID("question" + strconv.Itoa(i))
		}
	}()).AnyTimes()

	repository.EXPECT().Update(ctx, gomock.Any(), 1).DoAndReturn(
		func(ctx context.Context, questionnaire domain.Questionnaire, version int) domain.Error {
			assert.Equal(t, 1, version)
			assert.Equal(t, domain.QuestionnaireID("id"), questionnaire.ID())
			assert.Equal(t, []domain.QuestionItem{
				domain.NewQuestionItem(domain.NewQuestion("question1ID", "new question1"), true),
				domain.NewQuestionItem(domain.NewQuestion("question5", "question1"), true),
				domain.NewQuestionItem(domain.NewQuestion("question2ID", "new question2"), false),
			}, questionnaire.QuestionItems())
			return nil
		},
	)

	u := application.NewReplacingQuestionsUsecase(
		repository,
		questionIDProvider,
	)
	err := u.ReplaceQuestions(ctx, input)
	assert.NoError(t, err)
}
