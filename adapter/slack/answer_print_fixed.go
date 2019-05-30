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
	"strconv"
	"time"

	"github.com/nlopes/slack"
	"github.com/opensaasstudio/meerkat/domain"
	"gopkg.in/guregu/null.v3"
)

func (h AnsweringHandler) PrintFixed(
	ctx context.Context,
	channelID string,
	updateTargetID null.String,
	input AnsweringHandlerInput,
) domain.Error {
	blocks := make([]slack.Block, 0, 10)

	dividerBlock := slack.NewDividerBlock()

	plainText := func(text string) *slack.TextBlockObject {
		return slack.NewTextBlockObject("plain_text", text, false, false)
	}

	blocks = append(blocks, slack.NewSectionBlock(plainText("answered: "+input.QuestionnaireTitle), nil, nil))
	blocks = append(blocks, slack.NewContextBlock(strconv.FormatInt(time.Now().UnixNano(), 10)+"__questionnaireID", plainText("questionnaireID: "+input.QuestionnaireID)))

	for _, answer := range input.Answers {
		blocks = append(blocks, dividerBlock)

		blocks = append(blocks, slack.NewSectionBlock(plainText(answer.Question.Text+":"), nil, nil))
		blocks = append(blocks, slack.NewSectionBlock(plainText(answer.Value), nil, nil))
	}

	options := make([]slack.MsgOption, 0, 2)
	options = append(options, slack.MsgOptionBlocks(blocks...))
	if updateTargetID.Valid {
		options = append(options, slack.MsgOptionUpdate(updateTargetID.ValueOrZero()))
	}

	if _, _, err := h.slackClient.PostMessageContext(
		ctx,
		channelID,
		options...,
	); err != nil {
		return domain.ErrorUnknown(err)
	}
	return nil
}
