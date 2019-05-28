package slack

import (
	"strconv"

	"github.com/nlopes/slack"
)

type Schedules struct {
	WeekdayAndTimeSchedules        []WeekdayAndTimeSchedule
	YearMonthDayScheduleExceptions []YearMonthDayScheduleException
}

type WeekdayAndTimeSchedule struct {
	Hour     uint32
	Minute   uint32
	Sec      uint32
	Timezone int
	Mon      bool
	Tue      bool
	Wed      bool
	Thu      bool
	Fri      bool
	Sat      bool
	Sun      bool
}

type YearMonthDayScheduleException struct {
	Year     uint32
	Month    uint32
	Day      uint32
	Timezone int
}

func NumberOptions(left, right int) []*slack.OptionBlockObject {
	options := make([]*slack.OptionBlockObject, right-left+1)
	for i := 0; i < right-left+1; i++ {
		s := strconv.Itoa(left + i)
		options[i] = slack.NewOptionBlockObject(s, slack.NewTextBlockObject("plain_text", s, false, false))
	}
	return options
}
