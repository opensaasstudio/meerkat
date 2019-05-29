package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/hori-ryota/zaperr"
	"github.com/nlopes/slack/slackevents"
	"github.com/opensaasstudio/meerkat/domain"
	"go.uber.org/zap"
	"gopkg.in/guregu/null.v3"
)

func (h HTTPHandler) HandleEvent(w http.ResponseWriter, r *http.Request) {
	err := func() error {
		ctx := r.Context()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		defer r.Body.Close()

		event, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: h.slackVerificationToken}))
		if err != nil {
			return err
		}
		h.logger.Info("requested", zap.Any("event", event))

		switch eventData := event.Data.(type) {
		case *slackevents.EventsAPIURLVerificationEvent:
			if _, err := w.Write([]byte(eventData.Challenge)); err != nil {
				return err
			}
			return nil
		case *slackevents.EventsAPICallbackEvent:
			switch innerEvent := event.InnerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				h.logger.Info("mentioned", zap.Any("innerEvent", innerEvent))
				text := innerEvent.Text
				if strings.HasPrefix(text, "Reminder: ") {
					text = strings.TrimPrefix(text, "Reminder: ")
					text = strings.TrimSuffix(text, ".")
				}
				texts := strings.Split(text, " ")
				if len(texts) == 1 {
					return nil
				}
				commandName := texts[1]
				switch commandName {
				case "createQuestionnaire":
					callbackID := "EditingQuestionnaire_" + innerEvent.EventTimeStamp.String()
					input := EditingQuestionnaireHandlerInput{}
					if err := h.editingQuestionnaireHandler.RequestInput(
						ctx, innerEvent.Channel, null.String{}, callbackID, input,
					); err != nil {
						return err
					}
				case "editQuestionnaire":
					callbackID := "EditingQuestionnaire_" + innerEvent.EventTimeStamp.String()
					input := EditingQuestionnaireHandlerInput{}

					if len(texts) < 3 {
						return nil
					}
					questionnaireID := texts[2]
					existing, _, err := h.questionnaireSearcher.FindByID(ctx, domain.QuestionnaireID(questionnaireID))
					if err != nil {
						return err
					}
					input.FromDomainObject(existing)
					if err := h.paramStore.Store(context.TODO(), callbackID, input, 30*time.Minute); err != nil {
						return err
					}

					if err := h.editingQuestionnaireHandler.RequestInput(
						ctx, innerEvent.Channel, null.String{}, callbackID, input,
					); err != nil {
						return err
					}
				case "createAnswerer":
					callbackID := "CreatingAnswerer_" + innerEvent.EventTimeStamp.String()
					input := CreatingAnswererHandlerInput{}
					if err := h.creatingAnswererHandler.RequestInput(
						ctx, innerEvent.Channel, null.String{}, callbackID, input,
					); err != nil {
						return err
					}
				case "addAnswerer":
					callbackID := "AddingAnswerer_" + innerEvent.EventTimeStamp.String()
					input := AddingAnswererHandlerInput{}
					if err := h.addingAnswererHandler.RequestInput(
						ctx, innerEvent.Channel, null.String{}, callbackID, input,
					); err != nil {
						return err
					}
				}
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
