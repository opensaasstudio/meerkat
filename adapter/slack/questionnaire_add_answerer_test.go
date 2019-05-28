package slack_test

import (
	"context"
	"testing"

	"github.com/opensaasstudio/meerkat/adapter/slack"
	"github.com/opensaasstudio/meerkat/application"
	"github.com/stretchr/testify/assert"
)

func TestAddingAnswererHandler(t *testing.T) {
	t.Run("HandleAddingAnswerer", func(t *testing.T) {
		for _, tt := range []struct {
			name       string
			input      slack.AddingAnswererHandlerInput
			actionName string
			value      string
			assert     func(t *testing.T, input slack.AddingAnswererHandlerInput)
		}{
			{
				name:       "input questionnaireID",
				actionName: "questionnaireid",
				value:      "questionnaire id",
				assert: func(t *testing.T, input slack.AddingAnswererHandlerInput) {
					assert.Equal(t, "questionnaire id", input.QuestionnaireID)
				},
			},
			{
				name:       "input answererID",
				actionName: "answererid",
				value:      "answerer id",
				assert: func(t *testing.T, input slack.AddingAnswererHandlerInput) {
					assert.Equal(t, "answerer id", input.AnswererID)
				},
			},
			{
				name:       "input channelID",
				actionName: "channelid",
				value:      "channel id",
				assert: func(t *testing.T, input slack.AddingAnswererHandlerInput) {
					assert.Equal(t, "channel id", input.ChannelID)
				},
			},
			{
				name:       "input userID",
				actionName: "userid",
				value:      "user id",
				assert: func(t *testing.T, input slack.AddingAnswererHandlerInput) {
					assert.Equal(t, "user id", input.UserID)
				},
			},
			{
				name:       "toggle needsMention false -> true",
				actionName: "needsmention",
				value:      "",
				assert: func(t *testing.T, input slack.AddingAnswererHandlerInput) {
					assert.True(t, input.NeedsMention)
				},
			},
			{
				name:       "toggle needsMention true -> false",
				input:      slack.AddingAnswererHandlerInput{NeedsMention: true},
				actionName: "needsmention",
				value:      "",
				assert: func(t *testing.T, input slack.AddingAnswererHandlerInput) {
					assert.False(t, input.NeedsMention)
				},
			},
		} {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				ctx := context.Background()
				h := slack.NewAddingAnswererHandler(nil, application.AddingAnswererUsecase{})

				got, err := h.HandleAddingAnswerer(
					ctx,
					tt.input,
					tt.actionName,
					tt.value,
				)
				assert.NoError(t, err)
				tt.assert(t, got)
			})
		}
	})
}
