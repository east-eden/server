package game

import (
	"context"

	logger "github.com/sirupsen/logrus"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"github.com/yokaiio/yokai_server/transport"
)

func (m *MsgHandler) handleQueryPlayerInfo(ctx context.Context, sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("query player info failed")
		return
	}

	acct.PushWrapHandler(func() {
		reply := &pbGame.M2C_QueryPlayerInfo{
			Error: 0,
		}

		if pl := m.g.pm.GetPlayerByAccount(acct); pl != nil {
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
}

func (m *MsgHandler) handleCreatePlayer(ctx context.Context, sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("create player failed")
		return
	}

	msg, ok := p.Body.(*pbGame.C2M_CreatePlayer)
	if !ok {
		logger.Warn("create player failed, recv message body error")
		return
	}

	pl, err := m.g.am.CreatePlayer(acct, msg.Name)

	acct.PushWrapHandler(func() {
		reply := &pbGame.M2C_CreatePlayer{
			RpcId: msg.RpcId,
			Error: 0,
		}

		if err != nil {
			reply.Error = -1
			reply.Message = err.Error()
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
}

func (m *MsgHandler) handleSelectPlayer(ctx context.Context, sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("select player failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_SelectPlayer)
	if !ok {
		logger.Warn("Select player failed, recv message body error")
		return
	}

	pl, err := m.g.am.SelectPlayer(acct, msg.Id)

	acct.PushWrapHandler(func() {
		reply := &pbGame.MS_SelectPlayer{
			ErrorCode: 0,
		}

		if err != nil {
			reply.ErrorCode = -1
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
}

func (m *MsgHandler) handleChangeExp(ctx context.Context, sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("change exp failed")
		return
	}

	pl := m.g.pm.GetPlayerByAccount(acct)
	if pl == nil {
		return
	}

	msg, ok := p.Body.(*pbGame.C2M_ChangeExp)
	if !ok {
		logger.Warn("change exp failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		pl.ChangeExp(msg.AddExp)

		// sync player info
		reply := &pbGame.M2C_ExpUpdate{
			Exp: pl.GetExp(),
		}

		acct.SendProtoMessage(reply)
	})
}

func (m *MsgHandler) handleChangeLevel(ctx context.Context, sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("change level failed")
		return
	}

	msg, ok := p.Body.(*pbGame.C2M_ChangeLevel)
	if !ok {
		logger.Warn("change level failed, recv message body error")
		return
	}

	pl := m.g.pm.GetPlayerByAccount(acct)
	if pl == nil {
		return
	}

	acct.PushWrapHandler(func() {
		pl.ChangeLevel(msg.AddLevel)

		// sync player info
		reply := &pbGame.M2C_LevelUpdate{
			Level: pl.GetLevel(),
		}

		acct.SendProtoMessage(reply)

		// sync account info to gate
		acct.Level = pl.GetLevel()
		m.g.rpcHandler.CallUpdateUserInfo(acct)
	})
}
