package cache

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"

	"bitbucket.org/funplus/server/utils"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	json "github.com/json-iterator/go"
	"github.com/nitishm/go-rejson"
	"github.com/nitishm/go-rejson/rjs"
	"github.com/urfave/cli/v2"
)

type MiniRedis struct {
	redisServer *miniredis.Miniredis
	redisCli    *redis.Client
	handler     *rejson.Handler
	utils.WaitGroupWrapper
}

func NewMiniRedis(ctx *cli.Context) *MiniRedis {
	redisAddr, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		redisAddr = ctx.String("redis_addr")
	}

	r := &MiniRedis{
		redisCli: redis.NewClient(&redis.Options{Addr: redisAddr}),
		handler:  rejson.NewReJSONHandler(),
	}

	r.handler.SetGoRedisClient(r.redisCli)

	r.redisServer = miniredis.NewMiniRedis()
	if err := r.redisServer.StartAddr(redisAddr); err != nil {
		log.Fatal(err)
	}

	return r
}

func (r *MiniRedis) SaveObject(prefix string, k interface{}, x interface{}) error {
	key := fmt.Sprintf("%s:%v", prefix, k)
	if _, err := r.handler.JSONSet(key, ".", x); err != nil {
		return fmt.Errorf("Redis.SaveObject failed: %w", err)
	}

	// update expire
	r.redisCli.Expire(key, ExpireTime)

	return nil
}

func (r *MiniRedis) SaveFields(prefix string, k interface{}, fields map[string]interface{}) error {
	key := fmt.Sprintf("%s:%v", prefix, k)
	for path, val := range fields {
		if _, err := r.handler.JSONSet(key, "."+path, val); err != nil {
			return fmt.Errorf("Redis.SaveFields path<%s> failed: %w", path, err)
		}
	}

	// update expire
	r.redisCli.Expire(key, ExpireTime)

	return nil
}

func (r *MiniRedis) LoadObject(prefix string, k interface{}, x interface{}) error {
	key := fmt.Sprintf("%s:%v", prefix, k)

	res, err := r.handler.JSONGet(key, ".", rjs.GETOptionNOESCAPE)
	if errors.Is(err, redis.Nil) {
		return ErrObjectNotFound
	}

	if err != nil {
		return err
	}

	// empty result
	if res == nil {
		return ErrObjectNotFound
	}

	decoder := json.NewDecoder(bytes.NewBuffer(res.([]byte)))
	decoder.UseNumber()
	err = decoder.Decode(x)
	// err = json.Unmarshal(res.([]byte), x)
	if err != nil {
		return err
	}

	// update expire
	r.redisCli.Expire(key, ExpireTime)

	return nil
}

func (r *MiniRedis) DeleteObject(prefix string, k interface{}) error {
	key := fmt.Sprintf("%s:%v", prefix, k)
	_, err := r.handler.JSONDel(key, ".")
	utils.ErrPrint(err, "redis delete object failed", k)

	// delete object index
	// if x.GetStoreIndex() == -1 {
	// 	return nil
	// }

	// zremKey := fmt.Sprintf("%s_index:%v", prefix, x.GetStoreIndex())
	// err := r.redisCli.ZRem(zremKey, key).Err()
	// if err != nil {
	// 	return fmt.Errorf("GoRedis.DeleteObject index failed: %w", err)
	// }

	return err
}

func (r *MiniRedis) DeleteFields(prefix string, k interface{}, fieldsName []string) error {
	key := fmt.Sprintf("%s:%v", prefix, k)
	for _, path := range fieldsName {
		if _, err := r.handler.JSONDel(key, "."+path); err != nil {
			return fmt.Errorf("Redis.SaveFields path<%s> failed: %w", path, err)
		}
	}

	return nil
}

func (r *MiniRedis) Exit() error {
	r.Wait()
	err := r.redisCli.Close()
	r.redisServer.Close()
	return err
}
