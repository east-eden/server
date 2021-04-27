package game

import (
	"context"
	"fmt"

	"bitbucket.org/funplus/server/define"
	pbCommon "bitbucket.org/funplus/server/proto/global/common"
	"bitbucket.org/funplus/server/services/game/player"
)

func (m *MsgRegister) handleQueryTokens(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	pl, err := m.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleQueryTokens.AccountExecute failed: %w", err)
	}

	reply := &pbCommon.S2C_TokenList{}
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
