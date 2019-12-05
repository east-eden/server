package game

import (
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/transport"
	pbClient "github.com/yokaiio/yokai_server/proto/client"
)

func (m *MsgHandler) handleClientTest(sock transport.Socket, p *transport.Message) {
	logger.WithFields(logger.Fields{
		"type": p.Type,
		"name": p.Name,
		"body": p.Body,
	}).Info("recv client test")
}

func (m *MsgHandler) handleClientLogon(sock transport.Socket, p *transport.Message) {
	msg, ok := p.Body.(*pbClient.MC_ClientLogon)
	if !ok {
		logger.Warn("Cannot assert value to message")
		return
	}

	client := m.g.cm.GetClientBySock(sock)
	if client != nil {
		logger.Warn("client had logon:", sock)
		return
	}

	client, err := m.g.cm.ClientLogon(msg.ClientId, msg.ClientName, sock)
	if err != nil {
		logger.WithFields(logger.Fields{
			"id":   msg.ClientId,
			"name": msg.ClientName,
			"sock": sock,
		}).Warn("add client failed")
		return
	}

	reply := &pbClient.MS_ClientLogon{}
	client.SendProtoMessage(reply)
}

func (m *MsgHandler) handleHeartBeat(sock transport.Socket, p *transport.Message) {
	if client := m.g.cm.GetClientBySock(sock); client != nil {
		if t := int32(time.Now().Unix()); t == -1 {
			logger.Warn("Heart beat get time err")
			return
		}

		client.HeartBeat()
	}
}

func (m *MsgHandler) handleClientConnected(sock transport.Socket, p *transport.Message) {
	if client := m.g.cm.GetClientBySock(sock); client != nil {
		clientID := p.Body.(*pbClient.MC_ClientConnected).ClientId
		logger.WithFields(logger.Fields{
			"client_id": clientID,
		}).Info("client connected")

		// todo after connected
	}
}

func (m *MsgHandler) handleClientDisconnect(sock transport.Socket, p *transport.Message) {
	m.g.cm.DisconnectClient(sock, "client disconnect initiativly")
}
