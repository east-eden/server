package game

import (
	"context"
	"time"

	logger "github.com/sirupsen/logrus"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
	"github.com/yokaiio/yokai_server/transport"
)

func (m *MsgHandler) handleAccountTest(ctx context.Context, sock transport.Socket, p *transport.Message) {
}

func (m *MsgHandler) handleAccountLogon(ctx context.Context, sock transport.Socket, p *transport.Message) {
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

	acct, err := m.g.am.AccountLogon(ctx, msg.UserId, msg.AccountId, msg.AccountName, sock)
	if err != nil {
		logger.WithFields(logger.Fields{
			"user_id": msg.UserId,
			"id":      msg.AccountId,
			"name":    msg.AccountName,
			"sock":    sock,
		}).Warn("add account failed")
		return
	}

	acct.PushWrapHandler(func() {
		reply := &pbAccount.M2C_AccountLogon{
			RpcId:       msg.RpcId,
			UserId:      acct.UserId,
			AccountId:   acct.ID,
			PlayerId:    -1,
			PlayerName:  "",
			PlayerLevel: 0,
		}

		pl := m.g.am.GetPlayerByAccount(acct)
		if pl != nil {
			reply.PlayerId = pl.GetID()
			reply.PlayerName = pl.GetName()
			reply.PlayerLevel = pl.GetLevel()
		}

		acct.SendProtoMessage(reply)
	})
}

func (m *MsgHandler) handleHeartBeat(ctx context.Context, sock transport.Socket, p *transport.Message) {
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

		acct.PushAsyncHandler(func() {
			acct.HeartBeat(msg.RpcId)
		})
	}
}

func (m *MsgHandler) handleAccountConnected(ctx context.Context, sock transport.Socket, p *transport.Message) {
	if acct := m.g.am.GetAccountBySock(sock); acct != nil {
		accountID := p.Body.(*pbAccount.MC_AccountConnected).AccountId
		logger.WithFields(logger.Fields{
			"account_id": accountID,
		}).Info("account connected")

		// todo after connected
	}
}

func (m *MsgHandler) handleAccountDisconnect(ctx context.Context, sock transport.Socket, p *transport.Message) {
	m.g.am.DisconnectAccountBySock(sock, "account disconnect initiativly")
}
