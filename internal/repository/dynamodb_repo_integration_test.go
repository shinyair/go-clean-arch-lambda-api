//go:build integration
// +build integration

package repository_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"local.com/go-clean-lambda/internal/domain"
)

// https://mickey.dev/posts/go-build-tags-testing/

const (
	writeBatchSize int = 25
	readBatchSize  int = 25
)

var ddb struct {
	dcPath string
	client *dynamodb.DynamoDB
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	err := buildDockerComposePath()
	if err != nil {
		log.Fatal(err)
	}
	err = startDdbLocal()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("docker compose started")
	ddb.client = buildDdbClient()
	fmt.Println("dynamodb client inited")
	err = createDdbTables(ddb.client)
	if err != nil {
		fmt.Printf("failed to create dynamodb tables. caused by: %v \n", err)
		err2 := shutDdbLocal()
		if err2 != nil {
			fmt.Printf("failed to shutdown dynamodb local. caused by: %v \n", err2)
		}
		log.Fatal(err)
	}
	fmt.Println("dynamodb tables created")
	fmt.Println("dynamodb local setup completed")
}

func buildDockerComposePath() error {
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}
	configDir, err := filepath.Abs(currDir + "../../../test/dynamodb-local")
	if err != nil {
		return err
	}
	ddb.dcPath = configDir
	fmt.Printf("ddb local docker compose config: %s \n", configDir)
	return nil
}

func teardown() {
	err := shutDdbLocal()
	if err != nil {
		fmt.Printf("failed to shutdown dynamodb local, caused by: %v \n", err)
	} else {
		fmt.Println("dynamodb local teardown completed")
	}
}

func startDdbLocal() error {
	// another solution to run docker: use testcontainers
	// problem: import toooo many changes in go.sum
	// refs
	//   - github.com/testcontainers/testcontainers-go
	//   - https://golang.testcontainers.org/features/docker_compose/
	cmd := exec.Command("docker-compose", "up", "-d", "--force-recreate")
	cmd.Dir = ddb.dcPath
	cmdOutput, err := cmd.Output()
	fmt.Printf("%s std:\n    %s \n", cmd.String(), string(cmdOutput))
	return err
}

func shutDdbLocal() error {
	cmd := exec.Command("docker-compose", "down")
	cmd.Dir = ddb.dcPath
	cmdOutput, err := cmd.Output()
	fmt.Printf("%s std:\n    %s \n", cmd.String(), string(cmdOutput))
	return err
}

func buildDdbClient() *dynamodb.DynamoDB {
	awsopt := session.Options{
		Config: aws.Config{
			Region:      aws.String("local"),
			Endpoint:    aws.String("http://localhost:8000"),
			Credentials: credentials.NewStaticCredentials("dummy", "dummy", "dummy"),
		},
	}
	awssess := session.Must(session.NewSessionWithOptions(awsopt))
	return dynamodb.New(awssess)
}

func createDdbTables(client *dynamodb.DynamoDB) error {
	retry := 0
	maxRetry := 10
	var err error
	for err != nil && retry <= maxRetry {
		_, err = client.DescribeEndpoints(&dynamodb.DescribeEndpointsInput{})
		retry++
	}
	if err != nil {
		return err
	}
	_, err = client.CreateTable(buildCreateDummyTableInput())
	if err != nil {
		return err
	}
	return nil
}

func saveDdbItems[T domain.Dummy](
	tableName string,
	entities []*T,
	toDBItemFunc func(item *T) (map[string]*dynamodb.AttributeValue, error),
) error {
	writeReqs := []*dynamodb.WriteRequest{}
	writeInputs := []*dynamodb.BatchWriteItemInput{}
	for _, entity := range entities {
		dav, err := toDBItemFunc(entity)
		if err != nil {
			return err
		}
		putreq := &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: dav,
			},
		}
		writeReqs = append(writeReqs, putreq)
		if len(writeReqs) >= 25 {
			writeInputs = append(writeInputs, &dynamodb.BatchWriteItemInput{
				RequestItems: map[string][]*dynamodb.WriteRequest{
					tableName: writeReqs,
				},
			})
			writeReqs = []*dynamodb.WriteRequest{}
		}
	}
	if len(writeReqs) > 0 {
		writeInputs = append(writeInputs, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				tableName: writeReqs,
			},
		})
		writeReqs = []*dynamodb.WriteRequest{}
	}
	for _, writeInput := range writeInputs {
		resp, err := ddb.client.BatchWriteItemWithContext(context.TODO(), writeInput)
		if err != nil {
			return err
		}
		if resp == nil || len(resp.UnprocessedItems) > 0 {
			return errors.New("no response got or have unprocessed items")
		}
	}
	return nil
}

func loadDdbItems[T domain.Dummy](
	tableName string,
	keys []*T,
	toDBKeyFunc func(item *T) map[string]*dynamodb.AttributeValue,
	toEntityFunc func(dav map[string]*dynamodb.AttributeValue) (*T, error),
) ([]*T, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	batchGetInputs := []*dynamodb.BatchGetItemInput{}
	ddbKeys := []map[string]*dynamodb.AttributeValue{}
	for _, key := range keys {
		ddbKeys = append(ddbKeys, toDBKeyFunc(key))
		if len(ddbKeys) >= 100 {
			batchGetInputs = append(batchGetInputs, &dynamodb.BatchGetItemInput{
				RequestItems: map[string]*dynamodb.KeysAndAttributes{
					tableName: {Keys: ddbKeys},
				},
			})
			ddbKeys = []map[string]*dynamodb.AttributeValue{}
		}
	}
	if len(ddbKeys) > 0 {
		batchGetInputs = append(batchGetInputs, &dynamodb.BatchGetItemInput{
			RequestItems: map[string]*dynamodb.KeysAndAttributes{
				tableName: {Keys: ddbKeys},
			},
		})
	}
	batchGetItems := []map[string]*dynamodb.AttributeValue{}
	for _, batchGetInput := range batchGetInputs {
		batchGetOutput, err := ddb.client.BatchGetItemWithContext(
			context.TODO(),
			batchGetInput,
		)
		if err != nil {
			return nil, err
		}
		items := batchGetOutput.Responses[tableName]
		batchGetItems = append(batchGetItems, items...)
	}
	results := []*T{}
	for _, item := range batchGetItems {
		result, err := toEntityFunc(item)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
