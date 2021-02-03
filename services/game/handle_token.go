package game

import (
	"context"
	"errors"
	"fmt"

	"bitbucket.org/east-eden/server/define"
	pbGame "bitbucket.org/east-eden/server/proto/server/game"
	"bitbucket.org/east-eden/server/services/game/player"
	"bitbucket.org/east-eden/server/transport"
)

func (m *MsgHandler) handleAddToken(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGame.C2M_AddToken)
	if !ok {
		return errors.New("handleAddToken failed: recv message body error")
	}

	m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		pl, err := m.g.am.GetPlayerByAccount(acct)
		if err != nil {
			return fmt.Errorf("handleAddToken.AccountExecute failed: %w", err)
		}

		err = pl.TokenManager().TokenInc(msg.Type, msg.Value)
		if err != nil {
			return fmt.Errorf("handleAddToken.AccountExecute failed: %w", err)
		}

		reply := &pbGame.M2C_TokenList{}
		for n := 0; n < define.Token_End; n++ {
			v, err := pl.TokenManager().GetToken(int32(n))
			if err != nil {
				return fmt.Errorf("handleAddToken.AccountExecute failed: %w", err)
			}

			t := &pbGame.Token{
				Type:    v.ID,
				Value:   v.Value,
				MaxHold: v.MaxHold,
			}
			reply.Tokens = append(reply.Tokens, t)
		}
		acct.SendProtoMessage(reply)
		return nil
	})

	return nil
}

func (m *MsgHandler) handleQueryTokens(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		pl, err := m.g.am.GetPlayerByAccount(acct)
		if err != nil {
			return fmt.Errorf("handleQueryTokens.AccountExecute failed: %w", err)
		}

		reply := &pbGame.M2C_TokenList{}
		for n := 0; n < define.Token_End; n++ {
			v, err := pl.TokenManager().GetToken(int32(n))
			if err != nil {
				return fmt.Errorf("handleQueryTokens.AccountExecute failed: %w", err)
			}

			t := &pbGame.Token{
				Type:    v.ID,
				Value:   v.Value,
				MaxHold: v.MaxHold,
			}
			reply.Tokens = append(reply.Tokens, t)
		}
		acct.SendProtoMessage(reply)
		return nil
	})

	return nil
}
