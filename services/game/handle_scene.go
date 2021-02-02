package game

import (
	"context"
	"errors"
	"fmt"

	pbGame "bitbucket.org/east-eden/server/proto/game"
	"bitbucket.org/east-eden/server/services/game/player"
	"bitbucket.org/east-eden/server/transport"
)

func (m *MsgHandler) handleStartStageCombat(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGame.C2M_StartStageCombat)
	if !ok {
		return errors.New("handleStartStageCombat failed: recv message body error")
	}

	m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		pl, err := m.g.am.GetPlayerByAccount(acct)
		if err != nil {
			return fmt.Errorf("handleStartStageCombat.AccountExecute failed: %w", err)
		}

		reply := &pbGame.M2C_StartStageCombat{
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
	})

	return nil
}
