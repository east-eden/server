package game

import (
	"context"
	"errors"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/game/player"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"github.com/yokaiio/yokai_server/transport"
)

func (m *MsgHandler) handleAddToken(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGame.C2M_AddToken)
	if !ok {
		return errors.New("handleAddToken failed: recv message body error")
	}

	m.g.am.AccountLaterHandle(sock, func(acct *player.Account) {
		pl := m.g.am.GetPlayerByAccount(acct)
		if pl == nil {
			return
		}

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

	return nil
}

func (m *MsgHandler) handleQueryTokens(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	m.g.am.AccountLaterHandle(sock, func(acct *player.Account) {
		pl := m.g.am.GetPlayerByAccount(acct)
		if pl == nil {
			return
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

	return nil
}
