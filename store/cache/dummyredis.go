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

func (r *DummyRedis) SaveObject(prefix string, k interface{}, x interface{}) error {
	return nil
}

func (r *DummyRedis) SaveFields(prefix string, k interface{}, fields map[string]interface{}) error {
	return nil
}

func (r *DummyRedis) SaveMarshaledObject(prefix string, k interface{}, x interface{}) error {
	return nil
}

func (r *DummyRedis) LoadObject(prefix string, k interface{}, x interface{}) error {
	return ErrNoResult
}

func (r *DummyRedis) DeleteObject(prefix string, k interface{}) error {
	return nil
}

func (r *DummyRedis) DeleteFields(prefix string, k interface{}, fieldsName []string) error {
	return nil
}

func (r *DummyRedis) Exit() error {
	return nil
}
