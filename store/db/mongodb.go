package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type MongoDB struct {
	c        *mongo.Client
	db       *mongo.Database
	mapColls map[string]*mongo.Collection
	sync.RWMutex
	utils.WaitGroupWrapper
}

func NewMongoDB(ctx *cli.Context) DB {
	m := &MongoDB{
		mapColls: make(map[string]*mongo.Collection),
	}

	mongoCtx, _ := context.WithTimeout(ctx, 5*time.Second)
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
		cursor.Decode(&result)

		delete(needsCreated, result["name"].(string))
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
			log.Fatal().
				Str("coll_name", coll.Name()).
				Str("index_name", indexName).
				Err(err).
				Msg("create index failed")
		}
	}

	m.Lock()
	m.mapColls[name] = coll
	m.Unlock()

	return nil
}

func (m *MongoDB) LoadObject(tblName, key string, value interface{}, x DBObjector) error {
	coll := m.getCollection(tblName)
	if coll == nil {
		coll = m.db.Collection(tblName)
	}

	filter := bson.D{}
	if len(key) > 0 && value != nil {
		filter = append(filter, bson.E{key, value})
	}

	ctx, _ := context.WithTimeout(context.Background(), DatabaseLoadTimeout)
	res := coll.FindOne(ctx, filter)
	if res.Err() == nil {
		res.Decode(x)
		return nil
	}

	// load success with no result
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return ErrNoResult
	}

	return res.Err()
}

func (m *MongoDB) LoadArray(tblName string, key string, storeIndex int64, pool *sync.Pool) ([]interface{}, error) {
	coll := m.getCollection(tblName)
	if coll == nil {
		coll = m.db.Collection(tblName)
	}

	filter := bson.D{}
	if len(key) > 0 && storeIndex != -1 {
		filter = append(filter, bson.E{key, storeIndex})
	}

	list := make([]interface{}, 0)
	ctx, _ := context.WithTimeout(context.Background(), DatabaseLoadTimeout)
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return list, err
	}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		item := pool.Get()
		if err := cur.Decode(item); err != nil {
			log.Warn().
				Err(err).
				Msg("LoadArray decode item failed")
			continue
		}

		list = append(list, item.(DBObjector))
		item.(DBObjector).AfterLoad()
	}

	return list, nil
}

func (m *MongoDB) SaveObject(tblName string, x DBObjector) error {
	coll := m.getCollection(tblName)
	if coll == nil {
		coll = m.db.Collection(tblName)
	}

	ctx, _ := context.WithTimeout(context.Background(), DatabaseUpdateTimeout)
	filter := bson.D{{"_id", x.GetObjID()}}
	update := bson.D{{"$set", x}}
	op := options.Update().SetUpsert(true)

	if _, err := coll.UpdateOne(ctx, filter, update, op); err != nil {
		return fmt.Errorf("MongoDB.SaveObject failed: %w", err)
	}

	return nil
}

func (m *MongoDB) SaveFields(tblName string, x DBObjector, fields map[string]interface{}) error {
	coll := m.getCollection(tblName)
	if coll == nil {
		coll = m.db.Collection(tblName)
	}

	ctx, _ := context.WithTimeout(context.Background(), DatabaseUpdateTimeout)
	filter := bson.D{{"_id", x.GetObjID()}}

	values := bson.D{}
	for key, value := range fields {
		values = append(values, bson.E{key, value})
	}

	update := &bson.D{{"$set", values}}
	op := options.Update().SetUpsert(true)
	if _, err := coll.UpdateOne(ctx, filter, update, op); err != nil {
		return fmt.Errorf("MongoDB.SaveFields failed: %w", err)
	}

	return nil
}

func (m *MongoDB) DeleteObject(tblName string, x DBObjector) error {
	coll := m.getCollection(tblName)
	if coll == nil {
		coll = m.db.Collection(tblName)
	}

	ctx, _ := context.WithTimeout(context.Background(), DatabaseUpdateTimeout)
	filter := bson.D{{"_id", x.GetObjID()}}
	if _, err := coll.DeleteOne(ctx, filter); err != nil {
		return fmt.Errorf("MongoDB.DeleteObject failed: %w", err)
	}

	return nil
}

func (m *MongoDB) Exit() {
	m.Wait()
	m.c.Disconnect(nil)
}
