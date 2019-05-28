package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"github.com/hori-ryota/zaperr"
	"github.com/opensaasstudio/meerkat/domain"
)

type NotificationTargetStore struct {
	db        *dynamo.DB
	tableName string
}

func NewNotificationTargetStore(dynamoDBClient *dynamodb.DynamoDB, tableName string) *NotificationTargetStore {
	return &NotificationTargetStore{
		db:        dynamo.NewFromIface(dynamoDBClient),
		tableName: tableName,
	}
}

type notificationTargetItem struct {
	Version int
	domain.NotificationTargetValue
}

func (r *NotificationTargetStore) CreateTable(ctx context.Context) domain.Error {
	if err := r.db.CreateTable(r.tableName, notificationTargetItem{}).OnDemand(true).RunWithContext(ctx); err != nil {
		return domain.ErrorUnknown(err)
	}
	return nil
}

func (r *NotificationTargetStore) Create(ctx context.Context, notificationTarget domain.NotificationTarget) domain.Error {
	_, _, err := r.FindByID(ctx, notificationTarget.ID())
	if err == nil {
		return domain.ErrorBadRequest(zaperr.New("duplicated error"))
	}
	if !err.IsNotFound() {
		return err
	}
	if err := r.db.Table(r.tableName).Put(notificationTargetItem{
		Version:                 1,
		NotificationTargetValue: notificationTarget.Dump(),
	}).RunWithContext(ctx); err != nil {
		return domain.ErrorUnknown(err)
	}
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
	if err := r.db.Table(r.tableName).Put(notificationTargetItem{
		Version:                 version + 1,
		NotificationTargetValue: notificationTarget.Dump(),
	}).RunWithContext(ctx); err != nil {
		return domain.ErrorUnknown(err)
	}
	return nil
}

func (r *NotificationTargetStore) Delete(ctx context.Context, notificationTarget domain.NotificationTarget) domain.Error {
	_, _, err := r.FindByID(ctx, notificationTarget.ID())
	if err != nil {
		return err
	}
	if err := r.db.Table(r.tableName).Delete("ID", notificationTarget.ID()).RunWithContext(ctx); err != nil {
		return domain.ErrorUnknown(err)
	}
	return nil
}

func (r *NotificationTargetStore) FindByID(ctx context.Context, id domain.NotificationTargetID) (notificationTarget domain.NotificationTarget, version int, derr domain.Error) {
	var item notificationTargetItem
	if err := r.db.Table(r.tableName).Get("ID", id).OneWithContext(ctx, &item); err != nil {
		if err == dynamo.ErrNotFound {
			return nil, 0, domain.ErrorNotFound(err)
		}
		return nil, 0, domain.ErrorUnknown(err)
	}

	return domain.RestoreNotificationTargetFromDumped(item.NotificationTargetValue), item.Version, nil
}

func (r *NotificationTargetStore) SearchByQuestionnaireIDAndAnswererID(ctx context.Context, questionnaireID domain.QuestionnaireID, answererID domain.AnswererID) ([]domain.NotificationTarget, domain.Error) {
	var items []notificationTargetItem
	if err := r.db.Table(r.tableName).Scan().Filter("QuestionnaireID = ?", questionnaireID).Filter("AnswererID = ?", answererID).AllWithContext(ctx, &items); err != nil {
		if err == dynamo.ErrNotFound {
			return nil, domain.ErrorNotFound(err)
		}
		return nil, domain.ErrorUnknown(err)
	}

	notificationTargets := make([]domain.NotificationTarget, len(items))
	for i, item := range items {
		notificationTargets[i] = domain.RestoreNotificationTargetFromDumped(item.NotificationTargetValue)
	}

	return notificationTargets, nil
}

func (r *NotificationTargetStore) SearchByQuestionnaireID(ctx context.Context, questionnaireID domain.QuestionnaireID) ([]domain.NotificationTarget, domain.Error) {
	var items []notificationTargetItem
	if err := r.db.Table(r.tableName).Scan().Filter("QuestionnaireID = ?", questionnaireID).AllWithContext(ctx, &items); err != nil {
		if err == dynamo.ErrNotFound {
			return nil, domain.ErrorNotFound(err)
		}
		return nil, domain.ErrorUnknown(err)
	}

	notificationTargets := make([]domain.NotificationTarget, len(items))
	for i, item := range items {
		notificationTargets[i] = domain.RestoreNotificationTargetFromDumped(item.NotificationTargetValue)
	}

	return notificationTargets, nil
}
