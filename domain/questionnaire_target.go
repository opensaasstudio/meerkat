package domain

type QuestionnaireTargetID string
type QuestionnaireTargetType string

//genconstructor
type QuestionnaireTarget struct {
	id                 QuestionnaireTargetID   `required:"" getter:""`
	questionnaireID    QuestionnaireID         `required:"" getter:""`
	targetType         QuestionnaireTargetType `required:"" getter:""`
	targetIDs          []string                `required:"" getter:"" setter:""`
	schedules          []Schedule              `required:"" getter:"" setter:""`
	scheduleExceptions []ScheduleException     `required:"" getter:"" setter:""`
}
