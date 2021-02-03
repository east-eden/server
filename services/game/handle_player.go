package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "bitbucket.org/east-eden/server/proto/global"
	"bitbucket.org/east-eden/server/services/game/player"
	"bitbucket.org/east-eden/server/transport"
	"github.com/prometheus/client_golang/prometheus"
)

func (m *MsgHandler) handleQueryPlayerInfo(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		m.timeHistogram.WithLabelValues("handleQueryPlayerInfo").Observe(v)
	}))

	m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		defer timer.ObserveDuration()

		reply := &pbGlobal.M2C_QueryPlayerInfo{Error: 0}
		if pl, err := m.g.am.GetPlayerByAccount(acct); err == nil {
			reply.Info = &pbGlobal.PlayerInfo{
				LiteInfo: &pbGlobal.LitePlayer{
					Id:        pl.GetID(),
					AccountId: pl.GetAccountID(),
					Name:      pl.GetName(),
					Exp:       pl.GetExp(),
					Level:     pl.GetLevel(),
				},

				HeroNums: int32(pl.HeroManager().GetHeroNums()),
				ItemNums: int32(pl.ItemManager().GetItemNums()),
			}
		}

		acct.SendProtoMessage(reply)
		return nil
	})

	return nil
}

func (m *MsgHandler) handleCreatePlayer(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2M_CreatePlayer)
	if !ok {
		return errors.New("handleCreatePlayer failed: recv message body error")
	}

	m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		reply := &pbGlobal.M2C_CreatePlayer{
			Error: 0,
		}

		pl, err := m.g.am.CreatePlayer(acct, msg.Name)
		if err != nil {
			reply.Error = -1
			reply.Message = err.Error()
			acct.SendProtoMessage(reply)
			return fmt.Errorf("handleCreatePlayer.AccountExecute failed: %w", err)
		}

		reply.Info = &pbGlobal.PlayerInfo{
			LiteInfo: &pbGlobal.LitePlayer{
				Id:        pl.GetID(),
				AccountId: pl.GetAccountID(),
				Name:      pl.GetName(),
				Exp:       pl.GetExp(),
				Level:     pl.GetLevel(),
			},
			HeroNums: int32(pl.HeroManager().GetHeroNums()),
			ItemNums: int32(pl.ItemManager().GetItemNums()),
		}

		acct.SendProtoMessage(reply)
		return nil
	})

	return nil
}

func (m *MsgHandler) handleSelectPlayer(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.MC_SelectPlayer)
	if !ok {
		return errors.New("handleSelectPlayer failed: recv message body error")
	}

	m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		pl, err := m.g.am.SelectPlayer(acct, msg.Id)
		reply := &pbGlobal.MS_SelectPlayer{
			ErrorCode: 0,
		}

		if err != nil {
			reply.ErrorCode = -1
			acct.SendProtoMessage(reply)
			return fmt.Errorf("handleSelectPlayer.AccountExecute failed: %w", err)
		}

		reply.Info = &pbGlobal.PlayerInfo{
			LiteInfo: &pbGlobal.LitePlayer{
				Id:        pl.GetID(),
				AccountId: pl.GetAccountID(),
				Name:      pl.GetName(),
				Exp:       pl.GetExp(),
				Level:     pl.GetLevel(),
			},
			HeroNums: int32(pl.HeroManager().GetHeroNums()),
			ItemNums: int32(pl.ItemManager().GetItemNums()),
		}

		acct.SendProtoMessage(reply)
		return nil
	})

	return nil
}

func (m *MsgHandler) handleChangeExp(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2M_ChangeExp)
	if !ok {
		return errors.New("handleChangeExp failed: recv message body error")
	}

	m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		pl, err := m.g.am.GetPlayerByAccount(acct)
		if err != nil {
			return fmt.Errorf("handleChangeExp.AccountExecute failed: %w", err)
		}

		pl.ChangeExp(msg.AddExp)

		// sync player info
		reply := &pbGlobal.M2C_ExpUpdate{
			Exp:   pl.GetExp(),
			Level: pl.GetLevel(),
		}

		acct.SendProtoMessage(reply)
		return nil
	})

	return nil
}

func (m *MsgHandler) handleChangeLevel(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2M_ChangeLevel)
	if !ok {
		return errors.New("handleChangeLevel failed: recv message body error")
	}

	m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		pl, err := m.g.am.GetPlayerByAccount(acct)
		if err != nil {
			return fmt.Errorf("handleChangeLevel.AccountExecute failed: %w", err)
		}

		pl.ChangeLevel(msg.AddLevel)

		// sync player info
		reply := &pbGlobal.M2C_ExpUpdate{
			Exp:   pl.GetExp(),
			Level: pl.GetLevel(),
		}

		acct.SendProtoMessage(reply)

		// sync account info to gate
		acct.Level = pl.GetLevel()
		if _, err := m.g.rpcHandler.CallUpdateUserInfo(acct); err != nil {
			return err
		}

		return nil
	})

	return nil
}

func (m *MsgHandler) handleSyncPlayerInfo(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	fn := func(acct *player.Account) error {
		pl, err := m.g.am.GetPlayerByAccount(acct)
		if err != nil {
			return fmt.Errorf("handleSyncPlayerInfo.AccountExecute failed: %w", err)
		}

		_, err = m.g.rpcHandler.CallSyncPlayerInfo(acct.UserId, &pl.LitePlayer)
		if err != nil {
			return fmt.Errorf("handleSyncPlayerInfo.AccountExecute failed: %w", err)
		}

		acct.SendProtoMessage(&pbGlobal.M2C_SyncPlayerInfo{})

		return nil
	}

	m.g.am.AccountExecute(sock, fn)
	return nil
}

func (m *MsgHandler) handlePublicSyncPlayerInfo(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	fn := func(acct *player.Account) error {
		pl, err := m.g.am.GetPlayerByAccount(acct)
		if err != nil {
			return fmt.Errorf("handlePublicSyncPlayerInfo.AccountExecute failed: %w", err)
		}

		err = m.g.pubSub.PubSyncPlayerInfo(ctx, &pl.LitePlayer)
		if err != nil {
			return fmt.Errorf("handlePublicSyncPlayerInfo.AccountExecute failed: %w", err)
		}

		acct.SendProtoMessage(&pbGlobal.M2C_PublicSyncPlayerInfo{})

		return nil
	}

	m.g.am.AccountExecute(sock, fn)
	return nil
}
