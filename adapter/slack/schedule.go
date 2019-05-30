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
	"strconv"

	"github.com/nlopes/slack"
	"github.com/opensaasstudio/meerkat/domain"
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

func RestoreSchedulesFromDomainObject(s domain.Schedule) Schedules {
	schedules := Schedules{}
	switch s := s.(type) {
	case domain.WeekdayAndTimeSchedule:
		schedules.WeekdayAndTimeSchedules = append(schedules.WeekdayAndTimeSchedules, WeekdayAndTimeSchedule{
			Hour:     s.Hour(),
			Minute:   s.Minute(),
			Sec:      s.Sec(),
			Timezone: s.LocOffset(),
			Mon:      s.Mon(),
			Tue:      s.Tue(),
			Wed:      s.Wed(),
			Thu:      s.Thu(),
			Fri:      s.Fri(),
			Sat:      s.Sat(),
			Sun:      s.Sun(),
		})
	case domain.Schedules:
		for _, s := range s.Schedules() {
			schedules = schedules.Merge(RestoreSchedulesFromDomainObject(s))
		}
		for _, s := range s.Exceptions() {
			switch s := s.(type) {
			case domain.YearMonthDayScheduleException:
				schedules.YearMonthDayScheduleExceptions = append(schedules.YearMonthDayScheduleExceptions, YearMonthDayScheduleException{
					Year:     s.Year(),
					Month:    uint32(s.Month()),
					Day:      s.Day(),
					Timezone: s.LocOffset(),
				})
			}
		}
	}
	return schedules
}

func (s Schedules) Merge(t Schedules) Schedules {
	return Schedules{
		WeekdayAndTimeSchedules:        append(s.WeekdayAndTimeSchedules, t.WeekdayAndTimeSchedules...),
		YearMonthDayScheduleExceptions: append(s.YearMonthDayScheduleExceptions, t.YearMonthDayScheduleExceptions...),
	}
}
