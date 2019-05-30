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
	"github.com/opensaasstudio/meerkat/domain"
	"go.uber.org/zap"
)

//genconstructor
type HTTPHandler struct {
	slackVerificationToken      string                       `required:""`
	logger                      *zap.Logger                  `required:""`
	questionnaireSearcher       domain.QuestionnaireSearcher `required:""`
	answererSearcher            domain.AnswererSearcher      `required:""`
	editingQuestionnaireHandler EditingQuestionnaireHandler  `required:""`
	creatingAnswererHandler     CreatingAnswererHandler      `required:""`
	addingAnswererHandler       AddingAnswererHandler        `required:""`
	answeringHandler            AnsweringHandler             `required:""`
	paramStore                  ParamStore                   `required:""`
}
