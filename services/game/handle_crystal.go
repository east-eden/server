package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/player"
	"bitbucket.org/funplus/server/transport"
)

func (m *MsgHandler) handleAddCrystal(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_AddCrystal)
	if !ok {
		return errors.New("handleAddCrystal failed: recv message body error")
	}

	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleAddCrystal failed: %w", err)
	}

	if err := pl.CrystalManager().AddCrystalByTypeId(msg.TypeId); err != nil {
		return fmt.Errorf("handleAddCrystal failed: %w", err)
	}

	return nil
}

func (m *MsgHandler) handleDelCrystal(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_DelCrystal)
	if !ok {
		return errors.New("handleDelCrystal failed: recv message body error")
	}

	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleDelCrystal failed: %w", err)
	}

	if err := pl.CrystalManager().DeleteCrystal(msg.Id); err != nil {
		return fmt.Errorf("handleDelCrystal.AccountExecute failed: %w", err)
	}

	return nil
}

func (m *MsgHandler) handleQueryCrystals(ctx context.Context, acct *player.Account, p *transport.Message) error {
	_, ok := p.Body.(*pbGlobal.C2S_QueryCrystals)
	if !ok {
		return errors.New("handleQueryCrystals failed: recv message body error")
	}

	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleQueryCrystals.AccountExecute failed: %w", err)
	}

	rList := pl.CrystalManager().GetCrystalList()
	reply := &pbGlobal.S2C_CrystalList{}

	for _, v := range rList {
		reply.Crystals = append(reply.Crystals, &pbGlobal.Crystal{
			Id:         v.GetOptions().Id,
			TypeId:     int32(v.GetOptions().TypeId),
			EquipObjId: v.GetOptions().EquipObj,
		})
	}

	acct.SendProtoMessage(reply)
	return nil
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
		return fmt.Errorf("handleTakeoffCrystal.AccountExecute failed: %w", err)
	}

	if err := pl.HeroManager().TakeoffCrystal(msg.HeroId, msg.Pos); err != nil {
		return fmt.Errorf("handleTakeoffCrystal.AccountExecute failed: %w", err)
	}

	return nil
}
