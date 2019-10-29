package game

import (
	"bytes"
	"context"
	"encoding/binary"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/hellodudu/Ultimate/iface"
	"github.com/hellodudu/Ultimate/utils"
	"github.com/hellodudu/yokai_server/game/define"
	logger "github.com/sirupsen/logrus"
)

type ClientPeersInfo struct {
	ID   uint32 `gorm:"type:int(10);primary_key;column:id;default:0;not null"`
	Name string `gorm:"type:varchar(32);column:name;default:'';not null"`
	c    *TcpCon
	cm   *ClientMgr
}

type Client struct {
	peerInfo *ClientPeersInfo

	tHeartBeat *time.Timer // connection heart beat
	tTimeOut   *time.Timer // connection time out
	ctx        context.Context
	cancel     context.CancelFunc
	chw        chan uint32

	mapPlayer map[int64]*pbGame.CrossPlayerInfo
	mapGuild  map[int64]*pbGame.CrossGuildInfo

	chDBInit chan struct{}
	chStop   chan struct{}
}

func NewClient(peerInfo *ClientPeersInfo) *Client {
	client := &Client{
		peerInfo:   peerInfo,
		tHeartBeat: time.NewTimer(time.Duration(global.ClientHeartBeatSec) * time.Second),
		tTimeOut:   time.NewTimer(time.Duration(global.ClientConTimeOutSec) * time.Second),
		mapPlayer:  make(map[int64]*pbGame.CrossPlayerInfo),
		mapGuild:   make(map[int64]*pbGame.CrossGuildInfo),
		chDBInit:   make(chan struct{}, 1),
		chStop:     make(chan struct{}, 1),
	}

	client.ctx, client.cancel = context.WithCancel(peerInfo.cm.ctx)

	client.loadFromDB()
	return client
}

func (Client) TableName() string {
	return "Client"
}

func (w *Client) GetID() uint32 {
	return w.ID
}

func (w *Client) GetName() string {
	return w.Name
}

func (w *Client) GetCon() iface.ITCPConn {
	return w.con
}

func (w *Client) loadFromDB() {
	w.ds.DB().First(w)

	w.chDBInit <- struct{}{}
}

func (w *Client) Stop() {
	w.tHeartBeat.Stop()
	w.tTimeOut.Stop()
	w.cancel()
	<-w.chStop
	close(w.chStop)
	close(w.chDBInit)
	w.con.Close()
}

func (client *Client) Run() {
	<-client.chDBInit

	for {
		select {
		// context canceled
		case <-client.ctx.Done():
			logger.WithFields(logger.Fields{
				"id": w.GetID(),
			}).Info("Client context done!")
			w.chStop <- struct{}{}
			return

		// connecting timeout
		case <-w.tTimeOut.C:
			w.chw <- w.ID

		// Heart Beat
		case <-w.tHeartBeat.C:
			msg := &pbClient.MUW_TestConnect{}
			w.SendProtoMessage(msg)
			w.tHeartBeat.Reset(time.Duration(global.ClientHeartBeatSec) * time.Second)
		}
	}
}

func (w *Client) ResetTestConnect() {
	w.tHeartBeat.Reset(time.Duration(global.ClientHeartBeatSec) * time.Second)
	w.tTimeOut.Reset(time.Duration(global.ClientConTimeOutSec) * time.Second)
}

func (w *Client) SendProtoMessage(p proto.Message) {
	// reply message length = 4bytes size + 8bytes size BaseNetMsg + 2bytes message_name size + message_name + proto_data
	out, err := proto.Marshal(p)
	if err != nil {
		logger.Warn(err)
		return
	}

	typeName := proto.MessageName(p)
	baseMsg := &define.BaseNetMsg{}
	msgSize := binary.Size(baseMsg) + 2 + len(typeName) + len(out)
	baseMsg.ID = utils.Crc32("MUW_DirectProtoMsg")
	baseMsg.Size = uint32(msgSize)

	var resp []byte = make([]byte, 4+msgSize)
	binary.LittleEndian.PutUint32(resp[:4], uint32(msgSize))
	binary.LittleEndian.PutUint32(resp[4:8], baseMsg.ID)
	binary.LittleEndian.PutUint32(resp[8:12], baseMsg.Size)
	binary.LittleEndian.PutUint16(resp[12:12+2], uint16(len(typeName)))
	copy(resp[14:14+len(typeName)], []byte(typeName))
	copy(resp[14+len(typeName):], out)

	if _, err := w.con.Write(resp); err != nil {
		logger.Warn("send proto msg error:", err)
		return
	}
}

func (w *Client) SendTransferMessage(data []byte) {
	resp := make([]byte, 4+len(data))
	binary.LittleEndian.PutUint32(resp[:4], uint32(len(data)))
	copy(resp[4:], data)

	if _, err := w.con.Write(resp); err != nil {
		logger.Warn("send transfer msg error:", err)
		return
	}

	// for testing disconnected from Client server
	transferMsg := &define.TransferNetMsg{}
	byTransferMsg := make([]byte, binary.Size(transferMsg))

	copy(byTransferMsg, data[:binary.Size(transferMsg)])
	buf := &bytes.Buffer{}
	if _, err := buf.Write(byTransferMsg); err != nil {
		return
	}

	// get top 4 bytes messageid
	if err := binary.Read(buf, binary.LittleEndian, transferMsg); err != nil {
		return
	}
}
