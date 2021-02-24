package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/player"
	"github.com/east-eden/server/transport"
)

func (m *MsgHandler) handleAddRune(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_AddRune)
	if !ok {
		return errors.New("handleAddRune failed: recv message body error")
	}

	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleAddRune.AccountExecute failed: %w", err)
	}

	if err := pl.RuneManager().AddRuneByTypeID(msg.TypeId); err != nil {
		return fmt.Errorf("handleAddRune.AccountExecute failed: %w", err)
	}

	return nil
}

func (m *MsgHandler) handleDelRune(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_DelRune)
	if !ok {
		return errors.New("handleDelRune failed: recv message body error")
	}

	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleDelRune.AccountExecute failed: %w", err)
	}

	if err := pl.RuneManager().DeleteRune(msg.Id); err != nil {
		return fmt.Errorf("handleDelRune.AccountExecute failed: %w", err)
	}

	return nil
}

func (m *MsgHandler) handleQueryRunes(ctx context.Context, acct *player.Account, p *transport.Message) error {
	_, ok := p.Body.(*pbGlobal.C2S_QueryRunes)
	if !ok {
		return errors.New("handleQueryRunes failed: recv message body error")
	}

	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleQueryRunes.AccountExecute failed: %w", err)
	}

	rList := pl.RuneManager().GetRuneList()
	reply := &pbGlobal.S2C_RuneList{}

	for _, v := range rList {
		reply.Runes = append(reply.Runes, &pbGlobal.Rune{
			Id:         v.GetOptions().Id,
			TypeId:     int32(v.GetOptions().TypeId),
			EquipObjId: v.GetOptions().EquipObj,
		})
	}

	acct.SendProtoMessage(reply)
	return nil
}

func (m *MsgHandler) handlePutonRune(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_PutonRune)
	if !ok {
		return errors.New("handlePutonRune failed: recv message body error")
	}

	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handlePutonRune.AccountExecute failed: %w", err)
	}

	if err := pl.HeroManager().PutonRune(msg.HeroId, msg.RuneId); err != nil {
		return fmt.Errorf("handlePutonRune.AccountExecute failed: %w", err)
	}

	return nil
}

func (m *MsgHandler) handleTakeoffRune(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_TakeoffRune)
	if !ok {
		return errors.New("handleTakeoffRune failed: recv message body error")
	}

	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleTakeoffRune.AccountExecute failed: %w", err)
	}

	if err := pl.HeroManager().TakeoffRune(msg.HeroId, msg.Pos); err != nil {
		return fmt.Errorf("handleTakeoffRune.AccountExecute failed: %w", err)
	}

	return nil
}
