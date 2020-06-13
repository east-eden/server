package cache

import (
	"context"
	"sync"

	"github.com/urfave/cli/v2"
)

// CacheObjector save and load with all structure
type CacheObjector interface {
	GetObjID() interface{}
}

type Cache interface {
	SaveObject(prefix string, x CacheObjector) error
	SaveFields(prefix string, x CacheObjector, fields map[string]interface{}) error
	LoadObject(prefix string, value interface{}, x CacheObjector) error
	LoadArray(prefix string, pool *sync.Pool) ([]interface{}, error)
	DeleteObject(prefix string, x CacheObjector) error
	Exit(context.Context) error
}

func NewCache(ctx *cli.Context) Cache {
	return NewRedis(ctx)
}
