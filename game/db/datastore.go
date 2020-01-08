package db

import (
	"context"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/internal/define"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Datastore struct {
	c      *mongo.Client
	db     *mongo.Database
	ctx    context.Context
	cancel context.CancelFunc

	global *define.TableGlobal
}

func NewDatastore(ctx *cli.Context) *Datastore {
	ds := &Datastore{
		global: &define.TableGlobal{
			ID:        ctx.Int("game_id"),
			TimeStamp: int(time.Now().Unix()),
		},
	}

	ds.ctx, ds.cancel = context.WithCancel(ctx)

	mongoCtx, _ := context.WithTimeout(ds.ctx, 10*time.Second)
	var err error
	if ds.c, err = mongo.Connect(mongoCtx, options.Client().ApplyURI(ctx.String("db_dsn"))); err != nil {
		logger.Fatal("NewDatastore failed:", err, "with dsn:", ctx.String("db_dsn"))
		return nil
	}

	ds.db = ds.c.Database(ctx.String("database"))

	ds.initDatastore()
	return ds
}

func (ds *Datastore) Database() *mongo.Database {
	return ds.db
}

func (ds *Datastore) initDatastore() {
	ds.loadGlobal()
}

func (ds *Datastore) loadGlobal() {
	collection := ds.db.Collection(ds.global.TableName())
	filter := bson.D{{"_id", ds.global.ID}}
	update := bson.D{{"_id", ds.global.ID}, {"timestamp", ds.global.TimeStamp}}
	op := options.FindOneAndUpdate().SetUpsert(true)
	res := collection.FindOneAndUpdate(ds.ctx, filter, update, op)
	res.Decode(ds.global)
}

func (ds *Datastore) Run() error {
	for {
		select {
		case <-ds.ctx.Done():
			ds.Exit()
			logger.Info("Datastore context done...")
			return nil
		default:
			t := time.Now()
			d := time.Since(t)
			time.Sleep(time.Second - d)
		}
	}
}

func (ds *Datastore) Exit() {
	ds.c.Disconnect(ds.ctx)
}
