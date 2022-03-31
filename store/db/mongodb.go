package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/east-eden/server/utils"
	"github.com/hellodudu/channelwriter"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

var (
	ErrCollectionNotFound   = errors.New("collection not found")
	ErrBulkWriteInvalidType = errors.New("bulk write with invalid type")

	BulkWriteFlushLatency = time.Second * 2 // bulk write每隔两秒写入mongodb
)

type Collection struct {
	*mongo.Collection
	*channelwriter.ChannelWriter
}

func NewCollection(coll *mongo.Collection) *Collection {
	c := &Collection{
		Collection: coll,
	}

	c.ChannelWriter = channelwriter.NewChannelWriter(
		channelwriter.WithLogger(log.Logger),
		channelwriter.WithFlushHandler(c.flush),
	)

	return c
}

func (c *Collection) Write(p any) error {
	model, ok := p.(mongo.WriteModel)
	if !ok {
		return ErrBulkWriteInvalidType
	}

	c.ChannelWriter.Write(model)
	return nil
}

func (c *Collection) flush(datas []any) error {
	ctx, cancel := context.WithTimeout(context.Background(), DatabaseBulkWriteTimeout)
	defer cancel()

	models := make([]mongo.WriteModel, 0, len(datas))
	for _, data := range datas {
		model, ok := data.(mongo.WriteModel)
		if !ok {
			return ErrBulkWriteInvalidType
		}

		models = append(models, model)
	}
	res, err := c.Collection.BulkWrite(ctx, models)
	_ = utils.ErrCheck(err, "BulkWrite failed when Collection.Flush", datas, res)
	return err
}

func (c *Collection) Exit() {
	c.Stop()
}

type MongoDB struct {
	c        *mongo.Client
	db       *mongo.Database
	mapColls map[string]*Collection
	sync.RWMutex
}

func NewMongoDB(ctx *cli.Context) DB {
	m := &MongoDB{
		mapColls: make(map[string]*Collection),
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

func (m *MongoDB) getCollection(name string) *Collection {
	m.RLock()
	defer m.RUnlock()

	return m.mapColls[name]
}

func (m *MongoDB) setCollection(name string, coll *mongo.Collection) *Collection {
	m.Lock()
	defer m.Unlock()

	c := NewCollection(coll)
	m.mapColls[name] = c
	return c
}

func (m *MongoDB) GetCollection(name string) *Collection {
	m.RLock()
	coll, ok := m.mapColls[name]
	m.RUnlock()

	if !ok {
		coll = m.setCollection(name, m.db.Collection(name))
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

	_ = m.setCollection(name, coll)
	return nil
}

func (m *MongoDB) FindOne(ctx context.Context, colName string, filter any, result any) error {
	coll := m.GetCollection(colName)
	if coll == nil {
		return ErrCollectionNotFound
	}

	// timeout control
	subCtx, cancel := utils.WithTimeoutContext(ctx, DatabaseLoadTimeout)
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

func (m *MongoDB) Find(ctx context.Context, colName string, filter any) (map[string]any, error) {
	coll := m.GetCollection(colName)
	if coll == nil {
		return nil, ErrCollectionNotFound
	}

	// timeout control
	subCtx, cancel := utils.WithTimeoutContext(ctx, DatabaseLoadTimeout)
	defer cancel()

	cur, err := coll.Find(subCtx, filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)
	var docs []map[string]any
	err = cur.All(context.Background(), &docs)
	if !utils.ErrCheck(err, "cursor All failed when MongoDB.Find", filter) {
		return nil, err
	}

	result := make(map[string]any, len(docs))
	for _, v := range docs {
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		result[fmt.Sprintf("%d", v["_id"])] = data
	}

	return result, nil
}

func (m *MongoDB) InsertOne(ctx context.Context, colName string, insert any) error {
	coll := m.GetCollection(colName)
	if coll == nil {
		return ErrCollectionNotFound
	}

	// timeout control
	subCtx, cancel := utils.WithTimeoutContext(ctx, DatabaseWriteTimeout)
	defer cancel()

	if _, err := coll.InsertOne(subCtx, insert); err != nil {
		return fmt.Errorf("MongoDB.InsertOne failed: %w", err)
	}

	return nil
}

func (m *MongoDB) InsertMany(ctx context.Context, colName string, inserts []any) error {
	coll := m.GetCollection(colName)
	if coll == nil {
		return ErrCollectionNotFound
	}

	// timeout control
	subCtx, cancel := utils.WithTimeoutContext(ctx, DatabaseWriteTimeout)
	defer cancel()

	if _, err := coll.InsertMany(subCtx, inserts); err != nil {
		return fmt.Errorf("MongoDB.InsertOne failed: %w", err)
	}

	return nil
}

func (m *MongoDB) UpdateOne(ctx context.Context, colName string, filter any, update any, opts ...*options.UpdateOptions) error {
	coll := m.GetCollection(colName)
	if coll == nil {
		return ErrCollectionNotFound
	}

	// timeout control
	subCtx, cancel := utils.WithTimeoutContext(ctx, DatabaseWriteTimeout)
	defer cancel()

	if _, err := coll.UpdateOne(subCtx, filter, update, opts...); err != nil {
		return fmt.Errorf("MongoDB.UpdateOne failed: %w", err)
	}

	return nil
}

func (m *MongoDB) DeleteOne(ctx context.Context, colName string, filter any) error {
	coll := m.GetCollection(colName)
	if coll == nil {
		return ErrCollectionNotFound
	}

	// timeout control
	subCtx, cancel := utils.WithTimeoutContext(ctx, DatabaseWriteTimeout)
	defer cancel()

	_, err := coll.DeleteOne(subCtx, filter)
	return err
}

func (m *MongoDB) BulkWrite(ctx context.Context, colName string, model any) error {
	coll := m.GetCollection(colName)
	if coll == nil {
		return ErrCollectionNotFound
	}

	wm, ok := model.(mongo.WriteModel)
	if !ok {
		return ErrBulkWriteInvalidType
	}

	return coll.Write(wm)
}

func (m *MongoDB) Flush() {
	m.Lock()
	defer m.Unlock()
	for _, c := range m.mapColls {
		c.Flush()
	}
}

func (m *MongoDB) Exit() {
	var wg sync.WaitGroup
	m.Lock()
	for _, c := range m.mapColls {
		coll := c
		wg.Add(1)
		go func() {
			coll.Exit()
			wg.Done()
		}()
	}
	m.Unlock()

	wg.Wait()
	err := m.c.Disconnect(context.Background())
	utils.ErrPrint(err, "mongodb disconnect failed")
}
