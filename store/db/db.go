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
	MigrateTable(tblName string, indexNames ...string) error
	SaveObject(tblName string, k interface{}, x interface{}) error
	SaveFields(tblName string, k interface{}, fields map[string]interface{}) error
	LoadObject(tblName, key string, value interface{}, x interface{}) error
	// LoadArray(tblName, key string, storeIndex int64, pool *sync.Pool) ([]interface{}, error)
	DeleteObject(tblName string, k interface{}) error
	DeleteFields(tblName string, k interface{}, fieldsName []string) error
	Exit()
}

func NewDB(ctx *cli.Context) DB {
	return NewMongoDB(ctx)
}
