package store

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/utils"
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
				c, err := redis.DialTimeout("tcp", ctx.String("redis_address"), time.Second*2, time.Second*5, time.Second*5)
				if err != nil {
					panic(err.Error())
				}
				return c, err
			},
		},
	}
}

func (r *Redis) Do(commandName string, args ...interface{}) (interface{}, error) {
	c := r.pool.Get()
	if c.Err() != nil {
		return nil, c.Err()
	}

	return c.Do(commandName, args)
}

func (r *Redis) DoAsync(commandName string, cb RedisDoCallback, args ...interface{}) {
	r.Wrap(func() {
		c := r.pool.Get()
		if c.Err() != nil {
			cb(nil, c.Err())
			return
		}

		cb(c.Do(commandName, args))
	})
}

func (r *Redis) Send(commandName string, args ...interface{}) (redis.Conn, error) {
	c := r.pool.Get()
	if c.Err() != nil {
		return c, c.Err()
	}

	return c, c.Send(commandName, args)
}

func (r *Redis) Flush(con redis.Conn) error {
	return con.Flush()
}
