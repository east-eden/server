package game

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	logger "github.com/sirupsen/logrus"
)

type ClientMgr struct {
	mapClient      sync.Map
	mapConn        sync.Map
	g              *Game
	wg             sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
	chKickClientID chan uint32
}

func NewClientMgr(game *Game) *ClientMgr {
	cm := &ClientMgr{
		g:              game,
		chKickClientID: make(chan uint32, game.opts.ClientConnectMax),
	}

	cm.ctx, cm.cancel = context.WithCancel(game.ctx.Background())
	cm.g.db.orm.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(Client{})

	return cm, nil
}

func (cm *ClientMgr) Stop() {
	cm.mapClient.Range(func(_, v interface{}) bool {
		v.(*Client).Stop()
		return true
	})

	cm.cancel()
	close(cm.chStop)
	close(cm.chKickClientID)
}

func (cm *ClientMgr) AddClient(id uint32, name string, c *TCPConn) (*Client, error) {
	if int32(id) == -1 {
		return nil, errors.New("add world id invalid!")
	}

	if _, ok := cm.mapClient.Load(id); ok {
		cm.KickClient(id, "AddClient")
	}

	var numConn uint32
	cm.mapConn.Range(func(_, _ interface{}) bool {
		numConn++
		return true
	})

	if numConn >= uint32(cm.g.opts.ClientConnectMax) {
		return nil, errors.New("reach game server's max client connect num!")
	}

	// new client
	peerInfo := &ClientPeersInfo{
		ID:        id,
		Name:      name,
		c:         c,
		chTimeout: cm.chKickClientID,
	}
	w := NewClient(peerInfo, cm.chKickClientID)
	cm.mapClient.Store(w.GetID(), w)
	cm.mapConn.Store(w.GetCon(), w)
	logger.Info(fmt.Sprintf("add world <id:%d, name:%s, con:%v> success!", w.GetID(), w.GetName(), w.GetCon()))

	// world run
	go w.Run()

	w.SetLastConTime(int(time.Now().Unix()))
	cm.ds.DB().Save(w)

	return w, nil
}

func (cm *ClientMgr) GetClientByID(id uint32) iface.IClient {
	worldID := cm.getClientRefID(id)
	v, ok := cm.mapClient.Load(worldID)
	if !ok {
		return nil
	}

	return v.(*world)
}

func (cm *ClientMgr) GetClientByCon(con iface.ITCPConn) iface.IClient {
	v, ok := cm.mapConn.Load(con)
	if !ok {
		return nil
	}

	return v.(*world)
}

func (cm *ClientMgr) DisconnectClient(con iface.ITCPConn) {
	v, ok := cm.mapConn.Load(con)
	if !ok {
		return
	}

	world, ok := v.(*world)
	if !ok {
		return
	}

	logger.WithFields(logger.Fields{
		"id": world.GetID(),
	}).Warn("Client disconnected!")
	world.Stop()

	cm.mapClient.Delete(world.GetID())
	cm.mapConn.Delete(con)
}

func (cm *ClientMgr) KickClient(id uint32, reason string) {
	v, ok := cm.mapClient.Load(id)
	if !ok {
		return
	}

	world, ok := v.(*world)
	if !ok {
		return
	}

	logger.WithFields(logger.Fields{
		"id":     world.GetID(),
		"reason": reason,
	}).Warn("Client was kicked!")

	world.Stop()
	cm.mapConn.Delete(world.GetCon())
	cm.mapClient.Delete(world.GetID())
}

func (cm *ClientMgr) BroadCast(msg proto.Message) {
	cm.mapClient.Range(func(_, v interface{}) bool {
		if world, ok := v.(*world); ok {
			world.SendProtoMessage(msg)
		}
		return true
	})
}

func (cm *ClientMgr) Run() {
	for {
		select {
		case <-cm.ctx.Done():
			logger.Print("world session context done!")
			cm.chStop <- struct{}{}
			return
		case wid := <-cm.chKickClientID:
			cm.KickClient(wid, "time out")
		}
	}
}
