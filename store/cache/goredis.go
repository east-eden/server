package cache

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"bitbucket.org/funplus/server/utils"
	"github.com/go-redis/redis"
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

func (r *GoRedis) SaveObject(prefix string, k interface{}, x interface{}) error {
	key := fmt.Sprintf("%s:%v", prefix, k)
	data, err := json.Marshal(x)
	if pass := utils.ErrCheck(err, "json marshal failed when goredis SaveObject", key, x); !pass {
		return err
	}

	_, err = r.redisCli.Set(key, data, ExpireTime).Result()
	utils.ErrPrint(err, "goredis set failed", key)

	return err
}

func (r *GoRedis) SaveHashObject(prefix string, k interface{}, field interface{}, x interface{}) error {
	key := fmt.Sprintf("%s:%v", prefix, k)
	data, err := json.Marshal(x)
	if pass := utils.ErrCheck(err, "json marshal failed when goredis SaveObject", key, x); !pass {
		return err
	}

	f := fmt.Sprintf("%v", field)
	_, err = r.redisCli.HSet(key, f, data).Result()
	utils.ErrPrint(err, "goredis hset failed", key)

	return err
}

func (r *GoRedis) SaveHashAll(prefix string, k interface{}, fields map[string]interface{}) error {
	key := fmt.Sprintf("%s:%v", prefix, k)

	_, err := r.redisCli.HMSet(key, fields).Result()
	utils.ErrPrint(err, "goredis hmset failed", key, fields)

	// update expire
	r.redisCli.Expire(key, ExpireTime)

	return err
}

func (r *GoRedis) LoadObject(prefix string, k interface{}, x interface{}) error {
	key := fmt.Sprintf("%s:%v", prefix, k)

	data, err := r.redisCli.Get(key).Bytes()
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
	_, err = r.redisCli.Expire(key, ExpireTime).Result()
	return err
}

func (r *GoRedis) LoadHashAll(prefix string, keyValue interface{}) (interface{}, error) {
	key := fmt.Sprintf("%s:%v", prefix, keyValue)

	m, err := r.redisCli.HGetAll(key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, ErrNoResult
	}

	if len(m) == 0 {
		return nil, ErrNoResult
	}

	result := make(map[string]interface{}, len(m))
	for k, v := range m {
		result[k] = []byte(v)
	}

	// update expire
	_, err = r.redisCli.Expire(key, ExpireTime).Result()
	return result, err
}

func (r *GoRedis) DeleteObject(prefix string, k interface{}) error {
	key := fmt.Sprintf("%s:%v", prefix, k)
	_, err := r.redisCli.Del(key).Result()
	utils.ErrPrint(err, "redis delete object failed", k)

	return err
}

func (r *GoRedis) DeleteHashObject(prefix string, k interface{}, field interface{}) error {
	key := fmt.Sprintf("%s:%v", prefix, k)
	f := fmt.Sprintf("%v", field)
	_, err := r.redisCli.HDel(key, f).Result()
	utils.ErrPrint(err, "redis delete object failed", k)

	return err
}

func (r *GoRedis) Exit() error {
	return r.redisCli.Close()
}
