package domain_test

import (
	"testing"
	"time"

	"github.com/opensaasstudio/meerkat/domain"
	"github.com/stretchr/testify/assert"
)

func TestWeekdayAndTimeSchedule(t *testing.T) {
	scheduleImplemented(t, domain.NewWeekdayAndTimeSchedule(
		10, 30, 0, 0,
		true,
		true,
		true,
		true,
		true,
		false,
		false,
	))

	t.Run("NextTime", func(t *testing.T) {
		t.Run("UTC", func(t *testing.T) {
			s := domain.NewWeekdayAndTimeSchedule(
				10, 30, 0, 0,
				true,
				true,
				true,
				true,
				true,
				false,
				false,
			)
			for _, tt := range []struct {
				name string
				base time.Time
				want time.Time
			}{
				{
					name: "push the clock forward",
					base: time.Date(2019, 5, 7, 9, 0, 0, 0, time.UTC),
					want: time.Date(2019, 5, 7, 10, 30, 0, 0, time.UTC),
				},
				{
					name: "push the date forward",
					base: time.Date(2019, 5, 7, 10, 30, 0, 0, time.UTC),
					want: time.Date(2019, 5, 8, 10, 30, 0, 0, time.UTC),
				},
				{
					name: "push the clock forward: with location",
					base: time.Date(2019, 5, 7, 11, 0, 0, 0, time.FixedZone("+1", 1*60*60)),
					want: time.Date(2019, 5, 7, 10, 30, 0, 0, time.UTC),
				},
				{
					name: "push the date forward: with location",
					base: time.Date(2019, 5, 7, 9, 30, 0, 0, time.FixedZone("-1", -1*60*60)),
					want: time.Date(2019, 5, 8, 10, 30, 0, 0, time.UTC),
				},
				{
					name: "skip sat and sun",
					base: time.Date(2019, 5, 3, 10, 30, 0, 0, time.UTC),
					want: time.Date(2019, 5, 6, 10, 30, 0, 0, time.UTC),
				},
			} {
				tt := tt
				t.Run(tt.name, func(t *testing.T) {
					assert.WithinDuration(t, tt.want, s.NextTime(tt.base), 0)
				})
			}
		})

		t.Run("+1", func(t *testing.T) {
			s := domain.NewWeekdayAndTimeSchedule(
				10, 30, 0, 1,
				true,
				true,
				true,
				true,
				true,
				false,
				false,
			)
			for _, tt := range []struct {
				name string
				base time.Time
				want time.Time
			}{
				{
					name: "push the clock forward",
					base: time.Date(2019, 5, 7, 9, 0, 0, 0, time.UTC),
					want: time.Date(2019, 5, 7, 9, 30, 0, 0, time.UTC),
				},
				{
					name: "push the date forward",
					base: time.Date(2019, 5, 7, 9, 30, 0, 0, time.UTC),
					want: time.Date(2019, 5, 8, 9, 30, 0, 0, time.UTC),
				},
				{
					name: "push the clock forward: with location",
					base: time.Date(2019, 5, 7, 10, 0, 0, 0, time.FixedZone("+1", 1*60*60)),
					want: time.Date(2019, 5, 7, 9, 30, 0, 0, time.UTC),
				},
				{
					name: "push the date forward: with location",
					base: time.Date(2019, 5, 7, 8, 30, 0, 0, time.FixedZone("-1", -1*60*60)),
					want: time.Date(2019, 5, 8, 9, 30, 0, 0, time.UTC),
				},
				{
					name: "skip sat and sun",
					base: time.Date(2019, 5, 3, 9, 30, 0, 0, time.UTC),
					want: time.Date(2019, 5, 6, 9, 30, 0, 0, time.UTC),
				},
			} {
				tt := tt
				t.Run(tt.name, func(t *testing.T) {
					assert.WithinDuration(t, tt.want, s.NextTime(tt.base), 0)
				})
			}
		})
	})

	t.Run("PrevTime", func(t *testing.T) {
		s := domain.NewWeekdayAndTimeSchedule(
			10, 30, 0, 0,
			true,
			true,
			true,
			true,
			true,
			false,
			false,
		)
		for _, tt := range []struct {
			name string
			base time.Time
			want time.Time
		}{
			{
				name: "turn the clock back",
				base: time.Date(2019, 5, 7, 11, 0, 0, 0, time.UTC),
				want: time.Date(2019, 5, 7, 10, 30, 0, 0, time.UTC),
			},
			{
				name: "turn the date back",
				base: time.Date(2019, 5, 9, 10, 30, 0, 0, time.UTC),
				want: time.Date(2019, 5, 8, 10, 30, 0, 0, time.UTC),
			},
			{
				name: "turn the clock back: with location",
				base: time.Date(2019, 5, 7, 10, 0, 0, 0, time.FixedZone("-1", -1*60*60)),
				want: time.Date(2019, 5, 7, 10, 30, 0, 0, time.UTC),
			},
			{
				name: "turn the date back: with location",
				base: time.Date(2019, 5, 7, 11, 30, 0, 0, time.FixedZone("+1", 1*60*60)),
				want: time.Date(2019, 5, 6, 10, 30, 0, 0, time.UTC),
			},
			{
				name: "skip sat and sun",
				base: time.Date(2019, 5, 6, 10, 30, 0, 0, time.UTC),
				want: time.Date(2019, 5, 3, 10, 30, 0, 0, time.UTC),
			},
		} {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				assert.WithinDuration(t, tt.want, s.PrevTime(tt.base), 0)
			})
		}
	})
}
