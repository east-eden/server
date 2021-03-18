package store

import (
	"errors"
	"fmt"
	"sync"

	"bitbucket.org/funplus/server/store/cache"
	"bitbucket.org/funplus/server/store/db"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

// store find no result
var ErrNoResult = errors.New("store return no result")

// global store variables
var (
	gs Store
)

type StoreInfo struct {
	tp      int
	tblName string
	keyName string
}

type Store interface {
	InitCompleted() bool
	Exit()
	SetCache(cache.Cache)
	SetDB(db.DB)
	AddStoreInfo(tp int, tblName, keyName string)
	MigrateDbTable(tblName string, indexNames ...string) error
	LoadObject(storeType int, key interface{}, x interface{}) error
	LoadArray(storeType int, keyName string, keyValue interface{}, x interface{}) error
	SaveFields(storeType int, k interface{}, fields map[string]interface{}) error
	SaveObject(storeType int, k interface{}, x interface{}) error
	SaveMarshaledObject(storeType int, k interface{}, x interface{}) error
	DeleteObject(storeType int, k interface{}) error
	DeleteFields(storeType int, k interface{}, fieldsName []string) error
}

// defStore combines memory, cache and database
type defStore struct {
	cache    cache.Cache
	db       db.DB
	once     sync.Once
	done     bool
	infoList map[int]*StoreInfo
	sync.Mutex
}

func NewStore(ctx *cli.Context) Store {
	s := &defStore{}
	s.init(ctx)
	gs = s
	return gs
}

func GetStore() Store {
	return gs
}

func (s *defStore) init(ctx *cli.Context) {
	s.once.Do(func() {
		s.cache = cache.NewCache(ctx)
		s.db = db.NewDB(ctx)
		s.done = true
		s.infoList = make(map[int]*StoreInfo)
	})
}

func (s *defStore) InitCompleted() bool {
	return s.done
}

func (s *defStore) Exit() {
	s.cache.Exit()
	s.db.Exit()
	log.Info().Msg("store exit...")
}

func (s *defStore) SetCache(c cache.Cache) {
	s.cache = c
}

func (s *defStore) SetDB(db db.DB) {
	s.db = db
}

func (s *defStore) AddStoreInfo(tp int, tblName, keyName string) {
	s.Lock()
	defer s.Unlock()

	info := &StoreInfo{tp: tp, tblName: tblName, keyName: keyName}
	s.infoList[tp] = info
}

func (s *defStore) MigrateDbTable(tblName string, indexNames ...string) error {
	if !s.InitCompleted() {
		return errors.New("store didn't init")
	}

	return s.db.MigrateTable(tblName, indexNames...)
}

// LoadObject loads object from cache at first, if didn't hit, it will search from database. it neither search nor save with memory.
func (s *defStore) LoadObject(storeType int, key interface{}, x interface{}) error {
	if !s.InitCompleted() {
		return errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return fmt.Errorf("Store LoadObject: invalid store type %d", storeType)
	}

	// search in cache, if hit, store it in memory
	err := s.cache.LoadObject(info.tblName, key, x)
	if err == nil {
		return nil
	}

	// search in database, if hit, store it in both memory and cache
	err = s.db.LoadObject(info.tblName, info.keyName, key, x)
	if err == nil {
		return s.cache.SaveObject(info.tblName, key, x)
	}

	if errors.Is(err, db.ErrNoResult) {
		return ErrNoResult
	}

	return err
}

// LoadArray loads object from cache at first, if didn't hit, it will search from database. it neither search nor save with memory.
func (s *defStore) LoadArray(storeType int, keyName string, keyValue interface{}, x interface{}) error {
	if !s.InitCompleted() {
		return errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return fmt.Errorf("Store LoadArray: invalid store type %d", storeType)
	}

	// todo load array from cache
	// search in cache, if hit, store it in memory
	// err := s.cache.LoadObject(info.tblName, key, x)
	// if err == nil {
	// 	return nil
	// }

	// search in database, if hit, store it in both memory and cache
	err := s.db.LoadArray(info.tblName, keyName, keyValue, x)
	if err == nil {
		// todo save cache array
		// return s.cache.SaveObject(info.tblName, key, x)
	}

	if errors.Is(err, db.ErrNoResult) {
		return ErrNoResult
	}

	return err
}

// SaveFields save fields to cache and database with async call. it won't save to memory
func (s *defStore) SaveFields(storeType int, k interface{}, fields map[string]interface{}) error {
	if !s.InitCompleted() {
		return errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return fmt.Errorf("Store SaveFields: invalid store type %d", storeType)
	}

	// save into cache
	errCache := s.cache.SaveFields(info.tblName, k, fields)

	// save into database
	errDb := s.db.SaveFields(info.tblName, k, fields)

	if errCache != nil {
		return errCache
	}

	return errDb
}

// SaveObject save object cache and database with async call. it won't save to memory
func (s *defStore) SaveObject(storeType int, k interface{}, x interface{}) error {
	if !s.InitCompleted() {
		return errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return fmt.Errorf("Store SaveObject: invalid store type %d", storeType)
	}

	// save into cache
	errCache := s.cache.SaveObject(info.tblName, k, x)

	// save into database
	errDb := s.db.SaveObject(info.tblName, k, x)

	if errCache != nil {
		return errCache
	}

	return errDb
}

// SaveMarshaledObject save object cache and database with async call. it won't save to memory
func (s *defStore) SaveMarshaledObject(storeType int, k interface{}, x interface{}) error {
	if !s.InitCompleted() {
		return errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return fmt.Errorf("Store SaveObject: invalid store type %d", storeType)
	}

	// save into cache
	errCache := s.cache.SaveMarshaledObject(info.tblName, k, x)

	// save into database
	errDb := s.db.SaveObject(info.tblName, k, x)

	if errCache != nil {
		return errCache
	}

	return errDb
}

// DeleteObject delete object cache and database with async call. it won't delete from memory
func (s *defStore) DeleteObject(storeType int, k interface{}) error {
	if !s.InitCompleted() {
		return errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return fmt.Errorf("Store DeleteObject: invalid store type %d", storeType)
	}

	// delete from cache
	errCache := s.cache.DeleteObject(info.tblName, k)

	// delete from database
	errDb := s.db.DeleteObject(info.tblName, k)

	if errCache != nil {
		return errCache
	}

	return errDb
}

func (s *defStore) DeleteFields(storeType int, k interface{}, fieldsName []string) error {
	if !s.InitCompleted() {
		return errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return fmt.Errorf("Store DeleteFields: invalid store type %d", storeType)
	}

	// delete fields from cache
	errCache := s.cache.DeleteFields(info.tblName, k, fieldsName)

	// delete fields from database
	errDb := s.db.DeleteFields(info.tblName, k, fieldsName)

	if errCache != nil {
		return errCache
	}

	return errDb
}
