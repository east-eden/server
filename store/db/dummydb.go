package db

import "context"

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

func (m *DummyDB) Find(ctx context.Context, colName string, filter interface{}) (interface{}, error) {
	return nil, nil
}

func (m *DummyDB) UpdateOne(ctx context.Context, colName string, filter interface{}, update interface{}) error {
	return nil
}

func (m *DummyDB) DeleteOne(ctx context.Context, colName string, filter interface{}) error {
	return nil
}

func (m *DummyDB) Exit() {
}
