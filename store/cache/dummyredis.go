package cache

import (
	"os"

	"github.com/urfave/cli/v2"
)

type DummyRedis struct {
}

func NewDummyRedis(ctx *cli.Context) *DummyRedis {
	_, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		_ = ctx.String("redis_addr")
	}

	r := &DummyRedis{}

	return r
}

func (r *DummyRedis) SaveObject(prefix string, k any, x any) error {
	return nil
}

func (r *DummyRedis) SaveHashObject(prefix string, k any, field any, x any) error {
	return nil
}

func (r *DummyRedis) SaveHashAll(prefix string, k any, fields map[string]any) error {
	return nil
}

func (r *DummyRedis) SaveFields(prefix string, k any, fields map[string]any) error {
	return nil
}

func (r *DummyRedis) LoadObject(prefix string, k any, x any) error {
	return ErrNoResult
}

func (r *DummyRedis) LoadHashAll(prefix string, keyValue any) (any, error) {
	return nil, ErrNoResult
}

func (r *DummyRedis) DeleteObject(prefix string, k any) error {
	return nil
}

func (r *DummyRedis) DeleteHashObject(prefix string, k any, fields any) error {
	return nil
}

func (r *DummyRedis) Exit() error {
	return nil
}
