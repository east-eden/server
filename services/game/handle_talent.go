package game

import (
	"context"
	"errors"
	"fmt"

	pbGame "github.com/east-eden/server/proto/game"
	"github.com/east-eden/server/services/game/player"
	"github.com/east-eden/server/transport"
)

func (m *MsgHandler) handleAddTalent(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGame.C2M_AddTalent)
	if !ok {
		return errors.New("handleAddTalent failed: recv message body error")
	}

	return m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		pl, err := m.g.am.GetPlayerByAccount(acct)
		if err != nil {
			return fmt.Errorf("handleAddTalent.AccountExecute failed: %w", err)
		}

		blade, err := pl.BladeManager().GetBlade(msg.BladeId)
		if err != nil {
			return fmt.Errorf("handleAddTalent.AccountExecute failed: %w", err)
		}

		err = blade.GetTalentManager().AddTalent(int(msg.TalentId))
		if err != nil {
			return fmt.Errorf("Account.ExecutorHandler failed: %w", err)
		}

		list := blade.GetTalentManager().GetTalentList()
		reply := &pbGame.M2C_TalentList{
			BladeId: blade.GetOptions().Id,
			Talents: make([]*pbGame.Talent, 0, len(list)),
		}

		for _, v := range list {
			reply.Talents = append(reply.Talents, &pbGame.Talent{
				Id: int32(v.Id),
			})
		}

		acct.SendProtoMessage(reply)
		return nil
	})
}

func (m *MsgHandler) handleQueryTalents(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGame.C2M_QueryTalents)
	if !ok {
		return errors.New("handleQueryTalents failed: recv message body error")
	}

	return m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		pl, err := m.g.am.GetPlayerByAccount(acct)
		if err != nil {
			return fmt.Errorf("handleQueryTalents.AccountExecute failed: %w", err)
		}

		blade, err := pl.BladeManager().GetBlade(msg.BladeId)
		if err != nil {
			return fmt.Errorf("handleQueryTalents.AccountExecute failed: %w", err)
		}

		list := blade.GetTalentManager().GetTalentList()
		reply := &pbGame.M2C_TalentList{
			BladeId: msg.BladeId,
			Talents: make([]*pbGame.Talent, 0, len(list)),
		}

		for _, v := range list {
			reply.Talents = append(reply.Talents, &pbGame.Talent{
				Id: int32(v.Id),
			})
		}

		acct.SendProtoMessage(reply)
		return nil
	})
}
