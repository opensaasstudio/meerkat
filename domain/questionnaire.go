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
	"errors"
	"time"
)

type QuestionnaireID string

//genconstructor
type Questionnaire struct {
	id            QuestionnaireID `required:"" getter:""`
	title         string          `required:"" getter:"" setter:"OverwriteTitle"`
	questionItems []QuestionItem  `required:"" getter:"" setter:"ReplaceQuestions"`
	schedule      Schedule        `getter:"" setter:""`
	postTargets   []PostTarget    `getter:"" setter:""`
	lastExecuted  time.Time       `getter:"" setter:""`
	// TODO Authors
}

//genconstructor
type QuestionItem struct {
	question Question `required:"" getter:""`
	required bool     `required:"" getter:""`
}

func (q Questionnaire) Answer(answerID AnswerID, questionID QuestionID, answererID AnswererID, value string) (Answer, Error) {
	questionItem, err := q.chooseQuestion(questionID)
	if err != nil {
		return Answer{}, err
	}
	if questionItem.Question().ValidatingFunc() != nil {
		if err := questionItem.Question().ValidatingFunc()(value); err != nil {
			return Answer{}, err
		}
	}
	return NewAnswer(
		answerID,
		q.ID(),
		questionID,
		answererID,
		time.Now(),
		value,
	), nil
}

func (q Questionnaire) chooseQuestion(questionID QuestionID) (QuestionItem, Error) {
	for _, item := range q.QuestionItems() {
		if item.Question().ID() == questionID {
			return item, nil
		}
	}
	//errcode QuestionIsNotFound,questionID QuestionID
	return QuestionItem{}, ErrorBadRequest(errors.New("question is not found"), QuestionIsNotFoundError(questionID))
}

type QuestionItemValue struct {
	Question QuestionValue
	Required bool
}

type QuestionnaireValue struct {
	ID            QuestionnaireID `dynamo:",hash"`
	Title         string
	QuestionItems []QuestionItemValue
	Schedule      ScheduleValue
	PostTargets   []PostTargetValue
	LastExecuted  time.Time
	NextTime      time.Time
	NextTimeUnix  int64
}

func (m Questionnaire) Dump() QuestionnaireValue {
	questionItems := make([]QuestionItemValue, len(m.QuestionItems()))
	for i, q := range m.QuestionItems() {
		questionItems[i] = QuestionItemValue{
			Question: q.Question().Dump(),
			Required: q.Required(),
		}
	}
	postTargets := make([]PostTargetValue, len(m.PostTargets()))
	for i, t := range m.PostTargets() {
		postTargets[i] = t.Dump()
	}
	return QuestionnaireValue{
		ID:            m.ID(),
		Title:         m.Title(),
		QuestionItems: questionItems,
		Schedule:      m.Schedule().Dump(),
		PostTargets:   postTargets,
		LastExecuted:  m.LastExecuted(),
		NextTime:      m.Schedule().NextTime(m.LastExecuted()),
		NextTimeUnix:  m.Schedule().NextTime(m.LastExecuted()).Unix(),
	}
}

func RestoreQuestionnaireFromDumped(v QuestionnaireValue) Questionnaire {
	questionItems := make([]QuestionItem, len(v.QuestionItems))
	for i, q := range v.QuestionItems {
		questionItems[i] = NewQuestionItem(
			RestoreQuestionFromDumped(q.Question),
			q.Required,
		)
	}
	postTargets := make([]PostTarget, len(v.PostTargets))
	for i, t := range v.PostTargets {
		postTargets[i] = RestorePostTargetFromDumpled(t)
	}
	m := NewQuestionnaire(
		v.ID,
		v.Title,
		questionItems,
	)
	m.SetSchedule(RestoreScheduleFromDumped(v.Schedule))
	m.SetPostTargets(postTargets)
	m.SetLastExecuted(v.LastExecuted)
	return m
}
