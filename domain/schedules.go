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

import "time"

const ScheduleKindSchedules = "schedules"

type Schedules struct {
	scheduleKind ScheduleKind        `getter:""`
	schedules    []Schedule          `getter:""`
	exceptions   []ScheduleException `getter:""`
}

func NewSchedules(
	baseSchedule Schedule,
	additionalSchedules []Schedule,
	exceptions []ScheduleException,
) Schedules {
	ss := make([]Schedule, 0, 1+len(additionalSchedules))
	ss = append(ss, baseSchedule)
	ss = append(ss, additionalSchedules...)
	return Schedules{
		scheduleKind: ScheduleKindSchedules,
		schedules:    ss,
		exceptions:   exceptions,
	}
}

func (s Schedules) NextTime(baseTime time.Time) time.Time {
	var earliest time.Time
	for i, ss := range s.schedules {
		nextTime := ss.NextTime(baseTime)
		for s.needsIgnore(nextTime) {
			nextTime = ss.NextTime(nextTime)
		}
		if i != 0 && earliest.Before(nextTime) {
			continue
		}
		earliest = nextTime
	}
	return earliest
}

func (s Schedules) needsIgnore(t time.Time) bool {
	for _, e := range s.exceptions {
		if e.NeedsIgnore(t) {
			return true
		}
	}
	return false
}

func (s Schedules) PrevTime(baseTime time.Time) time.Time {
	var latest time.Time
	for i, ss := range s.schedules {
		prevTime := ss.PrevTime(baseTime)
		for s.needsIgnore(prevTime) {
			prevTime = ss.PrevTime(prevTime)
		}
		if i != 0 && latest.After(prevTime) {
			continue
		}
		latest = prevTime
	}
	return latest
}

func (s Schedules) Dump() ScheduleValue {
	v := ScheduleValue{
		ScheduleKind: s.ScheduleKind(),
		Schedules:    make([]ScheduleValue, len(s.schedules)),
		Exceptions:   make([]ScheduleExceptionValue, len(s.exceptions)),
	}
	for i, s := range s.schedules {
		v.Schedules[i] = s.Dump()
	}
	for i, s := range s.exceptions {
		v.Exceptions[i] = s.Dump()
	}
	return v
}

func RestoreSchedulesFromDumped(v ScheduleValue) Schedules {
	schedules := make([]Schedule, len(v.Schedules))
	for i, s := range v.Schedules {
		schedules[i] = RestoreScheduleFromDumped(s)
	}
	exceptions := make([]ScheduleException, len(v.Exceptions))
	for i, s := range v.Exceptions {
		exceptions[i] = RestoreScheduleExceptionFromDumped(s)
	}
	return NewSchedules(schedules[0], schedules[1:], exceptions)
}
