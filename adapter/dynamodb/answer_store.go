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

package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"github.com/hori-ryota/zaperr"
	"github.com/opensaasstudio/meerkat/domain"
)

type AnswerStore struct {
	db        *dynamo.DB
	tableName string
}

func NewAnswerStore(dynamoDBClient *dynamodb.DynamoDB, tableName string) *AnswerStore {
	return &AnswerStore{
		db:        dynamo.NewFromIface(dynamoDBClient),
		tableName: tableName,
	}
}

type answerItem struct {
	Version int
	domain.AnswerValue
}

func (r *AnswerStore) CreateTable(ctx context.Context) domain.Error {
	if err := r.db.CreateTable(r.tableName, answerItem{}).OnDemand(true).RunWithContext(ctx); err != nil {
		return domain.ErrorUnknown(err)
	}
	return nil
}

func (r *AnswerStore) Create(ctx context.Context, answer domain.Answer) domain.Error {
	_, _, err := r.FindByID(ctx, answer.ID())
	if err == nil {
		return domain.ErrorBadRequest(zaperr.New("duplicated error"))
	}
	if !err.IsNotFound() {
		return err
	}
	if err := r.db.Table(r.tableName).Put(answerItem{
		Version:     1,
		AnswerValue: answer.Dump(),
	}).RunWithContext(ctx); err != nil {
		return domain.ErrorUnknown(err)
	}
	return nil
}

func (r *AnswerStore) Update(ctx context.Context, answer domain.Answer, version int) domain.Error {
	_, storedVersion, err := r.FindByID(ctx, answer.ID())
	if err != nil {
		return err
	}
	if storedVersion != version {
		return domain.ErrorBadRequest(zaperr.New("version mismatch"))
	}
	if err := r.db.Table(r.tableName).Put(answerItem{
		Version:     version + 1,
		AnswerValue: answer.Dump(),
	}).RunWithContext(ctx); err != nil {
		return domain.ErrorUnknown(err)
	}
	return nil
}

func (r *AnswerStore) Delete(ctx context.Context, answer domain.Answer) domain.Error {
	_, _, err := r.FindByID(ctx, answer.ID())
	if err != nil {
		return err
	}
	if err := r.db.Table(r.tableName).Delete("ID", answer.ID()).RunWithContext(ctx); err != nil {
		return domain.ErrorUnknown(err)
	}
	return nil
}

func (r *AnswerStore) FindByID(ctx context.Context, id domain.AnswerID) (answer domain.Answer, version int, derr domain.Error) {
	var item answerItem
	if err := r.db.Table(r.tableName).Get("ID", id).OneWithContext(ctx, &item); err != nil {
		if err == dynamo.ErrNotFound {
			return domain.Answer{}, 0, domain.ErrorNotFound(err)
		}
		return domain.Answer{}, 0, domain.ErrorUnknown(err)
	}

	return domain.RestoreAnswerFromDumped(item.AnswerValue), item.Version, nil
}

func (r *AnswerStore) FetchAll(ctx context.Context) ([]domain.Answer, domain.Error) {
	var items []answerItem
	if err := r.db.Table(r.tableName).Scan().AllWithContext(ctx, &items); err != nil {
		if err == dynamo.ErrNotFound {
			return nil, domain.ErrorNotFound(err)
		}
		return nil, domain.ErrorUnknown(err)
	}

	answers := make([]domain.Answer, len(items))
	for i, item := range items {
		answers[i] = domain.RestoreAnswerFromDumped(item.AnswerValue)
	}

	return answers, nil
}
