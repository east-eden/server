package game

import (
	"context"
	"errors"
	"fmt"

	"bitbucket.org/east-eden/server/define"
	pbGlobal "bitbucket.org/east-eden/server/proto/global"
	"bitbucket.org/east-eden/server/services/game/player"
	"bitbucket.org/east-eden/server/transport"
)

func (m *MsgHandler) handleAddToken(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_AddToken)
	if !ok {
		return errors.New("handleAddToken failed: recv message body error")
	}
	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleAddToken.AccountExecute failed: %w", err)
	}

	err = pl.TokenManager().TokenInc(msg.Type, msg.Value)
	if err != nil {
		return fmt.Errorf("handleAddToken.AccountExecute failed: %w", err)
	}

	reply := &pbGlobal.S2C_TokenList{}
	for n := 0; n < define.Token_End; n++ {
		v, err := pl.TokenManager().GetToken(int32(n))
		if err != nil {
			return fmt.Errorf("handleAddToken.AccountExecute failed: %w", err)
		}

		reply.Tokens = append(reply.Tokens, v)
	}
	acct.SendProtoMessage(reply)
	return nil
}

func (m *MsgHandler) handleQueryTokens(ctx context.Context, acct *player.Account, p *transport.Message) error {
	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleQueryTokens.AccountExecute failed: %w", err)
	}

	reply := &pbGlobal.S2C_TokenList{}
	for n := 0; n < define.Token_End; n++ {
		v, err := pl.TokenManager().GetToken(int32(n))
		if err != nil {
			return fmt.Errorf("handleQueryTokens.AccountExecute failed: %w", err)
		}

		reply.Tokens = append(reply.Tokens, v)
	}
	acct.SendProtoMessage(reply)
	return nil
}
