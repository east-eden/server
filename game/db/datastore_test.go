package db

import (
	"context"
	"testing"
)

func TestDatastore(t *testing.T) {
	dsn := "mongodb://localhost:27017"
	database := "db_game"

	ds := newDatastore(context.Background(), dsn, database, 9999)
	if ds == nil {
		t.Errorf("new datastore failed")
	}

	go ds.Run()
}
