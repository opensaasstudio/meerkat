package slack_test

import (
	"context"
	"testing"

	"github.com/opensaasstudio/meerkat/adapter/slack"
	"github.com/opensaasstudio/meerkat/application"
	"github.com/stretchr/testify/assert"
)

func TestCreatingQuestionnaireHandler(t *testing.T) {
	t.Run("HandleCreatingQuestionnaire", func(t *testing.T) {
		for _, tt := range []struct {
			name       string
			input      slack.CreatingQuestionnaireHandlerInput
			actionName string
			value      string
			assert     func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput)
		}{
			{
				name:       "input title",
				actionName: "title",
				value:      "questionnaire title",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.Equal(t, "questionnaire title", input.Title)
				},
			},
			{
				name:       "input new question",
				actionName: "question_0_text",
				value:      "question_0 text",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.Equal(t, "question_0 text", input.Questions[0].Text)
				},
			},
			{
				name:       "append new question",
				actionName: "appendquestion",
				value:      "",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.Len(t, input.Questions, 1)
				},
			},
			{
				name: "input existing question",
				input: slack.CreatingQuestionnaireHandlerInput{
					Questions: []slack.Question{
						{Text: "question_0 text"},
						{Text: "question_1 text"},
						{Text: "question_2 text"},
					},
				},
				actionName: "question_1_text",
				value:      "new question_1 text",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.Equal(t, "new question_1 text", input.Questions[1].Text)
				},
			},
			{
				name: "remove existing question",
				input: slack.CreatingQuestionnaireHandlerInput{
					Questions: []slack.Question{
						{Text: "question_0 text"},
						{Text: "question_1 text"},
						{Text: "question_2 text"},
					},
				},
				actionName: "question_1_remove",
				value:      "",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					if assert.Len(t, input.Questions, 2) {
						assert.Equal(t, "question_0 text", input.Questions[0].Text)
						assert.Equal(t, "question_2 text", input.Questions[1].Text)
					}
				},
			},
			{
				name: "move up question",
				input: slack.CreatingQuestionnaireHandlerInput{
					Questions: []slack.Question{
						{Text: "question_0 text"},
						{Text: "question_1 text"},
						{Text: "question_2 text"},
					},
				},
				actionName: "question_1_moveup",
				value:      "",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.Equal(t, "question_1 text", input.Questions[0].Text)
					assert.Equal(t, "question_0 text", input.Questions[1].Text)
				},
			},
			{
				name: "move down question",
				input: slack.CreatingQuestionnaireHandlerInput{
					Questions: []slack.Question{
						{Text: "question_0 text"},
						{Text: "question_1 text"},
						{Text: "question_2 text"},
					},
				},
				actionName: "question_1_movedown",
				value:      "",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.Equal(t, "question_2 text", input.Questions[1].Text)
					assert.Equal(t, "question_1 text", input.Questions[2].Text)
				},
			},
			{
				name: "toggle required on question",
				input: slack.CreatingQuestionnaireHandlerInput{
					Questions: []slack.Question{
						{Text: "question_0 text"},
					},
				},
				actionName: "question_0_togglerequired",
				value:      "",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.True(t, input.Questions[0].Required)
				},
			},
			{
				name: "toggle required off question",
				input: slack.CreatingQuestionnaireHandlerInput{
					Questions: []slack.Question{
						{Text: "question_0 text", Required: true},
					},
				},
				actionName: "question_0_togglerequired",
				value:      "",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.False(t, input.Questions[0].Required)
				},
			},
			{
				name:       "append WeekdayAndTimeSchedule",
				actionName: "appendschedule",
				value:      "weekdayandtime",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.NotEmpty(t, input.Schedules.WeekdayAndTimeSchedules)
				},
			},
			{
				name: "WeekdayAndTimeSchedule input hour",
				input: slack.CreatingQuestionnaireHandlerInput{
					Schedules: slack.Schedules{WeekdayAndTimeSchedules: make([]slack.WeekdayAndTimeSchedule, 1)},
				},
				actionName: "schedule_0_weekdayandtime_hour",
				value:      "1",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.Equal(t, uint32(1), input.Schedules.WeekdayAndTimeSchedules[0].Hour)
				},
			},
			{
				name: "WeekdayAndTimeSchedule input minute",
				input: slack.CreatingQuestionnaireHandlerInput{
					Schedules: slack.Schedules{WeekdayAndTimeSchedules: make([]slack.WeekdayAndTimeSchedule, 1)},
				},
				actionName: "schedule_0_weekdayandtime_minute",
				value:      "1",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.Equal(t, uint32(1), input.Schedules.WeekdayAndTimeSchedules[0].Minute)
				},
			},
			{
				name: "WeekdayAndTimeSchedule input sec",
				input: slack.CreatingQuestionnaireHandlerInput{
					Schedules: slack.Schedules{WeekdayAndTimeSchedules: make([]slack.WeekdayAndTimeSchedule, 1)},
				},
				actionName: "schedule_0_weekdayandtime_sec",
				value:      "1",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.Equal(t, uint32(1), input.Schedules.WeekdayAndTimeSchedules[0].Sec)
				},
			},
			{
				name: "WeekdayAndTimeSchedule input timezone",
				input: slack.CreatingQuestionnaireHandlerInput{
					Schedules: slack.Schedules{WeekdayAndTimeSchedules: make([]slack.WeekdayAndTimeSchedule, 1)},
				},
				actionName: "schedule_0_weekdayandtime_timezone",
				value:      "1",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.Equal(t, 1, input.Schedules.WeekdayAndTimeSchedules[0].Timezone)
				},
			},
			{
				name: "WeekdayAndTimeSchedule toggle mon: false -> true",
				input: slack.CreatingQuestionnaireHandlerInput{
					Schedules: slack.Schedules{WeekdayAndTimeSchedules: make([]slack.WeekdayAndTimeSchedule, 1)},
				},
				actionName: "schedule_0_weekdayandtime_mon",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.True(t, input.Schedules.WeekdayAndTimeSchedules[0].Mon)
				},
			},
			{
				name: "WeekdayAndTimeSchedule toggle mon: true -> false",
				input: slack.CreatingQuestionnaireHandlerInput{
					Schedules: slack.Schedules{WeekdayAndTimeSchedules: []slack.WeekdayAndTimeSchedule{{Mon: true}}},
				},
				actionName: "schedule_0_weekdayandtime_mon",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.False(t, input.Schedules.WeekdayAndTimeSchedules[0].Mon)
				},
			},
			{
				name: "WeekdayAndTimeSchedule remove",
				input: slack.CreatingQuestionnaireHandlerInput{
					Schedules: slack.Schedules{WeekdayAndTimeSchedules: []slack.WeekdayAndTimeSchedule{
						{Mon: true, Tue: true},
						{Mon: true},
					}},
				},
				actionName: "schedule_0_weekdayandtime_remove",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.Len(t, input.Schedules.WeekdayAndTimeSchedules, 1)
					assert.False(t, input.Schedules.WeekdayAndTimeSchedules[0].Tue)
				},
			},
			{
				name:       "append YearMonthDayScheduleException",
				actionName: "appendscheduleexception",
				value:      "yearmonthday",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.NotEmpty(t, input.Schedules.YearMonthDayScheduleExceptions)
				},
			},
			{
				name: "YearMonthDayScheduleException input yearMonthDay",
				input: slack.CreatingQuestionnaireHandlerInput{
					Schedules: slack.Schedules{YearMonthDayScheduleExceptions: make([]slack.YearMonthDayScheduleException, 1)},
				},
				actionName: "scheduleexception_0_yearmonthday_yearmonthday",
				value:      "2019-05-20",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.Equal(t, uint32(2019), input.Schedules.YearMonthDayScheduleExceptions[0].Year)
					assert.Equal(t, uint32(5), input.Schedules.YearMonthDayScheduleExceptions[0].Month)
					assert.Equal(t, uint32(20), input.Schedules.YearMonthDayScheduleExceptions[0].Day)
				},
			},
			{
				name: "YearMonthDayScheduleException input timezone",
				input: slack.CreatingQuestionnaireHandlerInput{
					Schedules: slack.Schedules{YearMonthDayScheduleExceptions: make([]slack.YearMonthDayScheduleException, 1)},
				},
				actionName: "scheduleexception_0_yearmonthday_timezone",
				value:      "1",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.Equal(t, 1, input.Schedules.YearMonthDayScheduleExceptions[0].Timezone)
				},
			},
			{
				name: "YearMonthDayScheduleException remove",
				input: slack.CreatingQuestionnaireHandlerInput{
					Schedules: slack.Schedules{YearMonthDayScheduleExceptions: []slack.YearMonthDayScheduleException{
						{Year: 2019},
						{Year: 2018},
					}},
				},
				actionName: "scheduleexception_0_yearmonthday_remove",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.Len(t, input.Schedules.YearMonthDayScheduleExceptions, 1)
					assert.Equal(t, uint32(2018), input.Schedules.YearMonthDayScheduleExceptions[0].Year)
				},
			},
			{
				name:       "append SlackPostTarget",
				actionName: "appendposttarget",
				value:      "slack",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.NotEmpty(t, input.PostTargets.SlackPostTargets)
				},
			},
			{
				name: "SlackPostTarget input channelid",
				input: slack.CreatingQuestionnaireHandlerInput{
					PostTargets: slack.PostTargets{SlackPostTargets: make([]slack.SlackPostTarget, 1)},
				},
				actionName: "posttarget_0_slack_channelid",
				value:      "channel1",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.Equal(t, "channel1", input.PostTargets.SlackPostTargets[0].ChannelID)
				},
			},
			{
				name: "SlackPostTarget remove",
				input: slack.CreatingQuestionnaireHandlerInput{
					PostTargets: slack.PostTargets{SlackPostTargets: []slack.SlackPostTarget{
						{ChannelID: "channel1"},
						{ChannelID: "channel2"},
					}},
				},
				actionName: "posttarget_0_slack_remove",
				assert: func(t *testing.T, input slack.CreatingQuestionnaireHandlerInput) {
					assert.Len(t, input.PostTargets.SlackPostTargets, 1)
					assert.Equal(t, "channel2", input.PostTargets.SlackPostTargets[0].ChannelID)
				},
			},
		} {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				ctx := context.Background()
				h := slack.NewCreatingQuestionnaireHandler(nil, application.CreatingQuestionnaireUsecase{})

				got, err := h.HandleCreatingQuestionnaire(
					ctx,
					tt.input,
					tt.actionName,
					tt.value,
				)
				assert.NoError(t, err)
				tt.assert(t, got)
			})
		}
	})
}
