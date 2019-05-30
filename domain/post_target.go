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

type PostTargetKind string

const (
	PostTargetKindBase  PostTargetKind = "base"
	PostTargetKindSlack PostTargetKind = "slack"
)

type PostTargetID string

type PostTarget interface {
	ID() PostTargetID
	PostTargetKind() PostTargetKind
	Dump() PostTargetValue
}

//genconstructor
type PostTargetBase struct {
	id             PostTargetID   `required:"" getter:""`
	postTargetKind PostTargetKind `required:"PostTargetKindBase" getter:""`
}

//genconstructor
type PostTargetSlack struct {
	id             PostTargetID   `required:"" getter:""`
	postTargetKind PostTargetKind `required:"PostTargetKindSlack" getter:""`
	channelID      string         `required:"" getter:""`
}

type PostTargetValue struct {
	ID             PostTargetID
	PostTargetKind PostTargetKind

	ChannelID string
}

func (m PostTargetBase) Dump() PostTargetValue {
	return PostTargetValue{
		ID:             m.ID(),
		PostTargetKind: m.PostTargetKind(),
	}
}

func (m PostTargetSlack) Dump() PostTargetValue {
	return PostTargetValue{
		ID:             m.ID(),
		PostTargetKind: m.PostTargetKind(),
		ChannelID:      m.ChannelID(),
	}
}

func RestorePostTargetFromDumpled(v PostTargetValue) PostTarget {
	switch v.PostTargetKind {
	case PostTargetKindSlack:
		return RestorePostTargetSlackFromDumped(v)
	default:
		return RestorePostTargetBaseFromDumped(v)
	}
}

func RestorePostTargetBaseFromDumped(v PostTargetValue) PostTargetBase {
	return NewPostTargetBase(
		v.ID,
	)
}

func RestorePostTargetSlackFromDumped(v PostTargetValue) PostTargetSlack {
	return NewPostTargetSlack(
		v.ID,
		v.ChannelID,
	)
}
