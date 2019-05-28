package application

//go:generate mockgen -source $GOFILE -destination mock_$GOPACKAGE/${GOFILE}_mock.go

import (
	"context"

	"github.com/opensaasstudio/meerkat/domain"
)

type QuestionnaireRepository interface {
	Create(ctx context.Context, questionnaire domain.Questionnaire) domain.Error
	Update(ctx context.Context, questionnaire domain.Questionnaire, version int) domain.Error
	Delete(ctx context.Context, questionnaire domain.Questionnaire) domain.Error
	FindByID(ctx context.Context, id domain.QuestionnaireID) (questionnaire domain.Questionnaire, version int, derr domain.Error)
}

type AnswererRepository interface {
	Create(ctx context.Context, answerer domain.Answerer) domain.Error
	Update(ctx context.Context, answerer domain.Answerer, version int) domain.Error
	Delete(ctx context.Context, answerer domain.Answerer) domain.Error
	FindByID(ctx context.Context, id domain.AnswererID) (answerer domain.Answerer, version int, derr domain.Error)
}

type NotificationTargetRepository interface {
	Create(ctx context.Context, notificationTarget domain.NotificationTarget) domain.Error
	Update(ctx context.Context, notificationTarget domain.NotificationTarget, version int) domain.Error
	Delete(ctx context.Context, notificationTarget domain.NotificationTarget) domain.Error
}

type AnswerRepository interface {
	Create(ctx context.Context, answer domain.Answer) domain.Error
	Update(ctx context.Context, answer domain.Answer, version int) domain.Error
	Delete(ctx context.Context, answer domain.Answer) domain.Error
	FindByID(ctx context.Context, answerID domain.AnswerID) (answer domain.Answer, version int, derr domain.Error)
}
