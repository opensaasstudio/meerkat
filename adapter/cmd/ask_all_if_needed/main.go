package main

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	dynamodblib "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/kelseyhightower/envconfig"
	slacklib "github.com/nlopes/slack"
	"github.com/opensaasstudio/meerkat/adapter/dynamodb"
	"github.com/opensaasstudio/meerkat/adapter/slack"
	"github.com/opensaasstudio/meerkat/application"
	"github.com/opensaasstudio/meerkat/domain"
	"gopkg.in/guregu/null.v3"
)

type Config struct {
	SlackToken          string `envconfig:"SLACK_TOKEN" required:"true"`
	AWSRegion           string `envconfig:"AWS_REGION" default:"ap-northeast-1"`
	DynamoDBTablePrefix string `envconfig:"DYNAMODB_TABLE_PREFIX" default:"meerkat"`
}

func main() {
	// logger, err := zap.NewDevelopment()
	// if err != nil {
	// 	panic(err)
	// }

	var conf Config
	if err := envconfig.Process("MEERKAT", &conf); err != nil {
		panic(err)
	}

	slackClient := slacklib.New(
		conf.SlackToken,
		slacklib.OptionDebug(true),
	)
	var questionnaireRepository application.QuestionnaireRepository

	var questionnaireSearcher domain.QuestionnaireSearcher
	var notificationTargetSearcher domain.NotificationTargetSearcher

	dynamoDBClient := dynamodblib.New(session.New(), &aws.Config{
		Region: aws.String(conf.AWSRegion),
	})

	paramStore := dynamodb.NewParamStore(dynamoDBClient, conf.DynamoDBTablePrefix+"_param_store")
	questionnaireStore := dynamodb.NewQuestionnaireStore(dynamoDBClient, conf.DynamoDBTablePrefix+"_questionnaire")
	notificationTargetStore := dynamodb.NewNotificationTargetStore(dynamoDBClient, conf.DynamoDBTablePrefix+"_notificationtarget")

	questionnaireRepository = questionnaireStore

	questionnaireSearcher = questionnaireStore
	notificationTargetSearcher = notificationTargetStore

	answeringHandler := slack.NewAnsweringHandler(slackClient, application.AnsweringUsecase{})

	notificationService := application.NewNotificationService(func(
		ctx context.Context,
		target domain.NotificationTargetSlack,
		questionnaire domain.Questionnaire,
	) domain.Error {
		answers := make([]slack.Answer, len(questionnaire.QuestionItems()))
		for i, item := range questionnaire.QuestionItems() {
			answers[i] = slack.Answer{
				Question: slack.Question{
					ID:       string(item.Question().ID()),
					Text:     string(item.Question().Text()),
					Required: item.Required(),
				},
			}
		}
		input := slack.AnsweringHandlerInput{
			QuestionnaireID:    string(questionnaire.ID()),
			QuestionnaireTitle: questionnaire.Title(),
			AnswererID:         string(target.AnswererID()),
			Answers:            answers,
		}
		callbackID := "Answering_" + strconv.FormatInt(time.Now().UnixNano(), 10)
		if err := paramStore.Store(context.TODO(), callbackID, input, 30*time.Minute); err != nil {
			panic(err)
		}
		return answeringHandler.RequestInput(ctx, target.ChannelID(), null.String{}, callbackID, input)
	})

	lastExecutedRecorder := application.NewLastExecutedRecorder(
		questionnaireRepository,
	)

	askingUsecase := application.NewAskingUsecase(domain.NewAskingService(
		questionnaireSearcher,
		notificationTargetSearcher,
		notificationService,
		lastExecutedRecorder,
	))

	ctx := context.Background()

	if err := askingUsecase.AskAllIfNeeded(ctx); err != nil {
		panic(err)
	}
}
