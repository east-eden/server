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
	DatabaseWriteTimeout     = time.Second * 5
	DatabaseLoadTimeout      = time.Second * 5
	DatabaseBulkWriteTimeout = time.Second * 10
)

type DB interface {
	MigrateTable(colName string, indexNames ...string) error
	FindOne(ctx context.Context, colName string, filter any, result any) error
	Find(ctx context.Context, colName string, filter any) (map[string]any, error)
	InsertOne(ctx context.Context, colName string, insert any) error
	InsertMany(ctx context.Context, colName string, inserts []any) error
	UpdateOne(ctx context.Context, colName string, filter any, update any, opts ...*options.UpdateOptions) error
	DeleteOne(ctx context.Context, colName string, filter any) error
	BulkWrite(ctx context.Context, colName string, model any) error
	Flush()
	Exit()
}

func NewDB(ctx *cli.Context) DB {
	return NewMongoDB(ctx)
}
