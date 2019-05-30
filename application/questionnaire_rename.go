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
