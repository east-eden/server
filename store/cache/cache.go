package cache

import "github.com/urfave/cli/v2"

// CacheObjector save and load with all structure
type CacheObjector interface {
	GetObjID() interface{}
}

type Cache interface {
	SaveObject(prefix string, x CacheObjector) error
	LoadObject(prefix string, x CacheObjector) error
	DeleteObject(prefix string, x CacheObjector) error
}

func NewCache(ctx *cli.Context) Cache {
	return NewRedis(ctx)
}
