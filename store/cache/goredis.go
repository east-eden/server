package cache

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/east-eden/server/utils"
	"github.com/go-redis/redis/v8"
	json "github.com/json-iterator/go"
	"github.com/urfave/cli/v2"
)

type GoRedis struct {
	redisCli *redis.Client
}

func NewGoRedis(ctx *cli.Context) *GoRedis {
	redisAddr, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		redisAddr = ctx.String("redis_addr")
	}

	r := &GoRedis{
		redisCli: redis.NewClient(&redis.Options{Addr: redisAddr}),
	}

	return r
}

func (r *GoRedis) SaveObject(prefix string, k any, x any) error {
	key := fmt.Sprintf("%s:%v", prefix, k)
	data, err := json.Marshal(x)
	if !utils.ErrCheck(err, "json marshal failed when goredis SaveObject", key, x) {
		return err
	}

	_, err = r.redisCli.Set(context.Background(), key, data, ExpireTime).Result()
	utils.ErrPrint(err, "goredis set failed", key)

	return err
}

func (r *GoRedis) SaveHashObject(prefix string, k any, field any, x any) error {
	key := fmt.Sprintf("%s:%v", prefix, k)
	data, err := json.Marshal(x)
	if !utils.ErrCheck(err, "json marshal failed when goredis SaveObject", key, x) {
		return err
	}

	f := fmt.Sprintf("%v", field)
	_, err = r.redisCli.HSet(context.Background(), key, f, data).Result()
	utils.ErrPrint(err, "goredis hset failed", key)

	return err
}

func (r *GoRedis) SaveHashAll(prefix string, k any, fields map[string]any) error {
	key := fmt.Sprintf("%s:%v", prefix, k)

	_, err := r.redisCli.HMSet(context.Background(), key, fields).Result()
	utils.ErrPrint(err, "goredis hmset failed", key, fields)

	// update expire
	r.redisCli.Expire(context.Background(), key, ExpireTime)

	return err
}

func (r *GoRedis) LoadObject(prefix string, k any, x any) error {
	key := fmt.Sprintf("%s:%v", prefix, k)

	data, err := r.redisCli.Get(context.Background(), key).Bytes()
	if errors.Is(err, redis.Nil) {
		return ErrObjectNotFound
	}

	if err != nil {
		return err
	}

	decoder := json.NewDecoder(bytes.NewBuffer(data))
	decoder.UseNumber()
	err = decoder.Decode(x)
	if err != nil {
		return err
	}

	// update expire
	_, err = r.redisCli.Expire(context.Background(), key, ExpireTime).Result()
	return err
}

func (r *GoRedis) LoadHashAll(prefix string, keyValue any) (any, error) {
	key := fmt.Sprintf("%s:%v", prefix, keyValue)

	m, err := r.redisCli.HGetAll(context.Background(), key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, ErrNoResult
	}

	if len(m) == 0 {
		return nil, ErrNoResult
	}

	result := make(map[string]any, len(m))
	for k, v := range m {
		result[k] = []byte(v)
	}

	// update expire
	_, err = r.redisCli.Expire(context.Background(), key, ExpireTime).Result()
	return result, err
}

func (r *GoRedis) DeleteObject(prefix string, k any) error {
	key := fmt.Sprintf("%s:%v", prefix, k)
	_, err := r.redisCli.Del(context.Background(), key).Result()
	utils.ErrPrint(err, "redis delete object failed", k)

	return err
}

func (r *GoRedis) DeleteHashObject(prefix string, k any, field any) error {
	key := fmt.Sprintf("%s:%v", prefix, k)
	f := fmt.Sprintf("%v", field)
	_, err := r.redisCli.HDel(context.Background(), key, f).Result()
	utils.ErrPrint(err, "redis delete object failed", k)

	return err
}

func (r *GoRedis) Exit() error {
	return r.redisCli.Close()
}
