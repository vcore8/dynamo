package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Create struct {
	table table
}

func (create *Create) Run(ctx context.Context, item interface{}) (err error) {
	data, err := attributevalue.MarshalMap(item)
	if err != nil {
		return
	}

	_, err = create.table.db.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &create.table.name,
		Item:      data,
	})

	return
}
