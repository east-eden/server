package utils

import (
	"context"
	"sync"
	"time"

	"github.com/grafana/grafana/pkg/cmd/grafana-cli/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var expireNum = 1000

type CacheObjector interface {
	GetObjID() interface{}
	GetExpire() *time.Timer
	ResetExpire()
	StopExpire()
}

type CacheObjectNewFunc func() interface{}
type CacheDBLoadCB func(interface{})

type CacheLoader struct {
	mapObject sync.Map
	docField  string
	newFunc   CacheObjectNewFunc
	dbLoadCB  CacheDBLoadCB

	ctx    context.Context
	cancel context.CancelFunc

	coll      *mongo.Collection
	waitGroup WaitGroupWrapper

	chExpire chan interface{}
}

func NewCacheLoader(ctx context.Context, coll *mongo.Collection, docField string, newFunc CacheObjectNewFunc, dbCB CacheDBLoadCB) *CacheLoader {
	c := &CacheLoader{
		coll:     coll,
		docField: docField,
		newFunc:  newFunc,
		dbLoadCB: dbCB,
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

		// callback
		if c.dbLoadCB != nil {
			c.dbLoadCB(obj)
		}
		return obj.(CacheObjector)
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
				c.chExpire <- obj.GetObjID()
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
func (c *CacheLoader) Store(obj interface{}) {
	c.mapObject.Store(obj.(CacheObjector).GetObjID(), obj)
	c.beginTimeExpire(obj.(CacheObjector))
}

// get cache object, if not hit, load from database
func (c *CacheLoader) Load(key interface{}) CacheObjector {
	cache := c.LoadFromMemory(key)
	if cache == nil {
		cache = c.loadDBObject(key)
	}

	return cache
}

// delete cache, stop expire timer
func (c *CacheLoader) Delete(key interface{}) {
	cache := c.LoadFromMemory(key)
	if cache != nil {
		cache.StopExpire()
		c.mapObject.Delete(key)
	}
}

// only load from memory, usually you should only use CacheLoader.Load()
func (c *CacheLoader) LoadFromMemory(key interface{}) CacheObjector {
	v, ok := c.mapObject.Load(key)
	if ok {
		v.(CacheObjector).ResetExpire()
		return v.(CacheObjector)
	}

	return nil
}

// only load from database, usually you should only use CacheLoader.Load()
func (c *CacheLoader) LoadFromDB(key interface{}) CacheObjector {
	return c.loadDBObject(key)
}

func (c *CacheLoader) PureLoadFromDB(key interface{}) []CacheObjector {
	cur, err := c.coll.Find(c.ctx, bson.D{{c.docField, key}})
	if err != nil {
		logger.Warn("PureLoadFromDB failed:", err)
		return []CacheObjector{}
	}

	ret := make([]CacheObjector, 0)
	for cur.Next(c.ctx) {
		obj := c.newFunc()
		if err := cur.Decode(&obj); err != nil {
			logger.Warn("decode failed when call PureLoadFromDB:", err)
			continue
		}

		// callback
		if c.dbLoadCB != nil {
			c.dbLoadCB(obj)
		}

		ret = append(ret, obj.(CacheObjector))
	}

	return ret
}
