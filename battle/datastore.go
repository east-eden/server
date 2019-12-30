package battle

import (
	"context"
	"time"

	_ "github.com/go-sql-driver/mysql"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/internal/define"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Datastore struct {
	c      *mongo.Client
	db     *mongo.Database
	ctx    context.Context
	cancel context.CancelFunc
	b      *Battle

	tb *define.TableBattle
}

func NewDatastore(battle *Battle, ctx *cli.Context) *Datastore {
	ds := &Datastore{
		b: battle,
		tb: &define.TableBattle{
			ID:        ctx.Int("battle_id"),
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

func (ds *Datastore) initDatastore() {
	ds.loadBattle()
}

func (ds *Datastore) loadBattle() {

	collection := ds.db.Collection(ds.tb.TableName())
	filter := bson.D{{"_id", ds.tb.ID}}
	replace := bson.D{{"_id", ds.tb.ID}, {"timestamp", ds.tb.TimeStamp}}
	op := options.FindOneAndReplace().SetUpsert(true)
	res := collection.FindOneAndReplace(ds.ctx, filter, replace, op)
	res.Decode(ds.tb)

	logger.Info("datastore load table battle success:", ds.tb)
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
