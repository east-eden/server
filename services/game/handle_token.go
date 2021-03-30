package game

import (
	"context"
	"fmt"

	"github.com/east-eden/server/define"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/player"
	"github.com/east-eden/server/transport"
)

func (m *MsgRegister) handleQueryTokens(ctx context.Context, acct *player.Account, p *transport.Message) error {
	pl, err := m.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleQueryTokens.AccountExecute failed: %w", err)
	}

	reply := &pbGlobal.S2C_TokenList{}
	for n := define.Token_Begin; n < define.Token_End; n++ {
		v, err := pl.TokenManager().GetToken(int32(n))
		if err != nil {
			return fmt.Errorf("handleQueryTokens.AccountExecute failed: %w", err)
		}

		reply.Tokens = append(reply.Tokens, v)
	}
	acct.SendProtoMessage(reply)
	return nil
}
