package slack

import (
	"context"
	"time"

	"github.com/nlopes/slack"
	"github.com/opensaasstudio/meerkat/application"
	"github.com/opensaasstudio/meerkat/domain"
)

//genconstructor
type CreatingQuestionnaireHandler struct {
	slackClient *slack.Client                            `required:""`
	usecase     application.CreatingQuestionnaireUsecase `required:""`
}

type Question struct {
	ID       string
	Text     string
	Required bool
}

type CreatingQuestionnaireHandlerInput struct {
	Title       string
	Questions   []Question
	Schedules   Schedules
	PostTargets PostTargets
}

func (p CreatingQuestionnaireHandlerInput) ToUsecaseInput() application.CreatingQuestionnaireUsecaseInput {
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

// func (p *CreatingQuestionnaireHandlerInput) FromUsecaseInput(up application.CreatingQuestionnaireUsecaseInput) {
// 	p.Questions = make([]Question, len(up.QuestionItems()))
// 	for i, q := range up.QuestionItems() {
// 		p.Questions[i] = Question{
// 			ID:       string(q.Question().ID()),
// 			Text:     q.Question().Text(),
// 			Required: q.Required(),
// 		}
// 	}
// 	p.Title = up.Title()
// }

func (h CreatingQuestionnaireHandler) Execute(
	ctx context.Context,
	input CreatingQuestionnaireHandlerInput,
) domain.Error {
	_, err := h.usecase.Create(
		ctx,
		application.AdminDescriptor{},     // TODO
		application.WorkspaceDescriptor{}, // TODO
		input.ToUsecaseInput(),
	)
	return err
}
