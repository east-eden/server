package utils

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CacheObjector interface {
	GetID() interface{}
	GetExpire() *time.Timer
	ResetExpire()
}

type CacheObjectNewFunc func() CacheObjector

type CacheLoader struct {
	mapObject sync.Map
	docField  string
	newFunc   CacheObjectNewFunc

	ctx    context.Context
	cancel context.CancelFunc

	coll      *mongo.Collection
	waitGroup WaitGroupWrapper

	chExpire chan interface{}
}

func NewCacheLoader(ctx context.Context, docField string, expireNum int, newFunc CacheObjectNewFunc) *CacheLoader {
	c := &CacheLoader{
		docField: docField,
		newFunc:  newFunc,
		chExpire: make(chan interface{}, expireNum),
	}

	c.ctx, c.cancel = context.WithCancel(ctx)

	c.waitGroup.Wrap(func() {
		c.run()
	})

	return c
}

func (c *CacheLoader) loadDBObject(key interface{}) CacheObjector {
	res := c.coll.FindOne(c.ctx, bson.D{{c.docField, key}})
	if res.Err() == nil {
		obj := c.newFunc()
		res.Decode(obj)

		c.Store(obj)
		return obj
	}

	return nil
}

func (c *CacheLoader) beginTimeExpire(obj CacheObjector) {
	// memcache time expired
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return

			case <-obj.GetExpire().C:
				c.chExpire <- obj.GetID()
			}
		}
	}()
}

func (c *CacheLoader) run() error {
	for {
		select {
		case <-c.ctx.Done():
			return nil

		// memcache time expired
		case id := <-c.chExpire:
			c.mapObject.Delete(id)
		}
	}

	return nil
}

// save cache object and begin count down timer
func (c *CacheLoader) Store(obj CacheObjector) {
	c.mapObject.Store(obj.GetID(), obj)
	c.beginTimeExpire(obj)
}

// get cache object, if not hit, load from database
func (c *CacheLoader) Load(key interface{}) CacheObjector {
	v, ok := c.mapObject.Load(key)
	if ok {
		v.(CacheObjector).ResetExpire()
		return v.(CacheObjector)
	}

	return c.loadDBObject(key)
}
