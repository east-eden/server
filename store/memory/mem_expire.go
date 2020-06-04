package memory

import (
	"context"
	"errors"
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
	AfterLoad()
	AfterDelete()
}

type MemExpire struct {
	mapObject sync.Map
	pool      sync.Pool

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

func (m *MemExpireManager) AddMemExpire(ctx context.Context, tp int, newFn func() interface{}) error {
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

func (m *MemExpireManager) GetMemExpire(tp int) *MemExpire {
	m.RLock()
	defer m.RUnlock()

	return m.mapMemExpire[tp]
}

func (m *MemExpireManager) LoadObject(memType int, key interface{}) (MemObjector, error) {
	memExpire := m.GetMemExpire(memType)
	if memExpire == nil {
		return nil, fmt.Errorf("invalid memory expire type %d", memType)
	}

	x, ok := memExpire.Load(key)
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

func newMemExpire(newFunc func() interface{}) *MemExpire {
	c := &MemExpire{
		chExpire: make(chan interface{}, MemExpireParallelNum),
	}

	c.pool.New = newFunc
	return c
}

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
		v.(MemObjector).AfterLoad()
		return v.(MemObjector), ok
	}

	return c.pool.Get().(MemObjector), false
}

// delete cache, stop expire timer
func (c *MemExpire) Delete(key interface{}) {
	if x, ok := c.Load(key); ok {
		c.mapObject.Delete(key)
		x.AfterDelete()
		c.pool.Put(x)
	}
}
