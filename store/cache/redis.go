package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	RedisConnectTimeout = time.Minute * 2
	RedisReadTimeout    = time.Minute * 5
	RedisWriteTimeout   = time.Minute * 5
)

type RedisDoCallback func(interface{}, error)

type Redis struct {
	pool *redis.Pool
	utils.WaitGroupWrapper
	mapRejsonHandler map[redis.Conn]*rejson.Handler
	sync.RWMutex
}

func NewRedis(ctx *cli.Context) *Redis {
	return &Redis{
		pool: &redis.Pool{
			MaxIdle:   80,
			MaxActive: 12000,
			Dial: func() (redis.Conn, error) {
				c, err := redis.DialTimeout("tcp", ctx.String("redis_addr"), RedisConnectTimeout, RedisReadTimeout, RedisWriteTimeout)
				if err != nil {
					panic(err.Error())
				}
				return c, err
			},
		},
		mapRejsonHandler: make(map[redis.Conn]*rejson.Handler),
	}
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

func (r *Redis) SaveObject(prefix string, x CacheObjector) error {
	_, handler := r.getRejsonHandler()
	if handler == nil {
		return errors.New("can't get rejson handler")
	}

	key := fmt.Sprintf("%s:%v", prefix, x.GetObjID())

	r.Wrap(func() {
		if _, err := handler.JSONSet(key, ".", x); err != nil {
			logger.WithFields(logger.Fields{
				"object": x,
				"error":  err,
			}).Error("redis save object failed")
		}
	})

	return nil
}

func (r *Redis) SaveFields(prefix string, x CacheObjector, fields map[string]interface{}) error {
	_, handler := r.getRejsonHandler()
	if handler == nil {
		return errors.New("can't get rejson handler")
	}

	key := fmt.Sprintf("%s:%v", prefix, x.GetObjID())

	r.Wrap(func() {
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
	_, handler := r.getRejsonHandler()
	if handler == nil {
		return errors.New("can't get rejson handler")
	}

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

func (r *Redis) LoadArray(prefix string, pool *sync.Pool) ([]interface{}, error) {
	c, handler := r.getRejsonHandler()
	if handler == nil {
		return nil, errors.New("can't get rejson handler")
	}

	// scan all keys
	var (
		cursor int64
		items  []string
	)
	results := make([]string, 0)

	for {
		values, err := redis.Values(c.Do("SCAN", cursor, "MATCH", prefix+":*"))
		if err != nil {
			return nil, err
		}

		values, err = redis.Scan(values, &cursor, &items)
		if err != nil {
			return nil, err
		}

		results = append(results, items...)

		if cursor == 0 {
			break
		}
	}

	reply := make([]interface{}, 0)
	for _, key := range results {
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
	_, handler := r.getRejsonHandler()
	if handler == nil {
		return errors.New("can't get rejson handler")
	}

	key := fmt.Sprintf("%s:%v", prefix, x.GetObjID())

	_, err := handler.JSONDel(key, ".")
	return err
}

func (r *Redis) Exit(ctx context.Context) error {
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
