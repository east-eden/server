package store

import (
	"context"
	"errors"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/store/cache"
	"github.com/yokaiio/yokai_server/store/db"
)

const (
	StoreType_Begin = iota
	StoreType_User  = iota - 1
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
	"player",
	"player",
	"item",
	"hero",
	"blade",
	"token",
	"rune",
	"talent",
}

var (
	defaultStore = &Store{
		cache: nil,
		db:    nil,
		init:  false,
	}
)

// StoreObjector save and load with all structure
type StoreObjector interface {
	GetObjID() interface{}
	AfterLoad()
	TableName() string
}

// Store combines memory, cache and database
type Store struct {
	cache cache.Cache
	db    db.DB
	init  bool
}

func InitStore(ctx *cli.Context) {
	if !defaultStore.init {
		defaultStore.cache = cache.NewCache(ctx)
		defaultStore.db = db.NewDB(ctx)
		defaultStore.init = true
	}
}

func GetStore() *Store {
	return defaultStore
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
	s.cache.Exit(ctx)
	s.db.Exit(ctx)
	logger.Info("store exit...")
}

func (s *Store) MigrateDbTable(tblName string, indexNames ...string) error {
	if !s.init {
		return errors.New("store didn't init")
	}

	return s.db.MigrateTable(tblName, indexNames...)
}

// LoadObject loads object from cache at first, if didn't hit, it will search from database. it neither search nor save with memory.
func (s *Store) LoadObject(memType int, key string, value interface{}, x StoreObjector) error {
	if !s.init {
		return errors.New("store didn't init")
	}

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

func (s *Store) LoadArray(memType int, key string, value interface{}, pool *sync.Pool) ([]interface{}, error) {
	if !s.init {
		return nil, errors.New("store didn't init")
	}

	if memType < StoreType_Begin || memType >= StoreType_End {
		return nil, errors.New("memory type invalid")
	}

	cacheList, err := s.cache.LoadArray(StoreTypeNames[memType], pool)
	if err == nil {
		for _, val := range cacheList {
			val.(StoreObjector).AfterLoad()
		}

		return cacheList, nil
	}

	dbList, err := s.db.LoadArray(StoreTypeNames[memType], key, value, pool)
	if err == nil {
		for _, val := range dbList {
			val.(StoreObjector).AfterLoad()
			s.cache.SaveObject(StoreTypeNames[memType], val.(cache.CacheObjector))
		}
	}

	return dbList, err
}

// SaveFields save fields to cache and database with async call. it won't save to memory
func (s *Store) SaveFields(memType int, x StoreObjector, fields map[string]interface{}) error {
	if !s.init {
		return errors.New("store didn't init")
	}

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

// SaveObject save object cache and database with async call. it won't save to memory
func (s *Store) SaveObject(memType int, x StoreObjector) error {
	if !s.init {
		return errors.New("store didn't init")
	}

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
	if !s.init {
		return errors.New("store didn't init")
	}

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
