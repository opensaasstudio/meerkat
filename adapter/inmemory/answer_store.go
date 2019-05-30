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

	"github.com/hori-ryota/zaperr"
	"github.com/opensaasstudio/meerkat/domain"
)

type AnswerStore struct {
	items *sync.Map
}

func NewAnswerStore() *AnswerStore {
	return &AnswerStore{
		items: new(sync.Map),
	}
}

type answerItem struct {
	version int
	value   domain.Answer
}

func (r *AnswerStore) Create(ctx context.Context, answer domain.Answer) domain.Error {
	_, _, err := r.FindByID(ctx, answer.ID())
	if err == nil {
		return domain.ErrorBadRequest(zaperr.New("duplicated error"))
	}
	if !err.IsNotFound() {
		return err
	}
	r.items.Store(answer.ID(), answerItem{
		version: 1,
		value:   answer,
	})
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
	r.items.Store(answer.ID(), answerItem{
		version: version + 1,
		value:   answer,
	})
	return nil
}

func (r *AnswerStore) Delete(ctx context.Context, answer domain.Answer) domain.Error {
	_, _, err := r.FindByID(ctx, answer.ID())
	if err != nil {
		return err
	}
	r.items.Delete(answer.ID())
	return nil
}

func (r *AnswerStore) FindByID(ctx context.Context, id domain.AnswerID) (answer domain.Answer, version int, derr domain.Error) {
	item, ok := r.items.Load(id)
	if !ok {
		return domain.Answer{}, 0, domain.ErrorNotFound(zaperr.New("notfound"))
	}
	return item.(answerItem).value, item.(answerItem).version, nil
}

func (r *AnswerStore) FetchAll(ctx context.Context) ([]domain.Answer, domain.Error) {
	list := make([]domain.Answer, 0, 10)
	r.items.Range(func(key interface{}, item interface{}) bool {
		q := item.(answerItem).value
		list = append(list, q)
		return true
	})
	sort.Slice(list, func(i, j int) bool {
		return string(list[i].ID()) < string(list[j].ID())
	})
	return list, nil
}
