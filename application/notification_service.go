package application

import (
	"context"

	"github.com/hori-ryota/zaperr"
	"github.com/opensaasstudio/meerkat/domain"
)

//genconstructor
type NotificationService struct {
	notifySlack func(ctx context.Context, target domain.NotificationTargetSlack, questionnaire domain.Questionnaire) domain.Error `required:""`
}

func (s NotificationService) Notify(ctx context.Context, notificationTarget domain.NotificationTarget, questionnaire domain.Questionnaire) domain.Error {
	switch nt := notificationTarget.(type) {
	case domain.NotificationTargetSlack:
		return s.notifySlack(ctx, nt, questionnaire)
	}
	return domain.ErrorUnknown(zaperr.New("unknown notification target type"))
}
