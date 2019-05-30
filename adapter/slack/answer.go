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

package slack

import (
	"context"

	"github.com/nlopes/slack"
	"github.com/opensaasstudio/meerkat/application"
	"github.com/opensaasstudio/meerkat/domain"
)

//genconstructor
type AnsweringHandler struct {
	slackClient *slack.Client                `required:""`
	usecase     application.AnsweringUsecase `required:""`
}

type Answer struct {
	Question Question
	Value    string
}

type AnsweringHandlerInput struct {
	QuestionnaireID    string
	QuestionnaireTitle string
	AnswererID         string
	Answers            []Answer
}

func (p AnsweringHandlerInput) ToUsecaseInput() application.AnsweringUsecaseInput {
	answers := make([]application.AnswerInputValue, len(p.Answers))
	for i := range p.Answers {
		answers[i] = application.NewAnswerInputValue(
			domain.QuestionID(p.Answers[i].Question.ID),
			p.Answers[i].Value,
		)
	}
	return application.NewAnsweringUsecaseInput(
		domain.QuestionnaireID(p.QuestionnaireID),
		domain.AnswererID(p.AnswererID),
		answers,
	)
}

func (h AnsweringHandler) Execute(
	ctx context.Context,
	input AnsweringHandlerInput,
) domain.Error {
	return h.usecase.Answer(
		ctx,
		input.ToUsecaseInput(),
	)
}
