package game

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/transport"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

func (m *MsgHandler) handleAddItem(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.ID(),
			"account_name": acct.Name(),
		}).Warn("add item failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_AddItem)
	if !ok {
		logger.Warn("Add Item failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		acct.Player().ItemManager().AddItemByTypeID(msg.TypeId, 1)
		list := acct.Player().ItemManager().GetItemList()
		reply := &pbGame.MS_ItemList{Items: make([]*pbGame.Item, 0, len(list))}
		for _, v := range list {
			i := &pbGame.Item{
				Id:     v.GetID(),
				TypeId: v.GetTypeID(),
			}
			reply.Items = append(reply.Items, i)
		}
		acct.SendProtoMessage(reply)
	})

}

func (m *MsgHandler) handleDelItem(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.ID(),
			"account_name": acct.Name(),
		}).Warn("delete item failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_DelItem)
	if !ok {
		logger.Warn("Delete item failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		item := acct.Player().ItemManager().GetItem(msg.Id)
		if item == nil {
			logger.Warn("Delete item failed, non-existing item_id:", msg.Id)
			return
		}

		// clear hero's equip id before delete item
		equipObjID := item.GetEquipObj()
		if equipObjID != -1 {
			acct.Player().HeroManager().TakeoffEquip(equipObjID, item.Entry().EquipPos)
		}

		// delete item
		acct.Player().ItemManager().DelItem(msg.Id)

		// reply to client
		list := acct.Player().ItemManager().GetItemList()
		reply := &pbGame.MS_ItemList{Items: make([]*pbGame.Item, 0, len(list))}
		for _, v := range list {
			i := &pbGame.Item{
				Id:     v.GetID(),
				TypeId: v.GetTypeID(),
			}
			reply.Items = append(reply.Items, i)
		}
		acct.SendProtoMessage(reply)
	})
}

func (m *MsgHandler) handleQueryItems(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.ID(),
			"account_name": acct.Name(),
		}).Warn("query items failed")
		return
	}

	acct.PushWrapHandler(func() {
		list := acct.Player().ItemManager().GetItemList()
		reply := &pbGame.MS_ItemList{Items: make([]*pbGame.Item, 0, len(list))}
		for _, v := range list {
			i := &pbGame.Item{
				Id:     v.GetID(),
				TypeId: v.GetTypeID(),
			}
			reply.Items = append(reply.Items, i)
		}
		acct.SendProtoMessage(reply)
	})
}

func (m *MsgHandler) handlePutonEquip(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.ID(),
			"account_name": acct.Name(),
		}).Warn("Puton equip failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_PutonEquip)
	if !ok {
		logger.Warn("Puton equip failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		equip := acct.Player().ItemManager().GetItem(msg.EquipId)
		if equip == nil {
			logger.Warn("Puton equip failed, non-existing item:", msg.EquipId)
			return
		}

		hero := acct.Player().HeroManager().GetHero(msg.HeroId)
		if hero == nil {
			logger.Warn("Puton equip failed, non-existing hero:", msg.HeroId)
			return
		}

		if err := acct.Player().HeroManager().PutonEquip(msg.HeroId, msg.EquipId, equip.Entry().EquipPos); err != nil {
			logger.Warn(err)
			return
		}

		acct.Player().ItemManager().SetItemEquiped(msg.EquipId, msg.HeroId)

		reply := &pbGame.MS_HeroEquips{
			HeroId: msg.HeroId,
			Equips: make([]*pbGame.Item, 0),
		}

		equips := hero.GetEquips()
		for _, v := range equips {
			if v == -1 {
				continue
			}

			it := acct.Player().ItemManager().GetItem(v)
			i := &pbGame.Item{
				Id:     v,
				TypeId: it.GetTypeID(),
			}
			reply.Equips = append(reply.Equips, i)
		}
		acct.SendProtoMessage(reply)
	})

}

func (m *MsgHandler) handleTakeoffEquip(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.ID(),
			"account_name": acct.Name(),
		}).Warn("Takeoff equip failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_TakeoffEquip)
	if !ok {
		logger.Warn("Takeoff equip failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		hero := acct.Player().HeroManager().GetHero(msg.HeroId)
		if hero == nil {
			logger.Warn("Takeoff equip failed, non-existing hero:", msg.HeroId)
			return
		}

		equipID := hero.GetEquip(msg.Pos)
		if err := acct.Player().HeroManager().TakeoffEquip(msg.HeroId, msg.Pos); err != nil {
			logger.Warn(err)
			return
		}

		acct.Player().ItemManager().SetItemUnEquiped(equipID)

		reply := &pbGame.MS_HeroEquips{
			HeroId: msg.HeroId,
			Equips: make([]*pbGame.Item, 0),
		}

		equips := hero.GetEquips()
		for _, v := range equips {
			if v == -1 {
				continue
			}

			it := acct.Player().ItemManager().GetItem(v)
			i := &pbGame.Item{
				Id:     v,
				TypeId: it.GetTypeID(),
			}
			reply.Equips = append(reply.Equips, i)
		}
		acct.SendProtoMessage(reply)
	})
}

func (m *MsgHandler) handleQueryHeroEquips(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.ID(),
			"account_name": acct.Name(),
		}).Warn("Query hero equips failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_QueryHeroEquips)
	if !ok {
		logger.Warn("Query hero equips failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		hero := acct.Player().HeroManager().GetHero(msg.HeroId)
		if hero == nil {
			logger.Warn("Query hero equips failed, non-existing hero_id:", msg.HeroId)
			return
		}

		reply := &pbGame.MS_HeroEquips{
			HeroId: msg.HeroId,
			Equips: make([]*pbGame.Item, 0),
		}

		equips := hero.GetEquips()
		for _, v := range equips {
			if v == -1 {
				continue
			}

			it := acct.Player().ItemManager().GetItem(v)
			i := &pbGame.Item{
				Id:     v,
				TypeId: it.GetTypeID(),
			}
			reply.Equips = append(reply.Equips, i)
		}
		acct.SendProtoMessage(reply)
	})
}
