package cache

import (
	"context"
	"errors"
	"sync"

	"github.com/urfave/cli/v2"
)

// cache find no result
var ErrNoResult = errors.New("cache return no result")

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
	Exit(context.Context) error
}

func NewCache(ctx *cli.Context) Cache {
	return NewRedis(ctx)
}
