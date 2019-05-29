package slack

import (
	"context"
	"fmt"
	"strconv"

	"github.com/nlopes/slack"
	"github.com/opensaasstudio/meerkat/domain"
	"gopkg.in/guregu/null.v3"
)

func (h EditingQuestionnaireHandler) RequestInput(
	ctx context.Context,
	channelID string,
	updateTargetID null.String,
	callbackID string,
	input EditingQuestionnaireHandlerInput,
) domain.Error {
	blocks := make([]slack.Block, 0, 10)

	dividerBlock := slack.NewDividerBlock()
	plainText := func(text string) *slack.TextBlockObject {
		return slack.NewTextBlockObject("plain_text", text, false, false)
	}
	primaryIfTrue := func(button *slack.ButtonBlockElement, b bool) *slack.ButtonBlockElement {
		if b {
			button.WithStyle(slack.StylePrimary)
		}
		return button
	}

	applyInitialOption := func(block *slack.SelectBlockElement, value string) *slack.SelectBlockElement {
		if value == "" {
			return block
		}
		block.InitialOption = slack.NewOptionBlockObject(value, plainText(value))
		return block
	}
	applyInitialDate := func(block *slack.DatePickerBlockElement, value string) *slack.DatePickerBlockElement {
		block.InitialDate = value
		return block
	}
	applyInitialConversation := func(block *slack.SelectBlockElement, value string) *slack.SelectBlockElement {
		if value == "" {
			return block
		}
		block.InitialConversation = value
		return block
	}

	blocks = append(blocks, slack.NewActionBlock(
		callbackID+"__title",
		applyInitialOption(slack.NewOptionsSelectBlockElement(
			slack.OptTypeExternal,
			plainText("input title here"),
			callbackID+"__title",
		), input.Title),
	))

	blocks = append(blocks, dividerBlock)

	for i := range input.Questions {
		actionIDPrefix := callbackID + "__" + "question_" + strconv.Itoa(i)
		blocks = append(blocks, slack.NewActionBlock(
			callbackID+"__"+actionIDPrefix,
			applyInitialOption(slack.NewOptionsSelectBlockElement(
				slack.OptTypeExternal,
				plainText("input question"),
				actionIDPrefix+"_text",
			), input.Questions[i].Text),
			slack.NewButtonBlockElement(actionIDPrefix+"_moveup", "", plainText("↑")),
			slack.NewButtonBlockElement(actionIDPrefix+"_movedown", "", plainText("↓")),
			slack.NewButtonBlockElement(actionIDPrefix+"_remove", "", plainText("×")),
			primaryIfTrue(slack.NewButtonBlockElement(actionIDPrefix+"_togglerequired", "", plainText("required")), input.Questions[i].Required),
		))
	}

	blocks = append(blocks, slack.NewActionBlock(
		callbackID+"__appendquestion",
		slack.NewButtonBlockElement(callbackID+"__appendquestion", "", plainText("append question")),
	))

	for i, s := range input.Schedules.WeekdayAndTimeSchedules {
		actionIDPrefix := callbackID + "__" + "schedule_" + strconv.Itoa(i) + "_weekdayandtime"
		blocks = append(blocks, dividerBlock)
		blocks = append(blocks, slack.NewSectionBlock(plainText("WeekdayAndTimeSchedule"), nil, slack.NewAccessory(
			slack.NewButtonBlockElement(actionIDPrefix+"_remove", "", plainText("×")),
		)))
		blocks = append(blocks, slack.NewSectionBlock(plainText("Hour"), nil, slack.NewAccessory(
			applyInitialOption(slack.NewOptionsSelectBlockElement(
				slack.OptTypeStatic,
				plainText("hour"),
				actionIDPrefix+"_hour",
				NumberOptions(0, 23)...,
			), strconv.Itoa(int(s.Hour))),
		)))
		blocks = append(blocks, slack.NewSectionBlock(plainText("Minute"), nil, slack.NewAccessory(
			applyInitialOption(slack.NewOptionsSelectBlockElement(
				slack.OptTypeStatic,
				plainText("minute"),
				actionIDPrefix+"_minute",
				NumberOptions(0, 59)...,
			), strconv.Itoa(int(s.Minute))),
		)))
		blocks = append(blocks, slack.NewSectionBlock(plainText("Sec"), nil, slack.NewAccessory(
			applyInitialOption(slack.NewOptionsSelectBlockElement(
				slack.OptTypeStatic,
				plainText("sec"),
				actionIDPrefix+"_sec",
				NumberOptions(0, 59)...,
			), strconv.Itoa(int(s.Sec))),
		)))
		blocks = append(blocks, slack.NewSectionBlock(plainText("Timezone"), nil, slack.NewAccessory(
			applyInitialOption(slack.NewOptionsSelectBlockElement(
				slack.OptTypeStatic,
				plainText("timezone"),
				actionIDPrefix+"_timezone",
				NumberOptions(-12, 13)...,
			), strconv.Itoa(int(s.Timezone))),
		)))
		blocks = append(blocks, slack.NewActionBlock(
			callbackID+"__"+actionIDPrefix,
			primaryIfTrue(slack.NewButtonBlockElement(actionIDPrefix+"_mon", "", plainText("Mon")), s.Mon),
			primaryIfTrue(slack.NewButtonBlockElement(actionIDPrefix+"_tue", "", plainText("Tue")), s.Tue),
			primaryIfTrue(slack.NewButtonBlockElement(actionIDPrefix+"_wed", "", plainText("Wed")), s.Wed),
			primaryIfTrue(slack.NewButtonBlockElement(actionIDPrefix+"_thu", "", plainText("Thu")), s.Thu),
			primaryIfTrue(slack.NewButtonBlockElement(actionIDPrefix+"_fri", "", plainText("Fri")), s.Fri),
			primaryIfTrue(slack.NewButtonBlockElement(actionIDPrefix+"_sat", "", plainText("Sat")), s.Sat),
			primaryIfTrue(slack.NewButtonBlockElement(actionIDPrefix+"_sun", "", plainText("Sun")), s.Sun),
		))
	}

	for i, s := range input.Schedules.YearMonthDayScheduleExceptions {
		actionIDPrefix := callbackID + "__" + "scheduleexception_" + strconv.Itoa(i) + "_yearmonthday"
		blocks = append(blocks, dividerBlock)
		blocks = append(blocks, slack.NewSectionBlock(plainText("YearMonthDayScheduleException"), nil, slack.NewAccessory(
			slack.NewButtonBlockElement(actionIDPrefix+"_remove", "", plainText("×")),
		)))
		blocks = append(blocks, slack.NewSectionBlock(plainText("YearMonthDay"), nil, slack.NewAccessory(
			applyInitialDate(slack.NewDatePickerBlockElement(
				actionIDPrefix+"_yearmonthday",
			), fmt.Sprintf("%04d-%02d-%02d", s.Year, s.Month, s.Day)),
		)))
		blocks = append(blocks, slack.NewSectionBlock(plainText("Timezone"), nil, slack.NewAccessory(
			applyInitialOption(slack.NewOptionsSelectBlockElement(
				slack.OptTypeStatic,
				plainText("timezone"),
				actionIDPrefix+"_timezone",
				NumberOptions(-12, 13)...,
			), strconv.Itoa(int(s.Timezone))),
		)))
	}

	blocks = append(blocks, dividerBlock)
	blocks = append(blocks, slack.NewActionBlock(
		callbackID+"__appendschedule",
		slack.NewOptionsSelectBlockElement(
			slack.OptTypeStatic,
			plainText("append schedule"),
			callbackID+"__appendschedule",
			slack.NewOptionBlockObject("weekdayandtime", plainText("WeekdayAndTimeSchedule")),
		),
		slack.NewOptionsSelectBlockElement(
			slack.OptTypeStatic,
			plainText("append schedule exception"),
			callbackID+"__appendscheduleexception",
			slack.NewOptionBlockObject("yearmonthday", plainText("YearMonthDayScheduleException")),
		),
	))

	for i, s := range input.PostTargets.SlackPostTargets {
		actionIDPrefix := callbackID + "__" + "posttarget_" + strconv.Itoa(i) + "_slack"
		blocks = append(blocks, dividerBlock)
		blocks = append(blocks, slack.NewSectionBlock(plainText("Slack"), nil, slack.NewAccessory(
			slack.NewButtonBlockElement(actionIDPrefix+"_remove", "", plainText("×")),
		)))
		blocks = append(blocks, slack.NewSectionBlock(plainText("Select Post Channel"), nil, slack.NewAccessory(
			applyInitialConversation(slack.NewOptionsSelectBlockElement(
				slack.OptTypeConversations,
				plainText("select post channel"),
				actionIDPrefix+"_channelid",
			), s.ChannelID),
		)))
	}

	blocks = append(blocks, dividerBlock)
	blocks = append(blocks, slack.NewActionBlock(
		callbackID+"__appendposttarget",
		slack.NewOptionsSelectBlockElement(
			slack.OptTypeStatic,
			plainText("append postTarget"),
			callbackID+"__appendposttarget",
			slack.NewOptionBlockObject("slack", plainText("Slack")),
		),
	))

	filled := true
	if input.Title == "" {
		filled = false
	}
	if len(input.Questions) == 0 {
		filled = false
	}
	for _, q := range input.Questions {
		if q.Text == "" {
			filled = false
			break
		}
	}

	if len(input.Schedules.WeekdayAndTimeSchedules) == 0 {
		filled = false
	}

	for _, t := range input.PostTargets.SlackPostTargets {
		if t.ChannelID == "" {
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
	// TODO cancel
	// button := slack.NewButtonBlockElement(callbackID+"__cancel", "", plainText("cancel"))
	// button.WithStyle(slack.StyleDanger)
	// button.Confirm = slack.NewConfirmationBlockObject(
	// 	plainText("ok?"),
	// 	plainText("ok?"),
	// 	plainText("ok"),
	// 	plainText("no"),
	// )
	// blocks = append(blocks, slack.NewActionBlock(
	// 	callbackID+"__cancel",
	// 	button,
	// ))

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
