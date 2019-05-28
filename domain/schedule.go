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
