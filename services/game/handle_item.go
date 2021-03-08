package game

import (
	"context"
	"errors"
	"fmt"

	"bitbucket.org/funplus/server/define"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/item"
	"bitbucket.org/funplus/server/services/game/player"
	"bitbucket.org/funplus/server/transport"
)

func (m *MsgHandler) handleAddItem(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_AddItem)
	if !ok {
		return errors.New("handleAddItem failed: recv message body error")
	}
	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleAddItem.AccountExecute failed: %w", err)
	}

	var num int32 = 1
	if !pl.ItemManager().CanAddItem(msg.TypeId, num) {
		return fmt.Errorf("cannot add item<%d> num<%d>", msg.TypeId, num)
	}

	if err := pl.ItemManager().AddItemByTypeId(msg.TypeId, 1); err != nil {
		return fmt.Errorf("handleAddItem.AccountExecute failed: %w", err)
	}

	return nil
}

func (m *MsgHandler) handleDelItem(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_DelItem)
	if !ok {
		return errors.New("handleDelItem failed: recv message body error")
	}
	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleDelItem.AccountExecute failed: %w", err)
	}

	it, err := pl.ItemManager().GetItem(msg.Id)
	if err != nil {
		return fmt.Errorf("handleDelItem.AccountExecute failed: %w", err)
	}

	// clear hero's equip id before delete item
	if it.GetType() == define.Item_TypeEquip {
		equip := it.(*item.Equip)
		equipObjID := equip.GetEquipObj()
		if equipObjID != -1 {
			if err := pl.HeroManager().TakeoffEquip(equipObjID, equip.EquipEnchantEntry.EquipPos); err != nil {
				return fmt.Errorf("TakeoffEquip failed: %w", err)
			}
		}
	}

	// delete item
	return pl.ItemManager().DeleteItem(msg.Id)
}

func (m *MsgHandler) handleUseItem(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_UseItem)
	if !ok {
		return errors.New("handleUseItem failed: recv message body error")
	}
	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleUseItem.AccountExecute failed: %w", err)
	}

	if err := pl.ItemManager().UseItem(msg.ItemId); err != nil {
		return fmt.Errorf("handleUseItem.AccountExecute failed: %w", err)
	}

	return nil
}

func (m *MsgHandler) handleQueryItems(ctx context.Context, acct *player.Account, p *transport.Message) error {
	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleQueryItems.AccountExecute failed: %w", err)
	}

	reply := &pbGlobal.S2C_ItemList{}
	list := pl.ItemManager().GetItemList()
	for _, v := range list {
		i := &pbGlobal.Item{
			Id:     v.Opts().Id,
			TypeId: int32(v.Opts().TypeId),
		}
		reply.Items = append(reply.Items, i)
	}
	acct.SendProtoMessage(reply)
	return nil
}

func (m *MsgHandler) handlePutonEquip(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_PutonEquip)
	if !ok {
		return errors.New("handlePutonEquip failed: recv message body error")
	}
	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handlePutonEquip.AccountExecute failed: %w", err)
	}

	if err := pl.HeroManager().PutonEquip(msg.HeroId, msg.EquipId); err != nil {
		return fmt.Errorf("handlePutonEquip.AccountExecute failed: %w", err)
	}

	return nil
}

func (m *MsgHandler) handleTakeoffEquip(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_TakeoffEquip)
	if !ok {
		return errors.New("handleTakeoffEquip failed: recv message body error")
	}
	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleTakeoffEquip failed: %w", err)
	}

	return pl.HeroManager().TakeoffEquip(msg.HeroId, msg.Pos)
}

func (m *MsgHandler) handleEquipPromote(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_EquipPromote)
	if !ok {
		return errors.New("handleEquipPromote failed: recv message body error")
	}

	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleEquipPromote failed: %w", err)
	}

	return pl.ItemManager().EquipPromote(msg.ItemId)
}

func (m *MsgHandler) handleEquipLevelup(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_EquipLevelup)
	if !ok {
		return errors.New("handleEquipLevelup failed: recv message body error")
	}
	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleEquipLevelup.AccountExecute failed: %w", err)
	}

	return pl.ItemManager().EquipLevelup(msg.GetItemId(), msg.GetStuffItems(), msg.GetExpItems())
}

func (m *MsgHandler) handlePutonCrystal(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_PutonCrystal)
	if !ok {
		return errors.New("handlePutonCrystal failed: recv message body error")
	}

	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handlePutonCrystal failed: %w", err)
	}

	if err := pl.HeroManager().PutonCrystal(msg.HeroId, msg.CrystalId); err != nil {
		return fmt.Errorf("handlePutonCrystal.AccountExecute failed: %w", err)
	}

	return nil
}

func (m *MsgHandler) handleTakeoffCrystal(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_TakeoffCrystal)
	if !ok {
		return errors.New("handleTakeoffCrystal failed: recv message body error")
	}

	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleTakeoffCrystal failed: %w", err)
	}

	if err := pl.HeroManager().TakeoffCrystal(msg.HeroId, msg.Pos); err != nil {
		return fmt.Errorf("handleTakeoffCrystal failed: %w", err)
	}

	return nil
}

func (m *MsgHandler) handleCrystalLevelup(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_CrystalLevelup)
	if !ok {
		return errors.New("handleCrystalLevelup failed: recv message body error")
	}

	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleCrystalLevelup failed: %w", err)
	}

	return pl.ItemManager().CrystalLevelup(msg.GetItemId(), msg.GetStuffItems(), msg.GetExpItems())
}
