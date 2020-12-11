package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/go-redis/redis"
	"github.com/nitishm/go-rejson"
	"github.com/nitishm/go-rejson/rjs"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/east-eden/server/utils"
)

type GoRedis struct {
	redisCli *redis.Client
	handler  *rejson.Handler
	utils.WaitGroupWrapper
}

func NewGoRedis(ctx *cli.Context) *GoRedis {
	redisAddr, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		redisAddr = ctx.String("redis_addr")
	}

	r := &GoRedis{
		redisCli: redis.NewClient(&redis.Options{Addr: redisAddr}),
		handler:  rejson.NewReJSONHandler(),
	}

	r.handler.SetGoRedisClient(r.redisCli)

	return r
}

func (r *GoRedis) SaveObject(prefix string, x CacheObjector) error {
	key := fmt.Sprintf("%s:%v", prefix, x.GetObjID())
	if _, err := r.handler.JSONSet(key, ".", x); err != nil {
		return fmt.Errorf("Redis.SaveObject failed: %w", err)
	}

	// save object index
	if x.GetStoreIndex() != -1 {
		zaddKey := fmt.Sprintf("%s_index:%v", prefix, x.GetStoreIndex())

		err := r.redisCli.ZAdd(zaddKey, redis.Z{Score: 0, Member: key}).Err()
		if err != nil {
			return fmt.Errorf("Redis.SaveObject Index failed: %w", err)
		}
	}

	return nil
}

func (r *GoRedis) SaveFields(prefix string, x CacheObjector, fields map[string]interface{}) error {
	key := fmt.Sprintf("%s:%v", prefix, x.GetObjID())
	for path, val := range fields {
		if _, err := r.handler.JSONSet(key, "."+path, val); err != nil {
			return fmt.Errorf("Redis.SaveFields path<%s> failed: %w", path, err)
		}
	}

	return nil
}

func (r *GoRedis) LoadObject(prefix string, value interface{}, x CacheObjector) error {
	key := fmt.Sprintf("%s:%v", prefix, value)

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

	err = json.Unmarshal(res.([]byte), x)
	if err != nil {
		return err
	}

	return nil
}

func (r *GoRedis) LoadArray(prefix string, ownerId int64, pool *sync.Pool) ([]interface{}, error) {
	zKey := fmt.Sprintf("%s_index:%d", prefix, ownerId)
	keys, err := r.redisCli.ZRange(zKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, ErrNoResult
	}

	reply := make([]interface{}, 0)
	for _, key := range keys {
		res, err := r.handler.JSONGet(key, ".", rjs.GETOptionNOESCAPE)
		if err != nil && !errors.Is(err, redis.Nil) {
			return reply, err
		}

		// empty result
		if res == nil {
			continue
		}

		x := pool.Get()
		if err := json.Unmarshal(res.([]byte), x); err != nil {
			return reply, err
		}

		reply = append(reply, x)
	}

	return reply, nil
}

func (r *GoRedis) DeleteObject(prefix string, x CacheObjector) error {
	key := fmt.Sprintf("%s:%v", prefix, x.GetObjID())
	if _, err := r.handler.JSONDel(key, "."); err != nil {
		log.Error().
			Int64("obj_id", x.GetObjID()).
			Int64("store_idx", x.GetStoreIndex()).
			Err(err).
			Msg("redis delete object failed")
	}

	// delete object index
	if x.GetStoreIndex() == -1 {
		return nil
	}

	zremKey := fmt.Sprintf("%s_index:%v", prefix, x.GetStoreIndex())
	err := r.redisCli.ZRem(zremKey, key).Err()
	if err != nil {
		return fmt.Errorf("GoRedis.DeleteObject index failed: %w", err)
	}

	return nil
}

func (r *GoRedis) Exit() error {
	r.Wait()
	return nil
}
