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
