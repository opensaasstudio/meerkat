package slack

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/nlopes/slack"
	"github.com/opensaasstudio/meerkat/domain"
	"gopkg.in/guregu/null.v3"
)

func (h EditingQuestionnaireHandler) PrintFixed(
	ctx context.Context,
	channelID string,
	updateTargetID null.String,
	input EditingQuestionnaireHandlerInput,
	questionnaireID string,
) domain.Error {
	blocks := make([]slack.Block, 0, 10)

	dividerBlock := slack.NewDividerBlock()

	plainText := func(text string) *slack.TextBlockObject {
		return slack.NewTextBlockObject("plain_text", text, false, false)
	}

	blocks = append(blocks, slack.NewSectionBlock(plainText("questionnaire is created: "+input.Title), nil, nil))
	blocks = append(blocks, slack.NewContextBlock(strconv.FormatInt(time.Now().UnixNano(), 10)+"__questionnaireID", plainText("questionnaireID: "+questionnaireID)))

	blocks = append(blocks, dividerBlock)

	for i := range input.Questions {
		text := input.Questions[i].Text
		if input.Questions[i].Required {
			text = text + " [required]"
		}
		blocks = append(blocks, slack.NewSectionBlock(plainText(text), nil, nil))
	}
	for _, s := range input.Schedules.WeekdayAndTimeSchedules {
		blocks = append(blocks, dividerBlock)
		text := fmt.Sprintf("WeekdayAndTimeSchedule: %02d:%02d:%02d (tz=%02d)", s.Hour, s.Minute, s.Sec, s.Timezone)
		if s.Mon {
			text += ", Mon"
		}
		if s.Tue {
			text += ", Tue"
		}
		if s.Wed {
			text += ", Wed"
		}
		if s.Thu {
			text += ", Thu"
		}
		if s.Fri {
			text += ", Fri"
		}
		if s.Sat {
			text += ", Sat"
		}
		if s.Sun {
			text += ", Sun"
		}
		blocks = append(blocks, slack.NewSectionBlock(plainText(text), nil, nil))
	}
	for _, t := range input.PostTargets.SlackPostTargets {
		blocks = append(blocks, dividerBlock)
		text := fmt.Sprintf("Slack PostTarget: %s", t.ChannelID)
		blocks = append(blocks, slack.NewSectionBlock(plainText(text), nil, nil))
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
