package game

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/transport"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

func (m *MsgHandler) handleAddRune(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("Add Rune failed")
		return
	}

	pl := m.g.pm.GetPlayerByAccount(acct)
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

func (m *MsgHandler) handleDelRune(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("Del Rune failed")
		return
	}

	pl := m.g.pm.GetPlayerByAccount(acct)
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

func (m *MsgHandler) handleQueryRunes(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("Query Runes failed")
		return
	}

	pl := m.g.pm.GetPlayerByAccount(acct)
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
				Id:         v.ID,
				TypeId:     v.TypeID,
				EquipObjId: v.EquipObj,
			})
		}

		acct.SendProtoMessage(reply)
	})
}
