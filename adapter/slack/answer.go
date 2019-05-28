package slack

import (
	"context"

	"github.com/nlopes/slack"
	"github.com/opensaasstudio/meerkat/application"
	"github.com/opensaasstudio/meerkat/domain"
)

//genconstructor
type AnsweringHandler struct {
	slackClient *slack.Client                `required:""`
	usecase     application.AnsweringUsecase `required:""`
}

type Answer struct {
	Question Question
	Value    string
}

type AnsweringHandlerInput struct {
	QuestionnaireID    string
	QuestionnaireTitle string
	AnswererID         string
	Answers            []Answer
}

func (p AnsweringHandlerInput) ToUsecaseInput() application.AnsweringUsecaseInput {
	answers := make([]application.AnswerInputValue, len(p.Answers))
	for i := range p.Answers {
		answers[i] = application.NewAnswerInputValue(
			domain.QuestionID(p.Answers[i].Question.ID),
			p.Answers[i].Value,
		)
	}
	return application.NewAnsweringUsecaseInput(
		domain.QuestionnaireID(p.QuestionnaireID),
		domain.AnswererID(p.AnswererID),
		answers,
	)
}

func (h AnsweringHandler) Execute(
	ctx context.Context,
	input AnsweringHandlerInput,
) domain.Error {
	return h.usecase.Answer(
		ctx,
		input.ToUsecaseInput(),
	)
}
