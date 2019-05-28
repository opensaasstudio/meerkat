package main

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	dynamodblib "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/kelseyhightower/envconfig"
	slacklib "github.com/nlopes/slack"
	"github.com/opensaasstudio/meerkat/adapter/authorization"
	"github.com/opensaasstudio/meerkat/adapter/dynamodb"
	"github.com/opensaasstudio/meerkat/adapter/inmemory"
	"github.com/opensaasstudio/meerkat/adapter/slack"
	"github.com/opensaasstudio/meerkat/application"
	"github.com/opensaasstudio/meerkat/domain"
	"go.uber.org/zap"
)

type Config struct {
	SlackToken             string `envconfig:"SLACK_TOKEN" required:"true"`
	SlackVerificationToken string `envconfig:"SLACK_VARIFICATION_TOKEN" required:"true"`
	SlackListenAddr        string `envconfig:"SLACK_LISTEN_ADDR" required:"true"`
	AWSRegion              string `envconfig:"AWS_REGION" default:"ap-northeast-1"`
	UseInmemory            bool   `envconfig:"USE_INMEMORY" default:"false"`
	DynamoDBTablePrefix    string `envconfig:"DYNAMODB_TABLE_PREFIX" default:"meerkat"`
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	var conf Config
	if err := envconfig.Process("MEERKAT", &conf); err != nil {
		panic(err)
	}

	ulidProvider := application.NewULIDProvider()

	slackClient := slacklib.New(
		conf.SlackToken,
		slacklib.OptionDebug(true),
	)

	var paramStore slack.ParamStore

	var questionnaireRepository application.QuestionnaireRepository
	var notificationTargetRepository application.NotificationTargetRepository
	var answererRepository application.AnswererRepository
	var answerRepository application.AnswerRepository

	var questionnaireSearcher domain.QuestionnaireSearcher
	// var notificationTargetSearcher domain.NotificationTargetSearcher
	var answererSearcher domain.AnswererSearcher

	if conf.UseInmemory {
		paramStore = slack.NewInmemoryStore()

		questionnaireStore := inmemory.NewQuestionnaireStore()
		notificationTargetStore := inmemory.NewNotificationTargetStore()
		answererStore := inmemory.NewAnswererStore()
		answerStore := inmemory.NewAnswerStore()

		questionnaireRepository = questionnaireStore
		notificationTargetRepository = notificationTargetStore
		answererRepository = answererStore
		answerRepository = answerStore

		questionnaireSearcher = questionnaireStore
		// notificationTargetSearcher = notificationTargetStore
		answererSearcher = answererStore
	} else {
		dynamoDBClient := dynamodblib.New(session.New(), &aws.Config{
			Region: aws.String(conf.AWSRegion),
		})

		paramStore = dynamodb.NewParamStore(dynamoDBClient, conf.DynamoDBTablePrefix+"_param_store")

		questionnaireStore := dynamodb.NewQuestionnaireStore(dynamoDBClient, conf.DynamoDBTablePrefix+"_questionnaire")
		notificationTargetStore := dynamodb.NewNotificationTargetStore(dynamoDBClient, conf.DynamoDBTablePrefix+"_notificationtarget")
		answererStore := dynamodb.NewAnswererStore(dynamoDBClient, conf.DynamoDBTablePrefix+"_answerer")
		answerStore := dynamodb.NewAnswerStore(dynamoDBClient, conf.DynamoDBTablePrefix+"_answer")

		questionnaireRepository = questionnaireStore
		notificationTargetRepository = notificationTargetStore
		answererRepository = answererStore
		answerRepository = answerStore

		questionnaireSearcher = questionnaireStore
		// notificationTargetSearcher = notificationTargetStore
		answererSearcher = answererStore
	}

	creatingQuestionnaireUsecase := application.NewCreatingQuestionnaireUsecase(
		questionnaireRepository,
		ulidProvider,
		ulidProvider,
		authorization.NewService(),
	)

	creatingAnswererUsecase := application.NewCreatingAnswererUsecase(
		answererRepository,
		ulidProvider,
	)

	addingAnswererUsecase := application.NewAddingAnswererUsecase(
		notificationTargetRepository,
		ulidProvider,
	)

	answeringUsecase := application.NewAnsweringUsecase(
		questionnaireRepository,
		answererRepository,
		answerRepository,
		ulidProvider,
		application.NewPostingService(slack.NewPoster(slackClient).Post),
	)

	httpHandler := slack.NewHTTPHandler(
		conf.SlackVerificationToken,
		logger,
		questionnaireSearcher,
		answererSearcher,
		slack.NewCreatingQuestionnaireHandler(slackClient, creatingQuestionnaireUsecase),
		slack.NewCreatingAnswererHandler(slackClient, creatingAnswererUsecase),
		slack.NewAddingAnswererHandler(slackClient, addingAnswererUsecase),
		slack.NewAnsweringHandler(slackClient, answeringUsecase),
		paramStore,
	)

	http.HandleFunc("/slack/suggestion", httpHandler.HandleSuggestion)
	http.HandleFunc("/slack/event", httpHandler.HandleEvent)
	http.HandleFunc("/slack/interactive_component", httpHandler.HandleInteractiveComponent)

	err = http.ListenAndServe(conf.SlackListenAddr, nil)
	if err != nil {
		panic(err)
	}
}
