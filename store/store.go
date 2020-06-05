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
	ExpireType_Begin = iota
	ExpireType_User  = iota - 1
	ExpireType_LiteAccount
	ExpireType_Account
	ExpireType_LitePlayer
	ExpireType_Player

	ExpireType_End
)

var ExpireTypeNames = [ExpireType_End]string{"user", "account", "account", "player", "player"}

// StoreObjector save and load with all structure
type StoreObjector interface {
	GetObjID() interface{}
	GetExpire() *time.Timer
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
func (s *Store) LoadObject(memType int, filter string, key interface{}) (StoreObjector, error) {
	if memType < ExpireType_Begin || memType >= ExpireType_End {
		return nil, errors.New("memory type invalid")
	}

	// load memory object will search object in memory, if not hit, it will return an object which allocated by memory's pool.
	x, err := s.mem.LoadObject(memType, key)
	if err == nil {
		return x.(StoreObjector), nil
	}

	// then search in cache, if hit, store it in memory
	err = s.cache.LoadObject(key, x)
	if err == nil {
		s.mem.SaveObject(memType, x)
		return x.(StoreObjector), nil
	}

	logger.WithFields(logger.Fields{
		"memory_type": memType,
		"key":         key,
		"error":       err,
	}).Info("load cache object failed")

	// finally search in database, if hit, store it in both memory and cache
	err = s.db.LoadObject(filter, key, x.(db.DBObjector))
	if err == nil {
		s.mem.SaveObject(memType, x)
		s.cache.SaveObject(ExpireTypeNames[memType], x)
		return x.(StoreObjector), nil
	}

	// if all load failed, release x to memory pool
	s.mem.ReleaseObject(memType, x)
	return nil, err
}

func (s *Store) LoadObjectArrayFromDB(tblName, filter string, key interface{}, pool *sync.Pool) ([]db.DBObjector, error) {
	return s.db.LoadObjectArray(tblName, filter, key, pool)
}

// SaveObject save object into memory, save into cache and database with async call.
func (s *Store) SaveObject(memType int, x StoreObjector) error {
	if memType < ExpireType_Begin || memType >= ExpireType_End {
		return errors.New("memory type invalid")
	}

	// save into memory
	errMem := s.mem.SaveObject(memType, x)

	// save into cache
	errCache := s.cache.SaveObject(ExpireTypeNames[memType], x)

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
