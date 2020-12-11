package store

import (
	"errors"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/east-eden/server/store/cache"
	"github.com/east-eden/server/store/db"
)

// store find no result
var ErrNoResult = errors.New("store return no result")

// global store variables
var (
	defaultStore = &Store{
		cache: nil,
		db:    nil,
		init:  false,
	}
)

// StoreObjector save and load with all structure
type StoreObjector interface {
	GetObjID() int64
	GetStoreIndex() int64
	AfterLoad() error
}

type StoreInfo struct {
	tp        int
	tblName   string
	keyName   string
	indexName string
}

// Store combines memory, cache and database
type Store struct {
	cache    cache.Cache
	db       db.DB
	init     bool
	infoList map[int]*StoreInfo
	sync.Mutex
}

func InitStore(ctx *cli.Context) {
	if !defaultStore.init {
		defaultStore.cache = cache.NewCache(ctx)
		defaultStore.db = db.NewDB(ctx)
		defaultStore.init = true
		defaultStore.infoList = make(map[int]*StoreInfo)
	}
}

func GetStore() *Store {
	return defaultStore
}

func (s *Store) Exit() {
	s.cache.Exit()
	s.db.Exit()
	log.Info().Msg("store exit...")
}

func (s *Store) AddStoreInfo(tp int, tblName, keyName, indexName string) {
	s.Lock()
	defer s.Unlock()

	info := &StoreInfo{tp: tp, tblName: tblName, keyName: keyName, indexName: indexName}
	s.infoList[tp] = info
}

func (s *Store) MigrateDbTable(tblName string, indexNames ...string) error {
	if !s.init {
		return errors.New("store didn't init")
	}

	return s.db.MigrateTable(tblName, indexNames...)
}

// LoadObject loads object from cache at first, if didn't hit, it will search from database. it neither search nor save with memory.
func (s *Store) LoadObject(storeType int, keyValue interface{}, x StoreObjector) error {
	if !s.init {
		return errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return fmt.Errorf("Store LoadObject: invalid store type %d", storeType)
	}

	// search in cache, if hit, store it in memory
	err := s.cache.LoadObject(info.tblName, keyValue, x)
	if err == nil {
		return x.(StoreObjector).AfterLoad()
	}

	// search in database, if hit, store it in both memory and cache
	err = s.db.LoadObject(info.tblName, info.keyName, keyValue, x.(db.DBObjector))
	if err == nil {
		s.cache.SaveObject(info.tblName, x)
		return x.(StoreObjector).AfterLoad()
	}

	if errors.Is(err, db.ErrNoResult) {
		return ErrNoResult
	}

	return err
}

func (s *Store) LoadArray(storeType int, storeIndex int64, pool *sync.Pool) ([]interface{}, error) {
	if !s.init {
		return nil, errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return nil, fmt.Errorf("Store LoadArray: invalid store type %d", storeType)
	}

	cacheList, err := s.cache.LoadArray(info.tblName, storeIndex, pool)
	if err == nil {
		for _, val := range cacheList {
			if err := val.(StoreObjector).AfterLoad(); err != nil {
				return cacheList, err
			}
		}

		return cacheList, nil
	}

	dbList, err := s.db.LoadArray(info.tblName, info.indexName, storeIndex, pool)
	if err == nil {
		for _, val := range dbList {
			if err := val.(StoreObjector).AfterLoad(); err != nil {
				return dbList, err
			}
			s.cache.SaveObject(info.tblName, val.(cache.CacheObjector))
		}
	}

	return dbList, err
}

// SaveFields save fields to cache and database with async call. it won't save to memory
func (s *Store) SaveFields(storeType int, x StoreObjector, fields map[string]interface{}) error {
	if !s.init {
		return errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return fmt.Errorf("Store SaveFields: invalid store type %d", storeType)
	}

	// save into cache
	errCache := s.cache.SaveFields(info.tblName, x, fields)

	// save into database
	errDb := s.db.SaveFields(info.tblName, x, fields)

	if errCache != nil {
		return errCache
	}

	return errDb
}

// SaveObject save object cache and database with async call. it won't save to memory
func (s *Store) SaveObject(storeType int, x StoreObjector) error {
	if !s.init {
		return errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return fmt.Errorf("Store SaveObject: invalid store type %d", storeType)
	}

	// save into cache
	errCache := s.cache.SaveObject(info.tblName, x)

	// save into database
	errDb := s.db.SaveObject(info.tblName, x)

	if errCache != nil {
		return errCache
	}

	return errDb
}

// DeleteObject delete object cache and database with async call. it won't delete from memory
func (s *Store) DeleteObject(storeType int, x StoreObjector) error {
	if !s.init {
		return errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return fmt.Errorf("Store DeleteObject: invalid store type %d", storeType)
	}

	// delete from cache
	errCache := s.cache.DeleteObject(info.tblName, x)

	// delete from database
	errDb := s.db.DeleteObject(info.tblName, x)

	if errCache != nil {
		return errCache
	}

	return errDb
}
