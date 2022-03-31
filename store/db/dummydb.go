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

func (m *DummyDB) FindOne(ctx context.Context, colName string, filter any, result any) error {
	return nil
}

func (m *DummyDB) Find(ctx context.Context, colName string, filter any) (map[string]any, error) {
	return nil, nil
}

func (m *DummyDB) InsertOne(ctx context.Context, colName string, insert any) error {
	return nil
}

func (m *DummyDB) InsertMany(ctx context.Context, colName string, inserts []any) error {
	return nil
}

func (m *DummyDB) UpdateOne(ctx context.Context, colName string, filter any, update any, opts ...*options.UpdateOptions) error {
	return nil
}

func (m *DummyDB) DeleteOne(ctx context.Context, colName string, filter any) error {
	return nil
}

func (m *DummyDB) BulkWrite(ctx context.Context, colName string, model any) error {
	return nil
}

func (m *DummyDB) Flush() {

}

func (m *DummyDB) Exit() {
}
