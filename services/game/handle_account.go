package game

import (
	"context"
	"errors"
	"fmt"

	pbAccount "bitbucket.org/east-eden/server/proto/account"
	"bitbucket.org/east-eden/server/services/game/player"
	"bitbucket.org/east-eden/server/transport"
	"github.com/prometheus/client_golang/prometheus"
)

func (m *MsgHandler) handleAccountTest(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	return nil
}

func (m *MsgHandler) handleAccountLogon(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbAccount.C2M_AccountLogon)
	if !ok {
		return errors.New("handleAccountLogon failed: cannot assert value to message")
	}

	err := m.g.am.AccountLogon(ctx, msg.UserId, msg.AccountId, msg.AccountName, sock)
	if err != nil {
		return fmt.Errorf("handleAccountLogon failed: %w", err)
	}

	m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		reply := &pbAccount.M2C_AccountLogon{
			RpcId:       msg.RpcId,
			UserId:      acct.UserId,
			AccountId:   acct.ID,
			PlayerId:    -1,
			PlayerName:  "",
			PlayerLevel: 0,
		}

		if p := acct.GetPlayer(); p != nil {
			reply.PlayerId = p.GetID()
			reply.PlayerName = p.GetName()
			reply.PlayerLevel = p.GetLevel()
		}

		acct.SendProtoMessage(reply)
		return nil
	})

	return nil
}

func (m *MsgHandler) handleHeartBeat(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		m.timeHistogram.WithLabelValues("handleHeartBeat").Observe(v)
	}))

	msg, ok := p.Body.(*pbAccount.C2M_HeartBeat)
	if !ok {
		return errors.New("handleHeartBeat failed: cannot assert value to message")
	}

	m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		defer timer.ObserveDuration()
		acct.HeartBeat(msg.RpcId)
		return nil
	})

	return nil
}

// client disconnect
func (m *MsgHandler) handleAccountDisconnect(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	return player.ErrAccountDisconnect
}
