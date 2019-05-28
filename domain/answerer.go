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
