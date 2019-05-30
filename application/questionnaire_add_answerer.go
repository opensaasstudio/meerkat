// Copyright 2019 The meerkat Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
