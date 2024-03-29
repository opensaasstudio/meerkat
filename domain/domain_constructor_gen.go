// Code generated by go-genconstructor; DO NOT EDIT.

package domain

import (
	"time"
)

func NewAdmin(
	id AdminID,
	name string,
) Admin {
	return Admin{
		id:   id,
		name: name,
	}
}

func NewAnswer(
	id AnswerID,
	questionnaireID QuestionnaireID,
	questionID QuestionID,
	answererID AnswererID,
	answeredAt time.Time,
	value string,
) Answer {
	return Answer{
		id:              id,
		questionnaireID: questionnaireID,
		questionID:      questionID,
		answererID:      answererID,
		answeredAt:      answeredAt,
		value:           value,
	}
}

func NewAnswerer(
	id AnswererID,
	name string,
) Answerer {
	return Answerer{
		id:   id,
		name: name,
	}
}

func NewAskingService(
	questionnaireSearcher QuestionnaireSearcher,
	notificationTargetSearcher NotificationTargetSearcher,
	notificationService NotificationService,
	lastExecutedRecorder LastExecutedRecorder,
) AskingService {
	return AskingService{
		questionnaireSearcher:      questionnaireSearcher,
		notificationTargetSearcher: notificationTargetSearcher,
		notificationService:        notificationService,
		lastExecutedRecorder:       lastExecutedRecorder,
		nowFunc:                    time.Now,
	}
}

func NewNotificationTargetBase(
	id NotificationTargetID,
	questionnaireID QuestionnaireID,
	answererID AnswererID,
) NotificationTargetBase {
	return NotificationTargetBase{
		id:                     id,
		questionnaireID:        questionnaireID,
		answererID:             answererID,
		notificationTargetKind: NotificationTargetKindBase,
	}
}

func NewNotificationTargetSlack(
	id NotificationTargetID,
	questionnaireID QuestionnaireID,
	answererID AnswererID,
	channelID string,
	userID string,
) NotificationTargetSlack {
	return NotificationTargetSlack{
		id:                     id,
		questionnaireID:        questionnaireID,
		answererID:             answererID,
		notificationTargetKind: NotificationTargetKindSlack,
		channelID:              channelID,
		userID:                 userID,
	}
}

func NewPostTargetBase(
	id PostTargetID,
) PostTargetBase {
	return PostTargetBase{
		id:             id,
		postTargetKind: PostTargetKindBase,
	}
}

func NewPostTargetSlack(
	id PostTargetID,
	channelID string,
) PostTargetSlack {
	return PostTargetSlack{
		id:             id,
		postTargetKind: PostTargetKindSlack,
		channelID:      channelID,
	}
}

func NewQuestion(
	id QuestionID,
	text string,
) Question {
	return Question{
		id:   id,
		text: text,
	}
}

func NewQuestionnaire(
	id QuestionnaireID,
	title string,
	questionItems []QuestionItem,
) Questionnaire {
	return Questionnaire{
		id:            id,
		title:         title,
		questionItems: questionItems,
	}
}

func NewQuestionItem(
	question Question,
	required bool,
) QuestionItem {
	return QuestionItem{
		question: question,
		required: required,
	}
}

func NewQuestionnaireTarget(
	id QuestionnaireTargetID,
	questionnaireID QuestionnaireID,
	targetType QuestionnaireTargetType,
	targetIDs []string,
	schedules []Schedule,
	scheduleExceptions []ScheduleException,
) QuestionnaireTarget {
	return QuestionnaireTarget{
		id:                 id,
		questionnaireID:    questionnaireID,
		targetType:         targetType,
		targetIDs:          targetIDs,
		schedules:          schedules,
		scheduleExceptions: scheduleExceptions,
	}
}

func NewYearMonthDayScheduleException(
	year uint32,
	month time.Month,
	day uint32,
	locOffset int,
) YearMonthDayScheduleException {
	return YearMonthDayScheduleException{
		scheduleExceptionKind: ScheduleExceptionKindYearMonthDay,
		year:                  year,
		month:                 month,
		day:                   day,
		locOffset:             locOffset,
	}
}

func NewWeekdayAndTimeSchedule(
	hour uint32,
	minute uint32,
	sec uint32,
	locOffset int,
	mon bool,
	tue bool,
	wed bool,
	thu bool,
	fri bool,
	sat bool,
	sun bool,
) WeekdayAndTimeSchedule {
	return WeekdayAndTimeSchedule{
		scheduleKind: ScheduleKindWeekdayAndTime,
		hour:         hour,
		minute:       minute,
		sec:          sec,
		locOffset:    locOffset,
		mon:          mon,
		tue:          tue,
		wed:          wed,
		thu:          thu,
		fri:          fri,
		sat:          sat,
		sun:          sun,
	}
}
