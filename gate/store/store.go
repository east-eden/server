package store

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"gopkg.in/mgo.v2/bson"
)

// Store combines cache and database
type Store struct {
	c     *mongo.Client
	db    *mongo.Database
	cache *Redis

	mapColls map[string]*mongo.Collection
	sync.RWMutex
}

func NewStore(ctx *cli.Context) *Store {
	s := &Store{
		cache:    NewRedis(ctx),
		mapColls: make(map[string]*mongo.Collection),
	}

	mongoCtx, _ := context.WithTimeout(ctx, 5*time.Second)
	dsn, ok := os.LookupEnv("DB_DSN")
	if !ok {
		dsn = ctx.String("db_dsn")
	}

	var err error
	if s.c, err = mongo.Connect(mongoCtx, options.Client().ApplyURI(dsn)); err != nil {
		logger.Fatal("NewStore failed:", err, "with dsn:", dsn)
		return nil
	}

	s.db = s.c.Database(ctx.String("database"))
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
	s.c.Disconnect(ctx)
	logger.Info("store exit...")
}

func (s *Store) MigrateCollection(name string, indexNames ...string) error {
	s.RLock()
	if _, ok := s.mapColls[name]; ok {
		s.RUnlock()
		return fmt.Errorf("duplicate collection %s", name)
	}
	s.RUnlock()

	coll := s.db.Collection(name)

	// check index
	idx := coll.Indexes()

	opts := options.ListIndexes().SetMaxTime(3 * time.Second)
	cursor, err := idx.List(context.Background(), opts)
	if err != nil {
		log.Fatal(err)
	}

	needsCreated := make(map[string]struct{})
	for _, indexName := range indexNames {
		needsCreated[indexName] = struct{}{}
	}

	for cursor.Next(context.Background()) {
		var result bson.M
		cursor.Decode(&result)

		delete(needsCreated, result["name"])
	}

	// no index needs to be created
	if len(needsCreated) == 0 {
		return nil
	}

	// create index
	for indexName, _ := range needsCreated {
		if _, err := coll.Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys:    bsonx.Doc{{indexName, bsonx.Int32(1)}},
				Options: options.Index().SetName(indexName),
			},
		); err != nil {
			logger.Fatalf("collection<%s> create index<%s> failed:%s", coll.Name(), indexName, err)
		}
	}

	s.Lock()
	s.mapColls[name] = coll
	s.Unlock()

	return nil
}

func (s *Store) GetCollection(name string) *mongo.Collection {
	s.RLock()
	defer s.RUnlock()

	return s.mapColls[name]
}

func (s *Store) CollectionUpdate(ctx context.Context, collName string, filter, update interface{}, opts ...*options.UpdateOptions) (*UpdateResult, error) {
	coll := s.GetCollection(collName)
	if coll == nil {
		return nil, errors.New("invalid collection name")
	}

	return coll.UpdateOne(ctx, filter, update, opts)
}

func (s *Store) CacheDo(commandName string, args ...interface{}) (interface{}, error) {
	return s.cache.Do(commandName, args)
}

func (s *Store) CacheDoAsync(commandName string, cb RedisDoCallback, args ...interface{}) {
	s.cache.DoAsync(commandName, cb, args)
}

//func (ds *Datastore) initDatastore(ctx context.Context) {
//ds.loadGate(ctx)
//}

//func (ds *Datastore) loadGate(ctx context.Context) {

//collection := ds.db.Collection(ds.tb.TableName())
//filter := bson.D{{"_id", ds.tb.ID}}
//replace := bson.D{{"_id", ds.tb.ID}, {"timestamp", ds.tb.TimeStamp}}
//op := options.FindOneAndReplace().SetUpsert(true)
//res := collection.FindOneAndReplace(ctx, filter, replace, op)
//res.Decode(ds.tb)

//logger.Info("datastore load table gate success:", ds.tb)
//}
