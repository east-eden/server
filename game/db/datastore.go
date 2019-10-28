package db

import (
	"context"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hellodudu/Ultimate/iface"
	"github.com/hellodudu/Ultimate/utils/global"
	"github.com/jinzhu/gorm"
	logger "github.com/sirupsen/logrus"
)

type Datastore struct {
	db     *gorm.DB
	ctx    context.Context
	cancel context.CancelFunc
	chStop chan struct{}

	// table
	global *iface.TableGlobal
}

func NewDatastore() (*Datastore, error) {
	datastore := &Datastore{
		chStop: make(chan struct{}, 1),
	}

	datastore.ctx, datastore.cancel = context.WithCancel(context.Background())

	// default use docker env value
	var mysqlAddr string
	if mysqlAddr = os.Getenv("MYSQL_ADDR"); len(mysqlAddr) == 0 {
		mysqlAddr = global.MysqlAddr
	}

	mysqlDSN := fmt.Sprintf("%s:%s@(%s:%s)/%s", global.MysqlUser, global.MysqlPwd, mysqlAddr, global.MysqlPort, global.MysqlDB)
	var err error
	datastore.db, err = gorm.Open("mysql", mysqlDSN)
	if err != nil {
		logger.Fatal(err)
		return nil, err
	}

	return datastore, nil
}

func (m *Datastore) DB() *gorm.DB {
	return m.db
}

func (m *Datastore) TableGlobal() *iface.TableGlobal {
	return m.global
}

func (m *Datastore) Run() {
	for {
		select {
		case <-m.ctx.Done():
			logger.Info("datastore context done!")
			m.chStop <- struct{}{}
			return
		}
	}

}

func (m *Datastore) Stop() {
	m.db.Close()
	m.cancel()
	<-m.chStop
	close(m.chStop)
}
