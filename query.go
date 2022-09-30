package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Query struct {
	table     table
	hashKey   string
	hashValue interface{}
	sortKey   string
	sortValue interface{}
	index     string
}

func (query *Query) Index(name string) *Query {
	query.index = name
	return query
}

func (query *Query) SortBy(name string, value interface{}) *Query {
	query.sortKey = name
	query.sortValue = value
	return query
}

func (query *Query) All(ctx context.Context, out interface{}) (err error) {
	input := &dynamodb.QueryInput{
		TableName: &query.table.name,
	}

	if query.index != "" {
		input.IndexName = &query.index
	}

	if query.hashKey != "" {
		filter := expression.Key(query.hashKey).Equal(expression.Value(query.hashValue))

		if query.sortKey != "" {
			filter = expression.KeyAnd(filter,
				expression.Key(query.sortKey).Equal(expression.Value(query.sortValue)),
			)
		}

		expr, err := expression.NewBuilder().WithKeyCondition(filter).Build()
		if err != nil {
			return err
		}

		input.ExpressionAttributeNames = expr.Names()
		input.ExpressionAttributeValues = expr.Values()
		input.KeyConditionExpression = expr.KeyCondition()
	}

	response, err := query.table.db.client.Query(ctx, input)
	if err != nil {
		return
	}

	return attributevalue.UnmarshalListOfMaps(response.Items, &out)
}

func (query *Query) Scan(ctx context.Context, out interface{}) (err error) {
	input := &dynamodb.ScanInput{
		TableName: &query.table.name,
	}

	if query.index != "" {
		input.IndexName = &query.index
	}

	if query.hashKey != "" {
		filter := expression.Name(query.hashKey).Equal(expression.Value(query.hashValue))

		expr, err := expression.NewBuilder().WithFilter(filter).Build()
		if err != nil {
			return err
		}

		input.ExpressionAttributeNames = expr.Names()
		input.ExpressionAttributeValues = expr.Values()
		input.FilterExpression = expr.Filter()
	}

	response, err := query.table.db.client.Scan(ctx, input)
	if err != nil {
		return
	}

	return attributevalue.UnmarshalListOfMaps(response.Items, &out)
}

func (query *Query) One(ctx context.Context, out interface{}) (err error) {
	selectedKeys := map[string]interface{}{
		query.hashKey: query.hashValue,
		query.sortKey: query.sortValue,
	}

	key, err := attributevalue.MarshalMap(selectedKeys)
	if err != nil {
		return
	}

	response, err := query.table.db.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &query.table.name,
		Key:       key,
	})
	if err != nil {
		return
	}

	if len(response.Item) == 0 {
		return NotFound
	}

	return attributevalue.UnmarshalMap(response.Item, &out)
}
