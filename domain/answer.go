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
