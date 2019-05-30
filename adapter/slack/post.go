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
	"fmt"
	"strconv"
	"time"

	"github.com/nlopes/slack"
	"github.com/opensaasstudio/meerkat/domain"
)

//genconstructor
type Poster struct {
	slackClient *slack.Client `required:""`
}

func (s Poster) Post(ctx context.Context, target domain.PostTargetSlack, questionnaire domain.Questionnaire, answerer domain.Answerer, answers []domain.Answer) domain.Error {

	blocks := make([]slack.Block, 0, 10)

	dividerBlock := slack.NewDividerBlock()

	plainText := func(text string) *slack.TextBlockObject {
		return slack.NewTextBlockObject("plain_text", text, false, false)
	}

	blocks = append(blocks, slack.NewSectionBlock(plainText(fmt.Sprintf(
		"Answered by %s: %s",
		answerer.Name(),
		questionnaire.Title(),
	)), nil, nil))
	blocks = append(blocks, slack.NewContextBlock(strconv.FormatInt(time.Now().UnixNano(), 10)+"__questionnaireID", plainText("questionnaireID: "+string(questionnaire.ID()))))

	for i, answer := range answers {
		blocks = append(blocks, dividerBlock)

		blocks = append(blocks, slack.NewSectionBlock(plainText(questionnaire.QuestionItems()[i].Question().Text()+":"), nil, nil))
		blocks = append(blocks, slack.NewSectionBlock(plainText(answer.Value()), nil, nil))
	}

	options := make([]slack.MsgOption, 0, 1)
	options = append(options, slack.MsgOptionBlocks(blocks...))

	if _, _, err := s.slackClient.PostMessageContext(
		ctx,
		target.ChannelID(),
		options...,
	); err != nil {
		return domain.ErrorUnknown(err)
	}
	return nil
}
