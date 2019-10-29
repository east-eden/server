package game

import (
	"context"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hellodudu/yokai_server/game/define"
	"github.com/jinzhu/gorm"
	logger "github.com/sirupsen/logrus"
)

type Datastore struct {
	orm    *gorm.DB
	ctx    context.Context
	cancel context.CancelFunc
	g      *Game

	global *define.TableGlobal
}

func NewDatastore(game *Game) (*Datastore, error) {
	db := &Datastore{
		g: game,
		global: &define.TableGlobal{
			ID:        game.opts.GameID,
			TimeStamp: int32(time.Now().Unix()),
		},
	}

	db.ctx, db.cancel = context.WithCancel(ctx)

	var err error
	db.orm, err = gorm.Open("mysql", game.opts.MysqlDSN)
	if err != nil {
		return nil, err
	}

	datastore.initDatastore()
	return datastore, nil
}

func (db *Datastore) initDatastore() {
	db.loadGlobal()
}

func (db *Datastore) loadGlobal() {

	db.orm.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(db.global)
	if db.orm.FirstOrCreate(db.global, db.global.ID).RecordNotFound() {
		db.orm.Create(db.global)
	}

	logger.Info("datastore loadGlobal success:", db.global)
}

func (db *Datastore) Run() error {
	for {
		select {
		case <-db.ctx.Done():
			db.Exit()
			return nil
		default:
			t := time.Now()
			d := time.Since(t)
			time.Sleep(time.Second - d)
		}
	}
}

func (db *Datastore) Exit() {
	logger.Info("datastore context done!")
	db.cancel()
	db.orm.Close()
}
