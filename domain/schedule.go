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

//go:generate mockgen -source $GOFILE -destination mock_$GOPACKAGE/${GOFILE}_mock.go

import "time"

type ScheduleKind string

type Schedule interface {
	NextTime(baseTime time.Time) time.Time
	PrevTime(baseTime time.Time) time.Time
	Dump() ScheduleValue
	ScheduleKind() ScheduleKind
}

type ScheduleExceptionKind string

type ScheduleException interface {
	NeedsIgnore(time.Time) bool
	Dump() ScheduleExceptionValue
	ScheduleExceptionKind() ScheduleExceptionKind
}

type ScheduleValue struct {
	ScheduleKind ScheduleKind

	Hour      uint32
	Minute    uint32
	Sec       uint32
	LocOffset int
	Mon       bool
	Tue       bool
	Wed       bool
	Thu       bool
	Fri       bool
	Sat       bool
	Sun       bool

	Schedules  []ScheduleValue
	Exceptions []ScheduleExceptionValue
}

type ScheduleExceptionValue struct {
	ScheduleExceptionKind ScheduleExceptionKind
	Year                  uint32
	Month                 time.Month
	Day                   uint32
	LocOffset             int
}

func RestoreScheduleFromDumped(v ScheduleValue) Schedule {
	switch v.ScheduleKind {
	case ScheduleKindSchedules:
		return RestoreSchedulesFromDumped(v)
	case ScheduleKindWeekdayAndTime:
		return RestoreScheduleWeekdayAndTimeFromDumpled(v)
	}
	return nil
}

func RestoreScheduleExceptionFromDumped(v ScheduleExceptionValue) ScheduleException {
	switch v.ScheduleExceptionKind {
	case ScheduleExceptionKindYearMonthDay:
		return RestoreScheudleExceptionYearMonthDayFromDumped(v)
	}
	return nil
}
