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

	"github.com/nlopes/slack"
	"github.com/opensaasstudio/meerkat/domain"
	"gopkg.in/guregu/null.v3"
)

func (h AnsweringHandler) RequestInput(
	ctx context.Context,
	channelID string,
	updateTargetID null.String,
	callbackID string,
	input AnsweringHandlerInput,
) domain.Error {
	blocks := make([]slack.Block, 0, 10)

	dividerBlock := slack.NewDividerBlock()
	plainText := func(text string) *slack.TextBlockObject {
		return slack.NewTextBlockObject("plain_text", text, false, false)
	}

	applyInitialOption := func(block *slack.SelectBlockElement, value string) *slack.SelectBlockElement {
		if value == "" {
			return block
		}
		block.InitialOption = slack.NewOptionBlockObject(value, plainText(value))
		return block
	}

	blocks = append(blocks, slack.NewSectionBlock(plainText(input.QuestionnaireTitle), nil, nil))
	blocks = append(blocks, slack.NewContextBlock(callbackID+"__questionnaireID", plainText("questionnaireID: "+input.QuestionnaireID)))

	blocks = append(blocks, dividerBlock)

	for i := range input.Answers {
		actionIDPrefix := callbackID + "__" + "answer_" + strconv.Itoa(i)
		text := input.Answers[i].Question.Text
		if input.Answers[i].Question.Required {
			text += "[required]"
		}
		blocks = append(blocks, slack.NewSectionBlock(plainText(text), nil, slack.NewAccessory(
			applyInitialOption(slack.NewOptionsSelectBlockElement(
				slack.OptTypeExternal,
				plainText("input answer"),
				actionIDPrefix+"_value",
			), input.Answers[i].Value),
		)))
	}

	filled := true
	for _, answer := range input.Answers {
		if answer.Question.Required && answer.Value == "" {
			filled = false
			break
		}
	}

	if filled {
		blocks = append(blocks, dividerBlock)
		button := slack.NewButtonBlockElement(callbackID+"__fix", "", plainText("fix"))
		button.WithStyle(slack.StylePrimary)
		button.Confirm = slack.NewConfirmationBlockObject(
			plainText("ok?"),
			plainText("ok?"),
			plainText("submit!"),
			plainText("cancel"),
		)
		blocks = append(blocks, slack.NewActionBlock(
			callbackID+"__fix",
			button,
		))
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
