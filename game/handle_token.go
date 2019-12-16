package game

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/transport"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

func (m *MsgHandler) handleAddToken(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("add token failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_AddToken)
	if !ok {
		logger.Warn("Add Token failed, recv message body error")
		return
	}

	cli.PushWrapHandler(func() {
		err := cli.Player().TokenManager().TokenInc(msg.Type, msg.Value)
		if err != nil {
			logger.Warn("token inc failed:", err)
		}

		cli.Player().TokenManager().Save()

		reply := &pbGame.MS_TokenList{Tokens: make([]*pbGame.Token, 0, define.Token_End)}
		for n := 0; n < define.Token_End; n++ {
			v, err := cli.Player().TokenManager().GetToken(int32(n))
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
		cli.SendProtoMessage(reply)
	})

}

func (m *MsgHandler) handleQueryTokens(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("query tokens failed")
		return
	}

	reply := &pbGame.MS_TokenList{Tokens: make([]*pbGame.Token, 0, define.Token_End)}
	for n := 0; n < define.Token_End; n++ {
		v, err := cli.Player().TokenManager().GetToken(int32(n))
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
	cli.SendProtoMessage(reply)
}
