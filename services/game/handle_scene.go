package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "bitbucket.org/east-eden/server/proto/global"
	"bitbucket.org/east-eden/server/services/game/player"
	"bitbucket.org/east-eden/server/transport"
)

func (m *MsgHandler) handleStartStageCombat(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_StartStageCombat)
	if !ok {
		return errors.New("handleStartStageCombat failed: recv message body error")
	}

	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleStartStageCombat.AccountExecute failed: %w", err)
	}

	reply := &pbGlobal.S2C_StartStageCombat{
		RpcId: msg.RpcId,
	}

	resp, err := m.g.rpcHandler.CallStartStageCombat(pl)
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
