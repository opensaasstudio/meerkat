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

import (
	"context"

	"github.com/opensaasstudio/meerkat/domain"
)

//genconstructor
type ReplacingQuestionsUsecase struct {
	repository         QuestionnaireRepository `required:""`
	questionIDProvider QuestionIDProvider      `required:""`
}

//genconstructor
type ReplacingQuestionsUsecaseInput struct {
	questionnaireID domain.QuestionnaireID `required:"" getter:""`
	questionItems   []QuestionItem         `required:"" getter:""`
}

func (u ReplacingQuestionsUsecase) ReplaceQuestions(
	ctx context.Context,
	input ReplacingQuestionsUsecaseInput,
) domain.Error {
	questionnaire, version, err := u.repository.FindByID(ctx, "id")
	if err != nil {
		return err
	}
	questions := make([]domain.QuestionItem, len(input.QuestionItems()))
	for i, q := range input.QuestionItems() {
		id := q.Question().ID()
		if id == "" {
			id = u.questionIDProvider.NewQuestionID()
		}
		questions[i] = domain.NewQuestionItem(
			domain.NewQuestion(id, q.Question().Text()),
			q.Required(),
		)
	}
	questionnaire.ReplaceQuestions(questions)
	return u.repository.Update(ctx, questionnaire, version)
}
