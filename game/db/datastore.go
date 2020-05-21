package db

import (
	"context"
	"os"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Datastore struct {
	c  *mongo.Client
	db *mongo.Database
	utils.WaitGroupWrapper
}

func newDatastore(ctx context.Context, dsn string, database string, gameId int) *Datastore {
	ds := &Datastore{}

	mongoCtx, _ := context.WithTimeout(ctx, 10*time.Second)

	var err error
	if ds.c, err = mongo.Connect(mongoCtx, options.Client().ApplyURI(dsn)); err != nil {
		logger.Fatal("NewDatastore failed:", err, "with dsn:", dsn)
		return nil
	}

	ds.db = ds.c.Database(database)

	global := &define.TableGlobal{
		ID:        gameId,
		TimeStamp: int(time.Now().Unix()),
	}
	ds.loadGlobal(ctx, global)
	return ds
}

func NewDatastore(ctx *cli.Context) *Datastore {
	dsn, ok := os.LookupEnv("DB_DSN")
	if !ok {
		dsn = ctx.String("db_dsn")
	}

	return newDatastore(ctx, dsn, ctx.String("database"), ctx.Int("game_id"))
}

func (ds *Datastore) Database() *mongo.Database {
	return ds.db
}

func (ds *Datastore) loadGlobal(ctx context.Context, global *define.TableGlobal) {
	collection := ds.db.Collection(global.TableName())
	filter := bson.D{{"_id", global.ID}}
	update := bson.D{{"_id", global.ID}, {"timestamp", global.TimeStamp}}
	op := options.FindOneAndUpdate().SetUpsert(true)
	res := collection.FindOneAndUpdate(ctx, filter, update, op)
	res.Decode(global)
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
	ds.Wait()
	ds.c.Disconnect(ctx)
	logger.Info("datastore exit...")
}
