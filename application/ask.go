package application

import (
	"context"

	"github.com/opensaasstudio/meerkat/domain"
)

//genconstructor
type AskingUsecase struct {
	askingService domain.AskingService `required:""`
}

func (u AskingUsecase) AskAllIfNeeded(ctx context.Context) domain.Error {
	return u.askingService.AskAllIfNeeded(ctx)
}
