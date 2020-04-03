package game

import (
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/transport"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
)

func (m *MsgHandler) handleAccountTest(sock transport.Socket, p *transport.Message) {
}

func (m *MsgHandler) handleAccountLogon(sock transport.Socket, p *transport.Message) {
	msg, ok := p.Body.(*pbAccount.C2M_AccountLogon)
	if !ok {
		logger.Warn("Cannot assert value to message")
		return
	}

	acct := m.g.am.GetAccountBySock(sock)
	if acct != nil {
		logger.Warn("account had logon:", sock)
		return
	}

	acct, err := m.g.am.AccountLogon(msg.UserId, msg.AccountId, msg.AccountName, sock)
	if err != nil {
		logger.WithFields(logger.Fields{
			"user_id": msg.UserId,
			"id":      msg.AccountId,
			"name":    msg.AccountName,
			"sock":    sock,
		}).Warn("add account failed")
		return
	}

	reply := &pbAccount.M2C_AccountLogon{
		RpcId:       msg.RpcId,
		UserId:      acct.UserID,
		AccountId:   acct.ID,
		PlayerId:    -1,
		PlayerName:  "",
		PlayerLevel: 0,
	}

	pl := m.g.pm.GetPlayerByAccount(acct)
	if pl != nil {
		reply.PlayerId = pl.GetID()
		reply.PlayerName = pl.GetName()
		reply.PlayerLevel = pl.GetLevel()
	}

	acct.SendProtoMessage(reply)
}

func (m *MsgHandler) handleHeartBeat(sock transport.Socket, p *transport.Message) {
	msg, ok := p.Body.(*pbAccount.C2M_HeartBeat)
	if !ok {
		logger.Warn("Cannot assert value to message")
		return
	}

	if acct := m.g.am.GetAccountBySock(sock); acct != nil {
		if t := int32(time.Now().Unix()); t == -1 {
			logger.Warn("Heart beat get time err")
			return
		}

		acct.HeartBeat(msg.RpcId)
	}
}

func (m *MsgHandler) handleAccountConnected(sock transport.Socket, p *transport.Message) {
	if acct := m.g.am.GetAccountBySock(sock); acct != nil {
		accountID := p.Body.(*pbAccount.MC_AccountConnected).AccountId
		logger.WithFields(logger.Fields{
			"account_id": accountID,
		}).Info("account connected")

		// todo after connected
	}
}

func (m *MsgHandler) handleAccountDisconnect(sock transport.Socket, p *transport.Message) {
	m.g.am.DisconnectAccountBySock(sock, "account disconnect initiativly")
}
