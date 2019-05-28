package dynamodb

import (
	"bytes"
	"context"
	"encoding/gob"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/opensaasstudio/meerkat/domain"
)

type ParamStore struct {
	tableName      string
	dynamoDBClient *dynamodb.DynamoDB
}

func NewParamStore(
	dynamoDBClient *dynamodb.DynamoDB,
	dynamoDBTable string,
) ParamStore {
	return ParamStore{
		tableName:      dynamoDBTable,
		dynamoDBClient: dynamoDBClient,
	}
}

func (s *ParamStore) CreateTable(ctx context.Context) domain.Error {
	if _, err := s.dynamoDBClient.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		TableName:   aws.String(s.tableName),
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("key"),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("key"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
		},
	}); err != nil {
		return domain.ErrorUnknown(err)
	}
	return nil
}

func (s ParamStore) Store(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	b := new(bytes.Buffer)
	err := gob.NewEncoder(b).Encode(value)
	if err != nil {
		return err
	}
	_, err = s.dynamoDBClient.PutItemWithContext(
		ctx,
		&dynamodb.PutItemInput{
			TableName: aws.String(s.tableName),
			Item: map[string]*dynamodb.AttributeValue{
				"key": {
					S: aws.String(key),
				},
				"value": {
					B: b.Bytes(),
				},
				"ttl": {
					N: aws.String(strconv.FormatInt(time.Now().Add(expiration).Unix(), 10)),
				},
			},
		},
	)
	return err
}

func (s ParamStore) Restore(ctx context.Context, key string, valuePtr interface{}) error {
	output, err := s.dynamoDBClient.GetItemWithContext(
		ctx,
		&dynamodb.GetItemInput{
			TableName: aws.String(s.tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"key": {
					S: aws.String(key),
				},
			},
		},
	)
	if err != nil {
		if err, ok := err.(awserr.Error); ok {
			if err.Code() == dynamodb.ErrCodeResourceNotFoundException {
				return nil
			}
		}
		return err
	}
	v, ok := output.Item["value"]
	if !ok {
		return nil
	}
	return gob.NewDecoder(bytes.NewReader(v.B)).Decode(valuePtr)
}
