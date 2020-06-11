package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	logger "github.com/sirupsen/logrus"
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
		logger.Fatal("new mongodb failed:", err, "with dsn:", dsn)
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
		log.Fatal(err)
	}

	needsCreated := make(map[string]struct{})
	for _, indexName := range indexNames {
		needsCreated[indexName] = struct{}{}
	}

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
			logger.Fatalf("collection<%s> create index<%s> failed:%s", coll.Name(), indexName, err)
		}
	}

	m.Lock()
	m.mapColls[name] = coll
	m.Unlock()

	return nil
}

func (m *MongoDB) LoadObject(key string, value interface{}, x DBObjector) error {
	coll := m.getCollection(x.TableName())
	if coll == nil {
		coll = m.db.Collection(x.TableName())
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

	return res.Err()
}

func (m *MongoDB) LoadArray(tblName string, key string, value interface{}, pool *sync.Pool) ([]interface{}, error) {
	coll := m.getCollection(tblName)
	if coll == nil {
		coll = m.db.Collection(tblName)
	}

	filter := bson.D{}
	if len(key) > 0 && value != nil {
		filter = append(filter, bson.E{key, value})
	}

	list := make([]interface{}, 0)
	ctx, _ := context.WithTimeout(context.Background(), DatabaseLoadTimeout)
	cur, err := coll.Find(ctx, filter)
	defer cur.Close(ctx)
	if err != nil {
		return list, err
	}

	for cur.Next(ctx) {
		item := pool.Get()
		if err := cur.Decode(item); err != nil {
			logger.Warn("item decode failed:", err)
			continue
		}

		list = append(list, item.(DBObjector))
		item.(DBObjector).AfterLoad()
	}

	return list, nil
}

func (m *MongoDB) SaveObject(x DBObjector) error {
	coll := m.getCollection(x.TableName())
	if coll == nil {
		coll = m.db.Collection(x.TableName())
	}

	ctx, _ := context.WithTimeout(context.Background(), DatabaseUpdateTimeout)
	filter := bson.D{{"_id", x.GetObjID()}}
	update := bson.D{{"$set", x}}
	op := options.Update().SetUpsert(true)

	m.Wrap(func() {
		if _, err := coll.UpdateOne(ctx, filter, update, op); err != nil {
			logger.WithFields(logger.Fields{
				"collection": coll.Name(),
				"filter":     filter,
				"error":      err,
			}).Warning("mongodb save object failed")
		}
	})

	return nil
}

func (m *MongoDB) SaveFields(x DBObjector, fields map[string]interface{}) error {
	coll := m.getCollection(x.TableName())
	if coll == nil {
		coll = m.db.Collection(x.TableName())
	}

	ctx, _ := context.WithTimeout(context.Background(), DatabaseUpdateTimeout)
	filter := bson.D{{"_id", x.GetObjID()}}

	values := bson.D{}
	for key, value := range fields {
		values = append(values, bson.E{key, value})
	}

	update := &bson.D{{"$set", values}}
	op := options.Update().SetUpsert(true)

	m.Wrap(func() {
		if _, err := coll.UpdateOne(ctx, filter, update, op); err != nil {
			logger.WithFields(logger.Fields{
				"collection": coll.Name(),
				"filter":     filter,
				"fields":     fields,
				"error":      err,
			}).Warning("mongodb save fields failed")
		}
	})

	return nil
}

func (m *MongoDB) DeleteObject(x DBObjector) error {
	coll := m.getCollection(x.TableName())
	if coll == nil {
		coll = m.db.Collection(x.TableName())
	}

	ctx, _ := context.WithTimeout(context.Background(), DatabaseUpdateTimeout)
	filter := bson.D{{"_id", x.GetObjID()}}

	m.Wrap(func() {
		if _, err := coll.DeleteOne(ctx, filter); err != nil {
			logger.WithFields(logger.Fields{
				"collection": coll.Name(),
				"filter":     filter,
				"error":      err,
			}).Warning("mongodb delete object failed")
		}
	})

	return nil
}

func (m *MongoDB) Exit(ctx context.Context) {
	m.c.Disconnect(ctx)
}
