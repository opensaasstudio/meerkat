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
