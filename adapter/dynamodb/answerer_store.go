package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"github.com/hori-ryota/zaperr"
	"github.com/opensaasstudio/meerkat/domain"
)

type AnswererStore struct {
	db        *dynamo.DB
	tableName string
}

func NewAnswererStore(dynamoDBClient *dynamodb.DynamoDB, tableName string) *AnswererStore {
	return &AnswererStore{
		db:        dynamo.NewFromIface(dynamoDBClient),
		tableName: tableName,
	}
}

type answererItem struct {
	Version int
	domain.AnswererValue
}

func (r *AnswererStore) CreateTable(ctx context.Context) domain.Error {
	if err := r.db.CreateTable(r.tableName, answererItem{}).OnDemand(true).RunWithContext(ctx); err != nil {
		return domain.ErrorUnknown(err)
	}
	return nil
}

func (r *AnswererStore) Create(ctx context.Context, answerer domain.Answerer) domain.Error {
	_, _, err := r.FindByID(ctx, answerer.ID())
	if err == nil {
		return domain.ErrorBadRequest(zaperr.New("duplicated error"))
	}
	if !err.IsNotFound() {
		return err
	}
	if err := r.db.Table(r.tableName).Put(answererItem{
		Version:       1,
		AnswererValue: answerer.Dump(),
	}).RunWithContext(ctx); err != nil {
		return domain.ErrorUnknown(err)
	}
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
	if err := r.db.Table(r.tableName).Put(answererItem{
		Version:       version + 1,
		AnswererValue: answerer.Dump(),
	}).RunWithContext(ctx); err != nil {
		return domain.ErrorUnknown(err)
	}
	return nil
}

func (r *AnswererStore) Delete(ctx context.Context, answerer domain.Answerer) domain.Error {
	_, _, err := r.FindByID(ctx, answerer.ID())
	if err != nil {
		return err
	}
	if err := r.db.Table(r.tableName).Delete("ID", answerer.ID()).RunWithContext(ctx); err != nil {
		return domain.ErrorUnknown(err)
	}
	return nil
}

func (r *AnswererStore) FindByID(ctx context.Context, id domain.AnswererID) (answerer domain.Answerer, version int, derr domain.Error) {
	var item answererItem
	if err := r.db.Table(r.tableName).Get("ID", id).OneWithContext(ctx, &item); err != nil {
		if err == dynamo.ErrNotFound {
			return domain.Answerer{}, 0, domain.ErrorNotFound(err)
		}
		return domain.Answerer{}, 0, domain.ErrorUnknown(err)
	}

	return domain.RestoreAnswererFromDumped(item.AnswererValue), item.Version, nil
}

func (r *AnswererStore) FetchAll(ctx context.Context) ([]domain.Answerer, domain.Error) {
	var items []answererItem
	if err := r.db.Table(r.tableName).Scan().AllWithContext(ctx, &items); err != nil {
		if err == dynamo.ErrNotFound {
			return nil, domain.ErrorNotFound(err)
		}
		return nil, domain.ErrorUnknown(err)
	}

	answerers := make([]domain.Answerer, len(items))
	for i, item := range items {
		answerers[i] = domain.RestoreAnswererFromDumped(item.AnswererValue)
	}

	return answerers, nil
}
