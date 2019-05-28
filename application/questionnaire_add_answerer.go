package application

//go:generate mockgen -source $GOFILE -destination mock_$GOPACKAGE/${GOFILE}_mock.go

import (
	"context"

	"github.com/opensaasstudio/meerkat/domain"
)

//genconstructor
type AddingAnswererUsecase struct {
	repository                   NotificationTargetRepository `required:""`
	notificationTargetIDProvider NotificationTargetIDProvider `required:""`
}

//genconstructor
type AddingAnswererUsecaseInput struct {
	notificationTarget domain.NotificationTarget `required:"" getter:""`
}

func (u AddingAnswererUsecase) AddAnswerer(ctx context.Context, input AddingAnswererUsecaseInput) domain.Error {
	switch nt := input.NotificationTarget().(type) {
	case domain.NotificationTargetSlack:
		n := domain.NewNotificationTargetSlack(
			u.notificationTargetIDProvider.NewNotificationTargetID(),
			nt.QuestionnaireID(),
			nt.AnswererID(),
			nt.ChannelID(),
			nt.UserID(),
		)
		n.ToggleNeedsMention(nt.NeedsMention())
		return u.repository.Create(ctx, n)
	default:
		n := domain.NewNotificationTargetBase(
			u.notificationTargetIDProvider.NewNotificationTargetID(),
			nt.QuestionnaireID(),
			nt.AnswererID(),
		)
		return u.repository.Create(ctx, n)
	}
}

type NotificationTargetIDProvider interface {
	NewNotificationTargetID() domain.NotificationTargetID
}
