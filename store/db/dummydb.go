package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type DummyDB struct {
}

func NewDummyDB() DB {
	m := &DummyDB{}

	return m
}

// migrate collection
func (m *DummyDB) MigrateTable(name string, indexNames ...string) error {
	return nil
}

func (m *DummyDB) FindOne(ctx context.Context, colName string, filter interface{}, result interface{}) error {
	return nil
}

func (m *DummyDB) Find(ctx context.Context, colName string, filter interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func (m *DummyDB) InsertOne(ctx context.Context, colName string, insert interface{}) error {
	return nil
}

func (m *DummyDB) InsertMany(ctx context.Context, colName string, inserts []interface{}) error {
	return nil
}

func (m *DummyDB) UpdateOne(ctx context.Context, colName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) error {
	return nil
}

func (m *DummyDB) DeleteOne(ctx context.Context, colName string, filter interface{}) error {
	return nil
}

func (m *DummyDB) BulkWrite(ctx context.Context, colName string, model interface{}) error {
	return nil
}

func (m *DummyDB) Exit() {
}
