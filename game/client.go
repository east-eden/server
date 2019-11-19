package game

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/player"
	"github.com/yokaiio/yokai_server/internal/transport"
	"github.com/yokaiio/yokai_server/internal/utils"
	pbClient "github.com/yokaiio/yokai_server/proto/client"
)

type ClientPeersInfo struct {
	ID   int64  `gorm:"type:bigint(20);primary_key;column:id;default:0;not null"`
	Name string `gorm:"type:varchar(32);column:name;default:'';not null"`
	sock transport.Socket
	p    player.Player
}

type Client struct {
	peerInfo *ClientPeersInfo

	cm             *ClientManager
	ctx            context.Context
	cancel         context.CancelFunc
	waitGroup      utils.WaitGroupWrapper
	chw            chan uint32
	heartBeatTimer *time.Timer
}

func NewClient(cm *ClientManager, peerInfo *ClientPeersInfo) *Client {
	client := &Client{
		cm:             cm,
		peerInfo:       peerInfo,
		heartBeatTimer: time.NewTimer(cm.g.opts.HeartBeat),
	}

	client.ctx, client.cancel = context.WithCancel(cm.ctx)

	return client
}

func (Client) TableName() string {
	return "Client"
}

func (c *Client) GetID() int64 {
	return c.peerInfo.ID
}

func (c *Client) GetName() string {
	return c.peerInfo.Name
}

func (c *Client) GetSock() transport.Socket {
	return c.peerInfo.sock
}

func (c *Client) Main() error {
	c.loadFromDB()
	c.saveToDB()

	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal("Client Main() error:", err)
			}
			exitCh <- err
		})
	}

	c.waitGroup.Wrap(func() {
		exitFunc(c.Run())
	})

	return <-exitCh
}

func (c *Client) loadFromDB() {
	c.cm.g.db.orm.First(c.peerInfo)
}

func (c *Client) saveToDB() {
	c.cm.g.db.orm.Save(c.peerInfo)
}

func (c *Client) Exit() {
	c.heartBeatTimer.Stop()
	c.peerInfo.sock.Close()
}

func (c *Client) Run() error {
	for {
		select {
		// context canceled
		case <-c.ctx.Done():
			logger.WithFields(logger.Fields{
				"id": c.GetID(),
			}).Info("Client context done!")
			return nil

		// lost connection
		case <-c.heartBeatTimer.C:
			c.cm.DisconnectClient(c.peerInfo.sock, "timeout")
		}
	}
}

/*
msg Example:
	Type: transport.BodyProtobuf
	Name: yokai_client.MS_ClientLogon
	Body: protoBuf byte
*/
func (c *Client) SendProtoMessage(p proto.Message) {
	var msg transport.Message
	msg.Type = transport.BodyProtobuf
	msg.Name = proto.MessageName(p)
	msg.Body = p

	if err := c.peerInfo.sock.Send(&msg); err != nil {
		logger.Warn("send proto msg error:", err)
		return
	}
}

func (c *Client) HeartBeat() {
	c.heartBeatTimer.Reset(c.cm.g.opts.HeartBeat)

	reply := &pbClient.MS_HeartBeat{Timestamp: uint32(time.Now().Unix())}
	c.SendProtoMessage(reply)
}
