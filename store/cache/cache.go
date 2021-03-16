package cache

import (
	"errors"
	"time"

	"github.com/urfave/cli/v2"
)

// cache find no result
var (
	ErrNoResult       = errors.New("cache return no result")
	ErrObjectNotFound = errors.New("cache object not found")

	ExpireTime = 24 * time.Hour // 过期时间24小时
)

type Cache interface {
	SaveObject(prefix string, k interface{}, x interface{}) error
	SaveFields(prefix string, k interface{}, fields map[string]interface{}) error
	LoadObject(prefix string, k interface{}, x interface{}) error
	DeleteObject(prefix string, k interface{}) error
	DeleteFields(prefix string, k interface{}, fieldsName []string) error
	Exit() error
}

func NewCache(ctx *cli.Context) Cache {
	return NewDummyRedis(ctx)
}
