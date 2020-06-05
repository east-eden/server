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
	TableName() string
}

var (
	DatabaseUpdateTimeout = time.Second * 5
	DatabaseLoadTimeout   = time.Second * 5
)

type DB interface {
	MigrateTable(tblName string, indexNames ...string) error
	SaveObject(x DBObjector) error
	LoadObject(filter string, key interface{}, x DBObjector) error
	LoadObjectArray(tblName, filter string, key interface{}, pool *sync.Pool) ([]DBObjector, error)
	Exit(context.Context)
}

func NewDB(ctx *cli.Context) DB {
	return NewMongoDB(ctx)
}
