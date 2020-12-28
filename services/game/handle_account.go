package game

import (
	"context"
	"errors"
	"fmt"

	pbAccount "github.com/east-eden/server/proto/account"
	"github.com/east-eden/server/services/game/player"
	"github.com/east-eden/server/transport"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/rs/zerolog/log"
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

	return m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		reply := &pbAccount.M2C_AccountLogon{
			RpcId:       msg.RpcId,
			UserId:      acct.UserId,
			AccountId:   acct.ID,
			PlayerId:    -1,
			PlayerName:  "",
			PlayerLevel: 0,
		}

		if pl, err := m.g.am.GetPlayerByAccount(acct); err == nil {
			reply.PlayerId = pl.GetID()
			reply.PlayerName = pl.GetName()
			reply.PlayerLevel = pl.GetLevel()
		}
		acct.SendProtoMessage(reply)
		return nil
	})
}

func (m *MsgHandler) handleHeartBeat(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		m.timeHistogram.WithLabelValues("handleHeartBeat").Observe(v)
	}))

	msg, ok := p.Body.(*pbAccount.C2M_HeartBeat)
	if !ok {
		return errors.New("handleHeartBeat failed: cannot assert value to message")
	}

	return m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		defer timer.ObserveDuration()
		acct.HeartBeat(msg.RpcId)
		return nil
	})
}

// todo after account logon
func (m *MsgHandler) handleAccountConnected(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	if acct := m.g.am.GetAccountBySock(sock); acct != nil {
		accountId := p.Body.(*pbAccount.MC_AccountConnected).AccountId
		log.Info().
			Int64("account_id", accountId).
			Msg("account connected")

		// todo after connected
	}

	return nil
}

// client disconnect
func (m *MsgHandler) handleAccountDisconnect(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	return player.ErrAccountDisconnect
}
