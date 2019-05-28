package slack

import (
	"context"
	"fmt"

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
