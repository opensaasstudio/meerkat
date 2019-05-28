package slack

import (
	"context"
	"fmt"
	"strconv"

	"github.com/nlopes/slack"
	"github.com/opensaasstudio/meerkat/application"
	"github.com/opensaasstudio/meerkat/domain"
	"gopkg.in/guregu/null.v3"
)

//genconstructor
type AddingAnswererHandler struct {
	slackClient *slack.Client                     `required:""`
	usecase     application.AddingAnswererUsecase `required:""`
}

type AddingAnswererHandlerInput struct {
	QuestionnaireID string
	AnswererID      string
	ChannelID       string
	UserID          string
	NeedsMention    bool
}

func (p AddingAnswererHandlerInput) ToUsecaseInput() application.AddingAnswererUsecaseInput {
	t := domain.NewNotificationTargetSlack(
		"",
		domain.QuestionnaireID(p.QuestionnaireID),
		domain.AnswererID(p.AnswererID),
		p.ChannelID,
		p.UserID,
	)
	t.ToggleNeedsMention(p.NeedsMention)
	return application.NewAddingAnswererUsecaseInput(t)
}

func (h AddingAnswererHandler) HandleAddingAnswerer(
	ctx context.Context,
	input AddingAnswererHandlerInput,
	actionName string,
	value string,
) (AddingAnswererHandlerInput, domain.Error) {
	switch {
	case actionName == "questionnaireid":
		input.QuestionnaireID = value
		return input, nil
	case actionName == "answererid":
		input.AnswererID = value
		return input, nil
	case actionName == "channelid":
		input.ChannelID = value
		return input, nil
	case actionName == "userid":
		input.UserID = value
		return input, nil
	case actionName == "needsmention":
		input.NeedsMention = !input.NeedsMention
		return input, nil
	default:
		return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
	}
	return input, nil
}

func (h AddingAnswererHandler) RequestInput(
	ctx context.Context,
	channelID string,
	updateTargetID null.String,
	callbackID string,
	input AddingAnswererHandlerInput,
) domain.Error {
	blocks := make([]slack.Block, 0, 10)

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
	applyInitialUser := func(block *slack.SelectBlockElement, value string) *slack.SelectBlockElement {
		if value == "" {
			return block
		}
		block.InitialUser = value
		return block
	}
	applyInitialConversation := func(block *slack.SelectBlockElement, value string) *slack.SelectBlockElement {
		if value == "" {
			return block
		}
		block.InitialConversation = value
		return block
	}

	blocks = append(blocks, slack.NewSectionBlock(plainText("Select Questionnaire"), nil, slack.NewAccessory(
		applyInitialOption(slack.NewOptionsSelectBlockElement(
			slack.OptTypeExternal,
			plainText("select questionnaire"),
			callbackID+"__questionnaireid",
		), input.QuestionnaireID),
	)))
	blocks = append(blocks, slack.NewSectionBlock(plainText("Select Answerer"), nil, slack.NewAccessory(
		applyInitialOption(slack.NewOptionsSelectBlockElement(
			slack.OptTypeExternal,
			plainText("select answerer"),
			callbackID+"__answererid",
		), input.AnswererID),
	)))
	blocks = append(blocks, slack.NewSectionBlock(plainText("Select Notification Channel"), nil, slack.NewAccessory(
		applyInitialConversation(slack.NewOptionsSelectBlockElement(
			slack.OptTypeConversations,
			plainText("select notification channel"),
			callbackID+"__channelid",
		), input.ChannelID),
	)))
	blocks = append(blocks, slack.NewSectionBlock(plainText("Select Notification User"), nil, slack.NewAccessory(
		applyInitialUser(slack.NewOptionsSelectBlockElement(
			slack.OptTypeUser,
			plainText("select notification user"),
			callbackID+"__userid",
		), input.UserID),
	)))
	blocks = append(blocks, slack.NewSectionBlock(plainText("Toggle NeedsMention"), nil, slack.NewAccessory(
		primaryIfTrue(slack.NewButtonBlockElement(
			callbackID+"__needsmention",
			"",
			plainText("NeedsMention"),
		), input.NeedsMention),
	)))

	filled := true
	if input.QuestionnaireID == "" {
		filled = false
	}
	if input.AnswererID == "" {
		filled = false
	}
	if input.ChannelID == "" {
		filled = false
	}
	if input.UserID == "" {
		filled = false
	}

	if filled {
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

func (h AddingAnswererHandler) PrintFixed(
	ctx context.Context,
	channelID string,
	updateTargetID null.String,
	input AddingAnswererHandlerInput,
) domain.Error {
	blocks := make([]slack.Block, 0, 10)

	plainText := func(text string) *slack.TextBlockObject {
		return slack.NewTextBlockObject("plain_text", text, false, false)
	}

	blocks = append(blocks, slack.NewSectionBlock(plainText("answerer is added:"), nil, nil))
	blocks = append(blocks, slack.NewSectionBlock(plainText("questionnaireID: "+input.QuestionnaireID), nil, nil))
	blocks = append(blocks, slack.NewSectionBlock(plainText("answererID: "+input.AnswererID), nil, nil))
	blocks = append(blocks, slack.NewSectionBlock(plainText("channelID: "+input.ChannelID), nil, nil))
	blocks = append(blocks, slack.NewSectionBlock(plainText("userID: "+input.UserID), nil, nil))
	blocks = append(blocks, slack.NewSectionBlock(plainText("needsMention: "+strconv.FormatBool(input.NeedsMention)), nil, nil))

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

func (h AddingAnswererHandler) Execute(
	ctx context.Context,
	input AddingAnswererHandlerInput,
) domain.Error {
	return h.usecase.AddAnswerer(
		ctx,
		input.ToUsecaseInput(),
	)
}
