package dynamotest

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Music struct {
	Artist    string
	SongTitle string
	Attr1     string   `dynamo:"attr1"`
	Attr2     []string `dynamo:"attr2,set"`
}

func PrepareTestFixtures(fixturesPath string, config *aws.Config, deleteTables []*dynamodb.DeleteTableInput, createTables []*dynamodb.CreateTableInput) (*dynamodb.DynamoDB, error) {
	cDB, err := createrDB()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("get creater db err: %s", err))
	}

	for i, _ := range deleteTables {
		_, err := cDB.DeleteTable(deleteTables[i])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("delete table err: %s", err))
		}
	}

	for i, _ := range createTables {
		_, err := cDB.CreateTable(createTables[i])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("create table err: %s", err))
		}
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	db := dynamodb.New(sess)

	// db.PutItem(input)

	return db, nil
}

func createrDB() (*dynamodb.DynamoDB, error) {
	config := &aws.Config{
		Region:   aws.String(os.Getenv("DYNAMODB_REGION")),
		Endpoint: aws.String(os.Getenv("DYNAMODB_ENDPOINT")),
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return dynamodb.New(sess), nil
}
