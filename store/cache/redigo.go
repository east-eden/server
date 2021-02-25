package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"bitbucket.org/east-eden/server/utils"
	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson"
	"github.com/nitishm/go-rejson/rjs"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	RedisConnectTimeout = time.Second * 5
	RedisReadTimeout    = time.Second * 5
	RedisWriteTimeout   = time.Second * 5
)

type Redigo struct {
	addr string
	pool *redis.Pool
	utils.WaitGroupWrapper
	mapRejsonHandler map[redis.Conn]*rejson.Handler
	sync.RWMutex
}

func NewRedigo(ctx *cli.Context) *Redigo {
	redisAddr, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		redisAddr = ctx.String("redis_addr")
	}

	r := &Redigo{
		addr: redisAddr,
		pool: &redis.Pool{
			Wait:        true,
			MaxIdle:     500,
			MaxActive:   5000,
			IdleTimeout: time.Second * 300,
		},
		mapRejsonHandler: make(map[redis.Conn]*rejson.Handler),
	}

	r.pool.Dial = func() (redis.Conn, error) {
		c, err := redis.DialTimeout("tcp", r.addr, RedisConnectTimeout, RedisReadTimeout, RedisWriteTimeout)
		if err != nil {
			panic(err.Error())
		}
		return c, err
	}

	r.pool.TestOnBorrow = func(c redis.Conn, t time.Time) error {
		if time.Since(t) < time.Minute {
			return nil
		}

		_, err := c.Do("PING")
		return err
	}

	return r
}

// get rejson's handler by redigo.Conn, if not existing, create one.
func (r *Redigo) getRejsonHandler() (redis.Conn, *rejson.Handler) {
	r.Lock()
	defer r.Unlock()

	c := r.pool.Get()
	if c.Err() != nil {
		return c, nil
	}

	h, ok := r.mapRejsonHandler[c]
	if !ok {
		rh := rejson.NewReJSONHandler()
		rh.SetRedigoClient(c)
		r.mapRejsonHandler[c] = rh
		return c, rh
	}

	return c, h
}

// get rejson's handler by redigo.Conn, if not existing, create one.
func (r *Redigo) returnRejsonHandler(con redis.Conn) {
	if con == nil {
		return
	}

	r.Lock()
	defer r.Unlock()

	delete(r.mapRejsonHandler, con)
	con.Close()
}

func (r *Redigo) SaveObject(prefix string, k interface{}, x interface{}) error {
	con, handler := r.getRejsonHandler()
	if handler == nil {
		return fmt.Errorf("redis.SaveObject failed: %w", con.Err())
	}
	defer r.returnRejsonHandler(con)

	key := fmt.Sprintf("%s:%v", prefix, k)
	if _, err := handler.JSONSet(key, ".", x); err != nil {
		return fmt.Errorf("Redis.SaveObject failed: %w", err)
	}

	// update expire
	_, _ = con.Do("EXPIRE", key, ExpireTime/time.Second)

	return nil
}

func (r *Redigo) SaveFields(prefix string, k interface{}, fields map[string]interface{}) error {
	con, handler := r.getRejsonHandler()
	if handler == nil {
		return fmt.Errorf("redis.SaveFields failed: %w", con.Err())
	}
	defer r.returnRejsonHandler(con)

	key := fmt.Sprintf("%s:%v", prefix, k)
	for path, val := range fields {
		if _, err := handler.JSONSet(key, "."+path, val); err != nil {
			return fmt.Errorf("Redis.SaveFields path<%s> failed: %w", path, err)
		}
	}

	// update expire
	_, _ = con.Do("EXPIRE", key, ExpireTime/time.Second)

	return nil
}

func (r *Redigo) LoadObject(prefix string, value interface{}, x interface{}) error {
	con, handler := r.getRejsonHandler()
	if handler == nil {
		return fmt.Errorf("redis.LoadObject failed: %w", con.Err())
	}

	defer r.returnRejsonHandler(con)

	key := fmt.Sprintf("%s:%v", prefix, value)

	res, err := handler.JSONGet(key, ".", rjs.GETOptionNOESCAPE)
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

	// update expire
	_, _ = con.Do("EXPIRE", key, ExpireTime/time.Second)

	return nil
}

// deprecated
func (r *Redigo) LoadArray(prefix string, ownerId int64, pool *sync.Pool) ([]interface{}, error) {
	con, handler := r.getRejsonHandler()
	if handler == nil {
		return nil, fmt.Errorf("redis.LoadArray failed: %w", con.Err())
	}

	defer r.returnRejsonHandler(con)

	// scan all keys
	//var (
	//cursor int64
	//items  []string
	//)
	//results := make([]string, 0)

	//for {
	//values, err := redis.Values(c.Do("SCAN", cursor, "MATCH", prefix+":*"))
	//if err != nil {
	//return nil, err
	//}

	//values, err = redis.Scan(values, &cursor, &items)
	//if err != nil {
	//return nil, err
	//}

	//results = append(results, items...)

	//if cursor == 0 {
	//break
	//}
	//}

	zKey := fmt.Sprintf("%s_index:%d", prefix, ownerId)
	keys, err := redis.Strings(con.Do("ZRANGE", zKey, 0, -1))
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, ErrNoResult
	}

	reply := make([]interface{}, 0)
	for _, key := range keys {
		res, err := handler.JSONGet(key, ".", rjs.GETOptionNOESCAPE)
		if err != nil {
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

func (r *Redigo) DeleteObject(prefix string, k interface{}) error {
	con, handler := r.getRejsonHandler()
	if handler == nil {
		return fmt.Errorf("redis.DeleteObject failed:%w", con.Err())
	}

	defer r.returnRejsonHandler(con)

	key := fmt.Sprintf("%s:%v", prefix, k)
	if _, err := handler.JSONDel(key, "."); err != nil {
		log.Error().
			Interface("obj_id", k).
			Err(err).
			Msg("redis delete object failed")
	}

	// delete object index
	// if x.GetStoreIndex() == -1 {
	// 	return nil
	// }

	// zremKey := fmt.Sprintf("%s_index:%v", prefix, x.GetStoreIndex())
	// if _, err := con.Do("ZREM", zremKey, key); err != nil {
	// 	return fmt.Errorf("Redigo.DeleteObject index failed: %w", err)
	// }

	return nil
}

func (r *Redigo) DeleteFields(prefix string, k interface{}, fieldsName []string) error {
	con, handler := r.getRejsonHandler()
	if handler == nil {
		return fmt.Errorf("redis.DeleteObject failed:%w", con.Err())
	}

	defer r.returnRejsonHandler(con)

	key := fmt.Sprintf("%s:%v", prefix, k)
	for _, path := range fieldsName {
		if _, err := handler.JSONDel(key, "."+path); err != nil {
			return fmt.Errorf("Redis.SaveFields path<%s> failed: %w", path, err)
		}
	}

	return nil
}

func (r *Redigo) Exit() error {
	r.Wait()
	return r.pool.Close()
}
