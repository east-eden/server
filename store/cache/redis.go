package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson"
	"github.com/nitishm/go-rejson/rjs"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/utils"
)

var (
	RedisConnectTimeout = time.Second * 5
	RedisReadTimeout    = time.Second * 5
	RedisWriteTimeout   = time.Second * 5
)

type RedisDoCallback func(interface{}, error)

type Redis struct {
	addr string
	pool *redis.Pool
	utils.WaitGroupWrapper
	mapRejsonHandler map[redis.Conn]*rejson.Handler
	sync.RWMutex
}

func NewRedis(ctx *cli.Context) *Redis {
	redisAddr, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		redisAddr = ctx.String("redis_addr")
	}

	r := &Redis{
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
func (r *Redis) getRejsonHandler() (redis.Conn, *rejson.Handler) {
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
func (r *Redis) returnRejsonHandler(con redis.Conn) {
	if con == nil {
		return
	}

	r.Lock()
	defer r.Unlock()

	delete(r.mapRejsonHandler, con)
	con.Close()
}

func (r *Redis) SaveObject(prefix string, x CacheObjector) error {
	con, handler := r.getRejsonHandler()
	if handler == nil {
		return fmt.Errorf("redis.SaveObject failed: %w", con.Err())
	}

	key := fmt.Sprintf("%s:%v", prefix, x.GetObjID())

	r.Wrap(func() {
		defer r.returnRejsonHandler(con)

		if _, err := handler.JSONSet(key, ".", x); err != nil {
			logger.WithFields(logger.Fields{
				"object": x,
				"error":  err,
			}).Error("redis save object failed")
		}

		// save object index
		if x.GetStoreIndex() == -1 {
			return
		}

		zaddKey := fmt.Sprintf("%s_index:%v", prefix, x.GetStoreIndex())
		if _, err := con.Do("ZADD", zaddKey, 0, key); err != nil {
			logger.WithFields(logger.Fields{
				"object": x,
				"error":  err,
			}).Error("redis save object index failed")
		}
	})

	return nil
}

func (r *Redis) SaveFields(prefix string, x CacheObjector, fields map[string]interface{}) error {
	con, handler := r.getRejsonHandler()
	if handler == nil {
		return fmt.Errorf("redis.SaveFields failed: %w", con.Err())
	}

	key := fmt.Sprintf("%s:%v", prefix, x.GetObjID())

	r.Wrap(func() {
		defer r.returnRejsonHandler(con)

		for path, val := range fields {
			if _, err := handler.JSONSet(key, "."+path, val); err != nil {
				logger.WithFields(logger.Fields{
					"path":  "." + path,
					"value": val,
					"error": err,
				}).Error("redis save fields failed")
			}
		}
	})

	return nil
}

func (r *Redis) LoadObject(prefix string, value interface{}, x CacheObjector) error {
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
		return errors.New("cache object not found")
	}

	err = json.Unmarshal(res.([]byte), x)
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) LoadArray(prefix string, ownerId int64, pool *sync.Pool) ([]interface{}, error) {
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

func (r *Redis) DeleteObject(prefix string, x CacheObjector) error {
	con, handler := r.getRejsonHandler()
	if handler == nil {
		return fmt.Errorf("redis.DeleteObject failed:%w", con.Err())
	}

	key := fmt.Sprintf("%s:%v", prefix, x.GetObjID())

	r.Wrap(func() {
		defer r.returnRejsonHandler(con)

		if _, err := handler.JSONDel(key, "."); err != nil {
			logger.WithFields(logger.Fields{
				"object": x,
				"error":  err,
			}).Error("redis delete object failed")
		}

		// delete object index
		if x.GetStoreIndex() == -1 {
			return
		}

		zremKey := fmt.Sprintf("%s_index:%v", prefix, x.GetStoreIndex())
		if _, err := con.Do("ZREM", zremKey, key); err != nil {
			logger.WithFields(logger.Fields{
				"object": x,
				"error":  err,
			}).Error("redis delete object index failed")
		}
	})

	return nil
}

func (r *Redis) Exit() error {
	r.Wait()
	return r.pool.Close()
}

//func (r *Redis) Do(commandName string, args ...interface{}) (interface{}, error) {
//c := r.pool.Get()
//if c.Err() != nil {
//return nil, c.Err()
//}

//return c.Do(commandName, args...)
//}

//func (r *Redis) DoAsync(commandName string, cb RedisDoCallback, args ...interface{}) {
//r.Wrap(func() {
//c := r.pool.Get()
//if c.Err() != nil {
//cb(nil, c.Err())
//return
//}

//cb(c.Do(commandName, args...))
//})
//}

//func (r *Redis) Send(commandName string, args ...interface{}) (redis.Conn, error) {
//c := r.pool.Get()
//if c.Err() != nil {
//return c, c.Err()
//}

//return c, c.Send(commandName, args)
//}

//func (r *Redis) Flush(con redis.Conn) error {
//return con.Flush()
//}
