package game

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/micro/go-micro/transport"
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type ClientManager struct {
	mapClient sync.Map
	mapSocks  sync.Map
	g         *Game
	waitGroup utils.WaitGroupWrapper
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewClientManager(game *Game) *ClientManager {
	cm := &ClientManager{
		g: game,
	}

	cm.ctx, cm.cancel = context.WithCancel(game.ctx)
	cm.g.db.orm.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(ClientPeersInfo{})

	logger.Info("ClientManager Init OK ...")

	return cm
}

func (cm *ClientManager) Main() error {
	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal("ClientManager Main() error:", err)
			}
			exitCh <- err
		})
	}

	cm.waitGroup.Wrap(func() {
		exitFunc(cm.Run())
	})

	return <-exitCh
}

func (cm *ClientManager) Exit() {
	logger.Info("ClientManager context done...")
	cm.cancel()
	cm.waitGroup.Wait()
}

func (cm *ClientManager) AddClient(id int64, name string, sock transport.Socket) (*Client, error) {
	if id == -1 {
		return nil, errors.New("add world id invalid!")
	}

	if client, ok := cm.mapClient.Load(id); ok {
		// adding same client connection
		rc := client.(*Client)
		if rc.peerInfo.sock == sock {
			return rc, nil
		}

		// adding another connection client with existing client_id
		cm.DisconnectClient(rc.peerInfo.sock, "AddClient")
	}

	var numSocks uint32
	cm.mapSocks.Range(func(_, _ interface{}) bool {
		numSocks++
		return true
	})

	if numSocks >= uint32(cm.g.opts.ClientConnectMax) {
		return nil, errors.New("reach game server's max client connect num!")
	}

	// new client
	peerInfo := &ClientPeersInfo{
		ID:   id,
		Name: name,
		sock: sock,
		p:    cm.g.pm.NewPlayer(id, name),
	}

	client := NewClient(cm, peerInfo)
	cm.mapClient.Store(client.GetID(), client)
	cm.mapSocks.Store(sock, client)
	logger.Info(fmt.Sprintf("add client <id:%d, name:%s, sock:%v> success!", client.GetID(), client.GetName(), client.GetSock()))

	// client main
	cm.waitGroup.Wrap(func() {
		err := client.Main()
		if err != nil {
			logger.Info("client Main() return err:", err)
		}
		client.Exit()
		cm.mapSocks.Delete(client.peerInfo.sock)

		// maybe a new client connected with the same clientID
		if c, ok := cm.mapClient.Load(client.GetID()); ok {
			if c.(*Client).peerInfo.sock == client.peerInfo.sock {
				cm.mapClient.Delete(client.GetID())
			}
		}
	})

	return client, nil
}

func (cm *ClientManager) GetClientByID(id int64) *Client {
	v, ok := cm.mapClient.Load(id)
	if !ok {
		return nil
	}

	return v.(*Client)
}

func (cm *ClientManager) GetClientBySock(sock transport.Socket) *Client {
	v, ok := cm.mapSocks.Load(sock)
	if !ok {
		return nil
	}

	return v.(*Client)
}

func (cm *ClientManager) GetAllClients() []*Client {
	ret := make([]*Client, 0)
	cm.mapClient.Range(func(k, v interface{}) bool {
		c := v.(*Client)
		ret = append(ret, c)
		return true
	})

	return ret
}

func (cm *ClientManager) DisconnectClient(sock transport.Socket, reason string) {
	v, ok := cm.mapSocks.Load(sock)
	if !ok {
		return
	}

	client, ok := v.(*Client)
	if !ok {
		return
	}

	logger.WithFields(logger.Fields{
		"id":     client.GetID(),
		"reason": reason,
	}).Warn("Client disconnected!")

	client.cancel()
}

func (cm *ClientManager) BroadCast(msg proto.Message) {
	cm.mapClient.Range(func(_, v interface{}) bool {
		if client, ok := v.(*Client); ok {
			client.SendProtoMessage(msg)
		}
		return true
	})
}

func (cm *ClientManager) Run() error {
	for {
		select {
		case <-cm.ctx.Done():
			logger.Print("world session context done!")
			return nil
		}
	}

	return nil
}
