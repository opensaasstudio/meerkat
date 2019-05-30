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
