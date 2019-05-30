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
	"github.com/opensaasstudio/meerkat/domain/mock_domain"
	"github.com/stretchr/testify/assert"
)

func TestAnsweringUsecase(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		questionnaireRepository := mock_application.NewMockQuestionnaireRepository(ctrl)
		answererRepository := mock_application.NewMockAnswererRepository(ctrl)
		answerRepository := mock_application.NewMockAnswerRepository(ctrl)
		answerIDProvider := mock_application.NewMockAnswerIDProvider(ctrl)
		postingService := mock_domain.NewMockPostingService(ctrl)
		u := application.NewAnsweringUsecase(
			questionnaireRepository,
			answererRepository,
			answerRepository,
			answerIDProvider,
			postingService,
		)

		input := application.NewAnsweringUsecaseInput(
			"questionnaireID",
			"answererID",
			[]application.AnswerInputValue{
				application.NewAnswerInputValue("questionID", "answerValue"),
			},
		)

		questionnaireRepository.EXPECT().FindByID(ctx, domain.QuestionnaireID("questionnaireID")).Return(
			domain.NewQuestionnaire("questionnaireID", "title", []domain.QuestionItem{
				domain.NewQuestionItem(domain.NewQuestion("other questionID", ""), false),
				domain.NewQuestionItem(domain.NewQuestion("questionID", ""), false),
				domain.NewQuestionItem(domain.NewQuestion("other questionID", ""), false),
			}),
			1,
			nil,
		)
		answererRepository.EXPECT().FindByID(ctx, domain.AnswererID("answererID")).Return(
			domain.NewAnswerer("answererID", "answererName"),
			1,
			nil,
		)

		answerIDProvider.EXPECT().NewAnswerID().Return(domain.AnswerID("id"))

		answerRepository.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, answer domain.Answer) domain.Error {
				assert.Equal(t, domain.AnswerID("id"), answer.ID())
				assert.Equal(t, domain.QuestionnaireID("questionnaireID"), answer.QuestionnaireID())
				assert.Equal(t, domain.QuestionID("questionID"), answer.QuestionID())
				assert.Equal(t, domain.AnswererID("answererID"), answer.AnswererID())
				return nil
			},
		)

		postingService.EXPECT().PostAnswers(ctx, gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, questionnaire domain.Questionnaire, answerer domain.Answerer, answers []domain.Answer) domain.Error {
				assert.Len(t, answers, 1)
				return nil
			},
		)

		err := u.Answer(ctx, input)
		assert.NoError(t, err)
	})
}
