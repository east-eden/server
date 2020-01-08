package game

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/transport"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

func (m *MsgHandler) handleAddTalent(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("add talent failed")
		return
	}

	if acct.GetPlayer() == nil {
		return
	}

	msg, ok := p.Body.(*pbGame.MC_AddTalent)
	if !ok {
		logger.Warn("Add Talent failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		blade := acct.GetPlayer().BladeManager().GetBlade(msg.BladeId)
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
		reply := &pbGame.MS_TalentList{
			BladeId: blade.GetID(),
			Talents: make([]*pbGame.Talent, 0, len(list)),
		}

		for _, v := range list {
			reply.Talents = append(reply.Talents, &pbGame.Talent{
				Id: v.ID,
			})
		}

		acct.SendProtoMessage(reply)
	})
}

func (m *MsgHandler) handleQueryTalents(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("query talents failed")
		return
	}

	if acct.GetPlayer() == nil {
		return
	}

	msg, ok := p.Body.(*pbGame.MC_QueryTalents)
	if !ok {
		logger.Warn("Query Talents failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		blade := acct.GetPlayer().BladeManager().GetBlade(msg.BladeId)
		if blade == nil {
			logger.Warn("non-existing blade_id:", msg.BladeId)
			return
		}

		list := blade.TalentManager().GetTalentList()
		reply := &pbGame.MS_TalentList{
			BladeId: msg.BladeId,
			Talents: make([]*pbGame.Talent, 0, len(list)),
		}

		for _, v := range list {
			reply.Talents = append(reply.Talents, &pbGame.Talent{
				Id: v.ID,
			})
		}

		acct.SendProtoMessage(reply)
	})
}
