package application

//go:generate mockgen -source $GOFILE -destination mock_$GOPACKAGE/${GOFILE}_mock.go

import (
	"context"

	"github.com/opensaasstudio/meerkat/domain"
)

//genconstructor
type AnsweringUsecase struct {
	questionnaireRepository QuestionnaireRepository `required:""`
	answererRepository      AnswererRepository      `required:""`
	answerRepository        AnswerRepository        `required:""`
	answerIDProvider        AnswerIDProvider        `required:""`
	postingService          domain.PostingService   `required:""`
}

//genconstructor
type AnswerInputValue struct {
	questionID domain.QuestionID `required:"" getter:""`
	value      string            `required:"" getter:""`
}

//genconstructor
type AnsweringUsecaseInput struct {
	questionnaireID domain.QuestionnaireID `required:"" getter:""`
	answererID      domain.AnswererID      `required:"" getter:""`
	answers         []AnswerInputValue     `required:"" getter:""`
}

func (u AnsweringUsecase) Answer(ctx context.Context, input AnsweringUsecaseInput) domain.Error {
	questionnaire, _, err := u.questionnaireRepository.FindByID(ctx, input.QuestionnaireID())
	if err != nil {
		return err
	}
	answerer, _, err := u.answererRepository.FindByID(ctx, input.AnswererID())
	if err != nil {
		return err
	}
	answers := make([]domain.Answer, len(input.Answers()))
	for i, v := range input.Answers() {
		answer, err := questionnaire.Answer(
			u.answerIDProvider.NewAnswerID(),
			v.QuestionID(),
			input.AnswererID(),
			v.Value(),
		)
		if err != nil {
			return err
		}
		if err := u.answerRepository.Create(ctx, answer); err != nil {
			return err
		}
		answers[i] = answer
	}
	return u.postingService.PostAnswers(ctx, questionnaire, answerer, answers)
}

type AnswerIDProvider interface {
	NewAnswerID() domain.AnswerID
}
