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
	"go.mongodb.org/mongo-driver/bson"
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
	GetDB() db.DB
	Exit()
	AddStoreInfo(tp int, tblName, keyName string)
	MigrateDbTable(tblName string, indexNames ...string) error
	FindOne(storeType int, key interface{}, x interface{}) error
	FindAll(storeType int, keyName string, keyValue interface{}) (interface{}, error)
	UpdateOne(storeType int, k interface{}, x interface{}) error
	UpdateFields(storeType int, k interface{}, fields map[string]interface{}) error
	DeleteFields(storeType int, k interface{}, fields []string) error

	// deprecated
	SaveHashObjectFields(storeType int, k interface{}, field interface{}, x interface{}, fields map[string]interface{}) error
	SaveHashObject(storeType int, k interface{}, field interface{}, x interface{}) error
	DeleteObject(storeType int, k interface{}) error
	DeleteObjectFields(storeType int, k interface{}, x interface{}, fields []string) error
	DeleteHashObject(storeType int, k interface{}, field interface{}) error
	DeleteHashObjectFields(storeType int, k interface{}, field interface{}, x interface{}, fields []string) error
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

func (s *defStore) GetDB() db.DB {
	return s.db
}

func (s *defStore) Exit() {
	s.cache.Exit()
	s.db.Exit()
	log.Info().Msg("store exit...")
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

// FindOne loads object from cache at first, if didn't hit, it will search from database. it neither search nor save with memory.
func (s *defStore) FindOne(storeType int, key interface{}, x interface{}) error {
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
	filter := bson.D{}
	filter = append(filter, bson.E{Key: info.keyName, Value: key})
	err = s.db.FindOne(info.tblName, filter, x)
	if err == nil {
		return s.cache.SaveObject(info.tblName, key, x)
	}

	if errors.Is(err, db.ErrNoResult) {
		return ErrNoResult
	}

	return err
}

// FindAll loads object from cache at first, if didn't hit, it will search from database. it neither search nor save with memory.
func (s *defStore) FindAll(storeType int, keyName string, keyValue interface{}) (interface{}, error) {
	if !s.InitCompleted() {
		return nil, errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return nil, fmt.Errorf("Store LoadArray: invalid store type %d", storeType)
	}

	// search in cache, if hit, store it in memory
	result, err := s.cache.LoadHashAll(info.tblName, keyValue)
	if err == nil {
		return result, nil
	}

	// search in database, if hit, store it in both memory and cache
	filter := bson.D{}
	filter = append(filter, bson.E{Key: keyName, Value: keyValue})
	result, err = s.db.Find(info.tblName, filter)
	if err == nil {
		// todo save hash all
		return result, s.cache.SaveHashAll(info.tblName, keyValue, result.(map[string]interface{}))
	}

	if errors.Is(err, db.ErrNoResult) {
		return nil, ErrNoResult
	}

	return result, err
}

// UpdateFields save fields to cache and database with async call. it won't save to memory
func (s *defStore) UpdateFields(storeType int, k interface{}, fields map[string]interface{}) error {
	if !s.InitCompleted() {
		return errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return fmt.Errorf("Store SaveFields: invalid store type %d", storeType)
	}

	// save into cache
	// errCache := s.cache.SaveObject(info.tblName, k, x)

	// save into database
	filter := bson.D{}
	filter = append(filter, bson.E{Key: "_id", Value: k})
	update := bson.D{}
	updateValues := bson.D{}
	for updateKey, updateValue := range fields {
		updateValues = append(updateValues, bson.E{Key: updateKey, Value: updateValue})
	}
	update = append(update, bson.E{Key: "$set", Value: updateValues})
	errDb := s.db.UpdateOne(info.tblName, filter, update)

	// if errCache != nil {
	// 	return errCache
	// }

	return errDb
}

func (s *defStore) DeleteFields(storeType int, k interface{}, fields []string) error {
	if !s.InitCompleted() {
		return errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return fmt.Errorf("Store DeleteFields: invalid store type %d", storeType)
	}

	// delete fields from database
	filter := bson.D{}
	filter = append(filter, bson.E{Key: "_id", Value: k})
	update := bson.D{}
	updateValues := bson.D{}
	for _, keyName := range fields {
		updateValues = append(updateValues, bson.E{Key: keyName, Value: 1})
	}
	update = append(update, bson.E{Key: "$unset", Value: updateValues})
	errDb := s.db.UpdateOne(info.tblName, filter, update)

	return errDb
}

// SaveHashObjectFields save fields to cache and database with async call. it won't save to memory
func (s *defStore) SaveHashObjectFields(storeType int, k interface{}, field interface{}, x interface{}, fields map[string]interface{}) error {
	if !s.InitCompleted() {
		return errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return fmt.Errorf("Store SaveFields: invalid store type %d", storeType)
	}

	// save into cache
	errCache := s.cache.SaveHashObject(info.tblName, k, field, x)

	// save into database
	filter := bson.D{}
	filter = append(filter, bson.E{Key: "_id", Value: k})
	update := bson.D{}
	updateValues := bson.D{}
	for updateKey, updateValue := range fields {
		updateValues = append(updateValues, bson.E{Key: updateKey, Value: updateValue})
	}
	update = append(update, bson.E{Key: "$set", Value: updateValues})
	errDb := s.db.UpdateOne(info.tblName, filter, update)

	if errCache != nil {
		return errCache
	}

	return errDb
}

func (s *defStore) SaveHashObject(storeType int, k interface{}, field interface{}, x interface{}) error {
	if !s.InitCompleted() {
		return errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return fmt.Errorf("Store SaveHashObject: invalid store type %d", storeType)
	}

	// save into cache
	errCache := s.cache.SaveHashObject(info.tblName, k, field, x)

	// save into database
	filter := bson.D{}
	filter = append(filter, bson.E{Key: "_id", Value: field})
	update := bson.D{}
	update = append(update, bson.E{Key: "$set", Value: x})
	errDb := s.db.UpdateOne(info.tblName, filter, update)

	if errCache != nil {
		return errCache
	}

	return errDb
}

// UpdateOne save object cache and database with async call. it won't save to memory
func (s *defStore) UpdateOne(storeType int, k interface{}, x interface{}) error {
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
	filter := bson.D{}
	filter = append(filter, bson.E{Key: "_id", Value: k})
	update := bson.D{}
	update = append(update, bson.E{Key: "$set", Value: x})
	errDb := s.db.UpdateOne(info.tblName, filter, update)

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
	filter := bson.D{}
	filter = append(filter, bson.E{Key: "_id", Value: k})
	errDb := s.db.DeleteOne(info.tblName, filter)

	if errCache != nil {
		return errCache
	}

	return errDb
}

func (s *defStore) DeleteObjectFields(storeType int, k interface{}, x interface{}, fields []string) error {
	if !s.InitCompleted() {
		return errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return fmt.Errorf("Store DeleteFields: invalid store type %d", storeType)
	}

	// delete fields from cache
	errCache := s.cache.SaveObject(info.tblName, k, x)

	// delete fields from database
	filter := bson.D{}
	filter = append(filter, bson.E{Key: "_id", Value: k})
	update := bson.D{}
	updateValues := bson.D{}
	for _, keyName := range fields {
		updateValues = append(updateValues, bson.E{Key: keyName, Value: 1})
	}
	update = append(update, bson.E{Key: "$unset", Value: updateValues})
	errDb := s.db.UpdateOne(info.tblName, filter, update)

	if errCache != nil {
		return errCache
	}

	return errDb
}

// DeleteHashObject delete object cache and database with async call. it won't delete from memory
func (s *defStore) DeleteHashObject(storeType int, k interface{}, field interface{}) error {
	if !s.InitCompleted() {
		return errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return fmt.Errorf("Store DeleteObject: invalid store type %d", storeType)
	}

	// delete from cache
	errCache := s.cache.DeleteHashObject(info.tblName, k, field)

	// delete from database
	filter := bson.D{}
	filter = append(filter, bson.E{Key: "_id", Value: field})
	errDb := s.db.DeleteOne(info.tblName, filter)

	if errCache != nil {
		return errCache
	}

	return errDb
}

func (s *defStore) DeleteHashObjectFields(storeType int, k interface{}, field interface{}, x interface{}, fields []string) error {
	if !s.InitCompleted() {
		return errors.New("store didn't init")
	}

	info, ok := s.infoList[storeType]
	if !ok {
		return fmt.Errorf("Store DeleteHashObjectFields: invalid store type %d", storeType)
	}

	// delete fields from cache
	errCache := s.cache.SaveHashObject(info.tblName, k, field, x)

	// delete fields from database
	filter := bson.D{}
	filter = append(filter, bson.E{Key: "_id", Value: k})
	update := bson.D{}
	updateValues := bson.D{}
	for _, keyName := range fields {
		updateValues = append(updateValues, bson.E{Key: keyName, Value: 1})
	}
	update = append(update, bson.E{Key: "$unset", Value: updateValues})
	errDb := s.db.UpdateOne(info.tblName, filter, update)

	if errCache != nil {
		return errCache
	}

	return errDb
}
