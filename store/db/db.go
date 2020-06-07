package db

import (
	"context"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
)

// DBObjector save and load with all structure
type DBObjector interface {
	GetObjID() interface{}
	AfterLoad()
	TableName() string
}

var (
	DatabaseUpdateTimeout = time.Second * 5
	DatabaseLoadTimeout   = time.Second * 5
)

type DB interface {
	MigrateTable(tblName string, indexNames ...string) error
	SaveObject(x DBObjector) error
	SaveFields(x DBObjector, fields map[string]interface{}) error
	LoadObject(key string, value interface{}, x DBObjector) error
	LoadArray(tblName, key string, value interface{}, pool *sync.Pool) ([]DBObjector, error)
	DeleteObject(x DBObjector) error
	Exit(context.Context)
}

func NewDB(ctx *cli.Context) DB {
	return NewMongoDB(ctx)
}
