package domain

type NotificationTargetKind string

const (
	NotificationTargetKindBase  NotificationTargetKind = "base"
	NotificationTargetKindSlack NotificationTargetKind = "slack"
)

type NotificationTargetID string

type NotificationTarget interface {
	ID() NotificationTargetID
	QuestionnaireID() QuestionnaireID
	AnswererID() AnswererID
	NotificationTargetKind() NotificationTargetKind
	Dump() NotificationTargetValue
}

//genconstructor
type NotificationTargetBase struct {
	id                     NotificationTargetID   `required:"" getter:""`
	questionnaireID        QuestionnaireID        `required:"" getter:""`
	answererID             AnswererID             `required:"" getter:""`
	notificationTargetKind NotificationTargetKind `required:"NotificationTargetKindBase" getter:""`
}

//genconstructor
type NotificationTargetSlack struct {
	id                     NotificationTargetID   `required:"" getter:""`
	questionnaireID        QuestionnaireID        `required:"" getter:""`
	answererID             AnswererID             `required:"" getter:""`
	notificationTargetKind NotificationTargetKind `required:"NotificationTargetKindSlack" getter:""`
	channelID              string                 `required:"" getter:""`
	userID                 string                 `required:"" getter:""`
	needsMention           bool                   `getter:"" setter:"ToggleNeedsMention"`
}

type NotificationTargetValue struct {
	ID                     NotificationTargetID `dynamo:",hash"`
	QuestionnaireID        QuestionnaireID
	AnswererID             AnswererID
	NotificationTargetKind NotificationTargetKind

	ChannelID    string
	UserID       string
	NeedsMention bool
}

func (m NotificationTargetBase) Dump() NotificationTargetValue {
	return NotificationTargetValue{
		ID:                     m.ID(),
		QuestionnaireID:        m.QuestionnaireID(),
		AnswererID:             m.AnswererID(),
		NotificationTargetKind: m.NotificationTargetKind(),
	}
}

func (m NotificationTargetSlack) Dump() NotificationTargetValue {
	return NotificationTargetValue{
		ID:                     m.ID(),
		QuestionnaireID:        m.QuestionnaireID(),
		AnswererID:             m.AnswererID(),
		NotificationTargetKind: m.NotificationTargetKind(),
		ChannelID:              m.ChannelID(),
		UserID:                 m.UserID(),
		NeedsMention:           m.NeedsMention(),
	}
}

func RestoreNotificationTargetFromDumped(v NotificationTargetValue) NotificationTarget {
	switch v.NotificationTargetKind {
	case NotificationTargetKindSlack:
		return RestoreNotificationTargetSlackFromDumped(v)
	default:
		return RestoreNotificationTargetBaseFromDumped(v)
	}
}

func RestoreNotificationTargetBaseFromDumped(v NotificationTargetValue) NotificationTargetBase {
	return NewNotificationTargetBase(
		v.ID,
		v.QuestionnaireID,
		v.AnswererID,
	)
}

func RestoreNotificationTargetSlackFromDumped(v NotificationTargetValue) NotificationTargetSlack {
	m := NewNotificationTargetSlack(
		v.ID,
		v.QuestionnaireID,
		v.AnswererID,
		v.ChannelID,
		v.UserID,
	)
	m.ToggleNeedsMention(v.NeedsMention)
	return m
}
