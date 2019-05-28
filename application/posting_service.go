package application

import (
	"context"

	"github.com/hori-ryota/zaperr"
	"github.com/opensaasstudio/meerkat/domain"
)

//genconstructor
type PostingService struct {
	postSlack func(ctx context.Context, target domain.PostTargetSlack, questionnaire domain.Questionnaire, answerer domain.Answerer, answers []domain.Answer) domain.Error `required:""`
}

func (s PostingService) PostAnswers(ctx context.Context, questionnaire domain.Questionnaire, answerer domain.Answerer, answers []domain.Answer) domain.Error {
	for _, postTarget := range questionnaire.PostTargets() {
		switch pt := postTarget.(type) {
		case domain.PostTargetSlack:
			return s.postSlack(ctx, pt, questionnaire, answerer, answers)
		default:
			return domain.ErrorUnknown(zaperr.New("unknown posting target type"))
		}
	}
	return nil
}
