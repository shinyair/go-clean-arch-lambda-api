package dynamodbrepo

import (
	"context"
	"fmt"

	"local.com/go-clean-lambda/internal/domain"
	"local.com/go-clean-lambda/internal/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	FieldDummyPK string = "pk"
	FieldDummySK string = "sk"
)

// DummyDynamodbRepo
type DummyDynamodbRepo struct {
	tableName string
	client    *dynamodb.DynamoDB
}

// NewDummyDynamodbRepo
//  @param tableName
//  @param client
//  @return *DummyDynamodbRepo
func NewDummyDynamodbRepo(tableName string, client *dynamodb.DynamoDB) *DummyDynamodbRepo {
	return &DummyDynamodbRepo{
		tableName: tableName,
		client:    client,
	}
}

// GetById
//  @receiver repo
//  @param ctx
//  @param id
//  @return *domain.Dummy
//  @return error
func (repo *DummyDynamodbRepo) GetById(ctx context.Context, id string) (*domain.Dummy, error) {
	if len(id) == 0 {
		return nil, nil
	}
	data, err := repo.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key:       ToDummyDbKey(domain.ToKeyDummy(id)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get item from db. table: %s, id: %s. %w", repo.tableName, id, err)
	}
	logger.Debug("get from db. item: %s", logger.Pretty(data.Item))
	return ToDummyEntity(data.Item)
}

// Insert
//  @receiver repo
//  @param ctx
//  @param dummy
//  @return *domain.Dummy
//  @return error
func (repo *DummyDynamodbRepo) Insert(ctx context.Context, dummy *domain.Dummy) (*domain.Dummy, error) {
	if dummy == nil || len(dummy.ID) == 0 {
		return nil, nil
	}
	_, err := repo.client.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(repo.tableName),
		Item:      ToDummyDbItem(dummy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to put item to db. table: %s, item: %s. %w", repo.tableName, logger.Pretty(dummy), err)
	}
	logger.Debug("put to db. item: %s", logger.Pretty(dummy))
	return dummy, nil
}

// DeleteById
//  @receiver repo
//  @param ctx
//  @param id
//  @return *domain.Dummy
//  @return error
func (repo *DummyDynamodbRepo) DeleteById(ctx context.Context, id string) (*domain.Dummy, error) {
	if len(id) == 0 {
		return nil, nil
	}
	_, err := repo.client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(repo.tableName),
		Key:       ToDummyDbKey(domain.ToKeyDummy(id)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to delete item from db. table: %s, id: %s. %w", repo.tableName, id, err)
	}
	logger.Debug("delete from db. id: %s", id)
	return &domain.Dummy{
		ID: id,
	}, nil
}

// ToDbKey
//
//  @param dummy
//  @return map
func ToDummyDbKey(dummy *domain.Dummy) map[string]*dynamodb.AttributeValue {
	if dummy == nil {
		return make(map[string]*dynamodb.AttributeValue)
	}
	item := make(map[string]*dynamodb.AttributeValue)
	return addDummyKeys(item, dummy)
}

// ToDbItem
//
//  @param dummy
//  @return map
func ToDummyDbItem(dummy *domain.Dummy) map[string]*dynamodb.AttributeValue {
	if dummy == nil {
		return make(map[string]*dynamodb.AttributeValue)
	}
	// TODO: error handle
	item, _ := dynamodbattribute.MarshalMap(dummy)
	item = addDummyKeys(item, dummy)
	return item
}

// ToEntity
//
//  @param item
//  @return *domain.Dummy
//  @return error
func ToDummyEntity(item map[string]*dynamodb.AttributeValue) (*domain.Dummy, error) {
	if len(item) == 0 {
		return nil, nil
	}
	dummy := &domain.Dummy{}
	err := dynamodbattribute.UnmarshalMap(item, dummy)
	return dummy, err
}

// addKeys
// set pk and sk into db item map according to entity
//
//  @param item
//  @param dummy
//  @return map
func addDummyKeys(item map[string]*dynamodb.AttributeValue, dummy *domain.Dummy) map[string]*dynamodb.AttributeValue {
	item[FieldDummyPK] = &dynamodb.AttributeValue{S: aws.String("test")}
	item[FieldDummySK] = &dynamodb.AttributeValue{S: aws.String(dummy.ID)}
	return item
}
