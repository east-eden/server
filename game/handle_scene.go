package game

import (
	"context"
	"errors"

	"github.com/yokaiio/yokai_server/game/player"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"github.com/yokaiio/yokai_server/transport"
)

func (m *MsgHandler) handleStartStageCombat(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGame.C2M_StartStageCombat)
	if !ok {
		return errors.New("handleStartStageCombat failed: recv message body error")
	}

	m.g.am.AccountLaterHandle(sock, func(acct *player.Account) {
		pl := m.g.am.GetPlayerByAccount(acct)
		if pl == nil {
			return
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
	})

	return nil
}
