package memory

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
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
}

type MemExpire struct {
	mapObject sync.Map
	pool      *sync.Pool // external pool
	expire    time.Duration

	chExpire chan interface{}
}

type MemExpireManager struct {
	sync.RWMutex
	utils.WaitGroupWrapper
	mapMemExpire map[int]*MemExpire
}

func NewMemExpireManager() *MemExpireManager {
	return &MemExpireManager{
		mapMemExpire: make(map[int]*MemExpire),
	}
}

func (m *MemExpireManager) AddMemExpire(ctx context.Context, tp int, pool *sync.Pool, expire time.Duration) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.mapMemExpire[tp]; ok {
		return fmt.Errorf("init existing mem expire:%d", tp)
	}

	ex := newMemExpire(pool, expire)
	m.mapMemExpire[tp] = ex
	m.Wrap(func() {
		ex.Run(ctx)
	})
	return nil
}

func (m *MemExpireManager) GetMemExpire(tp int) *MemExpire {
	m.RLock()
	defer m.RUnlock()

	return m.mapMemExpire[tp]
}

func (m *MemExpireManager) LoadObject(memType int, value interface{}) (MemObjector, error) {
	memExpire := m.GetMemExpire(memType)
	if memExpire == nil {
		return nil, fmt.Errorf("invalid memory expire type %d", memType)
	}

	x, ok := memExpire.Load(value)
	if ok {
		return x, nil
	}

	return x, errors.New("memory object not found")
}

func (m *MemExpireManager) SaveObject(memType int, x MemObjector) error {
	memExpire := m.GetMemExpire(memType)
	if memExpire == nil {
		return fmt.Errorf("invalid memory expire type %d", memType)
	}

	memExpire.Store(x)
	return nil
}

func (m *MemExpireManager) DeleteObject(memType int, key interface{}) error {
	memExpire := m.GetMemExpire(memType)
	if memExpire == nil {
		return fmt.Errorf("invalid memory expire type %d", memType)
	}

	memExpire.Delete(key)
	return nil
}

func (m *MemExpireManager) ReleaseObject(memType int, x MemObjector) error {
	memExpire := m.GetMemExpire(memType)
	if memExpire == nil {
		return fmt.Errorf("invalid memory expire type %d", memType)
	}

	memExpire.Release(x)
	return nil
}

func newMemExpire(pool *sync.Pool, expire time.Duration) *MemExpire {
	c := &MemExpire{
		pool:     pool,
		expire:   expire,
		chExpire: make(chan interface{}, MemExpireParallelNum),
	}

	return c
}

func (c *MemExpire) beginTimeExpire(x MemObjector) {
	tm := x.GetExpire()
	if tm != nil {
		tm.Reset(c.expire + time.Second*time.Duration(rand.Intn(60)))
		return
	}

	// memcache time expired
	tm = time.NewTimer(c.expire)
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
		case key := <-c.chExpire:
			c.Delete(key)
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
	x, ok := c.mapObject.Load(key)
	if ok {
		c.beginTimeExpire(x.(MemObjector))
		return x.(MemObjector), ok
	}

	return c.pool.Get().(MemObjector), false
}

// delete cache, stop expire timer
func (c *MemExpire) Delete(key interface{}) {
	if x, ok := c.Load(key); ok {
		x.(MemObjector).GetExpire().Stop()
		c.mapObject.Delete(key)
		c.pool.Put(x)
	}
}

// release memory to pool
func (c *MemExpire) Release(x interface{}) {
	c.pool.Put(x)
}
