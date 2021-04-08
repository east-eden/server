package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"bitbucket.org/funplus/server/utils"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

var (
	ErrCollectionNotFound = errors.New("collection not found")
)

type MongoDB struct {
	c        *mongo.Client
	db       *mongo.Database
	mapColls map[string]*mongo.Collection
	sync.RWMutex
}

func NewMongoDB(ctx *cli.Context) DB {
	m := &MongoDB{
		mapColls: make(map[string]*mongo.Collection),
	}

	mongoCtx, cancel := context.WithTimeout(ctx.Context, 5*time.Second)
	defer cancel()
	dsn, ok := os.LookupEnv("DB_DSN")
	if !ok {
		dsn = ctx.String("db_dsn")
	}

	var err error
	if m.c, err = mongo.Connect(mongoCtx, options.Client().ApplyURI(dsn)); err != nil {
		log.Fatal().
			Str("dsn", dsn).
			Err(err).
			Msg("new mongodb failed")
		return nil
	}

	m.db = m.c.Database(ctx.String("database"))
	return m
}

func (m *MongoDB) getCollection(name string) *mongo.Collection {
	m.RLock()
	defer m.RUnlock()

	return m.mapColls[name]
}

func (m *MongoDB) setCollection(name string, coll *mongo.Collection) {
	m.Lock()
	defer m.Unlock()

	m.mapColls[name] = coll
}

func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	m.RLock()
	coll, ok := m.mapColls[name]
	m.RUnlock()

	if !ok {
		coll = m.db.Collection(name)
		m.setCollection(name, coll)
	}

	return coll
}

// migrate collection
func (m *MongoDB) MigrateTable(name string, indexNames ...string) error {
	if coll := m.getCollection(name); coll != nil {
		return fmt.Errorf("duplicate collection %s", name)
	}

	coll := m.db.Collection(name)

	// check index
	idx := coll.Indexes()

	opts := options.ListIndexes().SetMaxTime(3 * time.Second)
	cursor, err := idx.List(context.Background(), opts)
	if err != nil {
		log.Fatal().Err(err)
	}

	needsCreated := make(map[string]struct{})
	for _, indexName := range indexNames {
		needsCreated[indexName] = struct{}{}
	}

	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var result bson.M
		err := cursor.Decode(&result)
		utils.ErrPrint(err, "peek index failed", name, indexNames)

		delete(needsCreated, result["name"].(string))
	}

	// no index needs to be created
	if len(needsCreated) == 0 {
		return nil
	}

	// create index
	for indexName := range needsCreated {
		if _, err := coll.Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys:    bsonx.Doc{{Key: indexName, Value: bsonx.Int32(1)}},
				Options: options.Index().SetName(indexName),
			},
		); err != nil {
			log.Fatal().
				Str("coll_name", coll.Name()).
				Str("index_name", indexName).
				Err(err).
				Msg("create index failed")
		}
	}

	m.setCollection(name, coll)
	return nil
}

func (m *MongoDB) FindOne(ctx context.Context, colName string, filter interface{}, result interface{}) error {
	coll := m.GetCollection(colName)
	if coll == nil {
		return ErrCollectionNotFound
	}

	// timeout control
	subCtx, cancel := utils.WithTimeoutContext(ctx, DatabaseUpdateTimeout)
	defer cancel()

	res := coll.FindOne(subCtx, filter)
	if res.Err() == nil {
		err := res.Decode(result)
		utils.ErrPrint(err, "Decode failed when MongoDB.FindOne", colName, filter)
		return nil
	}

	// load success with no result
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return ErrNoResult
	}

	return res.Err()
}

func (m *MongoDB) Find(ctx context.Context, colName string, filter interface{}) (map[string]interface{}, error) {
	coll := m.GetCollection(colName)
	if coll == nil {
		return nil, ErrCollectionNotFound
	}

	// timeout control
	subCtx, cancel := utils.WithTimeoutContext(ctx, DatabaseUpdateTimeout)
	defer cancel()

	cur, err := coll.Find(subCtx, filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)
	var docs []map[string]interface{}
	err = cur.All(context.Background(), &docs)
	if !utils.ErrCheck(err, "cursor All failed when MongoDB.Find", filter) {
		return nil, err
	}

	result := make(map[string]interface{}, len(docs))
	for _, v := range docs {
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		result[fmt.Sprintf("%d", v["_id"])] = data
	}

	return result, nil
}

func (m *MongoDB) InsertOne(ctx context.Context, colName string, insert interface{}) error {
	coll := m.GetCollection(colName)
	if coll == nil {
		return ErrCollectionNotFound
	}

	// timeout control
	subCtx, cancel := utils.WithTimeoutContext(ctx, DatabaseUpdateTimeout)
	defer cancel()

	if _, err := coll.InsertOne(subCtx, insert); err != nil {
		return fmt.Errorf("MongoDB.InsertOne failed: %w", err)
	}

	return nil
}

func (m *MongoDB) InsertMany(ctx context.Context, colName string, inserts []interface{}) error {
	coll := m.GetCollection(colName)
	if coll == nil {
		return ErrCollectionNotFound
	}

	// timeout control
	subCtx, cancel := utils.WithTimeoutContext(ctx, DatabaseUpdateTimeout)
	defer cancel()

	if _, err := coll.InsertMany(subCtx, inserts); err != nil {
		return fmt.Errorf("MongoDB.InsertOne failed: %w", err)
	}

	return nil
}

func (m *MongoDB) UpdateOne(ctx context.Context, colName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) error {
	coll := m.GetCollection(colName)
	if coll == nil {
		return ErrCollectionNotFound
	}

	// timeout control
	subCtx, cancel := utils.WithTimeoutContext(ctx, DatabaseUpdateTimeout)
	defer cancel()

	if _, err := coll.UpdateOne(subCtx, filter, update, opts...); err != nil {
		return fmt.Errorf("MongoDB.UpdateOne failed: %w", err)
	}

	return nil
}

func (m *MongoDB) DeleteOne(ctx context.Context, colName string, filter interface{}) error {
	coll := m.GetCollection(colName)
	if coll == nil {
		return ErrCollectionNotFound
	}

	// timeout control
	subCtx, cancel := utils.WithTimeoutContext(ctx, DatabaseUpdateTimeout)
	defer cancel()

	_, err := coll.DeleteOne(subCtx, filter)
	return err
}

func (m *MongoDB) Exit() {
	err := m.c.Disconnect(context.Background())
	utils.ErrPrint(err, "mongodb disconnect failed")
}
