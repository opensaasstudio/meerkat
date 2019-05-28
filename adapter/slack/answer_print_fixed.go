package slack

import (
	"context"

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
