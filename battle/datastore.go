package battle

import (
	"context"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/battle/define"
)

type Datastore struct {
	orm    *gorm.DB
	ctx    context.Context
	cancel context.CancelFunc
	b      *Battle

	tb *define.TableBattle
}

func NewDatastore(battle *Battle) *Datastore {
	ds := &Datastore{
		b: battle,
		tb: &define.TableBattle{
			ID:        battle.opts.BattleID,
			TimeStamp: int(time.Now().Unix()),
		},
	}

	ds.ctx, ds.cancel = context.WithCancel(battle.ctx)

	var err error
	ds.orm, err = gorm.Open("mysql", battle.opts.MysqlDSN)
	if err != nil {
		logger.Fatal("NewDatastore failed:", err)
		return nil
	}

	ds.initDatastore()
	return ds
}

func (ds *Datastore) initDatastore() {
	ds.loadBattle()
}

func (ds *Datastore) loadBattle() {

	ds.orm.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(ds.tb)
	if ds.orm.FirstOrCreate(ds.tb, ds.tb.ID).RecordNotFound() {
		ds.orm.Create(ds.tb)
	}

	logger.Info("datastore loadBattle success:", ds.tb)
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
