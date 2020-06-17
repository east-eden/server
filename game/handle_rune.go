package game

import (
	"context"

	logger "github.com/sirupsen/logrus"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"github.com/yokaiio/yokai_server/transport"
)

func (m *MsgHandler) handleAddRune(ctx context.Context, sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("Add Rune failed")
		return
	}

	pl := m.g.am.GetPlayerByAccount(acct)
	if pl == nil {
		return
	}

	msg, ok := p.Body.(*pbGame.C2M_AddRune)
	if !ok {
		logger.Warn("Add Rune failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		if err := pl.RuneManager().AddRuneByTypeID(msg.TypeId); err != nil {
			logger.Warn(err)
			return
		}
	})
}

func (m *MsgHandler) handleDelRune(ctx context.Context, sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("Del Rune failed")
		return
	}

	pl := m.g.am.GetPlayerByAccount(acct)
	if pl == nil {
		return
	}

	msg, ok := p.Body.(*pbGame.C2M_DelRune)
	if !ok {
		logger.Warn("Del Rune failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		if err := pl.RuneManager().DeleteRune(msg.Id); err != nil {
			logger.Warn(err)
			return
		}
	})
}

func (m *MsgHandler) handleQueryRunes(ctx context.Context, sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("Query Runes failed")
		return
	}

	pl := m.g.am.GetPlayerByAccount(acct)
	if pl == nil {
		return
	}

	_, ok := p.Body.(*pbGame.C2M_QueryRunes)
	if !ok {
		logger.Warn("Query Runes failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		rList := pl.RuneManager().GetRuneList()
		reply := &pbGame.M2C_RuneList{}

		for _, v := range rList {
			reply.Runes = append(reply.Runes, &pbGame.Rune{
				Id:         v.GetOptions().Id,
				TypeId:     v.GetOptions().TypeId,
				EquipObjId: v.GetOptions().EquipObj,
			})
		}

		acct.SendProtoMessage(reply)
	})
}

func (m *MsgHandler) handlePutonRune(ctx context.Context, sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("puton rune failed")
		return
	}

	pl := m.g.am.GetPlayerByAccount(acct)
	if pl == nil {
		return
	}

	msg, ok := p.Body.(*pbGame.C2M_PutonRune)
	if !ok {
		logger.Warn("puton rune failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		if err := pl.HeroManager().PutonRune(msg.HeroId, msg.RuneId); err != nil {
			logger.Warn(err)
			return
		}
	})
}

func (m *MsgHandler) handleTakeoffRune(ctx context.Context, sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("takeoff rune failed")
		return
	}

	pl := m.g.am.GetPlayerByAccount(acct)
	if pl == nil {
		return
	}

	msg, ok := p.Body.(*pbGame.C2M_TakeoffRune)
	if !ok {
		logger.Warn("takeoff rune failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		if err := pl.HeroManager().TakeoffRune(msg.HeroId, msg.Pos); err != nil {
			logger.Warn(err)
			return
		}
	})
}
