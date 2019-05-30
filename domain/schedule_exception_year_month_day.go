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

package domain

import (
	"strconv"
	"time"
)

const ScheduleExceptionKindYearMonthDay = "yearmonthday"

//genconstructor
type YearMonthDayScheduleException struct {
	scheduleExceptionKind ScheduleExceptionKind `required:"ScheduleExceptionKindYearMonthDay" getter:""`
	year                  uint32                `required:"" getter:""`
	month                 time.Month            `required:"" getter:""`
	day                   uint32                `required:"" getter:""`
	locOffset             int                   `required:"" getter:""`
}

func (s YearMonthDayScheduleException) NeedsIgnore(t time.Time) bool {
	left := time.Date(int(s.year), s.month, int(s.day), 0, 0, 0, 0, time.FixedZone(strconv.Itoa(s.locOffset), s.locOffset*60*60))
	return !t.Before(left) && t.Before(left.AddDate(0, 0, 1))
}

func (s YearMonthDayScheduleException) Dump() ScheduleExceptionValue {
	return ScheduleExceptionValue{
		ScheduleExceptionKind: s.ScheduleExceptionKind(),
		Year:                  s.Year(),
		Month:                 s.Month(),
		Day:                   s.Day(),
		LocOffset:             s.LocOffset(),
	}
}

func RestoreScheudleExceptionYearMonthDayFromDumped(v ScheduleExceptionValue) YearMonthDayScheduleException {
	return NewYearMonthDayScheduleException(v.Year, v.Month, v.Day, v.LocOffset)
}
