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

type QuestionnaireRepository interface {
	Create(ctx context.Context, questionnaire domain.Questionnaire) domain.Error
	Update(ctx context.Context, questionnaire domain.Questionnaire, version int) domain.Error
	Delete(ctx context.Context, questionnaire domain.Questionnaire) domain.Error
	FindByID(ctx context.Context, id domain.QuestionnaireID) (questionnaire domain.Questionnaire, version int, derr domain.Error)
}

type AnswererRepository interface {
	Create(ctx context.Context, answerer domain.Answerer) domain.Error
	Update(ctx context.Context, answerer domain.Answerer, version int) domain.Error
	Delete(ctx context.Context, answerer domain.Answerer) domain.Error
	FindByID(ctx context.Context, id domain.AnswererID) (answerer domain.Answerer, version int, derr domain.Error)
}

type NotificationTargetRepository interface {
	Create(ctx context.Context, notificationTarget domain.NotificationTarget) domain.Error
	Update(ctx context.Context, notificationTarget domain.NotificationTarget, version int) domain.Error
	Delete(ctx context.Context, notificationTarget domain.NotificationTarget) domain.Error
}

type AnswerRepository interface {
	Create(ctx context.Context, answer domain.Answer) domain.Error
	Update(ctx context.Context, answer domain.Answer, version int) domain.Error
	Delete(ctx context.Context, answer domain.Answer) domain.Error
	FindByID(ctx context.Context, answerID domain.AnswerID) (answer domain.Answer, version int, derr domain.Error)
}
