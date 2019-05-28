package domain

import (
	"strconv"
	"time"
)

const ScheduleKindWeekdayAndTime = "weekdayandtime"

//genconstructor
type WeekdayAndTimeSchedule struct {
	scheduleKind ScheduleKind `required:"ScheduleKindWeekdayAndTime" getter:""`
	hour         uint32       `required:"" getter:""`
	minute       uint32       `required:"" getter:""`
	sec          uint32       `required:"" getter:""`
	locOffset    int          `required:"" getter:""`
	mon          bool         `required:"" getter:""`
	tue          bool         `required:"" getter:""`
	wed          bool         `required:"" getter:""`
	thu          bool         `required:"" getter:""`
	fri          bool         `required:"" getter:""`
	sat          bool         `required:"" getter:""`
	sun          bool         `required:"" getter:""`
}

func (s WeekdayAndTimeSchedule) NextTime(baseTime time.Time) time.Time {
	loc := time.FixedZone(strconv.Itoa(s.LocOffset()), s.LocOffset()*60*60)
	t := baseTime.In(loc)
	t = time.Date(t.Year(), t.Month(), t.Day(), int(s.hour), int(s.minute), int(s.sec), 0, loc)
	if !t.After(baseTime) {
		t = t.AddDate(0, 0, 1)
	}
	for s.needsIgnore(t) {
		t = t.AddDate(0, 0, 1)
	}

	return t
}

func (s WeekdayAndTimeSchedule) PrevTime(baseTime time.Time) time.Time {
	loc := time.FixedZone(strconv.Itoa(s.LocOffset()), s.LocOffset()*60*60)
	t := baseTime.In(loc)
	t = time.Date(t.Year(), t.Month(), t.Day(), int(s.hour), int(s.minute), int(s.sec), 0, loc)
	if !t.Before(baseTime) {
		t = t.AddDate(0, 0, -1)
	}
	for s.needsIgnore(t) {
		t = t.AddDate(0, 0, -1)
	}

	return t
}

func (s WeekdayAndTimeSchedule) needsIgnore(t time.Time) bool {
	switch t.In(time.FixedZone(strconv.Itoa(s.LocOffset()), s.LocOffset()*60*60)).Weekday() {
	case time.Monday:
		return !s.mon
	case time.Tuesday:
		return !s.tue
	case time.Wednesday:
		return !s.wed
	case time.Thursday:
		return !s.thu
	case time.Friday:
		return !s.fri
	case time.Saturday:
		return !s.sat
	case time.Sunday:
		return !s.sun
	}
	return true
}

func (s WeekdayAndTimeSchedule) Dump() ScheduleValue {
	return ScheduleValue{
		ScheduleKind: s.ScheduleKind(),
		Hour:         s.Hour(),
		Minute:       s.Minute(),
		Sec:          s.Sec(),
		LocOffset:    s.LocOffset(),
		Mon:          s.Mon(),
		Tue:          s.Tue(),
		Wed:          s.Wed(),
		Thu:          s.Thu(),
		Fri:          s.Fri(),
		Sat:          s.Sat(),
		Sun:          s.Sun(),
	}
}

func RestoreScheduleWeekdayAndTimeFromDumpled(v ScheduleValue) WeekdayAndTimeSchedule {
	return NewWeekdayAndTimeSchedule(
		v.Hour,
		v.Minute,
		v.Sec,
		v.LocOffset,
		v.Mon,
		v.Tue,
		v.Wed,
		v.Thu,
		v.Fri,
		v.Sat,
		v.Sun,
	)
}
