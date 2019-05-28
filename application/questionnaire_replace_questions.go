package application

import (
	"context"

	"github.com/opensaasstudio/meerkat/domain"
)

//genconstructor
type ReplacingQuestionsUsecase struct {
	repository         QuestionnaireRepository `required:""`
	questionIDProvider QuestionIDProvider      `required:""`
}

//genconstructor
type ReplacingQuestionsUsecaseInput struct {
	questionnaireID domain.QuestionnaireID `required:"" getter:""`
	questionItems   []QuestionItem         `required:"" getter:""`
}

func (u ReplacingQuestionsUsecase) ReplaceQuestions(
	ctx context.Context,
	input ReplacingQuestionsUsecaseInput,
) domain.Error {
	questionnaire, version, err := u.repository.FindByID(ctx, "id")
	if err != nil {
		return err
	}
	questions := make([]domain.QuestionItem, len(input.QuestionItems()))
	for i, q := range input.QuestionItems() {
		id := q.Question().ID()
		if id == "" {
			id = u.questionIDProvider.NewQuestionID()
		}
		questions[i] = domain.NewQuestionItem(
			domain.NewQuestion(id, q.Question().Text()),
			q.Required(),
		)
	}
	questionnaire.ReplaceQuestions(questions)
	return u.repository.Update(ctx, questionnaire, version)
}
