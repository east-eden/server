package game

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/transport"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

func (m *MsgHandler) handleQueryPlayerInfos(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.ID(),
			"account_name": acct.Name(),
		}).Warn("query player info failed")
		return
	}

	playerList := m.g.pm.GetPlayers(acct.ID())
	reply := &pbGame.MS_QueryPlayerInfos{
		Infos: make([]*pbGame.PlayerInfo, 0, len(playerList)),
	}

	for _, v := range playerList {
		info := &pbGame.PlayerInfo{
			LiteInfo: &pbGame.LitePlayer{
				Id:        v.GetID(),
				AccountId: v.GetAccountID(),
				Name:      v.GetName(),
				Exp:       v.GetExp(),
				Level:     v.GetLevel(),
			},

			HeroNums: int32(v.HeroManager().GetHeroNums()),
			ItemNums: int32(v.ItemManager().GetItemNums()),
		}

		reply.Infos = append(reply.Infos, info)
	}

	acct.SendProtoMessage(reply)
}

func (m *MsgHandler) handleCreatePlayer(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.ID(),
			"account_name": acct.Name(),
		}).Warn("create player failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_CreatePlayer)
	if !ok {
		logger.Warn("create player failed, recv message body error")
		return
	}

	pl, err := m.g.am.CreatePlayer(acct, msg.Name)
	reply := &pbGame.MS_CreatePlayer{
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
}

func (m *MsgHandler) handleSelectPlayer(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.ID(),
			"account_name": acct.Name(),
		}).Warn("select player failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_SelectPlayer)
	if !ok {
		logger.Warn("Select player failed, recv message body error")
		return
	}

	pl, err := m.g.am.SelectPlayer(acct, msg.Id)
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
}

func (m *MsgHandler) handleExpirePlayer(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.ID(),
			"account_name": acct.Name(),
		}).Warn("select player failed")
		return
	}

	if acct.Player() == nil {
		return
	}

	m.g.ExpirePlayer(acct.Player().GetID())
}

func (m *MsgHandler) handleChangeExp(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.ID(),
			"account_name": acct.Name(),
		}).Warn("change exp failed")
		return
	}

	if acct.Player() == nil {
		return
	}

	msg, ok := p.Body.(*pbGame.MC_ChangeExp)
	if !ok {
		logger.Warn("change exp failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		acct.Player().ChangeExp(msg.AddExp)

		// sync player info
		pl := acct.Player()
		reply := &pbGame.MS_QueryPlayerInfo{
			Info: &pbGame.PlayerInfo{
				LiteInfo: &pbGame.LitePlayer{
					Id:        pl.GetID(),
					AccountId: pl.GetAccountID(),
					Name:      pl.GetName(),
					Exp:       pl.GetExp(),
					Level:     pl.GetLevel(),
				},
				HeroNums: int32(pl.HeroManager().GetHeroNums()),
				ItemNums: int32(pl.ItemManager().GetItemNums()),
			},
		}

		acct.SendProtoMessage(reply)
	})
}

func (m *MsgHandler) handleChangeLevel(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.ID(),
			"account_name": acct.Name(),
		}).Warn("change level failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_ChangeLevel)
	if !ok {
		logger.Warn("change level failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		acct.Player().ChangeLevel(msg.AddLevel)

		// sync player info
		pl := acct.Player()
		reply := &pbGame.MS_QueryPlayerInfo{
			Info: &pbGame.PlayerInfo{
				LiteInfo: &pbGame.LitePlayer{
					Id:        pl.GetID(),
					AccountId: pl.GetAccountID(),
					Name:      pl.GetName(),
					Exp:       pl.GetExp(),
					Level:     pl.GetLevel(),
				},
				HeroNums: int32(pl.HeroManager().GetHeroNums()),
				ItemNums: int32(pl.ItemManager().GetItemNums()),
			},
		}

		acct.SendProtoMessage(reply)
	})
}
