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
	left := time.Date(int(s.year), s.month, int(s.day), 0, 0, 0, 0, time.FixedZone(strconv.Itoa(s.locOffset), s.locOffset))
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
