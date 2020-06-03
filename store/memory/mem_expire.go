package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/yokaiio/yokai_server/utils"
)

var (
	MemExpireTimeout     = time.Second * 10
	MemExpireParallelNum = 100
)

type MemObjector interface {
	GetObjID() interface{}
	GetExpire() *time.Timer
	ResetExpire()
	StopExpire()
}

type MemExpire struct {
	mapObject sync.Map
	pool      sync.Pool

	chExpire chan interface{}
}

type MemExpireManager struct {
	sync.RWMutex
	utils.WaitGroupWrapper
	mapMemExpire map[string]*MemExpire
}

func NewMemExpireManager() *MemExpireManager {
	return &MemExpireManager{
		mapMemExpire: make(map[string]*MemExpire),
	}
}

func (m *MemExpireManager) AddMemExpire(ctx context.Context, tp string, newFn func() interface{}) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.mapMemExpire[tp]; ok {
		return fmt.Errorf("init existing mem expire:", tp)
	}

	ex := newMemExpire(newFn)
	m.mapMemExpire[tp] = ex
	m.Wrap(func() {
		ex.Run(ctx)
	})
	return nil
}

func (m *MemExpireManager) GetMemExpire(tp string) *MemExpire {
	m.RLock()
	defer m.RUnlock()

	return m.mapMemExpire[tp]
}

func newMemExpire(newFunc func() interface{}) *MemExpire {
	c := &MemExpire{
		chExpire: make(chan interface{}, MemExpireParallelNum),
	}

	c.pool.New = newFunc
	return c
}

//func (c *MemExpire) loadDBObject(key interface{}) MemObjector {
//if c.store == nil {
//return nil
//}

//ctx, _ := context.WithTimeout(context.Background(), MemExpireTimeout)
//res := c.coll.FindOne(ctx, bson.D{{c.docField, key}})
//if res.Err() == nil {
//obj := c.pool.Get()
//res.Decode(obj)
//c.Store(obj)

//// callback
//if c.dbLoadCB != nil {
//c.dbLoadCB(obj)
//}
//return obj.(MemObjector)
//}

//return nil
//}

func (c *MemExpire) beginTimeExpire(x MemObjector) {
	// memcache time expired
	go func() {
		select {
		case <-x.GetExpire().C:
			c.chExpire <- x.GetObjID()
		}
	}()
}

func (c *MemExpire) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		// memcache time expired
		case id := <-c.chExpire:
			c.Delete(id)
		}
	}

	return nil
}

// save cache object and begin count down timer
func (c *MemExpire) Store(x interface{}) {
	c.mapObject.Store(x.(MemObjector).GetObjID(), x)
	c.beginTimeExpire(x.(MemObjector))
}

// get cache object, if not hit, load from database
func (c *MemExpire) Load(key interface{}) (MemObjector, bool) {
	v, ok := c.mapObject.Load(key)
	if ok {
		v.(MemObjector).ResetExpire()
		return v.(MemObjector), ok
	}

	return c.pool.Get().(MemObjector), false
}

// delete cache, stop expire timer
func (c *MemExpire) Delete(key interface{}) {
	if cache, ok := c.Load(key); ok {
		cache.StopExpire()
		c.mapObject.Delete(key)
		c.pool.Put(cache)
	}
}

// only load from database, usually you should only use MemExpire.Load()
//func (c *MemExpire) LoadFromDB(key interface{}) MemObjector {
//return c.loadDBObject(key)
//}

//func (c *MemExpire) PureLoadFromDB(key interface{}) []MemObjector {
//ret := make([]MemObjector, 0)
//if c.coll == nil {
//return ret
//}

//ctx, _ := context.WithTimeout(context.Background(), MemExpireTimeout)
//cur, err := c.coll.Find(ctx, bson.D{{c.docField, key}})
//if err != nil {
//logger.Warn("PureLoadFromDB failed:", err)
//return []MemObjector{}
//}

//for cur.Next(ctx) {
//obj := c.pool.Get()
//if err := cur.Decode(&obj); err != nil {
//logger.Warn("decode failed when call PureLoadFromDB:", err)
//continue
//}

//// callback
//if c.dbLoadCB != nil {
//c.dbLoadCB(obj)
//}

//ret = append(ret, obj.(MemObjector))
//}

//return ret
//}
