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
