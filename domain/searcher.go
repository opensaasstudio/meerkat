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

package domain

import (
	"context"
)

//go:generate mockgen -source $GOFILE -destination mock_$GOPACKAGE/${GOFILE}_mock.go

type NotificationTargetSearcher interface {
	SearchByQuestionnaireIDAndAnswererID(ctx context.Context, questionnaireID QuestionnaireID, answererID AnswererID) ([]NotificationTarget, Error)
	SearchByQuestionnaireID(ctx context.Context, questionnaireID QuestionnaireID) ([]NotificationTarget, Error)
}

type QuestionnaireSearcher interface {
	SearchExecutionNeeded(ctx context.Context) ([]Questionnaire, Error)
	FetchAll(ctx context.Context) ([]Questionnaire, Error)
	FindByID(ctx context.Context, id QuestionnaireID) (questionnaire Questionnaire, version int, derr Error)
}

type AnswererSearcher interface {
	FetchAll(ctx context.Context) ([]Answerer, Error)
}
