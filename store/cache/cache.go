package cache

import (
	"errors"
	"sync"

	"github.com/urfave/cli/v2"
)

// cache find no result
var (
	ErrNoResult       = errors.New("cache return no result")
	ErrObjectNotFound = errors.New("cache object not found")
)

// CacheObjector save and load with all structure
type CacheObjector interface {
	GetObjID() int64
	GetStoreIndex() int64
}

type Cache interface {
	SaveObject(prefix string, x CacheObjector) error
	SaveFields(prefix string, x CacheObjector, fields map[string]interface{}) error
	LoadObject(prefix string, value interface{}, x CacheObjector) error
	LoadArray(prefix string, ownerId int64, pool *sync.Pool) ([]interface{}, error)
	DeleteObject(prefix string, x CacheObjector) error
	Exit() error
}

func NewCache(ctx *cli.Context) Cache {
	return NewGoRedis(ctx)
}
