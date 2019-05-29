package domain

import (
	"context"
)

//go:generate mockgen -source $GOFILE -destination mock_$GOPACKAGE/${GOFILE}_mock.go

type NotificationTargetSearcher interface {
	SearchByQuestionnaireIDAndAnswererID(ctx context.Context, questionnaireID QuestionnaireID, answererID AnswererID) ([]NotificationTarget, Error)
	SearchByQuestionnaireID(ctx context.Context, questionnaireID QuestionnaireID) ([]NotificationTarget, Error)
}

type QuestionnaireSearcher interface {
	SearchExecutionNeeded(ctx context.Context) ([]Questionnaire, Error)
	FetchAll(ctx context.Context) ([]Questionnaire, Error)
	FindByID(ctx context.Context, id QuestionnaireID) (questionnaire Questionnaire, version int, derr Error)
}

type AnswererSearcher interface {
	FetchAll(ctx context.Context) ([]Answerer, Error)
}
