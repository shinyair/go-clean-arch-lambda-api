//go:build integration
// +build integration

package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"local.com/go-clean-lambda/internal/domain"
	"local.com/go-clean-lambda/internal/repository"
)

const (
	dummyTableName        string = "dummy_table_name"
	invalidDummyTableName string = "invalid_dummy_table_name"
)

func buildCreateDummyTableInput() *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(repository.FieldDummyPK),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
			{
				AttributeName: aws.String(repository.FieldDummySK),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(repository.FieldDummyPK),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String(repository.FieldDummySK),
				KeyType:       aws.String(dynamodb.KeyTypeRange),
			},
		},
		TableName: aws.String(dummyTableName),
	}
}

func TestDummyGetByIDWithBadRepoReturnError(t *testing.T) {
	assert := require.New(t)
	repo := repository.NewDummyDynamodbRepo(invalidDummyTableName, ddb.client)

	id := uuid.New().String()
	item, err := repo.GetByID(context.TODO(), id)

	msg := "get dummy entity with bad repo didin't fail"
	assert.Nil(item, msg, "returned item")
	assert.NotNil(err, msg, "error not found")
}

func TestDummyGetByIDWithIdReturnEntity(t *testing.T) {
	assert := require.New(t)
	msg := "failed to get entity by valid id"
	repo := repository.NewDummyDynamodbRepo(dummyTableName, ddb.client)

	random := uuid.New().String()
	expected := &domain.Dummy{
		ID:       random,
		Name:     fmt.Sprintf("test_name_%s", random),
		SomeAttr: fmt.Sprintf("test_some_attr_%s", random),
	}
	err1 := saveDdbItems(
		dummyTableName,
		[]*domain.Dummy{expected},
		repository.ToDummyDBItem,
	)
	if err1 != nil {
		t.Fatalf("%s. error happened when preparing necesarry data, %v", msg, err1)
	}

	actual, err2 := repo.GetByID(context.TODO(), expected.ID)
	assert.Nil(err2, msg, "found error")
	assert.Equal(expected, actual, msg, "wrong entity")
}

func TestDummyInsertWithEntityReturnEntity(t *testing.T) {
	assert := require.New(t)
	msg := "failed to insert valid entity"
	repo := repository.NewDummyDynamodbRepo(dummyTableName, ddb.client)

	random := uuid.New().String()
	expected := &domain.Dummy{
		ID:       random,
		Name:     fmt.Sprintf("test_name_%s", random),
		SomeAttr: fmt.Sprintf("test_some_attr_%s", random),
	}

	actual, err1 := repo.Insert(context.TODO(), expected)

	loaded, err2 := loadDdbItems(
		dummyTableName,
		[]*domain.Dummy{expected},
		repository.ToDummyDBKey,
		repository.ToDummyEntity,
	)
	if err2 != nil {
		t.Fatalf("%s. error happened when load db data, %v", msg, err2)
	}
	assert.Nil(err1, msg, "found error")
	assert.Equal(expected, actual, msg, "wrong actual return item")
	assert.Len(loaded, 1, msg, "wrong loaded item size")
	assert.Equal(expected, loaded[0], msg, "wrong actual db item")
}
