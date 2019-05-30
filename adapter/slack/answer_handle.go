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
