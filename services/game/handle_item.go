package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "bitbucket.org/east-eden/server/proto/global"
	"bitbucket.org/east-eden/server/services/game/player"
	"bitbucket.org/east-eden/server/transport"
)

func (m *MsgHandler) handleAddItem(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2M_AddItem)
	if !ok {
		return errors.New("handleAddItem failed: recv message body error")
	}

	m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		pl, err := m.g.am.GetPlayerByAccount(acct)
		if err != nil {
			return fmt.Errorf("handleAddItem.AccountExecute failed: %w", err)
		}

		if err := pl.ItemManager().AddItemByTypeID(msg.TypeId, 1); err != nil {
			return fmt.Errorf("handleAddItem.AccountExecute failed: %w", err)
		}

		return nil
	})

	return nil
}

func (m *MsgHandler) handleDelItem(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2M_DelItem)
	if !ok {
		return errors.New("handleDelItem failed: recv message body error")
	}

	m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		pl, err := m.g.am.GetPlayerByAccount(acct)
		if err != nil {
			return fmt.Errorf("handleDelItem.AccountExecute failed: %w", err)
		}

		item, err := pl.ItemManager().GetItem(msg.Id)
		if err != nil {
			return fmt.Errorf("handleDelItem.AccountExecute failed: %w", err)
		}

		// clear hero's equip id before delete item
		equipObjID := item.GetEquipObj()
		if equipObjID != -1 {
			if err := pl.HeroManager().TakeoffEquip(equipObjID, item.EquipEnchantEntry().EquipPos); err != nil {
				return fmt.Errorf("TakeoffEquip failed: %w", err)
			}
		}

		// delete item
		return pl.ItemManager().DeleteItem(msg.Id)
	})

	return nil
}

func (m *MsgHandler) handleUseItem(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2M_UseItem)
	if !ok {
		return errors.New("handleUseItem failed: recv message body error")
	}

	m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		pl, err := m.g.am.GetPlayerByAccount(acct)
		if err != nil {
			return fmt.Errorf("handleUseItem.AccountExecute failed: %w", err)
		}

		if err := pl.ItemManager().UseItem(msg.ItemId); err != nil {
			return fmt.Errorf("handleUseItem.AccountExecute failed: %w", err)
		}

		return nil
	})

	return nil
}

func (m *MsgHandler) handleQueryItems(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		pl, err := m.g.am.GetPlayerByAccount(acct)
		if err != nil {
			return fmt.Errorf("handleQueryItems.AccountExecute failed: %w", err)
		}

		reply := &pbGlobal.M2C_ItemList{}
		list := pl.ItemManager().GetItemList()
		for _, v := range list {
			i := &pbGlobal.Item{
				Id:     v.GetOptions().Id,
				TypeId: int32(v.GetOptions().TypeId),
			}
			reply.Items = append(reply.Items, i)
		}
		acct.SendProtoMessage(reply)
		return nil
	})

	return nil
}

func (m *MsgHandler) handlePutonEquip(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2M_PutonEquip)
	if !ok {
		return errors.New("handlePutonEquip failed: recv message body error")
	}

	m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		pl, err := m.g.am.GetPlayerByAccount(acct)
		if err != nil {
			return fmt.Errorf("handlePutonEquip.AccountExecute failed: %w", err)
		}

		if err := pl.HeroManager().PutonEquip(msg.HeroId, msg.EquipId); err != nil {
			return fmt.Errorf("handlePutonEquip.AccountExecute failed: %w", err)
		}

		return nil
	})

	return nil
}

func (m *MsgHandler) handleTakeoffEquip(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2M_TakeoffEquip)
	if !ok {
		return errors.New("handleTakeoffEquip failed: recv message body error")
	}

	m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		pl, err := m.g.am.GetPlayerByAccount(acct)
		if err != nil {
			return fmt.Errorf("handleTakeoffEquip.AccountExecute failed: %w", err)
		}

		if err := pl.HeroManager().TakeoffEquip(msg.HeroId, msg.Pos); err != nil {
			return fmt.Errorf("handleTakeoffEquip.AccountExecute failed: %w", err)
		}

		return nil
	})

	return nil
}
