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

type AnswererID string

//genconstructor
type Answerer struct {
	id                  AnswererID           `required:"" getter:""`
	name                string               `required:"" getter:"" setter:"Rename"`
	notificationTargets []NotificationTarget `getter:""`
}

type AnswererValue struct {
	ID                  AnswererID `dynamo:",hash"`
	Name                string
	NotificationTargets []NotificationTargetValue
}

func (m Answerer) Dump() AnswererValue {
	notificationTargets := make([]NotificationTargetValue, len(m.NotificationTargets()))
	for i, nt := range m.NotificationTargets() {
		notificationTargets[i] = nt.Dump()
	}
	return AnswererValue{
		ID:                  m.ID(),
		Name:                m.Name(),
		NotificationTargets: notificationTargets,
	}
}

func RestoreAnswererFromDumped(v AnswererValue) Answerer {
	m := NewAnswerer(
		v.ID,
		v.Name,
	)
	m.notificationTargets = make([]NotificationTarget, len(v.NotificationTargets))
	for i, nt := range v.NotificationTargets {
		m.notificationTargets[i] = RestoreNotificationTargetFromDumped(nt)
	}
	return m
}
