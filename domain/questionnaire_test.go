package domain_test

import (
	"errors"
	"testing"

	"github.com/opensaasstudio/meerkat/domain"
	"github.com/stretchr/testify/assert"
)

func TestQuestionnaire(t *testing.T) {
	t.Run("Answer", func(t *testing.T) {
		q := domain.NewQuestionnaire(
			"questionnaireID",
			"questionnaireTitle",
			[]domain.QuestionItem{
				domain.NewQuestionItem(
					domain.NewQuestion(
						"question1",
						"question1 text",
					),
					false,
				),
			},
		)
		answer, err := q.Answer("answer1", "question1", "answererID", "value")
		assert.NoError(t, err)
		assert.Equal(t, domain.QuestionnaireID("questionnaireID"), answer.QuestionnaireID())
		assert.Equal(t, domain.QuestionID("question1"), answer.QuestionID())
		assert.Equal(t, domain.AnswererID("answererID"), answer.AnswererID())
		assert.Equal(t, "value", answer.Value())

		t.Run("unknown questionID erorr", func(t *testing.T) {
			question := domain.NewQuestion(
				"question1",
				"question1 text",
			)
			question.SetValidatingFunc(func(currentText string) domain.Error {
				return domain.ErrorUnknown(errors.New("validation error"))
			})

			q := domain.NewQuestionnaire(
				"questionnaireID",
				"questionnaireTitle",
				[]domain.QuestionItem{
					domain.NewQuestionItem(
						question,
						false,
					),
				},
			)
			_, err := q.Answer("answer1", "unknown questionID", "answererID", "value")
			assert.Error(t, err)
		})

		t.Run("validation erorr", func(t *testing.T) {
			question := domain.NewQuestion(
				"question1",
				"question1 text",
			)
			question.SetValidatingFunc(func(currentText string) domain.Error {
				return domain.ErrorUnknown(errors.New("validation error"))
			})

			q := domain.NewQuestionnaire(
				"questionnaireID",
				"questionnaireTitle",
				[]domain.QuestionItem{
					domain.NewQuestionItem(
						question,
						false,
					),
				},
			)
			_, err := q.Answer("answer1", "question1", "answererID", "value")
			assert.Error(t, err)
		})
	})
}
