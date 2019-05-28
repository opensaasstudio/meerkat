package application

import (
	"context"
	"time"

	"github.com/opensaasstudio/meerkat/domain"
)

//genconstructor
type LastExecutedRecorder struct {
	quesitonnaireRepository QuestionnaireRepository `required:""`
}

func (s LastExecutedRecorder) RecordLastExecuted(ctx context.Context, questionnaire domain.Questionnaire, lastExecuted time.Time) domain.Error {
	q, version, err := s.quesitonnaireRepository.FindByID(ctx, questionnaire.ID())
	if err != nil {
		return err
	}
	q.SetLastExecuted(lastExecuted)
	return s.quesitonnaireRepository.Update(ctx, q, version)
}
