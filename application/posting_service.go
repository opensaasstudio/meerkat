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
type PostingService struct {
	postSlack func(ctx context.Context, target domain.PostTargetSlack, questionnaire domain.Questionnaire, answerer domain.Answerer, answers []domain.Answer) domain.Error `required:""`
}

func (s PostingService) PostAnswers(ctx context.Context, questionnaire domain.Questionnaire, answerer domain.Answerer, answers []domain.Answer) domain.Error {
	for _, postTarget := range questionnaire.PostTargets() {
		switch pt := postTarget.(type) {
		case domain.PostTargetSlack:
			return s.postSlack(ctx, pt, questionnaire, answerer, answers)
		default:
			return domain.ErrorUnknown(zaperr.New("unknown posting target type"))
		}
	}
	return nil
}
