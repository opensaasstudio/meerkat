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
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	dynamodblib "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/kelseyhightower/envconfig"
	"github.com/opensaasstudio/meerkat/adapter/dynamodb"
)

type Config struct {
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

	dynamoDBClient := dynamodblib.New(session.New(), &aws.Config{
		Region: aws.String(conf.AWSRegion),
	})

	paramStore := dynamodb.NewParamStore(dynamoDBClient, conf.DynamoDBTablePrefix+"_param_store")
	answerStore := dynamodb.NewAnswerStore(dynamoDBClient, conf.DynamoDBTablePrefix+"_answer")
	answererStore := dynamodb.NewAnswererStore(dynamoDBClient, conf.DynamoDBTablePrefix+"_answerer")
	notificationTargetStore := dynamodb.NewNotificationTargetStore(dynamoDBClient, conf.DynamoDBTablePrefix+"_notificationtarget")
	questionnaireStore := dynamodb.NewQuestionnaireStore(dynamoDBClient, conf.DynamoDBTablePrefix+"_questionnaire")

	ctx := context.Background()

	if err := paramStore.CreateTable(ctx); err != nil {
		panic(err)
	}
	if err := answerStore.CreateTable(ctx); err != nil {
		panic(err)
	}
	if err := answererStore.CreateTable(ctx); err != nil {
		panic(err)
	}
	if err := notificationTargetStore.CreateTable(ctx); err != nil {
		panic(err)
	}
	if err := questionnaireStore.CreateTable(ctx); err != nil {
		panic(err)
	}
}
