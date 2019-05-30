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

package application_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/opensaasstudio/meerkat/application"
	"github.com/opensaasstudio/meerkat/application/mock_application"
	"github.com/opensaasstudio/meerkat/domain"
	"github.com/stretchr/testify/assert"
)

func TestAddingAnswererUsecase(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := mock_application.NewMockNotificationTargetRepository(ctrl)
	notificationTargetIDProvider := mock_application.NewMockNotificationTargetIDProvider(ctrl)
	u := application.NewAddingAnswererUsecase(repository, notificationTargetIDProvider)

	input := application.NewAddingAnswererUsecaseInput(
		domain.NewNotificationTargetBase(
			"",
			"questionnaireID",
			"answererID",
		),
	)

	notificationTargetIDProvider.EXPECT().NewNotificationTargetID().Return(domain.NotificationTargetID("notificationTargetID"))

	repository.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, notificationTarget domain.NotificationTarget) domain.Error {
			assert.Equal(t, domain.NotificationTargetID("notificationTargetID"), notificationTarget.ID())
			assert.Equal(t, domain.QuestionnaireID("questionnaireID"), notificationTarget.QuestionnaireID())
			return nil
		},
	)

	err := u.AddAnswerer(ctx, input)
	assert.NoError(t, err)
}
