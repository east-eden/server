package game

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/game/player"
	"github.com/yokaiio/yokai_server/internal/transport"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type ClientManager struct {
	mapClient sync.Map
	mapSocks  sync.Map
	g         *Game
	waitGroup utils.WaitGroupWrapper
	ctx       context.Context
	cancel    context.CancelFunc

	clientConnectMax int
	clientTimeout    time.Duration
}

func NewClientManager(game *Game, ctx *cli.Context) *ClientManager {
	cm := &ClientManager{
		g:                game,
		clientConnectMax: ctx.Int("client_connect_max"),
		clientTimeout:    ctx.Duration("client_timeout"),
	}

	cm.ctx, cm.cancel = context.WithCancel(ctx)
	cm.g.ds.ORM().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(ClientInfo{})

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

func (cm *ClientManager) addClient(id int64, name string, sock transport.Socket) (*Client, error) {
	if id == -1 {
		return nil, errors.New("add world id invalid!")
	}

	var numSocks uint32
	cm.mapSocks.Range(func(_, _ interface{}) bool {
		numSocks++
		return true
	})

	if numSocks >= uint32(cm.clientConnectMax) {
		return nil, errors.New("reach game server's max client connect num!")
	}

	// new client
	info := &ClientInfo{
		ID:   id,
		Name: name,
		sock: sock,
	}

	client := NewClient(cm, info)
	cm.mapClient.Store(client.ID(), client)
	cm.mapSocks.Store(sock, client)
	logger.Info(fmt.Sprintf("add client <id:%d, name:%s, sock:%v> success!", client.ID(), client.Name(), client.Sock()))

	// client main
	cm.waitGroup.Wrap(func() {
		err := client.Main()
		if err != nil {
			logger.Info("client Main() return err:", err)
		}
		cm.mapSocks.Delete(client.info.sock)
		cm.mapClient.Delete(client.ID())
		client.Exit()

	})

	return client, nil
}

func (cm *ClientManager) ClientLogon(id int64, name string, sock transport.Socket) (*Client, error) {
	client, ok := cm.mapClient.Load(id)
	if ok {
		// return exist connection sock
		rc := client.(*Client)
		if rc.info.sock == sock {
			return rc, nil
		}

		// disconnect last client sock
		cm.DisconnectClient(rc.info.sock, "AddClient")
	}

	return cm.addClient(id, name, sock)
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
		"id":     client.ID(),
		"reason": reason,
	}).Warn("Client disconnected!")

	client.cancel()
}

func (cm *ClientManager) CreatePlayer(c *Client, name string) (player.Player, error) {
	m := cm.g.pm.GetPlayersByClientID(c.ID())
	for _, v := range m {
		if v.GetName() == name {
			return nil, fmt.Errorf("existed player name")
		}
	}

	p, err := cm.g.pm.CreatePlayer(c.ID(), name)
	if p != nil {
		c.info.p = p
	}

	return p, err
}

func (cm *ClientManager) SelectPlayer(c *Client, id int64) (player.Player, error) {
	playerList := cm.g.pm.GetPlayersByClientID(c.ID())
	for _, v := range playerList {
		if v.GetID() == id {
			return v, nil
		}
	}

	return nil, fmt.Errorf("select player with wrong id:", id)
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
