package domain_test

import (
	"testing"
	"time"

	"github.com/opensaasstudio/meerkat/domain"
	"github.com/stretchr/testify/assert"
)

func scheduleImplemented(t *testing.T, schedule domain.Schedule) {
	t.Run("schedule implemeneted", func(t *testing.T) {
		bt := time.Date(2019, 5, 1, 2, 3, 4, 5, time.UTC)
		assert.NotEqual(t, schedule.NextTime(bt), schedule.NextTime(schedule.NextTime(bt)))
		assert.NotEqual(t, schedule.PrevTime(bt), schedule.PrevTime(schedule.PrevTime(bt)))
	})
}
