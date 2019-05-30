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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/hori-ryota/zaperr"
	"github.com/nlopes/slack"
	"go.uber.org/zap"
)

type SuggestionRequestPayload struct {
	Token    string
	ActionID string `json:"action_id"`
	Value    string
}

func (h HTTPHandler) HandleSuggestion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := func() error {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		defer r.Body.Close()
		payload, err := url.QueryUnescape(string(body[len("payload="):]))
		if err != nil {
			return err
		}
		h.logger.Info("suggestion", zap.String("payload", payload))

		var requestBody SuggestionRequestPayload
		if err := json.Unmarshal([]byte(payload), &requestBody); err != nil {
			return err
		}
		if requestBody.Token != h.slackVerificationToken {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}

		w.Header().Add("Content-type", "application/json")
		switch {
		case strings.HasSuffix(requestBody.ActionID, "questionnaireid"):
			questionnaires, err := h.questionnaireSearcher.FetchAll(ctx)
			if err != nil {
				return err
			}
			options := make([]*slack.OptionBlockObject, len(questionnaires))
			for i, q := range questionnaires {
				options[i] = slack.NewOptionBlockObject(string(q.ID()), slack.NewTextBlockObject("plain_text", q.Title()+":"+string(q.ID()), true, false))
			}
			if err := json.NewEncoder(w).Encode(map[string][]*slack.OptionBlockObject{
				"options": options,
			}); err != nil {
				return err // TODO いい感じのレスポンス
			}
		case strings.HasSuffix(requestBody.ActionID, "answererid"):
			answerers, err := h.answererSearcher.FetchAll(ctx)
			if err != nil {
				return err
			}
			options := make([]*slack.OptionBlockObject, len(answerers))
			for i, q := range answerers {
				options[i] = slack.NewOptionBlockObject(string(q.ID()), slack.NewTextBlockObject("plain_text", q.Name()+":"+string(q.ID()), true, false))
			}
			if err := json.NewEncoder(w).Encode(map[string][]*slack.OptionBlockObject{
				"options": options,
			}); err != nil {
				return err // TODO いい感じのレスポンス
			}
		default:
			if err := json.NewEncoder(w).Encode(map[string][]*slack.OptionBlockObject{
				"options": []*slack.OptionBlockObject{
					slack.NewOptionBlockObject(requestBody.Value, slack.NewTextBlockObject("plain_text", requestBody.Value, true, false)),
				},
			}); err != nil {
				return err // TODO いい感じのレスポンス
			}
		}
		return nil
	}()
	if err != nil {
		h.logger.Error("error", zaperr.ToField(err))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "err: %+v", err)
	}
}
