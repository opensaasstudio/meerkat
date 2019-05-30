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

	"github.com/hori-ryota/zaperr"
	"github.com/opensaasstudio/meerkat/domain"
)

//genconstructor
type UpdatingQuestionnaireUsecase struct {
	repository           QuestionnaireRepository                   `required:""`
	questionIDProvider   QuestionIDProvider                        `required:""`
	authorizationService UpdatingQuestionnaireAuthorizationService `required:""`
}

//genconstructor
type UpdatingQuestionnaireUsecaseInput struct {
	questionnaireID domain.QuestionnaireID            `required:"" getter:""`
	creatingInput   CreatingQuestionnaireUsecaseInput `required:"" getter:""`
}

func (u UpdatingQuestionnaireUsecase) Update(
	ctx context.Context,
	admin AdminDescriptor,
	workspace WorkspaceDescriptor,
	input UpdatingQuestionnaireUsecaseInput,
) (domain.Questionnaire, domain.Error) {

	questionnaire, version, err := u.repository.FindByID(ctx, input.QuestionnaireID())
	if err != nil {
		return domain.Questionnaire{}, err
	}

	canUpdateQuestionnaire, err := u.authorizationService.CanUpdateQuestionnaire(ctx, questionnaire, admin, workspace)
	if err != nil {
		return domain.Questionnaire{}, err
	}
	if !canUpdateQuestionnaire {
		return domain.Questionnaire{}, domain.ErrorPermissionDenied(zaperr.New("forbidden"))
	}

	questions := make([]domain.QuestionItem, len(input.CreatingInput().QuestionItems()))
	for i, q := range input.CreatingInput().QuestionItems() {
		id := q.Question().ID()
		if id == "" {
			id = u.questionIDProvider.NewQuestionID()
		}
		questions[i] = domain.NewQuestionItem(
			domain.NewQuestion(id, q.Question().Text()),
			q.Required(),
		)
	}
	questionnaire.OverwriteTitle(input.CreatingInput().Title())
	questionnaire.ReplaceQuestions(questions)
	if input.CreatingInput().Schedule() != nil {
		questionnaire.SetSchedule(input.CreatingInput().Schedule())
	}
	if len(input.CreatingInput().PostTargets()) > 0 {
		questionnaire.SetPostTargets(input.CreatingInput().PostTargets())
	}

	err = u.repository.Update(ctx, questionnaire, version)
	if err != nil {
		return domain.Questionnaire{}, err
	}

	return questionnaire, nil
}

type UpdatingQuestionnaireAuthorizationService interface {
	CanUpdateQuestionnaire(
		ctx context.Context,
		questionnaire domain.Questionnaire,
		adminDescriptor AdminDescriptor,
		workspaceDescriptor WorkspaceDescriptor,
	) (bool, domain.Error)
}
