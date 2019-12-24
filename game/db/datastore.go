package db

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/internal/define"
)

type Datastore struct {
	c      *mongo.Client
	dbName string
	ctx    context.Context
	cancel context.CancelFunc

	global *define.TableGlobal
}

func NewDatastore(id uint16, ctx *cli.Context) *Datastore {
	ds := &Datastore{
		global: &define.TableGlobal{
			ID:        int(id),
			TimeStamp: int(time.Now().Unix()),
		},
		dbName: ctx.String("database"),
	}

	ds.ctx, ds.cancel = context.WithCancel(ctx)

	mongoCtx, _ := context.WithTimeout(ds.ctx, 10*time.Second)
	var err error
	if ds.c, err = mongo.Connect(mongoCtx, options.Client().ApplyURI(ctx.db_dsn)); err != nil {
		logger.Fatal("NewDatastore failed:", err, "with dsn:", ctx.String("db_dsn"))
		return nil
	}

	ds.initDatastore()
	return ds
}

func (ds *Datastore) Mongo() *mongo.Client {
	return ds.c
}

func (ds *Datastore) initDatastore() {
	ds.loadGlobal()
}

func (ds *Datastore) loadGlobal() {

	//ds.orm.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(ds.global)
	//if ds.orm.FirstOrCreate(ds.global, ds.global.ID).RecordNotFound() {
	//ds.orm.Create(ds.global)
	//}

	logger.Info("datastore loadGlobal success:", ds.global)
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
	ds.orm.Close()
}
