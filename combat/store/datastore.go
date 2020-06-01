package store

import (
	"context"
	"os"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Datastore struct {
	c      *mongo.Client
	db     *mongo.Database
	ctx    context.Context
	cancel context.CancelFunc
	utils.WaitGroupWrapper
}

func newDatastore(ctx context.Context, dsn string, database string, combatId int) *Datastore {
	ds := &Datastore{}
	ds.ctx, ds.cancel = context.WithCancel(ctx)

	mongoCtx, _ := context.WithTimeout(ctx, 10*time.Second)

	var err error
	if ds.c, err = mongo.Connect(mongoCtx, options.Client().ApplyURI(dsn)); err != nil {
		logger.Fatal("NewDatastore failed:", err, "with dsn:", dsn)
		return nil
	}

	ds.db = ds.c.Database(database)

	return ds
}

func NewDatastore(ctx *cli.Context) *Datastore {
	dsn, ok := os.LookupEnv("DB_DSN")
	if !ok {
		dsn = ctx.String("db_dsn")
	}

	return newDatastore(ctx.Context, dsn, ctx.String("database"), ctx.Int("combat_id"))
}

func (ds *Datastore) Database() *mongo.Database {
	return ds.db
}

func (ds *Datastore) Run() error {
	for {
		select {
		case <-ds.ctx.Done():
			ds.Exit()
			logger.Info("Datastore context done...")
			return nil
		}
	}
}

func (ds *Datastore) Exit() {
	ds.Wait()
	ds.c.Disconnect(ds.ctx)
}
