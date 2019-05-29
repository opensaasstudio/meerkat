package slack

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/opensaasstudio/meerkat/domain"
	"github.com/pkg/errors"
)

func (h EditingQuestionnaireHandler) HandleEditingQuestionnaire(
	ctx context.Context,
	input EditingQuestionnaireHandlerInput,
	actionName string,
	value string,
) (EditingQuestionnaireHandlerInput, domain.Error) {
	switch {
	case actionName == "title":
		input.Title = value
		return input, nil
	case actionName == "appendquestion":
		input.Questions = append(input.Questions, Question{})
		return input, nil
	case strings.HasPrefix(actionName, "question_"):
		// e.g. question_0_text
		ss := strings.SplitN(actionName, "_", 3)
		if len(ss) < 3 {
			return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
		}
		index, err := strconv.Atoi(ss[1])
		if err != nil {
			return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
		}
		switch ss[2] {
		case "text":
			if index < len(input.Questions) {
				input.Questions[index].Text = value
				return input, nil
			}
			input.Questions = append(input.Questions, Question{Text: value})
			return input, nil
		case "remove":
			if index >= len(input.Questions) {
				return input, nil
			}
			input.Questions = append(input.Questions[:index], input.Questions[index+1:]...)
			return input, nil
		case "moveup":
			if index >= len(input.Questions) || index == 0 {
				return input, nil
			}
			input.Questions[index-1], input.Questions[index] = input.Questions[index], input.Questions[index-1]
			return input, nil
		case "movedown":
			if index >= len(input.Questions)-1 {
				return input, nil
			}
			input.Questions[index], input.Questions[index+1] = input.Questions[index+1], input.Questions[index]
			return input, nil
		case "togglerequired":
			if index >= len(input.Questions) {
				return input, nil
			}
			input.Questions[index].Required = !input.Questions[index].Required
			return input, nil
		default:
			return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
		}
	case strings.HasPrefix(actionName, "schedule_"):
		// e.g. schedule_0_weekdayandtime_hour
		ss := strings.SplitN(actionName, "_", 4)
		if len(ss) < 4 {
			return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
		}
		index, err := strconv.Atoi(ss[1])
		if err != nil {
			return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
		}
		switch ss[2] {
		case "weekdayandtime":
			switch ss[3] {
			case "hour":
				v, err := strconv.Atoi(value)
				if err != nil {
					return input, domain.ErrorBadRequest(fmt.Errorf("unknown value %s", value))
				}
				if index < len(input.Schedules.WeekdayAndTimeSchedules) {
					input.Schedules.WeekdayAndTimeSchedules[index].Hour = uint32(v)
					return input, nil
				}
				return input, nil
			case "minute":
				v, err := strconv.Atoi(value)
				if err != nil {
					return input, domain.ErrorBadRequest(fmt.Errorf("unknown value %s", value))
				}
				if index < len(input.Schedules.WeekdayAndTimeSchedules) {
					input.Schedules.WeekdayAndTimeSchedules[index].Minute = uint32(v)
					return input, nil
				}
				return input, nil
			case "sec":
				v, err := strconv.Atoi(value)
				if err != nil {
					return input, domain.ErrorBadRequest(fmt.Errorf("unknown value %s", value))
				}
				if index < len(input.Schedules.WeekdayAndTimeSchedules) {
					input.Schedules.WeekdayAndTimeSchedules[index].Sec = uint32(v)
					return input, nil
				}
				return input, nil
			case "timezone":
				v, err := strconv.Atoi(value)
				if err != nil {
					return input, domain.ErrorBadRequest(fmt.Errorf("unknown value %s", value))
				}
				if index < len(input.Schedules.WeekdayAndTimeSchedules) {
					input.Schedules.WeekdayAndTimeSchedules[index].Timezone = v
					return input, nil
				}
				return input, nil
			case "mon":
				if index < len(input.Schedules.WeekdayAndTimeSchedules) {
					input.Schedules.WeekdayAndTimeSchedules[index].Mon = !input.Schedules.WeekdayAndTimeSchedules[index].Mon
					return input, nil
				}
				return input, nil
			case "tue":
				if index < len(input.Schedules.WeekdayAndTimeSchedules) {
					input.Schedules.WeekdayAndTimeSchedules[index].Tue = !input.Schedules.WeekdayAndTimeSchedules[index].Tue
					return input, nil
				}
				return input, nil
			case "wed":
				if index < len(input.Schedules.WeekdayAndTimeSchedules) {
					input.Schedules.WeekdayAndTimeSchedules[index].Wed = !input.Schedules.WeekdayAndTimeSchedules[index].Wed
					return input, nil
				}
				return input, nil
			case "thu":
				if index < len(input.Schedules.WeekdayAndTimeSchedules) {
					input.Schedules.WeekdayAndTimeSchedules[index].Thu = !input.Schedules.WeekdayAndTimeSchedules[index].Thu
					return input, nil
				}
				return input, nil
			case "fri":
				if index < len(input.Schedules.WeekdayAndTimeSchedules) {
					input.Schedules.WeekdayAndTimeSchedules[index].Fri = !input.Schedules.WeekdayAndTimeSchedules[index].Fri
					return input, nil
				}
				return input, nil
			case "sat":
				if index < len(input.Schedules.WeekdayAndTimeSchedules) {
					input.Schedules.WeekdayAndTimeSchedules[index].Sat = !input.Schedules.WeekdayAndTimeSchedules[index].Sat
					return input, nil
				}
				return input, nil
			case "sun":
				if index < len(input.Schedules.WeekdayAndTimeSchedules) {
					input.Schedules.WeekdayAndTimeSchedules[index].Sun = !input.Schedules.WeekdayAndTimeSchedules[index].Sun
					return input, nil
				}
				return input, nil
			case "remove":
				if index >= len(input.Schedules.WeekdayAndTimeSchedules) {
					return input, nil
				}
				input.Schedules.WeekdayAndTimeSchedules = append(input.Schedules.WeekdayAndTimeSchedules[:index], input.Schedules.WeekdayAndTimeSchedules[index+1:]...)
				return input, nil
			default:
				return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
			}
		default:
			return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
		}
	case strings.HasPrefix(actionName, "scheduleexception_"):
		// e.g. scheduleexception_0_yearmonthday_timezone
		ss := strings.SplitN(actionName, "_", 4)
		if len(ss) < 4 {
			return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
		}
		index, err := strconv.Atoi(ss[1])
		if err != nil {
			return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
		}
		switch ss[2] {
		case "yearmonthday":
			switch ss[3] {
			case "yearmonthday":
				if index < len(input.Schedules.YearMonthDayScheduleExceptions) {
					_, err := fmt.Sscanf(
						value,
						"%04d-%02d-%02d",
						&input.Schedules.YearMonthDayScheduleExceptions[index].Year,
						&input.Schedules.YearMonthDayScheduleExceptions[index].Month,
						&input.Schedules.YearMonthDayScheduleExceptions[index].Day,
					)
					if err != nil {
						return input, domain.ErrorBadRequest(errors.Wrapf(err, "unknown value %s", value))
					}
					return input, nil
				}
				return input, nil
			case "timezone":
				v, err := strconv.Atoi(value)
				if err != nil {
					return input, domain.ErrorBadRequest(fmt.Errorf("unknown value %s", value))
				}
				if index < len(input.Schedules.YearMonthDayScheduleExceptions) {
					input.Schedules.YearMonthDayScheduleExceptions[index].Timezone = v
					return input, nil
				}
				return input, nil
			case "remove":
				if index >= len(input.Schedules.YearMonthDayScheduleExceptions) {
					return input, nil
				}
				input.Schedules.YearMonthDayScheduleExceptions = append(input.Schedules.YearMonthDayScheduleExceptions[:index], input.Schedules.YearMonthDayScheduleExceptions[index+1:]...)
				return input, nil
			default:
				return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
			}
		default:
			return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
		}
	case actionName == "appendschedule":
		switch value {
		case "weekdayandtime":
			input.Schedules.WeekdayAndTimeSchedules = append(input.Schedules.WeekdayAndTimeSchedules, WeekdayAndTimeSchedule{
				Hour:     10,
				Minute:   0,
				Sec:      0,
				Timezone: 0,
				Mon:      true,
				Tue:      true,
				Wed:      true,
				Thu:      true,
				Fri:      true,
				Sat:      false,
				Sun:      false,
			})
		}
		return input, nil
	case actionName == "appendscheduleexception":
		switch value {
		case "yearmonthday":
			input.Schedules.YearMonthDayScheduleExceptions = append(input.Schedules.YearMonthDayScheduleExceptions, YearMonthDayScheduleException{
				Year:     uint32(time.Now().Year()),
				Month:    uint32(time.Now().Month()),
				Day:      uint32(time.Now().Day()),
				Timezone: 0,
			})
		}
		return input, nil
	case strings.HasPrefix(actionName, "posttarget_"):
		// e.g. posttarget_0_slack_channelid
		ss := strings.SplitN(actionName, "_", 4)
		if len(ss) < 4 {
			return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
		}
		index, err := strconv.Atoi(ss[1])
		if err != nil {
			return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
		}
		switch ss[2] {
		case "slack":
			switch ss[3] {
			case "channelid":
				if index < len(input.PostTargets.SlackPostTargets) {
					input.PostTargets.SlackPostTargets[index].ChannelID = value
					return input, nil
				}
				return input, nil
			case "remove":
				if index >= len(input.PostTargets.SlackPostTargets) {
					return input, nil
				}
				input.PostTargets.SlackPostTargets = append(input.PostTargets.SlackPostTargets[:index], input.PostTargets.SlackPostTargets[index+1:]...)
			default:
				return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
			}
		default:
			return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
		}
	case actionName == "appendposttarget":
		switch value {
		case "slack":
			input.PostTargets.SlackPostTargets = append(input.PostTargets.SlackPostTargets, SlackPostTarget{})
		}
		return input, nil
	default:
		return input, domain.ErrorBadRequest(fmt.Errorf("unknown actionName %s", actionName))
	}
	return input, nil
}
