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

package domain_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/opensaasstudio/meerkat/domain"
	"github.com/opensaasstudio/meerkat/domain/mock_domain"
	"github.com/stretchr/testify/assert"
)

func TestSchedules(t *testing.T) {
	var _ domain.Schedule = domain.Schedules{}

	t.Run("NextTime", func(t *testing.T) {
		t.Run("use earliest", func(t *testing.T) {

			baseTime := time.Date(2019, 5, 1, 2, 3, 4, 5, time.UTC)

			for _, tt := range []struct {
				name                      string
				baseScheduleResult        time.Time
				additionalScheduleResults []time.Time
				want                      time.Time
			}{
				{
					name:               "only baseSchedule",
					baseScheduleResult: time.Date(2019, 5, 2, 2, 3, 4, 5, time.UTC),
					want:               time.Date(2019, 5, 2, 2, 3, 4, 5, time.UTC),
				},
				{
					name:               "baseSchedule is earliest",
					baseScheduleResult: time.Date(2019, 5, 2, 2, 3, 4, 5, time.UTC),
					additionalScheduleResults: []time.Time{
						time.Date(2019, 5, 2, 3, 3, 4, 5, time.UTC),
						time.Date(2019, 5, 2, 4, 3, 4, 5, time.UTC),
					},
					want: time.Date(2019, 5, 2, 2, 3, 4, 5, time.UTC),
				},
				{
					name:               "one of additionalSchedules is earliest",
					baseScheduleResult: time.Date(2019, 5, 2, 5, 3, 4, 5, time.UTC),
					additionalScheduleResults: []time.Time{
						time.Date(2019, 5, 2, 3, 3, 4, 5, time.UTC), // <- earliest
						time.Date(2019, 5, 2, 4, 3, 4, 5, time.UTC),
					},
					want: time.Date(2019, 5, 2, 3, 3, 4, 5, time.UTC),
				},
			} {
				tt := tt
				t.Run(tt.name, func(t *testing.T) {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()

					baseSchedule := mock_domain.NewMockSchedule(ctrl)
					baseSchedule.EXPECT().NextTime(baseTime).Return(tt.baseScheduleResult)

					additionalSchedules := make([]domain.Schedule, len(tt.additionalScheduleResults))
					for i := range tt.additionalScheduleResults {
						schedule := mock_domain.NewMockSchedule(ctrl)
						schedule.EXPECT().NextTime(baseTime).Return(tt.additionalScheduleResults[i])
						additionalSchedules[i] = schedule
					}

					s := domain.NewSchedules(baseSchedule, additionalSchedules, nil)
					assert.Equal(t, tt.want, s.NextTime(baseTime))
				})
			}
		})
		t.Run("with exceptions", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			baseTime := time.Date(2019, 5, 1, 2, 3, 4, 5, time.UTC)

			baseSchedule := mock_domain.NewMockSchedule(ctrl)
			baseSchedule.EXPECT().NextTime(gomock.Any()).DoAndReturn(func(baseTime time.Time) time.Time {
				return baseTime.Add(time.Hour)
			}).AnyTimes()

			exception := mock_domain.NewMockScheduleException(ctrl)
			exception.EXPECT().NeedsIgnore(time.Date(2019, 5, 1, 3, 3, 4, 5, time.UTC)).Return(true)
			exception.EXPECT().NeedsIgnore(time.Date(2019, 5, 1, 4, 3, 4, 5, time.UTC)).Return(true)
			exception.EXPECT().NeedsIgnore(time.Date(2019, 5, 1, 5, 3, 4, 5, time.UTC)).Return(false)

			s := domain.NewSchedules(baseSchedule, nil, []domain.ScheduleException{exception})
			assert.Equal(t, time.Date(2019, 5, 1, 5, 3, 4, 5, time.UTC), s.NextTime(baseTime))
		})
	})

	t.Run("PrevTime", func(t *testing.T) {
		t.Run("use latest", func(t *testing.T) {

			baseTime := time.Date(2019, 5, 1, 2, 3, 4, 5, time.UTC)

			for _, tt := range []struct {
				name                      string
				baseScheduleResult        time.Time
				additionalScheduleResults []time.Time
				want                      time.Time
			}{
				{
					name:               "only baseSchedule",
					baseScheduleResult: time.Date(2019, 4, 30, 2, 3, 4, 5, time.UTC),
					want:               time.Date(2019, 4, 30, 2, 3, 4, 5, time.UTC),
				},
				{
					name:               "baseSchedule is latest",
					baseScheduleResult: time.Date(2019, 4, 30, 4, 3, 4, 5, time.UTC),
					additionalScheduleResults: []time.Time{
						time.Date(2019, 4, 30, 2, 3, 4, 5, time.UTC),
						time.Date(2019, 4, 30, 3, 3, 4, 5, time.UTC),
					},
					want: time.Date(2019, 4, 30, 4, 3, 4, 5, time.UTC),
				},
				{
					name:               "one of additionalSchedules is latest",
					baseScheduleResult: time.Date(2019, 4, 30, 3, 3, 4, 5, time.UTC),
					additionalScheduleResults: []time.Time{
						time.Date(2019, 4, 30, 5, 3, 4, 5, time.UTC), // <- latest
						time.Date(2019, 4, 30, 4, 3, 4, 5, time.UTC),
					},
					want: time.Date(2019, 4, 30, 5, 3, 4, 5, time.UTC),
				},
			} {
				tt := tt
				t.Run(tt.name, func(t *testing.T) {
					ctrl := gomock.NewController(t)
					defer ctrl.Finish()

					baseSchedule := mock_domain.NewMockSchedule(ctrl)
					baseSchedule.EXPECT().PrevTime(baseTime).Return(tt.baseScheduleResult)

					additionalSchedules := make([]domain.Schedule, len(tt.additionalScheduleResults))
					for i := range tt.additionalScheduleResults {
						schedule := mock_domain.NewMockSchedule(ctrl)
						schedule.EXPECT().PrevTime(baseTime).Return(tt.additionalScheduleResults[i])
						additionalSchedules[i] = schedule
					}

					s := domain.NewSchedules(baseSchedule, additionalSchedules, nil)
					assert.Equal(t, tt.want, s.PrevTime(baseTime))
				})
			}
		})
		t.Run("with exceptions", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			baseTime := time.Date(2019, 5, 1, 5, 3, 4, 5, time.UTC)

			baseSchedule := mock_domain.NewMockSchedule(ctrl)
			baseSchedule.EXPECT().PrevTime(gomock.Any()).DoAndReturn(func(baseTime time.Time) time.Time {
				return baseTime.Add(-time.Hour)
			}).AnyTimes()

			exception := mock_domain.NewMockScheduleException(ctrl)
			exception.EXPECT().NeedsIgnore(time.Date(2019, 5, 1, 4, 3, 4, 5, time.UTC)).Return(true)
			exception.EXPECT().NeedsIgnore(time.Date(2019, 5, 1, 3, 3, 4, 5, time.UTC)).Return(true)
			exception.EXPECT().NeedsIgnore(time.Date(2019, 5, 1, 2, 3, 4, 5, time.UTC)).Return(false)

			s := domain.NewSchedules(baseSchedule, nil, []domain.ScheduleException{exception})
			assert.Equal(t, time.Date(2019, 5, 1, 2, 3, 4, 5, time.UTC), s.PrevTime(baseTime))
		})
	})
}
