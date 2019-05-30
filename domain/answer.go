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

import "time"

type AnswerID string

//genconstructor
type Answer struct {
	id              AnswerID        `required:"" getter:""`
	questionnaireID QuestionnaireID `required:"" getter:""`
	questionID      QuestionID      `required:"" getter:""`
	answererID      AnswererID      `required:"" getter:""`
	answeredAt      time.Time       `required:"" getter:""`
	value           string          `required:"" getter:"" setter:"Overwrite"`
}

type AnswerValue struct {
	ID              AnswerID `dynamo:",hash"`
	QuestionnaireID QuestionnaireID
	QuestionID      QuestionID
	AnswererID      AnswererID
	AnsweredAt      time.Time
	Value           string
}

func (m Answer) Dump() AnswerValue {
	return AnswerValue{
		ID:              m.ID(),
		QuestionnaireID: m.QuestionnaireID(),
		QuestionID:      m.QuestionID(),
		AnswererID:      m.AnswererID(),
		AnsweredAt:      m.AnsweredAt(),
		Value:           m.Value(),
	}
}

func RestoreAnswerFromDumped(v AnswerValue) Answer {
	return NewAnswer(
		v.ID,
		v.QuestionnaireID,
		v.QuestionID,
		v.AnswererID,
		v.AnsweredAt,
		v.Value,
	)
}
