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
type RemovingAnswererUsecase struct {
	searcher   domain.NotificationTargetSearcher `required:""`
	repository NotificationTargetRepository      `required:""`
}

//genconstructor
type RemovingAnswererUsecaseInput struct {
	questionnaireID domain.QuestionnaireID `required:"" getter:""`
	answererID      domain.AnswererID      `required:"" getter:""`
}

func (u RemovingAnswererUsecase) RemoveAnswerer(ctx context.Context, input RemovingAnswererUsecaseInput) domain.Error {
	notificationTargets, err := u.searcher.SearchByQuestionnaireIDAndAnswererID(
		ctx,
		input.QuestionnaireID(),
		input.AnswererID(),
	)
	if err != nil {
		return err
	}
	for _, notificationTarget := range notificationTargets {
		if err := u.repository.Delete(ctx, notificationTarget); err != nil {
			return err
		}
	}
	return nil
}
