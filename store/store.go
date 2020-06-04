package store

import (
	"context"
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
	AfterLoad()
	AfterDelete()
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

func (s *Store) AddMemExpire(ctx context.Context, tp int, newFn func() interface{}) error {
	return s.mem.AddMemExpire(ctx, tp, newFn)
}

func (s *Store) MigrateDbTable(tblName string, indexNames ...string) error {
	return s.db.MigrateTable(tblName, indexNames...)
}

func (s *Store) SaveCacheObject(x cache.CacheObjector) error {
	return s.cache.SaveObject(x)
}

// LoadObject loads object from memory at first, if didn't hit, it will search from cache. if still find nothing, it will finally search from database.
func (s *Store) LoadObject(memType int, idxName string, key interface{}) (StoreObjector, error) {
	// load memory object will search object in memory, if not hit, it will return an object which allocated by memory's pool.
	x, err := s.mem.LoadObject(memType, key)
	if err == nil {
		return x.(StoreObjector), nil
	}

	// then search in cache, if hit, store it in memory
	err = s.cache.LoadObject(key, x)
	if err == nil {
		memExpire := s.mem.GetMemExpire(memType)
		memExpire.Store(x)
		return x.(StoreObjector), nil
	}

	logger.WithFields(logger.Fields{
		"memory_type": memType,
		"key":         key,
		"error":       err,
	}).Info("load cache object failed")

	// finally search in database, if hit, store it in both memory and cache
	err = s.db.LoadObject(idxName, key, x.(db.DBObjector))
	if err == nil {
		memExpire := s.mem.GetMemExpire(memType)
		memExpire.Store(x)
		s.cache.SaveObject(x)
		return x.(StoreObjector), nil
	}

	return nil, err
}

// SaveObject save object into memory, save into cache and database with async call.
func (s *Store) SaveObject(memType int, x StoreObjector) error {
	// save into memory
	errMem := s.mem.SaveObject(memType, x)

	// save into cache
	errCache := s.cache.SaveObject(x)

	// save into database
	errDb := s.db.SaveObject(x)

	if errMem != nil {
		return errMem
	}

	if errCache != nil {
		return errCache
	}

	return errDb
}
