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

//go:generate mockgen -source $GOFILE -destination mock_$GOPACKAGE/${GOFILE}_mock.go

import (
	"context"
	"time"
)

//genconstructor
type AskingService struct {
	questionnaireSearcher      QuestionnaireSearcher      `required:""`
	notificationTargetSearcher NotificationTargetSearcher `required:""`
	notificationService        NotificationService        `required:""`
	lastExecutedRecorder       LastExecutedRecorder       `required:""`
	nowFunc                    func() time.Time           `required:"time.Now" setter:"OverwriteNowFunc"`
}

func (s AskingService) AskAllIfNeeded(ctx context.Context) Error {
	questionnaires, err := s.questionnaireSearcher.SearchExecutionNeeded(ctx)
	if err != nil {
		return err
	}
	for _, q := range questionnaires {
		if err := s.AskQuestionnaire(ctx, q); err != nil {
			return err
		}
	}
	return nil
}

func (s AskingService) AskQuestionnaire(ctx context.Context, questionnaire Questionnaire) Error {
	nts, err := s.notificationTargetSearcher.SearchByQuestionnaireID(ctx, questionnaire.ID())
	if err != nil {
		return err
	}
	for _, nt := range nts {
		err := s.notificationService.Notify(ctx, nt, questionnaire)
		if err != nil {
			return err
		}
	}
	return s.lastExecutedRecorder.RecordLastExecuted(ctx, questionnaire, s.nowFunc())
}

type NotificationService interface {
	Notify(ctx context.Context, notificationTarget NotificationTarget, questionnaire Questionnaire) Error
}

type LastExecutedRecorder interface {
	RecordLastExecuted(ctx context.Context, questionnaire Questionnaire, lastExecuted time.Time) Error
}
