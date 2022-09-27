package dynamo

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DB struct {
	client *dynamodb.Client
}

func New(cfg aws.Config) *DB {
	client := dynamodb.NewFromConfig(cfg)

	return &DB{client}
}
