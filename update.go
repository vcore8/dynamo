package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Update struct {
	table      table
	hashKey    string
	hashValue  interface{}
	sortKey    string
	sortValue  interface{}
	removeExpr []string
	setExpr    map[string]interface{}
}

func (update *Update) SetRemoveExpr(attr string) *Update {
	update.removeExpr = append(update.removeExpr, attr)
	return update
}

func (update *Update) UpdateColumn(column string, value interface{}) *Update {
	if update.setExpr == nil {
		update.setExpr = map[string]interface{}{}
	}

	update.setExpr[column] = value
	return update
}

func (update *Update) SetSortKey(sort string, sortValue interface{}) *Update {
	update.sortKey = sort
	update.sortValue = sortValue
	return update
}

func (update *Update) Run(ctx context.Context) (err error) {
	primaryKey := map[string]interface{}{
		update.hashKey: update.hashValue,
	}

	if update.sortKey != "" {
		primaryKey[update.sortKey] = update.sortValue
	}

	pk, err := attributevalue.MarshalMap(primaryKey)
	if err != nil {
		return
	}

	builder := expression.NewBuilder()

	if len(update.removeExpr) > 0 || len(update.setExpr) > 0 {
		upd := expression.UpdateBuilder{}
		for _, rm := range update.removeExpr {
			upd = upd.Remove(expression.Name(rm))
		}
		for column, value := range update.setExpr {
			upd = upd.Set(expression.Name(column), expression.Value(value))
		}
		builder = builder.WithUpdate(upd)
	}

	expr, err := builder.Build()
	if err != nil {
		return
	}

	_, err = update.table.db.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 &update.table.name,
		Key:                       pk,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	})
	return
}
