package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/player"
)

func (m *MsgRegister) handleStartStageCombat(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_StartStageCombat)
	if !ok {
		return errors.New("handleStartStageCombat failed: recv message body error")
	}

	pl, err := m.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleStartStageCombat.AccountExecute failed: %w", err)
	}

	reply := &pbGlobal.S2C_StartStageCombat{
		RpcId: msg.RpcId,
	}

	resp, err := m.rpcHandler.CallStartStageCombat(pl)
	if err != nil {
		reply.Error = 1
		reply.Message = err.Error()
	}

	if resp != nil {
		reply.Result = resp.Result
	}

	pl.SendProtoMessage(reply)
	return nil
}
