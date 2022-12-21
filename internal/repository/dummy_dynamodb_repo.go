package repository

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	"local.com/go-clean-lambda/internal/domain"
	"local.com/go-clean-lambda/internal/logger"
)

const (
	FieldDummyPK string = "pk"
	FieldDummySK string = "sk"
)

// DummyDynamodbRepo.
type DummyDynamodbRepo struct {
	tableName string
	client    *dynamodb.DynamoDB
}

// NewDummyDynamodbRepo
//
//	@param tableName
//	@param client
//	@return *DummyDynamodbRepo
func NewDummyDynamodbRepo(tableName string, client *dynamodb.DynamoDB) *DummyDynamodbRepo {
	return &DummyDynamodbRepo{
		tableName: tableName,
		client:    client,
	}
}

// GetByID
//
//	@receiver repo
//	@param ctx
//	@param id
//	@return *domain.Dummy
//	@return error
func (repo *DummyDynamodbRepo) GetByID(ctx context.Context, id string) (*domain.Dummy, error) {
	if len(id) == 0 {
		return nil, nil
	}
	data, err := repo.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key:       ToDummyDBKey(domain.ToKeyDummy(id)),
	})
	if err != nil {
		// new an error to record stack from current position
		rootErr := errors.New(err.Error())
		return nil, errors.Wrapf(rootErr, "get db item error. table: %s, id: %s", repo.tableName, id)
	}
	logger.Debug("get from db. item: %s", logger.Pretty(data.Item))
	return ToDummyEntity(data.Item)
}

// Insert
//
//	@receiver repo
//	@param ctx
//	@param dummy
//	@return *domain.Dummy
//	@return error
func (repo *DummyDynamodbRepo) Insert(ctx context.Context, dummy *domain.Dummy) (*domain.Dummy, error) {
	if dummy == nil || len(dummy.ID) == 0 {
		return nil, nil
	}
	item, err := ToDummyDBItem(dummy)
	if err != nil {
		return nil, errors.Wrap(err, "failed build db item")
	}
	_, err = repo.client.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(repo.tableName),
		Item:      item,
	})
	if err != nil {
		rootErr := errors.New(err.Error())
		return nil, errors.Wrapf(rootErr, "put db item error. table: %s, item: %s", repo.tableName, logger.Pretty(dummy))
	}
	logger.Debug("put to db. item: %s", logger.Pretty(dummy))
	return dummy, nil
}

// DeleteByID
//
//	@receiver repo
//	@param ctx
//	@param id
//	@return *domain.Dummy
//	@return error
func (repo *DummyDynamodbRepo) DeleteByID(ctx context.Context, id string) (*domain.Dummy, error) {
	if len(id) == 0 {
		return nil, nil
	}
	_, err := repo.client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(repo.tableName),
		Key:       ToDummyDBKey(domain.ToKeyDummy(id)),
	})
	if err != nil {
		rootErr := errors.New(err.Error())
		return nil, errors.Wrapf(rootErr, "delete db item error. table: %s, id: %s", repo.tableName, id)
	}
	logger.Debug("delete from db. id: %s", id)
	return &domain.Dummy{
		ID: id,
	}, nil
}

// ToDummyDBKey
//
//	@param dummy
//	@return map
func ToDummyDBKey(dummy *domain.Dummy) map[string]*dynamodb.AttributeValue {
	if dummy == nil {
		return make(map[string]*dynamodb.AttributeValue)
	}
	item := make(map[string]*dynamodb.AttributeValue)
	return addDummyKeys(item, dummy)
}

// ToDummyDBItem
//
//	@param dummy
//	@return map
func ToDummyDBItem(dummy *domain.Dummy) (map[string]*dynamodb.AttributeValue, error) {
	if dummy == nil {
		return make(map[string]*dynamodb.AttributeValue), nil
	}
	item, err := dynamodbattribute.MarshalMap(dummy)
	if err != nil {
		rootErr := errors.New(err.Error())
		return nil, errors.Wrap(rootErr, "marshal dummy entity error")
	}
	item = addDummyKeys(item, dummy)
	return item, nil
}

// ToDummyEntity
//
//	@param item
//	@return *domain.Dummy
//	@return error
func ToDummyEntity(item map[string]*dynamodb.AttributeValue) (*domain.Dummy, error) {
	if len(item) == 0 {
		return nil, nil
	}
	dummy := &domain.Dummy{}
	err := dynamodbattribute.UnmarshalMap(item, dummy)
	if err != nil {
		rootErr := errors.New(err.Error())
		return nil, errors.Wrap(rootErr, "unmarshal dummy item error")
	}
	return dummy, nil
}

// addDummyKeys
// set pk and sk into db item map according to entity
//
//	@param item
//	@param dummy
//	@return map
func addDummyKeys(item map[string]*dynamodb.AttributeValue, dummy *domain.Dummy) map[string]*dynamodb.AttributeValue {
	item[FieldDummyPK] = &dynamodb.AttributeValue{S: aws.String("test")}
	item[FieldDummySK] = &dynamodb.AttributeValue{S: aws.String(dummy.ID)}
	return item
}
