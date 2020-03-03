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
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("add item failed")
		return
	}

	if acct.GetPlayer() == nil {
		return
	}

	msg, ok := p.Body.(*pbGame.MC_AddItem)
	if !ok {
		logger.Warn("Add Item failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		acct.GetPlayer().ItemManager().AddItemByTypeID(msg.TypeId, 1)
		list := acct.GetPlayer().ItemManager().GetItemList()
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
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("delete item failed")
		return
	}

	if acct.GetPlayer() == nil {
		return
	}

	msg, ok := p.Body.(*pbGame.MC_DelItem)
	if !ok {
		logger.Warn("Delete item failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		item := acct.GetPlayer().ItemManager().GetItem(msg.Id)
		if item == nil {
			logger.Warn("Delete item failed, non-existing item_id:", msg.Id)
			return
		}

		// clear hero's equip id before delete item
		equipObjID := item.GetEquipObj()
		if equipObjID != -1 {
			acct.GetPlayer().HeroManager().TakeoffEquip(equipObjID, item.EquipEnchantEntry().EquipPos)
		}

		// delete item
		acct.GetPlayer().ItemManager().DelItem(msg.Id)

		// reply to client
		list := acct.GetPlayer().ItemManager().GetItemList()
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
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("query items failed")
		return
	}

	if acct.GetPlayer() == nil {
		return
	}

	acct.PushWrapHandler(func() {
		list := acct.GetPlayer().ItemManager().GetItemList()
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
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("Puton equip failed")
		return
	}

	if acct.GetPlayer() == nil {
		return
	}

	msg, ok := p.Body.(*pbGame.MC_PutonEquip)
	if !ok {
		logger.Warn("Puton equip failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		equip := acct.GetPlayer().ItemManager().GetItem(msg.EquipId)
		if equip == nil {
			logger.Warn("Puton equip failed, non-existing item:", msg.EquipId)
			return
		}

		hero := acct.GetPlayer().HeroManager().GetHero(msg.HeroId)
		if hero == nil {
			logger.Warn("Puton equip failed, non-existing hero:", msg.HeroId)
			return
		}

		if err := acct.GetPlayer().HeroManager().PutonEquip(msg.HeroId, msg.EquipId, equip.EquipEnchantEntry().EquipPos); err != nil {
			logger.Warn(err)
			return
		}

		acct.GetPlayer().ItemManager().SetItemEquiped(msg.EquipId, msg.HeroId)

		reply := &pbGame.MS_HeroEquips{
			HeroId: msg.HeroId,
			Equips: make([]*pbGame.Item, 0),
		}

		equips := hero.GetEquips()
		for _, v := range equips {
			if v == -1 {
				continue
			}

			it := acct.GetPlayer().ItemManager().GetItem(v)
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
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("Takeoff equip failed")
		return
	}

	if acct.GetPlayer() == nil {
		return
	}

	msg, ok := p.Body.(*pbGame.MC_TakeoffEquip)
	if !ok {
		logger.Warn("Takeoff equip failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		hero := acct.GetPlayer().HeroManager().GetHero(msg.HeroId)
		if hero == nil {
			logger.Warn("Takeoff equip failed, non-existing hero:", msg.HeroId)
			return
		}

		equipID := hero.GetEquip(msg.Pos)
		if err := acct.GetPlayer().HeroManager().TakeoffEquip(msg.HeroId, msg.Pos); err != nil {
			logger.Warn(err)
			return
		}

		acct.GetPlayer().ItemManager().SetItemUnEquiped(equipID)

		reply := &pbGame.MS_HeroEquips{
			HeroId: msg.HeroId,
			Equips: make([]*pbGame.Item, 0),
		}

		equips := hero.GetEquips()
		for _, v := range equips {
			if v == -1 {
				continue
			}

			it := acct.GetPlayer().ItemManager().GetItem(v)
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
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("Query hero equips failed")
		return
	}

	if acct.GetPlayer() == nil {
		return
	}

	msg, ok := p.Body.(*pbGame.MC_QueryHeroEquips)
	if !ok {
		logger.Warn("Query hero equips failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		hero := acct.GetPlayer().HeroManager().GetHero(msg.HeroId)
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

			it := acct.GetPlayer().ItemManager().GetItem(v)
			i := &pbGame.Item{
				Id:     v,
				TypeId: it.GetTypeID(),
			}
			reply.Equips = append(reply.Equips, i)
		}
		acct.SendProtoMessage(reply)
	})
}
