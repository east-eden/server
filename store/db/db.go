package db

import (
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
	FindOne(colName string, filter interface{}, result interface{}) error
	Find(colName string, filter interface{}) (interface{}, error)
	UpdateOne(colName string, filter interface{}, update interface{}) error
	DeleteOne(colName string, filter interface{}) error
	Exit()
}

func NewDB(ctx *cli.Context) DB {
	return NewMongoDB(ctx)
}
