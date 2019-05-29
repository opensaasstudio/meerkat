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
