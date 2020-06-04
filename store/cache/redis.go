package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/utils"
)

var (
	RedisConnectTimeout = time.Second * 2
	RedisReadTimeout    = time.Second * 5
	RedisWriteTimeout   = time.Second * 5
)

type RedisDoCallback func(interface{}, error)

type Redis struct {
	pool *redis.Pool
	utils.WaitGroupWrapper
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
	}
}

func (r *Redis) SaveObject(prefix string, x CacheObjector) error {
	c := r.pool.Get()
	if c.Err() != nil {
		return c.Err()
	}

	key := fmt.Sprintf("%s:%v", prefix, x.GetObjID())

	r.Wrap(func() {
		if _, err := c.Do("HMSET", redis.Args{}.Add(key).AddFlat(x)...); err != nil {
			logger.WithFields(logger.Fields{
				"object": x,
				"error":  err,
			}).Error("redis save object failed")
		}
	})

	return nil
}

func (r *Redis) LoadObject(key interface{}, x CacheObjector) error {
	c := r.pool.Get()
	if c.Err() != nil {
		return c.Err()
	}

	val, err := redis.Values(c.Do("HGETALL", key))
	if err != nil {
		return err
	}

	if len(val) == 0 {
		return errors.New("empty array")
	}

	if err := redis.ScanStruct(val, x); err != nil {
		return err
	}

	return nil
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
