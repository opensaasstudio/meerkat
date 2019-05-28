package domain

//go:generate mockgen -source $GOFILE -destination mock_$GOPACKAGE/${GOFILE}_mock.go

import (
	"context"
)

type PostingService interface {
	PostAnswers(ctx context.Context, questionnaire Questionnaire, answerer Answerer, answers []Answer) Error
}
