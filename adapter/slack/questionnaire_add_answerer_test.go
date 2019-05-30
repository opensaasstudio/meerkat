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
