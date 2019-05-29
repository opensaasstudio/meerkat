package slack

import (
	"context"
	"time"

	"github.com/nlopes/slack"
	"github.com/opensaasstudio/meerkat/application"
	"github.com/opensaasstudio/meerkat/domain"
)

//genconstructor
type EditingQuestionnaireHandler struct {
	slackClient                  *slack.Client                            `required:""`
	creatingQuestionnaireUsecase application.CreatingQuestionnaireUsecase `required:""`
	updatingQuestionnaireUsecase application.UpdatingQuestionnaireUsecase `required:""`
}

type Question struct {
	ID       string
	Text     string
	Required bool
}

type EditingQuestionnaireHandlerInput struct {
	ID          string
	Title       string
	Questions   []Question
	Schedules   Schedules
	PostTargets PostTargets
}

func (p EditingQuestionnaireHandlerInput) ToCreatingUsecaseInput() application.CreatingQuestionnaireUsecaseInput {
	questions := make([]application.QuestionItem, len(p.Questions))
	for i, q := range p.Questions {
		questions[i] = application.NewQuestionItem(
			application.NewQuestion(q.Text), q.Required,
		)
	}

	scheduleList := make([]domain.Schedule, len(p.Schedules.WeekdayAndTimeSchedules))
	for i, s := range p.Schedules.WeekdayAndTimeSchedules {
		scheduleList[i] = domain.NewWeekdayAndTimeSchedule(
			s.Hour, s.Minute, s.Sec, s.Timezone,
			s.Mon, s.Tue, s.Wed, s.Thu, s.Fri, s.Sat, s.Sun,
		)
	}

	scheduleExceptionList := make([]domain.ScheduleException, len(p.Schedules.YearMonthDayScheduleExceptions))
	for i, s := range p.Schedules.YearMonthDayScheduleExceptions {
		scheduleExceptionList[i] = domain.NewYearMonthDayScheduleException(
			s.Year, time.Month(s.Month), s.Day, s.Timezone,
		)
	}

	schedule := domain.NewSchedules(scheduleList[0], scheduleList[1:], scheduleExceptionList)

	postTarget := make([]domain.PostTarget, len(p.PostTargets.SlackPostTargets))
	for i, s := range p.PostTargets.SlackPostTargets {
		postTarget[i] = domain.NewPostTargetSlack("", s.ChannelID)
	}

	input := application.NewCreatingQuestionnaireUsecaseInput(
		p.Title,
		questions,
	)
	input.SetSchedule(schedule)
	input.SetPostTargets(postTarget)
	return input
}

func (p EditingQuestionnaireHandlerInput) ToUpdatingUsecaseInput() application.UpdatingQuestionnaireUsecaseInput {
	return application.NewUpdatingQuestionnaireUsecaseInput(
		domain.QuestionnaireID(p.ID),
		p.ToCreatingUsecaseInput(),
	)
}

func (p *EditingQuestionnaireHandlerInput) FromDomainObject(questionnaire domain.Questionnaire) {
	p.ID = string(questionnaire.ID())
	p.Title = questionnaire.Title()
	p.Questions = make([]Question, len(questionnaire.QuestionItems()))
	for i, q := range questionnaire.QuestionItems() {
		p.Questions[i] = Question{
			ID:       string(q.Question().ID()),
			Text:     q.Question().Text(),
			Required: q.Required(),
		}
	}
	p.Schedules = RestoreSchedulesFromDomainObject(questionnaire.Schedule())
	for i := range questionnaire.PostTargets() {
		p.PostTargets = p.PostTargets.Merge(RestorePostTargetFromDomainObject(questionnaire.PostTargets()[i]))
	}
}

func (h EditingQuestionnaireHandler) Execute(
	ctx context.Context,
	input EditingQuestionnaireHandlerInput,
) (domain.Questionnaire, domain.Error) {
	if input.ID == "" {
		return h.creatingQuestionnaireUsecase.Create(
			ctx,
			application.AdminDescriptor{},     // TODO
			application.WorkspaceDescriptor{}, // TODO
			input.ToCreatingUsecaseInput(),
		)
	}
	return h.updatingQuestionnaireUsecase.Update(
		ctx,
		application.AdminDescriptor{},     // TODO
		application.WorkspaceDescriptor{}, // TODO
		input.ToUpdatingUsecaseInput(),
	)
}
