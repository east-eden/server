package game

import (
	"context"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/define"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"github.com/yokaiio/yokai_server/transport"
)

func (m *MsgHandler) handleAddToken(ctx context.Context, sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("add token failed")
		return
	}

	pl := m.g.am.GetPlayerByAccount(acct)
	if pl == nil {
		return
	}

	msg, ok := p.Body.(*pbGame.C2M_AddToken)
	if !ok {
		logger.Warn("Add Token failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		err := pl.TokenManager().TokenInc(msg.Type, msg.Value)
		if err != nil {
			logger.Warn("token inc failed:", err)
		}

		reply := &pbGame.M2C_TokenList{}
		for n := 0; n < define.Token_End; n++ {
			v, err := pl.TokenManager().GetToken(int32(n))
			if err != nil {
				logger.Warn("token get value failed:", err)
				return
			}

			t := &pbGame.Token{
				Type:    v.ID,
				Value:   v.Value,
				MaxHold: v.MaxHold,
			}
			reply.Tokens = append(reply.Tokens, t)
		}
		acct.SendProtoMessage(reply)
	})

}

func (m *MsgHandler) handleQueryTokens(ctx context.Context, sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("query tokens failed")
		return
	}

	pl := m.g.am.GetPlayerByAccount(acct)
	if pl == nil {
		return
	}

	acct.PushWrapHandler(func() {
		reply := &pbGame.M2C_TokenList{}
		for n := 0; n < define.Token_End; n++ {
			v, err := pl.TokenManager().GetToken(int32(n))
			if err != nil {
				logger.Warn("token get value failed:", err)
				return
			}

			t := &pbGame.Token{
				Type:    v.ID,
				Value:   v.Value,
				MaxHold: v.MaxHold,
			}
			reply.Tokens = append(reply.Tokens, t)
		}
		acct.SendProtoMessage(reply)
	})
}
