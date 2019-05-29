package slack

import "github.com/opensaasstudio/meerkat/domain"

type PostTargets struct {
	SlackPostTargets []SlackPostTarget
}

type SlackPostTarget struct {
	ChannelID string
}

func RestorePostTargetFromDomainObject(s domain.PostTarget) PostTargets {
	postTargets := PostTargets{}
	switch s := s.(type) {
	case domain.PostTargetSlack:
		postTargets.SlackPostTargets = append(postTargets.SlackPostTargets, SlackPostTarget{
			ChannelID: s.ChannelID(),
		})
	}
	return postTargets
}

func (s PostTargets) Merge(t PostTargets) PostTargets {
	return PostTargets{
		SlackPostTargets: append(s.SlackPostTargets, t.SlackPostTargets...),
	}
}
