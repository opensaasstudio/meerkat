package domain_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/opensaasstudio/meerkat/domain"
	"github.com/stretchr/testify/assert"
)

func TestYearMonthDaySchedleException(t *testing.T) {
	var s domain.ScheduleException = domain.NewYearMonthDayScheduleException(2019, 5, 1, 0)
	for _, tt := range []struct {
		t    time.Time
		want bool
	}{
		{
			t:    time.Date(2019, 4, 30, 23, 59, 59, 999999999, time.UTC),
			want: false,
		},
		{
			t:    time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC),
			want: true,
		},
		{
			t:    time.Date(2019, 5, 1, 23, 59, 59, 999999999, time.UTC),
			want: true,
		},
		{
			t:    time.Date(2019, 5, 2, 0, 0, 0, 0, time.UTC),
			want: false,
		},
		{
			t:    time.Date(2019, 5, 1, 0, 59, 59, 999999999, time.FixedZone("+1", 1*60*60)),
			want: false,
		},
		{
			t:    time.Date(2019, 4, 30, 23, 0, 0, 0, time.FixedZone("-1", -1*60*60)),
			want: true,
		},
		{
			t:    time.Date(2019, 5, 2, 0, 59, 59, 999999999, time.FixedZone("+1", 1*60*60)),
			want: true,
		},
		{
			t:    time.Date(2019, 5, 1, 23, 0, 0, 0, time.FixedZone("-1", -1*60*60)),
			want: false,
		},
	} {
		tt := tt
		t.Run(fmt.Sprint(tt.t), func(t *testing.T) {
			assert.Equal(t, tt.want, s.NeedsIgnore(tt.t))
		})
	}
}
