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

package slack

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"
)

type ParamStore interface {
	Store(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Restore(ctx context.Context, key string, valuePtr interface{}) error
}

type StoreValue struct {
	ExpiredAt time.Time
	Value     []byte
}

type InmemoryStore struct {
	v map[string]StoreValue
}

func NewInmemoryStore() InmemoryStore {
	return InmemoryStore{
		v: make(map[string]StoreValue, 100),
	}
}

func (s InmemoryStore) Store(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	b := new(bytes.Buffer)
	err := gob.NewEncoder(b).Encode(value)
	if err != nil {
		return err
	}
	s.v[key] = StoreValue{
		ExpiredAt: time.Now().Add(expiration),
		Value:     b.Bytes(),
	}
	return nil
}

func (s InmemoryStore) Restore(ctx context.Context, key string, valuePtr interface{}) error {
	v, ok := s.v[key]
	if !ok {
		return nil
	}
	if time.Now().After(v.ExpiredAt) {
		return nil
	}
	return gob.NewDecoder(bytes.NewReader(v.Value)).Decode(valuePtr)
}
