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

/*
- 上書き
*/

func TestOverwritingQuestionnaireTitleUsecase(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := mock_application.NewMockQuestionnaireRepository(ctrl)

	repository.EXPECT().FindByID(ctx, domain.QuestionnaireID("id")).Return(
		domain.NewQuestionnaire(
			"id",
			"title",
			nil,
		), 1, nil,
	)

	input := application.NewOverwritingQuestionnaireTitleUsecaseInput(
		"id",
		"new title",
	)

	repository.EXPECT().Update(ctx, gomock.Any(), 1).DoAndReturn(
		func(ctx context.Context, questionnaire domain.Questionnaire, version int) domain.Error {
			assert.Equal(t, 1, version)
			assert.Equal(t, domain.QuestionnaireID("id"), questionnaire.ID())
			assert.Equal(t, "new title", questionnaire.Title())
			return nil
		},
	)

	u := application.NewOverwritingQuestionnaireTitleUsecase(repository)
	err := u.OverwriteQuestionnaireTitle(ctx, input)
	assert.NoError(t, err)
}
