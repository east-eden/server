package game

import (
	"context"
	"errors"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/player"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"github.com/yokaiio/yokai_server/transport"
)

func (m *MsgHandler) handleAddItem(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGame.C2M_AddItem)
	if !ok {
		return errors.New("handleAddItem failed: recv message body error")
	}

	m.g.am.AccountLaterHandle(sock, func(acct *player.Account) {
		pl := m.g.am.GetPlayerByAccount(acct)
		if pl == nil {
			return
		}

		pl.ItemManager().AddItemByTypeID(msg.TypeId, 1)
	})

	return nil
}

func (m *MsgHandler) handleDelItem(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGame.C2M_DelItem)
	if !ok {
		return errors.New("handleDelItem failed: recv message body error")
	}

	m.g.am.AccountLaterHandle(sock, func(acct *player.Account) {
		pl := m.g.am.GetPlayerByAccount(acct)
		if pl == nil {
			return
		}

		item := pl.ItemManager().GetItem(msg.Id)
		if item == nil {
			logger.Warn("Delete item failed, non-existing item_id:", msg.Id)
			return
		}

		// clear hero's equip id before delete item
		equipObjID := item.GetEquipObj()
		if equipObjID != -1 {
			pl.HeroManager().TakeoffEquip(equipObjID, item.EquipEnchantEntry().EquipPos)
		}

		// delete item
		pl.ItemManager().DeleteItem(msg.Id)
	})

	return nil
}

func (m *MsgHandler) handleUseItem(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGame.C2M_UseItem)
	if !ok {
		return errors.New("handleUseItem failed: recv message body error")
	}

	m.g.am.AccountLaterHandle(sock, func(acct *player.Account) {
		pl := m.g.am.GetPlayerByAccount(acct)
		if pl == nil {
			return
		}

		pl.ItemManager().UseItem(msg.ItemId)
	})

	return nil
}

func (m *MsgHandler) handleQueryItems(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	m.g.am.AccountLaterHandle(sock, func(acct *player.Account) {
		pl := m.g.am.GetPlayerByAccount(acct)
		if pl == nil {
			return
		}

		reply := &pbGame.M2C_ItemList{}
		list := pl.ItemManager().GetItemList()
		for _, v := range list {
			i := &pbGame.Item{
				Id:     v.GetOptions().Id,
				TypeId: v.GetOptions().TypeId,
			}
			reply.Items = append(reply.Items, i)
		}
		acct.SendProtoMessage(reply)
	})

	return nil
}

func (m *MsgHandler) handlePutonEquip(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGame.C2M_PutonEquip)
	if !ok {
		return errors.New("handlePutonEquip failed: recv message body error")
	}

	m.g.am.AccountLaterHandle(sock, func(acct *player.Account) {
		pl := m.g.am.GetPlayerByAccount(acct)
		if pl == nil {
			return
		}

		if err := pl.HeroManager().PutonEquip(msg.HeroId, msg.EquipId); err != nil {
			logger.Warn(err)
			return
		}
	})

	return nil
}

func (m *MsgHandler) handleTakeoffEquip(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGame.C2M_TakeoffEquip)
	if !ok {
		return errors.New("handleTakeoffEquip failed: recv message body error")
	}

	m.g.am.AccountLaterHandle(sock, func(acct *player.Account) {
		pl := m.g.am.GetPlayerByAccount(acct)
		if pl == nil {
			return
		}

		if err := pl.HeroManager().TakeoffEquip(msg.HeroId, msg.Pos); err != nil {
			logger.Warn(err)
			return
		}
	})

	return nil
}
