package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/store/cache"
	"github.com/yokaiio/yokai_server/store/db"
	"github.com/yokaiio/yokai_server/store/memory"
)

// StoreObjector save and load with all structure
type StoreObjector interface {
	GetObjID() interface{}
	GetExpire() *time.Timer
	ResetExpire()
	StopExpire()
	TableName() string
}

// Store combines memory, cache and database
type Store struct {
	mem   *memory.MemExpireManager
	cache cache.Cache
	db    db.DB
}

func NewStore(ctx *cli.Context) *Store {
	s := &Store{
		mem:   memory.NewMemExpireManager(),
		cache: cache.NewCache(ctx),
		db:    db.NewDB(ctx),
	}

	return s
}

func (s *Store) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			logger.Info("store context done...")
			return nil
		}
	}
}

func (s *Store) Exit(ctx context.Context) {
	s.db.Exit(ctx)
	logger.Info("store exit...")
}

func (s *Store) AddMemExpire(ctx context.Context, tp string, newFn func() interface{}) error {
	return s.mem.AddMemExpire(ctx, tp, newFn)
}

func (s *Store) MigrateDbTable(tblName string, indexNames ...string) error {
	return s.db.MigrateTable(tblName, indexNames...)
}

//func (s *Store) CacheDo(commandName string, args ...interface{}) (interface{}, error) {
//return s.cache.Do(commandName, args)
//}

//func (s *Store) CacheDoAsync(commandName string, cb RedisDoCallback, args ...interface{}) {
//s.cache.DoAsync(commandName, cb, args)
//}

func (s *Store) SaveCacheObject(x cache.CacheObjector) error {
	return s.cache.SaveObject(x)
}

// LoadObject loads object from memory at first, if didn't hit, it will search from cache. if still find nothing, it will finally search from database.
func (s *Store) LoadObject(name string, idxName string, key interface{}) (StoreObjector, error) {
	// load from memory
	x, err := s.LoadMemoryObject(name, key)
	if err == nil {
		return x, nil
	}

	// then search in cache, if hit, store it in memory
	if err := s.LoadCacheObject(key, x); err == nil {
		memExpire := s.mem.GetMemExpire(name)
		memExpire.Store(x)
		return x, nil
	}

	logger.WithFields(logger.Fields{
		"name":  name,
		"key":   key,
		"error": err,
	}).Info("load cache object failed")

	// finally search in database, if hit, store it in both memory and cache
	if err := s.LoadDBObject(idxName, key, x); err == nil {
		memExpire := s.mem.GetMemExpire(name)
		memExpire.Store(x)
		s.cache.SaveObject(x)
		return x, nil
	}

	return nil, errors.New("cannot find object")
}

// LoadMemoryObject will search object in memory, if not hit, it will return an object which allocated by memory's pool.
func (s *Store) LoadMemoryObject(name string, key interface{}) (StoreObjector, error) {
	memExpire := s.mem.GetMemExpire(name)
	if memExpire == nil {
		return nil, fmt.Errorf("invalid memory expire type %s", name)
	}

	x, ok := memExpire.Load(key)
	if ok {
		return x.(StoreObjector), nil
	}

	return x.(StoreObjector), errors.New("memory object not found")
}

func (s *Store) LoadCacheObject(key interface{}, x cache.CacheObjector) error {
	return s.cache.LoadObject(key, x)
}

func (s *Store) LoadDBObject(idxName string, key interface{}, x db.DBObjector) error {
	return s.db.LoadObject(idxName, key, x)
}

// SaveObject save object into memory, save into cache and database with async call.
func (s *Store) SaveObject(name string, x StoreObjector) error {
	memExpire := s.mem.GetMemExpire(name)
	if memExpire == nil {
		return fmt.Errorf("invalid memory expire type %s", name)
	}

	// save into memory
	memExpire.Store(x)

	// save into cache
	errCache := s.cache.SaveObject(x)

	// save into database
	errDb := s.db.SaveObject(x)

	if errCache != nil {
		return errCache
	}

	return errDb
}
