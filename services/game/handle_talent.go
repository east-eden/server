package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/player"
	"github.com/east-eden/server/transport"
)

func (m *MsgHandler) handleAddTalent(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_AddTalent)
	if !ok {
		return errors.New("handleAddTalent failed: recv message body error")
	}

	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleAddTalent.AccountExecute failed: %w", err)
	}

	blade, err := pl.BladeManager().GetBlade(msg.BladeId)
	if err != nil {
		return fmt.Errorf("handleAddTalent.AccountExecute failed: %w", err)
	}

	err = blade.GetTalentManager().AddTalent(msg.TalentId)
	if err != nil {
		return fmt.Errorf("Account.ExecutorHandler failed: %w", err)
	}

	list := blade.GetTalentManager().GetTalentList()
	reply := &pbGlobal.S2C_TalentList{
		BladeId: blade.GetOptions().Id,
		Talents: make([]*pbGlobal.Talent, 0, len(list)),
	}

	for _, v := range list {
		reply.Talents = append(reply.Talents, &pbGlobal.Talent{
			Id: int32(v.Id),
		})
	}

	acct.SendProtoMessage(reply)
	return nil
}

func (m *MsgHandler) handleQueryTalents(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_QueryTalents)
	if !ok {
		return errors.New("handleQueryTalents failed: recv message body error")
	}

	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleQueryTalents.AccountExecute failed: %w", err)
	}

	blade, err := pl.BladeManager().GetBlade(msg.BladeId)
	if err != nil {
		return fmt.Errorf("handleQueryTalents.AccountExecute failed: %w", err)
	}

	list := blade.GetTalentManager().GetTalentList()
	reply := &pbGlobal.S2C_TalentList{
		BladeId: msg.BladeId,
		Talents: make([]*pbGlobal.Talent, 0, len(list)),
	}

	for _, v := range list {
		reply.Talents = append(reply.Talents, &pbGlobal.Talent{
			Id: int32(v.Id),
		})
	}

	acct.SendProtoMessage(reply)
	return nil
}
