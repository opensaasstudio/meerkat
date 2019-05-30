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
	"sync"

	"github.com/hori-ryota/zaperr"
	"github.com/opensaasstudio/meerkat/domain"
)

type NotificationTargetStore struct {
	items *sync.Map
}

func NewNotificationTargetStore() *NotificationTargetStore {
	return &NotificationTargetStore{
		items: new(sync.Map),
	}
}

type notificationTargetItem struct {
	version int
	value   domain.NotificationTarget
}

func (r *NotificationTargetStore) Create(ctx context.Context, notificationTarget domain.NotificationTarget) domain.Error {
	_, _, err := r.FindByID(ctx, notificationTarget.ID())
	if err == nil {
		return domain.ErrorBadRequest(zaperr.New("duplicated error"))
	}
	if !err.IsNotFound() {
		return err
	}
	r.items.Store(notificationTarget.ID(), notificationTargetItem{
		version: 1,
		value:   notificationTarget,
	})
	return nil
}

func (r *NotificationTargetStore) Update(ctx context.Context, notificationTarget domain.NotificationTarget, version int) domain.Error {
	_, storedVersion, err := r.FindByID(ctx, notificationTarget.ID())
	if err != nil {
		return err
	}
	if storedVersion != version {
		return domain.ErrorBadRequest(zaperr.New("version mismatch"))
	}
	r.items.Store(notificationTarget.ID(), notificationTargetItem{
		version: version + 1,
		value:   notificationTarget,
	})
	return nil
}

func (r *NotificationTargetStore) Delete(ctx context.Context, notificationTarget domain.NotificationTarget) domain.Error {
	_, _, err := r.FindByID(ctx, notificationTarget.ID())
	if err != nil {
		return err
	}
	r.items.Delete(notificationTarget.ID())
	return nil
}

func (r *NotificationTargetStore) FindByID(ctx context.Context, id domain.NotificationTargetID) (notificationTarget domain.NotificationTarget, version int, derr domain.Error) {
	item, ok := r.items.Load(id)
	if !ok {
		return nil, 0, domain.ErrorNotFound(zaperr.New("notfound"))
	}
	return item.(notificationTargetItem).value, item.(notificationTargetItem).version, nil
}

func (r *NotificationTargetStore) SearchByQuestionnaireIDAndAnswererID(ctx context.Context, questionnaireID domain.QuestionnaireID, answererID domain.AnswererID) ([]domain.NotificationTarget, domain.Error) {
	list := make([]domain.NotificationTarget, 0, 10)
	r.items.Range(func(key interface{}, item interface{}) bool {
		nt := item.(notificationTargetItem).value
		if nt.QuestionnaireID() == questionnaireID && nt.AnswererID() == answererID {
			list = append(list, nt)
		}
		return true
	})
	return list, nil
}

func (r *NotificationTargetStore) SearchByQuestionnaireID(ctx context.Context, questionnaireID domain.QuestionnaireID) ([]domain.NotificationTarget, domain.Error) {
	list := make([]domain.NotificationTarget, 0, 10)
	r.items.Range(func(key interface{}, item interface{}) bool {
		nt := item.(notificationTargetItem).value
		if nt.QuestionnaireID() == questionnaireID {
			list = append(list, nt)
		}
		return true
	})
	return list, nil
}
