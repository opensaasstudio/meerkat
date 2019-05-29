// Code generated by go-genconstructor; DO NOT EDIT.

package application

import (
	"context"

	"github.com/opensaasstudio/meerkat/domain"
)

func NewAnsweringUsecase(
	questionnaireRepository QuestionnaireRepository,
	answererRepository AnswererRepository,
	answerRepository AnswerRepository,
	answerIDProvider AnswerIDProvider,
	postingService domain.PostingService,
) AnsweringUsecase {
	return AnsweringUsecase{
		questionnaireRepository: questionnaireRepository,
		answererRepository:      answererRepository,
		answerRepository:        answerRepository,
		answerIDProvider:        answerIDProvider,
		postingService:          postingService,
	}
}

func NewAnswerInputValue(
	questionID domain.QuestionID,
	value string,
) AnswerInputValue {
	return AnswerInputValue{
		questionID: questionID,
		value:      value,
	}
}

func NewAnsweringUsecaseInput(
	questionnaireID domain.QuestionnaireID,
	answererID domain.AnswererID,
	answers []AnswerInputValue,
) AnsweringUsecaseInput {
	return AnsweringUsecaseInput{
		questionnaireID: questionnaireID,
		answererID:      answererID,
		answers:         answers,
	}
}

func NewCreatingAnswererUsecase(
	repository AnswererRepository,
	answererIDProvider AnswererIDProvider,
) CreatingAnswererUsecase {
	return CreatingAnswererUsecase{
		repository:         repository,
		answererIDProvider: answererIDProvider,
	}
}

func NewCreatingAnswererUsecaseInput(
	name string,
) CreatingAnswererUsecaseInput {
	return CreatingAnswererUsecaseInput{
		name: name,
	}
}

func NewAskingUsecase(
	askingService domain.AskingService,
) AskingUsecase {
	return AskingUsecase{
		askingService: askingService,
	}
}

func NewLastExecutedRecorder(
	quesitonnaireRepository QuestionnaireRepository,
) LastExecutedRecorder {
	return LastExecutedRecorder{
		quesitonnaireRepository: quesitonnaireRepository,
	}
}

func NewNotificationService(
	notifySlack func(ctx context.Context, target domain.NotificationTargetSlack, questionnaire domain.Questionnaire) domain.Error,
) NotificationService {
	return NotificationService{
		notifySlack: notifySlack,
	}
}

func NewPostingService(
	postSlack func(ctx context.Context, target domain.PostTargetSlack, questionnaire domain.Questionnaire, answerer domain.Answerer, answers []domain.Answer) domain.Error,
) PostingService {
	return PostingService{
		postSlack: postSlack,
	}
}

func NewAddingAnswererUsecase(
	repository NotificationTargetRepository,
	notificationTargetIDProvider NotificationTargetIDProvider,
) AddingAnswererUsecase {
	return AddingAnswererUsecase{
		repository:                   repository,
		notificationTargetIDProvider: notificationTargetIDProvider,
	}
}

func NewAddingAnswererUsecaseInput(
	notificationTarget domain.NotificationTarget,
) AddingAnswererUsecaseInput {
	return AddingAnswererUsecaseInput{
		notificationTarget: notificationTarget,
	}
}

func NewCreatingQuestionnaireUsecase(
	repository QuestionnaireRepository,
	questionnaireIDProvider QuestionnaireIDProvider,
	questionIDProvider QuestionIDProvider,
	authorizationService CreatingQuestionnaireAuthorizationService,
) CreatingQuestionnaireUsecase {
	return CreatingQuestionnaireUsecase{
		repository:              repository,
		questionnaireIDProvider: questionnaireIDProvider,
		questionIDProvider:      questionIDProvider,
		authorizationService:    authorizationService,
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

func NewQuestion(
	text string,
) Question {
	return Question{
		text: text,
	}
}

func NewCreatingQuestionnaireUsecaseInput(
	title string,
	questionItems []QuestionItem,
) CreatingQuestionnaireUsecaseInput {
	return CreatingQuestionnaireUsecaseInput{
		title:         title,
		questionItems: questionItems,
	}
}

func NewRemovingAnswererUsecase(
	searcher domain.NotificationTargetSearcher,
	repository NotificationTargetRepository,
) RemovingAnswererUsecase {
	return RemovingAnswererUsecase{
		searcher:   searcher,
		repository: repository,
	}
}

func NewRemovingAnswererUsecaseInput(
	questionnaireID domain.QuestionnaireID,
	answererID domain.AnswererID,
) RemovingAnswererUsecaseInput {
	return RemovingAnswererUsecaseInput{
		questionnaireID: questionnaireID,
		answererID:      answererID,
	}
}

func NewOverwritingQuestionnaireTitleUsecase(
	repository QuestionnaireRepository,
) OverwritingQuestionnaireTitleUsecase {
	return OverwritingQuestionnaireTitleUsecase{
		repository: repository,
	}
}

func NewOverwritingQuestionnaireTitleUsecaseInput(
	questionnaireID domain.QuestionnaireID,
	title string,
) OverwritingQuestionnaireTitleUsecaseInput {
	return OverwritingQuestionnaireTitleUsecaseInput{
		questionnaireID: questionnaireID,
		title:           title,
	}
}

func NewReplacingQuestionsUsecase(
	repository QuestionnaireRepository,
	questionIDProvider QuestionIDProvider,
) ReplacingQuestionsUsecase {
	return ReplacingQuestionsUsecase{
		repository:         repository,
		questionIDProvider: questionIDProvider,
	}
}

func NewReplacingQuestionsUsecaseInput(
	questionnaireID domain.QuestionnaireID,
	questionItems []QuestionItem,
) ReplacingQuestionsUsecaseInput {
	return ReplacingQuestionsUsecaseInput{
		questionnaireID: questionnaireID,
		questionItems:   questionItems,
	}
}

func NewUpdatingQuestionnaireUsecase(
	repository QuestionnaireRepository,
	questionIDProvider QuestionIDProvider,
	authorizationService UpdatingQuestionnaireAuthorizationService,
) UpdatingQuestionnaireUsecase {
	return UpdatingQuestionnaireUsecase{
		repository:           repository,
		questionIDProvider:   questionIDProvider,
		authorizationService: authorizationService,
	}
}

func NewUpdatingQuestionnaireUsecaseInput(
	questionnaireID domain.QuestionnaireID,
	creatingInput CreatingQuestionnaireUsecaseInput,
) UpdatingQuestionnaireUsecaseInput {
	return UpdatingQuestionnaireUsecaseInput{
		questionnaireID: questionnaireID,
		creatingInput:   creatingInput,
	}
}