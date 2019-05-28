package inmemory

import (
	"context"
	"sort"
	"sync"

	"github.com/hori-ryota/zaperr"
	"github.com/opensaasstudio/meerkat/domain"
)

type AnswererStore struct {
	items *sync.Map
}

func NewAnswererStore() *AnswererStore {
	return &AnswererStore{
		items: new(sync.Map),
	}
}

type answererItem struct {
	version int
	value   domain.Answerer
}

func (r *AnswererStore) Create(ctx context.Context, answerer domain.Answerer) domain.Error {
	_, _, err := r.FindByID(ctx, answerer.ID())
	if err == nil {
		return domain.ErrorBadRequest(zaperr.New("duplicated error"))
	}
	if !err.IsNotFound() {
		return err
	}
	r.items.Store(answerer.ID(), answererItem{
		version: 1,
		value:   answerer,
	})
	return nil
}

func (r *AnswererStore) Update(ctx context.Context, answerer domain.Answerer, version int) domain.Error {
	_, storedVersion, err := r.FindByID(ctx, answerer.ID())
	if err != nil {
		return err
	}
	if storedVersion != version {
		return domain.ErrorBadRequest(zaperr.New("version mismatch"))
	}
	r.items.Store(answerer.ID(), answererItem{
		version: version + 1,
		value:   answerer,
	})
	return nil
}

func (r *AnswererStore) Delete(ctx context.Context, answerer domain.Answerer) domain.Error {
	_, _, err := r.FindByID(ctx, answerer.ID())
	if err != nil {
		return err
	}
	r.items.Delete(answerer.ID())
	return nil
}

func (r *AnswererStore) FindByID(ctx context.Context, id domain.AnswererID) (answerer domain.Answerer, version int, derr domain.Error) {
	item, ok := r.items.Load(id)
	if !ok {
		return domain.Answerer{}, 0, domain.ErrorNotFound(zaperr.New("notfound"))
	}
	return item.(answererItem).value, item.(answererItem).version, nil
}

func (r *AnswererStore) FetchAll(ctx context.Context) ([]domain.Answerer, domain.Error) {
	list := make([]domain.Answerer, 0, 10)
	r.items.Range(func(key interface{}, item interface{}) bool {
		q := item.(answererItem).value
		list = append(list, q)
		return true
	})
	sort.Slice(list, func(i, j int) bool {
		return string(list[i].ID()) < string(list[j].ID())
	})
	return list, nil
}
