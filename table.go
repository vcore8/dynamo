package dynamo

import "context"

type table struct {
	name string
	db   *DB
}

func (db *DB) Table(name string) table {
	return table{
		name: name,
		db:   db,
	}
}

func (table table) Get(name string, value interface{}) *Query {
	q := &Query{
		table:     table,
		hashKey:   name,
		hashValue: value,
	}

	return q
}

func (table table) Update(hash string, hashValue interface{}) *Update {
	q := &Update{
		table:      table,
		hashKey:    hash,
		hashValue:  hashValue,
		removeExpr: []string{},
	}

	return q
}

func (table table) Create(ctx context.Context, item interface{}) error {
	tb := &Create{table: table}
	return tb.Run(ctx, item)
}
