package application

import (
	"context"

	"github.com/opensaasstudio/meerkat/domain"
)

//genconstructor
type RemovingAnswererUsecase struct {
	searcher   domain.NotificationTargetSearcher `required:""`
	repository NotificationTargetRepository      `required:""`
}

//genconstructor
type RemovingAnswererUsecaseInput struct {
	questionnaireID domain.QuestionnaireID `required:"" getter:""`
	answererID      domain.AnswererID      `required:"" getter:""`
}

func (u RemovingAnswererUsecase) RemoveAnswerer(ctx context.Context, input RemovingAnswererUsecaseInput) domain.Error {
	notificationTargets, err := u.searcher.SearchByQuestionnaireIDAndAnswererID(
		ctx,
		input.QuestionnaireID(),
		input.AnswererID(),
	)
	if err != nil {
		return err
	}
	for _, notificationTarget := range notificationTargets {
		if err := u.repository.Delete(ctx, notificationTarget); err != nil {
			return err
		}
	}
	return nil
}
