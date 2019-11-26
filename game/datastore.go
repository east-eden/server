package game

import (
	"context"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/game/define"
)

type Datastore struct {
	orm    *gorm.DB
	ctx    context.Context
	cancel context.CancelFunc
	g      *Game

	global *define.TableGlobal
}

func NewDatastore(game *Game, ctx *cli.Context) *Datastore {
	ds := &Datastore{
		g: game,
		global: &define.TableGlobal{
			ID:        game.ID,
			TimeStamp: int(time.Now().Unix()),
		},
	}

	ds.ctx, ds.cancel = context.WithCancel(ctx)

	var err error
	ds.orm, err = gorm.Open("mysql", ctx.String("db_dsn"))
	if err != nil {
		logger.Fatal("NewDatastore failed:", err, "with mysql dsn:", ctx.String("db_dsn"))
		return nil
	}

	ds.initDatastore()
	return ds
}

func (ds *Datastore) initDatastore() {
	ds.loadGlobal()
}

func (ds *Datastore) loadGlobal() {

	ds.orm.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(ds.global)
	if ds.orm.FirstOrCreate(ds.global, ds.global.ID).RecordNotFound() {
		ds.orm.Create(ds.global)
	}

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
