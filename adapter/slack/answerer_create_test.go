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

func TestCreatingAnswererHandler(t *testing.T) {
	t.Run("HandleCreatingAnswerer", func(t *testing.T) {
		for _, tt := range []struct {
			name       string
			input      slack.CreatingAnswererHandlerInput
			actionName string
			value      string
			assert     func(t *testing.T, input slack.CreatingAnswererHandlerInput)
		}{
			{
				name:       "input name",
				actionName: "name",
				value:      "answerer name",
				assert: func(t *testing.T, input slack.CreatingAnswererHandlerInput) {
					assert.Equal(t, "answerer name", input.Name)
				},
			},
		} {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				ctx := context.Background()
				h := slack.NewCreatingAnswererHandler(nil, application.CreatingAnswererUsecase{})

				got, err := h.HandleCreatingAnswerer(
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
