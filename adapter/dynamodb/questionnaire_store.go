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
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"github.com/hori-ryota/zaperr"
	"github.com/opensaasstudio/meerkat/domain"
)

type QuestionnaireStore struct {
	db        *dynamo.DB
	tableName string
	nowFunc   func() time.Time
}

func NewQuestionnaireStore(dynamoDBClient *dynamodb.DynamoDB, tableName string) *QuestionnaireStore {
	return &QuestionnaireStore{
		db:        dynamo.NewFromIface(dynamoDBClient),
		tableName: tableName,
		nowFunc:   time.Now,
	}
}

type questionnaireItem struct {
	Version int
	domain.QuestionnaireValue
}

func (r *QuestionnaireStore) CreateTable(ctx context.Context) domain.Error {
	if err := r.db.CreateTable(r.tableName, questionnaireItem{}).OnDemand(true).RunWithContext(ctx); err != nil {
		return domain.ErrorUnknown(err)
	}
	return nil
}

func (r *QuestionnaireStore) Create(ctx context.Context, questionnaire domain.Questionnaire) domain.Error {
	_, _, err := r.FindByID(ctx, questionnaire.ID())
	if err == nil {
		return domain.ErrorBadRequest(zaperr.New("duplicated error"))
	}
	if !err.IsNotFound() {
		return err
	}
	if err := r.db.Table(r.tableName).Put(questionnaireItem{
		Version:            1,
		QuestionnaireValue: questionnaire.Dump(),
	}).RunWithContext(ctx); err != nil {
		return domain.ErrorUnknown(err)
	}
	return nil
}

func (r *QuestionnaireStore) Update(ctx context.Context, questionnaire domain.Questionnaire, version int) domain.Error {
	_, storedVersion, err := r.FindByID(ctx, questionnaire.ID())
	if err != nil {
		return err
	}
	if storedVersion != version {
		return domain.ErrorBadRequest(zaperr.New("version mismatch"))
	}
	if err := r.db.Table(r.tableName).Put(questionnaireItem{
		Version:            version + 1,
		QuestionnaireValue: questionnaire.Dump(),
	}).RunWithContext(ctx); err != nil {
		return domain.ErrorUnknown(err)
	}
	return nil
}

func (r *QuestionnaireStore) Delete(ctx context.Context, questionnaire domain.Questionnaire) domain.Error {
	_, _, err := r.FindByID(ctx, questionnaire.ID())
	if err != nil {
		return err
	}
	if err := r.db.Table(r.tableName).Delete("ID", questionnaire.ID()).RunWithContext(ctx); err != nil {
		return domain.ErrorUnknown(err)
	}
	return nil
}

func (r *QuestionnaireStore) FindByID(ctx context.Context, id domain.QuestionnaireID) (questionnaire domain.Questionnaire, version int, derr domain.Error) {
	var item questionnaireItem
	if err := r.db.Table(r.tableName).Get("ID", id).OneWithContext(ctx, &item); err != nil {
		if err == dynamo.ErrNotFound {
			return domain.Questionnaire{}, 0, domain.ErrorNotFound(err)
		}
		return domain.Questionnaire{}, 0, domain.ErrorUnknown(err)
	}

	return domain.RestoreQuestionnaireFromDumped(item.QuestionnaireValue), item.Version, nil
}

func (r *QuestionnaireStore) SearchExecutionNeeded(ctx context.Context) ([]domain.Questionnaire, domain.Error) {
	var items []questionnaireItem
	if err := r.db.Table(r.tableName).Scan().Filter("NextTime <= ?", r.nowFunc()).AllWithContext(ctx, &items); err != nil {
		if err == dynamo.ErrNotFound {
			return nil, domain.ErrorNotFound(err)
		}
		return nil, domain.ErrorUnknown(err)
	}

	questionnaires := make([]domain.Questionnaire, len(items))
	for i, item := range items {
		questionnaires[i] = domain.RestoreQuestionnaireFromDumped(item.QuestionnaireValue)
	}

	return questionnaires, nil
}

func (r *QuestionnaireStore) FetchAll(ctx context.Context) ([]domain.Questionnaire, domain.Error) {
	var items []questionnaireItem
	if err := r.db.Table(r.tableName).Scan().AllWithContext(ctx, &items); err != nil {
		if err == dynamo.ErrNotFound {
			return nil, domain.ErrorNotFound(err)
		}
		return nil, domain.ErrorUnknown(err)
	}

	questionnaires := make([]domain.Questionnaire, len(items))
	for i, item := range items {
		questionnaires[i] = domain.RestoreQuestionnaireFromDumped(item.QuestionnaireValue)
	}

	return questionnaires, nil
}
