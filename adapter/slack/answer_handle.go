package slack

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/opensaasstudio/meerkat/domain"
)

func (h AnsweringHandler) HandleAnswering(
	ctx context.Context,
	input AnsweringHandlerInput,
	actionName string,
	value string,
) (AnsweringHandlerInput, domain.Error) {
	switch {
	case strings.HasPrefix(actionName, "answer_"):
		// e.g. answer_0_value
		ss := strings.SplitN(actionName, "_", 3)
		if len(ss) < 3 {
			return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
		}
		index, err := strconv.Atoi(ss[1])
		if err != nil {
			return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
		}
		switch ss[2] {
		case "value":
			if index < len(input.Answers) {
				input.Answers[index].Value = value
				return input, nil
			}
		default:
			return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
		}
	default:
		return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
	}
	return input, nil
}
