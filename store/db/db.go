package db

import (
	"context"
	"errors"
	"time"

	"github.com/urfave/cli/v2"
)

// db find no result
var ErrNoResult = errors.New("db return no result")

var (
	DatabaseUpdateTimeout = time.Second * 5
	DatabaseLoadTimeout   = time.Second * 5
)

type DB interface {
	MigrateTable(colName string, indexNames ...string) error
	FindOne(ctx context.Context, colName string, filter interface{}, result interface{}) error
	Find(ctx context.Context, colName string, filter interface{}) (interface{}, error)
	UpdateOne(ctx context.Context, colName string, filter interface{}, update interface{}) error
	DeleteOne(ctx context.Context, colName string, filter interface{}) error
	Exit()
}

func NewDB(ctx *cli.Context) DB {
	return NewMongoDB(ctx)
}
