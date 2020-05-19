package game

import (
	logger "github.com/sirupsen/logrus"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"github.com/yokaiio/yokai_server/transport"
)

func (m *MsgHandler) handleStartStageCombat(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("handleStartStageCombat failed")
		return
	}

	pl := m.g.pm.GetPlayerByAccount(acct)
	if pl == nil {
		return
	}

	msg, ok := p.Body.(*pbGame.C2M_StartStageCombat)
	if !ok {
		logger.Warn("handleStartStageCombat failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		reply := &pbGame.M2C_StartStageCombat{
			RpcId: msg.RpcId,
		}

		resp, err := m.g.rpcHandler.CallStartStageCombat(pl)
		if err != nil {
			reply.Error = 1
			reply.Message = err.Error()
		}

		reply.Result = resp.Result
		pl.SendProtoMessage(reply)
	})
}
