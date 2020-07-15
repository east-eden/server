package game

import (
	"context"
	"errors"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/player"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"github.com/yokaiio/yokai_server/transport"
)

func (m *MsgHandler) handleAddTalent(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGame.C2M_AddTalent)
	if !ok {
		return errors.New("handleAddTalent failed: recv message body error")
	}

	m.g.am.AccountLaterHandle(sock, func(acct *player.Account) {
		pl := m.g.am.GetPlayerByAccount(acct)
		if pl == nil {
			return
		}

		blade := pl.BladeManager().GetBlade(msg.BladeId)
		if blade == nil {
			logger.Warn("non-existing blade_id:", msg.BladeId)
			return
		}

		err := blade.TalentManager().AddTalent(msg.TalentId)
		if err != nil {
			logger.Warn("add talent failed:", err)
			return
		}

		list := blade.TalentManager().GetTalentList()
		reply := &pbGame.M2C_TalentList{
			BladeId: blade.GetOptions().Id,
			Talents: make([]*pbGame.Talent, 0, len(list)),
		}

		for _, v := range list {
			reply.Talents = append(reply.Talents, &pbGame.Talent{
				Id: v.Id,
			})
		}

		acct.SendProtoMessage(reply)
	})

	return nil
}

func (m *MsgHandler) handleQueryTalents(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGame.C2M_QueryTalents)
	if !ok {
		return errors.New("handleQueryTalents failed: recv message body error")
	}

	m.g.am.AccountLaterHandle(sock, func(acct *player.Account) {
		pl := m.g.am.GetPlayerByAccount(acct)
		if pl == nil {
			return
		}

		blade := pl.BladeManager().GetBlade(msg.BladeId)
		if blade == nil {
			logger.Warn("non-existing blade_id:", msg.BladeId)
			return
		}

		list := blade.TalentManager().GetTalentList()
		reply := &pbGame.M2C_TalentList{
			BladeId: msg.BladeId,
			Talents: make([]*pbGame.Talent, 0, len(list)),
		}

		for _, v := range list {
			reply.Talents = append(reply.Talents, &pbGame.Talent{
				Id: v.Id,
			})
		}

		acct.SendProtoMessage(reply)
	})

	return nil
}
