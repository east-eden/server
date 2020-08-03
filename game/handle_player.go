package game

import (
	"context"
	"errors"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/player"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"github.com/yokaiio/yokai_server/transport"
)

func (m *MsgHandler) handleQueryPlayerInfo(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	m.g.am.AccountLaterHandle(sock, func(acct *player.Account) {
		reply := &pbGame.M2C_QueryPlayerInfo{
			Error: 0,
		}

		if pl := m.g.am.GetPlayerByAccount(acct); pl != nil {
			reply.Info = &pbGame.PlayerInfo{
				LiteInfo: &pbGame.LitePlayer{
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
	})

	return nil
}

func (m *MsgHandler) handleCreatePlayer(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGame.C2M_CreatePlayer)
	if !ok {
		return errors.New("handleCreatePlayer failed: recv message body error")
	}

	m.g.am.AccountLaterHandle(sock, func(acct *player.Account) {
		pl, err := m.g.am.CreatePlayer(acct, msg.Name)
		reply := &pbGame.M2C_CreatePlayer{
			RpcId: msg.RpcId,
			Error: 0,
		}

		if err != nil {
			reply.Error = -1
			reply.Message = err.Error()
			logger.Warningf("handleCreatePlayer failed: account_id<%d>, error<%s>", acct.ID, err.Error())
		}

		if pl != nil {
			reply.Info = &pbGame.PlayerInfo{
				LiteInfo: &pbGame.LitePlayer{
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
	})

	return nil
}

func (m *MsgHandler) handleSelectPlayer(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGame.MC_SelectPlayer)
	if !ok {
		return errors.New("handleSelectPlayer failed: recv message body error")
	}

	m.g.am.AccountLaterHandle(sock, func(acct *player.Account) {
		pl, err := m.g.am.SelectPlayer(acct, msg.Id)
		reply := &pbGame.MS_SelectPlayer{
			ErrorCode: 0,
		}

		if err != nil {
			reply.ErrorCode = -1
			logger.Warningf("handleSelectPlayer failed: %v", err)
		}

		if pl != nil {
			reply.Info = &pbGame.PlayerInfo{
				LiteInfo: &pbGame.LitePlayer{
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
	})

	return nil
}

func (m *MsgHandler) handleChangeExp(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGame.C2M_ChangeExp)
	if !ok {
		return errors.New("handleChangeExp failed: recv message body error")
	}

	m.g.am.AccountLaterHandle(sock, func(acct *player.Account) {
		pl := m.g.am.GetPlayerByAccount(acct)
		if pl == nil {
			return
		}

		pl.ChangeExp(msg.AddExp)

		// sync player info
		reply := &pbGame.M2C_ExpUpdate{
			Exp:   pl.GetExp(),
			Level: pl.GetLevel(),
		}

		acct.SendProtoMessage(reply)
	})

	return nil
}

func (m *MsgHandler) handleChangeLevel(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGame.C2M_ChangeLevel)
	if !ok {
		return errors.New("handleChangeLevel failed: recv message body error")
	}

	m.g.am.AccountLaterHandle(sock, func(acct *player.Account) {
		pl := m.g.am.GetPlayerByAccount(acct)
		if pl == nil {
			return
		}

		pl.ChangeLevel(msg.AddLevel)

		// sync player info
		reply := &pbGame.M2C_ExpUpdate{
			Exp:   pl.GetExp(),
			Level: pl.GetLevel(),
		}

		acct.SendProtoMessage(reply)

		// sync account info to gate
		acct.Level = pl.GetLevel()
		m.g.rpcHandler.CallUpdateUserInfo(acct)
	})

	return nil
}
