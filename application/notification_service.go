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

	"github.com/hori-ryota/zaperr"
	"github.com/opensaasstudio/meerkat/domain"
)

//genconstructor
type NotificationService struct {
	notifySlack func(ctx context.Context, target domain.NotificationTargetSlack, questionnaire domain.Questionnaire) domain.Error `required:""`
}

func (s NotificationService) Notify(ctx context.Context, notificationTarget domain.NotificationTarget, questionnaire domain.Questionnaire) domain.Error {
	switch nt := notificationTarget.(type) {
	case domain.NotificationTargetSlack:
		return s.notifySlack(ctx, nt, questionnaire)
	}
	return domain.ErrorUnknown(zaperr.New("unknown notification target type"))
}
