// Copyright 2019 The meerkat Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	dynamodblib "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/kelseyhightower/envconfig"
	slacklib "github.com/nlopes/slack"
	"github.com/opensaasstudio/meerkat/adapter/authorization"
	"github.com/opensaasstudio/meerkat/adapter/dynamodb"
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

	creatingQuestionnaireUsecase := application.NewCreatingQuestionnaireUsecase(
		questionnaireRepository,
		ulidProvider,
		ulidProvider,
		authorization.NewService(),
	)
	updatingQuestionnaireUsecase := application.NewUpdatingQuestionnaireUsecase(
		questionnaireRepository,
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
		slack.NewEditingQuestionnaireHandler(slackClient, creatingQuestionnaireUsecase, updatingQuestionnaireUsecase),
		slack.NewCreatingAnswererHandler(slackClient, creatingAnswererUsecase),
		slack.NewAddingAnswererHandler(slackClient, addingAnswererUsecase),
		slack.NewAnsweringHandler(slackClient, answeringUsecase),
		paramStore,
	)

	lambda.Start(func(ctx context.Context, param LambdaRequest) (LambdaResponse, error) {
		w := &LambdaResponseWriter{
			Body:              new(bytes.Buffer),
			StatusCode:        http.StatusOK,
			MultiValueHeaders: http.Header{},
		}
		r, err := http.NewRequest(
			param.HTTPMethod,
			param.Path,
			ioutil.NopCloser(strings.NewReader(param.Body)),
		)
		if err != nil {
			return LambdaResponse{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
		switch param.Path {
		case "/slack/suggestion":
			httpHandler.HandleSuggestion(w, r)
		case "/slack/event":
			httpHandler.HandleEvent(w, r)
		case "/slack/interactive":
			httpHandler.HandleInteractiveComponent(w, r)
		default:
			return LambdaResponse{
				StatusCode: http.StatusNotFound,
			}, err
		}

		return w.ToResponse(), nil
	})

	err = http.ListenAndServe(conf.SlackListenAddr, nil)
	if err != nil {
		panic(err)
	}
}

type LambdaRequest struct {
	Body       string `json:"body"`
	HTTPMethod string `json:"httpMethod"`
	Path       string `json:"path"`
}

type LambdaResponse struct {
	Body              string              `json:"body"`
	StatusCode        int                 `json:"statusCode"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
	IsBase64Encoded   bool                `json:"isBase64Encoded"`
}

type LambdaResponseWriter struct {
	Body              *bytes.Buffer
	StatusCode        int
	MultiValueHeaders http.Header
}

func (w *LambdaResponseWriter) Header() http.Header {
	return w.MultiValueHeaders
}

func (w *LambdaResponseWriter) Write(b []byte) (int, error) {
	return w.Body.Write(b)
}

func (w *LambdaResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

func (w LambdaResponseWriter) ToResponse() LambdaResponse {
	return LambdaResponse{
		Body:              w.Body.String(),
		StatusCode:        w.StatusCode,
		MultiValueHeaders: w.MultiValueHeaders,
		IsBase64Encoded:   false,
	}
}
