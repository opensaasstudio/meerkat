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

type QuestionID string

//genconstructor
type Question struct {
	id             QuestionID                        `required:"" getter:""`
	text           string                            `required:"" getter:"" setter:"Overwrite"`
	suggestingFunc func(currentText string) []string `getter:"" setter:""` // TODO 設定できるようにする
	validatingFunc func(currentText string) Error    `getter:"" setter:""` // TODO 設定できるようにする
}

type QuestionValue struct {
	ID   QuestionID
	Text string
}

func (m Question) Dump() QuestionValue {
	return QuestionValue{
		ID:   m.ID(),
		Text: m.Text(),
	}
}

func RestoreQuestionFromDumped(v QuestionValue) Question {
	return NewQuestion(v.ID, v.Text)
}
