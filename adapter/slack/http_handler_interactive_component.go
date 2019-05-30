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

package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hori-ryota/zaperr"
	"github.com/opensaasstudio/meerkat/domain"
	"go.uber.org/zap"
	"gopkg.in/guregu/null.v3"
)

type Callback struct {
	Token   string
	Actions []CallbackAction

	Container struct {
		ChannelID string `json:"channel_id"`
		MessageTS string `json:"message_ts"`
	}
}

type CallbackAction struct {
	ActionID       string `json:"action_id"`
	BlockID        string `json:"block_id"`
	Value          string `json:"value"`
	SelectedOption struct {
		Value string
	} `json:"selected_option"`
	SelectedDate         string `json:"selected_date"`
	SelectedUser         string `json:"selected_user"`
	SelectedConversation string `json:"selected_conversation"`
	SelectedChannel      string `json:"selected_channel"`
}

func (h HTTPHandler) HandleInteractiveComponent(w http.ResponseWriter, r *http.Request) {
	err := func() error {
		ctx := r.Context()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		defer r.Body.Close()
		payload, err := url.QueryUnescape(string(body[len("payload="):]))
		if err != nil {
			return err
		}

		var callback Callback

		if err := json.Unmarshal([]byte(payload), &callback); err != nil {
			return err
		}
		if callback.Token != h.slackVerificationToken {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}
		h.logger.Info("requested", zap.Any("callback", callback))

		for _, action := range callback.Actions {
			ai := strings.SplitN(action.ActionID, "__", 2)
			if len(ai) < 2 {
				return domain.ErrorBadRequest(zaperr.New("invalid actionID " + action.ActionID))
			}

			callbackID := ai[0]
			actionName := ai[1]
			switch strings.SplitN(callbackID, "_", 2)[0] {
			case "EditingQuestionnaire":
				input := EditingQuestionnaireHandlerInput{}
				if err := h.paramStore.Restore(context.TODO(), callbackID, &input); err != nil {
					return err
				}
				if actionName == "fix" {
					editted, err := h.editingQuestionnaireHandler.Execute(ctx, input)
					if err != nil {
						return err
					}
					return h.editingQuestionnaireHandler.PrintFixed(
						ctx,
						callback.Container.ChannelID,
						null.StringFrom(callback.Container.MessageTS),
						input,
						string(editted.ID()),
					)
				}
				value := parseValue(action)

				input, err = h.editingQuestionnaireHandler.HandleEditingQuestionnaire(ctx, input, actionName, value)
				if err != nil {
					return err
				}
				if err := h.editingQuestionnaireHandler.RequestInput(
					ctx,
					callback.Container.ChannelID,
					null.StringFrom(callback.Container.MessageTS),
					callbackID,
					input,
				); err != nil {
					return err
				}
				if err := h.paramStore.Store(context.TODO(), callbackID, input, 30*time.Minute); err != nil {
					return err
				}
			case "CreatingAnswerer":
				input := CreatingAnswererHandlerInput{}
				if err := h.paramStore.Restore(context.TODO(), callbackID, &input); err != nil {
					return err
				}
				if actionName == "fix" {
					if err := h.creatingAnswererHandler.Execute(ctx, input); err != nil {
						return err
					}
					return h.creatingAnswererHandler.PrintFixed(
						ctx,
						callback.Container.ChannelID,
						null.StringFrom(callback.Container.MessageTS),
						input,
					)
				}
				value := parseValue(action)

				input, err = h.creatingAnswererHandler.HandleCreatingAnswerer(ctx, input, actionName, value)
				if err != nil {
					return err
				}
				if err := h.creatingAnswererHandler.RequestInput(
					ctx,
					callback.Container.ChannelID,
					null.StringFrom(callback.Container.MessageTS),
					callbackID,
					input,
				); err != nil {
					return err
				}
				if err := h.paramStore.Store(context.TODO(), callbackID, input, 30*time.Minute); err != nil {
					return err
				}
			case "AddingAnswerer":
				input := AddingAnswererHandlerInput{}
				if err := h.paramStore.Restore(context.TODO(), callbackID, &input); err != nil {
					return err
				}
				if actionName == "fix" {
					if err := h.addingAnswererHandler.Execute(ctx, input); err != nil {
						return err
					}
					return h.addingAnswererHandler.PrintFixed(
						ctx,
						callback.Container.ChannelID,
						null.StringFrom(callback.Container.MessageTS),
						input,
					)
				}
				value := parseValue(action)

				input, err = h.addingAnswererHandler.HandleAddingAnswerer(ctx, input, actionName, value)
				if err != nil {
					return err
				}
				if err := h.addingAnswererHandler.RequestInput(
					ctx,
					callback.Container.ChannelID,
					null.StringFrom(callback.Container.MessageTS),
					callbackID,
					input,
				); err != nil {
					return err
				}
				if err := h.paramStore.Store(context.TODO(), callbackID, input, 30*time.Minute); err != nil {
					return err
				}
			case "Answering":
				input := AnsweringHandlerInput{}
				if err := h.paramStore.Restore(context.TODO(), callbackID, &input); err != nil {
					return err
				}
				if actionName == "fix" {
					if err := h.answeringHandler.Execute(ctx, input); err != nil {
						return err
					}
					return h.answeringHandler.PrintFixed(
						ctx,
						callback.Container.ChannelID,
						null.StringFrom(callback.Container.MessageTS),
						input,
					)
				}
				value := parseValue(action)

				input, err = h.answeringHandler.HandleAnswering(ctx, input, actionName, value)
				if err != nil {
					return err
				}
				if err := h.answeringHandler.RequestInput(
					ctx,
					callback.Container.ChannelID,
					null.StringFrom(callback.Container.MessageTS),
					callbackID,
					input,
				); err != nil {
					return err
				}
				if err := h.paramStore.Store(context.TODO(), callbackID, input, 30*time.Minute); err != nil {
					return err
				}
			}
		}
		return nil
	}()
	if err != nil {
		h.logger.Error("error", zaperr.ToField(err))
		w.WriteHeader(http.StatusInternalServerError)
		h.handleError(w, r, err)
		fmt.Fprintf(w, "err: %+v", err)
	}
}

func parseValue(action CallbackAction) string {
	switch {
	case action.Value != "":
		return action.Value
	case action.SelectedOption.Value != "":
		return action.SelectedOption.Value
	case action.SelectedDate != "":
		return action.SelectedDate
	case action.SelectedUser != "":
		return action.SelectedUser
	case action.SelectedConversation != "":
		return action.SelectedConversation
	case action.SelectedChannel != "":
		return action.SelectedChannel
	}
	return ""
}
