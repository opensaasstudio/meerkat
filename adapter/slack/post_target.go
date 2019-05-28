package slack

type PostTargets struct {
	SlackPostTargets []SlackPostTarget
}

type SlackPostTarget struct {
	ChannelID string
}
