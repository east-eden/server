package store

import (
	"context"
	"errors"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/store/cache"
	"github.com/yokaiio/yokai_server/store/db"
	"github.com/yokaiio/yokai_server/store/memory"
)

const (
	StoreType_Begin = iota
	StoreType_User  = iota - 1
	StoreType_LiteAccount
	StoreType_Account
	StoreType_LitePlayer
	StoreType_Player
	StoreType_Item
	StoreType_Hero
	StoreType_Blade
	StoreType_Token
	StoreType_Rune
	StoreType_Talent

	StoreType_End
)

var StoreTypeNames = [StoreType_End]string{
	"user",
	"account",
	"account",
	"player",
	"player",
	"item",
	"hero",
	"blade",
	"token",
	"rune",
	"talent",
}

// StoreObjector save and load with all structure
type StoreObjector interface {
	GetObjID() interface{}
	GetExpire() *time.Timer
	AfterLoad()
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

func (s *Store) AddMemExpire(ctx context.Context, tp int, pool *sync.Pool, expire time.Duration) error {
	return s.mem.AddMemExpire(ctx, tp, pool, expire)
}

func (s *Store) MigrateDbTable(tblName string, indexNames ...string) error {
	return s.db.MigrateTable(tblName, indexNames...)
}

// LoadObject loads object from memory at first, if didn't hit, it will search from cache. if still find nothing, it will finally search from database.
func (s *Store) LoadObject(memType int, key string, value interface{}) (StoreObjector, error) {
	if memType < StoreType_Begin || memType >= StoreType_End {
		return nil, errors.New("memory type invalid")
	}

	// load memory object will search object in memory, if not hit, it will return an object which allocated by memory's pool.
	x, err := s.mem.LoadObject(memType, value)
	if err == nil {
		return x.(StoreObjector), nil
	}

	// then search in cache, if hit, store it in memory
	err = s.cache.LoadObject(StoreTypeNames[memType], value, x)
	if err == nil {
		s.mem.SaveObject(memType, x)
		x.(StoreObjector).AfterLoad()
		return x.(StoreObjector), nil
	}

	logger.WithFields(logger.Fields{
		"memory_type": memType,
		"key":         key,
		"error":       err,
	}).Info("load cache object failed")

	// finally search in database, if hit, store it in both memory and cache
	err = s.db.LoadObject(key, value, x.(db.DBObjector))
	if err == nil {
		s.mem.SaveObject(memType, x)
		s.cache.SaveObject(StoreTypeNames[memType], x)
		x.(StoreObjector).AfterLoad()
		return x.(StoreObjector), nil
	}

	// if all load failed, release x to memory pool
	s.mem.ReleaseObject(memType, x)
	return nil, err
}

// LoadObjectFromCacheAndDB loads object from cache at first, if didn't hit, it will search from database. it neither search nor save with memory.
func (s *Store) LoadObjectFromCacheAndDB(memType int, key string, value interface{}, x StoreObjector) error {
	if memType < StoreType_Begin || memType >= StoreType_End {
		return errors.New("memory type invalid")
	}

	// search in cache, if hit, store it in memory
	err := s.cache.LoadObject(StoreTypeNames[memType], value, x)
	if err == nil {
		x.(StoreObjector).AfterLoad()
		return nil
	}

	// search in database, if hit, store it in both memory and cache
	err = s.db.LoadObject(key, value, x.(db.DBObjector))
	if err == nil {
		s.cache.SaveObject(StoreTypeNames[memType], x)
		x.(StoreObjector).AfterLoad()
		return nil
	}

	return err
}

func (s *Store) LoadArrayFromCacheAndDB(memType int, key string, value interface{}, pool *sync.Pool) ([]db.DBObjector, error) {
	if memType < StoreType_Begin || memType >= StoreType_End {
		return nil, errors.New("memory type invalid")
	}

	// todo load from cache
	//s.cache.LoadArray(tblName, key, value, pool)

	return s.db.LoadArray(StoreTypeNames[memType], key, value, pool)
}

// SaveObject save object into memory, save into cache and database with async call.
func (s *Store) SaveObject(memType int, x StoreObjector) error {
	if memType < StoreType_Begin || memType >= StoreType_End {
		return errors.New("memory type invalid")
	}

	// save into memory
	errMem := s.mem.SaveObject(memType, x)

	// save into cache
	errCache := s.cache.SaveObject(StoreTypeNames[memType], x)

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

// SaveFieldsToCacheAndDB save fields to cache and database with async call. it won't save to memory
func (s *Store) SaveFieldsToCacheAndDB(memType int, x StoreObjector, fields map[string]interface{}) error {
	if memType < StoreType_Begin || memType >= StoreType_End {
		return errors.New("memory type invalid")
	}

	// save into cache
	errCache := s.cache.SaveFields(StoreTypeNames[memType], x, fields)

	// save into database
	errDb := s.db.SaveFields(x, fields)

	if errCache != nil {
		return errCache
	}

	return errDb
}

// SaveObjectToCacheAndDB save object cache and database with async call. it won't save to memory
func (s *Store) SaveObjectToCacheAndDB(memType int, x StoreObjector) error {
	if memType < StoreType_Begin || memType >= StoreType_End {
		return errors.New("memory type invalid")
	}

	// save into cache
	errCache := s.cache.SaveObject(StoreTypeNames[memType], x)

	// save into database
	errDb := s.db.SaveObject(x)

	if errCache != nil {
		return errCache
	}

	return errDb
}

// DeleteObjectFromCacheAndDB delete object cache and database with async call. it won't delete from memory
func (s *Store) DeleteObjectFromCacheAndDB(memType int, x StoreObjector) error {
	if memType < StoreType_Begin || memType >= StoreType_End {
		return errors.New("memory type invalid")
	}

	// delete from cache
	errCache := s.cache.DeleteObject(StoreTypeNames[memType], x)

	// delete from database
	errDb := s.db.DeleteObject(x)

	if errCache != nil {
		return errCache
	}

	return errDb
}
