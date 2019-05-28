package application

import (
	"context"

	"github.com/opensaasstudio/meerkat/domain"
)

//genconstructor
type OverwritingQuestionnaireTitleUsecase struct {
	repository QuestionnaireRepository `required:""`
}

//genconstructor
type OverwritingQuestionnaireTitleUsecaseInput struct {
	questionnaireID domain.QuestionnaireID `required:"" getter:""`
	title           string                 `required:"" getter:""`
}

func (u OverwritingQuestionnaireTitleUsecase) OverwriteQuestionnaireTitle(
	ctx context.Context,
	input OverwritingQuestionnaireTitleUsecaseInput,
) domain.Error {
	questionnaire, version, err := u.repository.FindByID(ctx, input.QuestionnaireID())
	if err != nil {
		return err
	}
	questionnaire.OverwriteTitle(input.Title())
	return u.repository.Update(ctx, questionnaire, version)
}
