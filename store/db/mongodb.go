package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"bitbucket.org/east-eden/server/utils"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
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

	m.Lock()
	m.mapColls[name] = coll
	m.Unlock()

	return nil
}

func (m *MongoDB) LoadObject(tblName, key string, value interface{}, x interface{}) error {
	coll := m.getCollection(tblName)
	if coll == nil {
		coll = m.db.Collection(tblName)
	}

	filter := bson.D{}
	if len(key) > 0 && value != nil {
		filter = append(filter, bson.E{Key: key, Value: value})
	}

	ctx, cancel := context.WithTimeout(context.Background(), DatabaseLoadTimeout)
	defer cancel()
	res := coll.FindOne(ctx, filter)
	if res.Err() == nil {
		err := res.Decode(x)
		utils.ErrPrint(err, "mongodb load object failed", tblName, key)
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
		filter = append(filter, bson.E{Key: key, Value: storeIndex})
	}

	list := make([]interface{}, 0)
	ctx, cancel := context.WithTimeout(context.Background(), DatabaseLoadTimeout)
	defer cancel()
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return list, err
	}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		item := pool.Get()
		err := cur.Decode(item)
		if !utils.ErrCheck(err, "mongodb LoadArray decode item failed") {
			continue
		}

		list = append(list, item.(DBObjector))
		err = item.(DBObjector).AfterLoad()
		utils.ErrPrint(err, "mongodb LoadArray AfterLoad failed")
	}

	return list, nil
}

func (m *MongoDB) SaveObject(tblName string, k interface{}, x interface{}) error {
	coll := m.getCollection(tblName)
	if coll == nil {
		coll = m.db.Collection(tblName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), DatabaseUpdateTimeout)
	defer cancel()
	filter := bson.D{{Key: "_id", Value: k}}
	update := bson.D{{Key: "$set", Value: x}}
	op := options.Update().SetUpsert(true)

	if _, err := coll.UpdateOne(ctx, filter, update, op); err != nil {
		return fmt.Errorf("MongoDB.SaveObject failed: %w", err)
	}

	return nil
}

func (m *MongoDB) SaveFields(tblName string, k interface{}, fields map[string]interface{}) error {
	coll := m.getCollection(tblName)
	if coll == nil {
		coll = m.db.Collection(tblName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), DatabaseUpdateTimeout)
	defer cancel()
	filter := bson.D{{Key: "_id", Value: k}}

	values := bson.D{}
	for key, value := range fields {
		values = append(values, bson.E{Key: key, Value: value})
	}

	update := &bson.D{{Key: "$set", Value: values}}
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

	ctx, cancel := context.WithTimeout(context.Background(), DatabaseUpdateTimeout)
	defer cancel()
	filter := bson.D{{Key: "_id", Value: x.GetObjID()}}
	if _, err := coll.DeleteOne(ctx, filter); err != nil {
		return fmt.Errorf("MongoDB.DeleteObject failed: %w", err)
	}

	return nil
}

func (m *MongoDB) Exit() {
	err := m.c.Disconnect(context.Background())
	utils.ErrPrint(err, "mongodb disconnect failed")
}
