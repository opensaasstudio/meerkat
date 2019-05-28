package slack

import (
	"context"
	"fmt"

	"github.com/nlopes/slack"
	"github.com/opensaasstudio/meerkat/application"
	"github.com/opensaasstudio/meerkat/domain"
	"gopkg.in/guregu/null.v3"
)

//genconstructor
type CreatingAnswererHandler struct {
	slackClient *slack.Client                       `required:""`
	usecase     application.CreatingAnswererUsecase `required:""`
}

type CreatingAnswererHandlerInput struct {
	Name string
}

func (p CreatingAnswererHandlerInput) ToUsecaseInput() application.CreatingAnswererUsecaseInput {
	return application.NewCreatingAnswererUsecaseInput(p.Name)
}

func (h CreatingAnswererHandler) HandleCreatingAnswerer(
	ctx context.Context,
	input CreatingAnswererHandlerInput,
	actionName string,
	value string,
) (CreatingAnswererHandlerInput, domain.Error) {
	switch {
	case actionName == "name":
		input.Name = value
		return input, nil
	default:
		return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
	}
	return input, nil
}

func (h CreatingAnswererHandler) RequestInput(
	ctx context.Context,
	channelID string,
	updateTargetID null.String,
	callbackID string,
	input CreatingAnswererHandlerInput,
) domain.Error {
	blocks := make([]slack.Block, 0, 10)

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

	blocks = append(blocks, slack.NewActionBlock(
		callbackID+"__name",
		applyInitialOption(slack.NewOptionsSelectBlockElement(
			slack.OptTypeExternal,
			plainText("input name here"),
			callbackID+"__name",
		), input.Name),
	))

	if input.Name != "" {
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

func (h CreatingAnswererHandler) PrintFixed(
	ctx context.Context,
	channelID string,
	updateTargetID null.String,
	input CreatingAnswererHandlerInput,
) domain.Error {
	blocks := make([]slack.Block, 0, 10)

	plainText := func(text string) *slack.TextBlockObject {
		return slack.NewTextBlockObject("plain_text", text, false, false)
	}

	blocks = append(blocks, slack.NewSectionBlock(plainText("answerer is created: "+input.Name), nil, nil))

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

func (h CreatingAnswererHandler) Execute(
	ctx context.Context,
	input CreatingAnswererHandlerInput,
) domain.Error {
	_, err := h.usecase.CreateAnswerer(
		ctx,
		input.ToUsecaseInput(),
	)
	return err
}
