package db

import (
	"context"
	"errors"
	"time"

	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// db find no result
var ErrNoResult = errors.New("db return no result")

var (
	DatabaseWriteTimeout = time.Second * 5
	DatabaseLoadTimeout  = time.Second * 5
)

type DB interface {
	MigrateTable(colName string, indexNames ...string) error
	FindOne(ctx context.Context, colName string, filter interface{}, result interface{}) error
	Find(ctx context.Context, colName string, filter interface{}) (map[string]interface{}, error)
	InsertOne(ctx context.Context, colName string, insert interface{}) error
	InsertMany(ctx context.Context, colName string, inserts []interface{}) error
	UpdateOne(ctx context.Context, colName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) error
	DeleteOne(ctx context.Context, colName string, filter interface{}) error
	BulkWrite(ctx context.Context, colName string, model interface{}) error
	Exit()
}

func NewDB(ctx *cli.Context) DB {
	return NewMongoDB(ctx)
}
