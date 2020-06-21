package db

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
)

// db find no result
var ErrNoResult = errors.New("db return no result")

// DBObjector save and load with all structure
type DBObjector interface {
	GetObjID() int64
	AfterLoad() error
}

var (
	DatabaseUpdateTimeout = time.Minute * 5
	DatabaseLoadTimeout   = time.Minute * 5
)

type DB interface {
	MigrateTable(tblName string, indexNames ...string) error
	SaveObject(tblName string, x DBObjector) error
	SaveFields(tblName string, x DBObjector, fields map[string]interface{}) error
	LoadObject(tblName, key string, value interface{}, x DBObjector) error
	LoadArray(tblName, key string, value interface{}, pool *sync.Pool) ([]interface{}, error)
	DeleteObject(tblName string, x DBObjector) error
	Exit(context.Context)
}

func NewDB(ctx *cli.Context) DB {
	return NewMongoDB(ctx)
}
