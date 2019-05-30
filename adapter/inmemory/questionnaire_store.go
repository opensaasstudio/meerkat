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

package inmemory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/hori-ryota/zaperr"
	"github.com/opensaasstudio/meerkat/domain"
)

type QuestionnaireStore struct {
	items   *sync.Map
	nowFunc func() time.Time
}

func NewQuestionnaireStore() *QuestionnaireStore {
	return &QuestionnaireStore{
		items:   new(sync.Map),
		nowFunc: time.Now,
	}
}

// OverwirteNowFunc is for test
func (r *QuestionnaireStore) OverwirteNowFunc(f func() time.Time) {
	r.nowFunc = f
}

type questionnaireItem struct {
	version int
	value   domain.Questionnaire
}

func (r *QuestionnaireStore) Create(ctx context.Context, questionnaire domain.Questionnaire) domain.Error {
	_, _, err := r.FindByID(ctx, questionnaire.ID())
	if err == nil {
		return domain.ErrorBadRequest(zaperr.New("duplicated error"))
	}
	if !err.IsNotFound() {
		return err
	}
	r.items.Store(questionnaire.ID(), questionnaireItem{
		version: 1,
		value:   questionnaire,
	})
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
	r.items.Store(questionnaire.ID(), questionnaireItem{
		version: version + 1,
		value:   questionnaire,
	})
	return nil
}

func (r *QuestionnaireStore) Delete(ctx context.Context, questionnaire domain.Questionnaire) domain.Error {
	_, _, err := r.FindByID(ctx, questionnaire.ID())
	if err != nil {
		return err
	}
	r.items.Delete(questionnaire.ID())
	return nil
}

func (r *QuestionnaireStore) FindByID(ctx context.Context, id domain.QuestionnaireID) (questionnaire domain.Questionnaire, version int, derr domain.Error) {
	item, ok := r.items.Load(id)
	if !ok {
		return domain.Questionnaire{}, 0, domain.ErrorNotFound(zaperr.New("notfound"))
	}
	return item.(questionnaireItem).value, item.(questionnaireItem).version, nil
}

func (r *QuestionnaireStore) SearchExecutionNeeded(ctx context.Context) ([]domain.Questionnaire, domain.Error) {
	list := make([]domain.Questionnaire, 0, 10)
	r.items.Range(func(key interface{}, item interface{}) bool {
		q := item.(questionnaireItem).value
		if q.Schedule() == nil {
			return true
		}
		if !q.Schedule().NextTime(q.LastExecuted()).After(r.nowFunc()) {
			list = append(list, q)
		}
		return true
	})
	return list, nil
}

func (r *QuestionnaireStore) FetchAll(ctx context.Context) ([]domain.Questionnaire, domain.Error) {
	list := make([]domain.Questionnaire, 0, 10)
	r.items.Range(func(key interface{}, item interface{}) bool {
		q := item.(questionnaireItem).value
		list = append(list, q)
		return true
	})
	sort.Slice(list, func(i, j int) bool {
		return string(list[i].ID()) < string(list[j].ID())
	})
	return list, nil
}
