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
