package db

import (
	"context"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/define"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Datastore struct {
	c  *mongo.Client
	db *mongo.Database

	tb *define.TableGate
}

func NewDatastore(ctx *cli.Context) *Datastore {
	ds := &Datastore{
		tb: &define.TableGate{
			ID:        ctx.Int("gate_id"),
			TimeStamp: int(time.Now().Unix()),
		},
	}

	mongoCtx, _ := context.WithTimeout(ctx, 10*time.Second)
	dsn, ok := os.LookupEnv("DB_DSN")
	if !ok {
		dsn = ctx.String("db_dsn")
	}

	var err error
	if ds.c, err = mongo.Connect(mongoCtx, options.Client().ApplyURI(dsn)); err != nil {
		logger.Fatal("NewDatastore failed:", err, "with dsn:", dsn)
		return nil
	}

	ds.db = ds.c.Database(ctx.String("database"))

	ds.initDatastore(ctx)
	return ds
}

func (ds *Datastore) initDatastore(ctx context.Context) {
	ds.loadGate(ctx)
}

func (ds *Datastore) loadGate(ctx context.Context) {

	collection := ds.db.Collection(ds.tb.TableName())
	filter := bson.D{{"_id", ds.tb.ID}}
	replace := bson.D{{"_id", ds.tb.ID}, {"timestamp", ds.tb.TimeStamp}}
	op := options.FindOneAndReplace().SetUpsert(true)
	res := collection.FindOneAndReplace(ctx, filter, replace, op)
	res.Decode(ds.tb)

	logger.Info("datastore load table gate success:", ds.tb)
}

func (ds *Datastore) Database() *mongo.Database {
	return ds.db
}

func (ds *Datastore) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			logger.Info("Datastore context done...")
			return nil
		}
	}
}

func (ds *Datastore) Exit(ctx context.Context) {
	ds.c.Disconnect(ctx)
	logger.Info("datastore exit...")
}
