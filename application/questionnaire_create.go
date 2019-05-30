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
type CreatingQuestionnaireUsecase struct {
	repository              QuestionnaireRepository                   `required:""`
	questionnaireIDProvider QuestionnaireIDProvider                   `required:""`
	questionIDProvider      QuestionIDProvider                        `required:""`
	authorizationService    CreatingQuestionnaireAuthorizationService `required:""`
}

//genconstructor
type QuestionItem struct {
	question Question `required:"" getter:""`
	required bool     `required:"" getter:""`
}

//genconstructor
type Question struct {
	id   domain.QuestionID `getter:"" setter:""`
	text string            `required:"" getter:""`
}

//genconstructor
type CreatingQuestionnaireUsecaseInput struct {
	title         string              `required:"" getter:""`
	questionItems []QuestionItem      `required:"" getter:""`
	schedule      domain.Schedule     `getter:"" setter:""`
	postTargets   []domain.PostTarget `getter:"" setter:""`
}

func (u CreatingQuestionnaireUsecase) Create(
	ctx context.Context,
	admin AdminDescriptor,
	workspace WorkspaceDescriptor,
	input CreatingQuestionnaireUsecaseInput,
) (domain.Questionnaire, domain.Error) {

	canCreateQuestionnaire, err := u.authorizationService.CanCreateQuestionnaire(ctx, admin, workspace)
	if err != nil {
		return domain.Questionnaire{}, err
	}
	if !canCreateQuestionnaire {
		return domain.Questionnaire{}, domain.ErrorPermissionDenied(zaperr.New("forbidden"))
	}

	questions := make([]domain.QuestionItem, len(input.QuestionItems()))
	for i, q := range input.QuestionItems() {
		questions[i] = domain.NewQuestionItem(
			domain.NewQuestion(u.questionIDProvider.NewQuestionID(), q.Question().Text()),
			q.Required(),
		)
	}
	questionnaire := domain.NewQuestionnaire(
		u.questionnaireIDProvider.NewQuestionnaireID(),
		input.Title(),
		questions,
	)
	if input.Schedule() != nil {
		questionnaire.SetSchedule(input.Schedule())
	}
	if len(input.PostTargets()) > 0 {
		questionnaire.SetPostTargets(input.PostTargets())
	}

	err = u.repository.Create(ctx, questionnaire)
	if err != nil {
		return domain.Questionnaire{}, err
	}

	return questionnaire, nil
}

type QuestionnaireIDProvider interface {
	NewQuestionnaireID() domain.QuestionnaireID
}

type QuestionIDProvider interface {
	NewQuestionID() domain.QuestionID
}

type CreatingQuestionnaireAuthorizationService interface {
	CanCreateQuestionnaire(
		ctx context.Context,
		adminDescriptor AdminDescriptor,
		workspaceDescriptor WorkspaceDescriptor,
	) (bool, domain.Error)
}
