package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/player"
	"bitbucket.org/funplus/server/transport"
)

func (m *MsgHandler) handleDelHero(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_DelHero)
	if !ok {
		return errors.New("handelDelHero failed: recv message body error")
	}
	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleDelHero.AccountExecute failed: %w", err)
	}

	pl.HeroManager().DelHero(msg.Id)
	list := pl.HeroManager().GetHeroList()
	reply := &pbGlobal.S2C_HeroList{}
	for _, v := range list {
		h := &pbGlobal.Hero{
			Id:             v.GetOptions().Id,
			TypeId:         int32(v.GetOptions().TypeId),
			Exp:            v.GetOptions().Exp,
			Level:          int32(v.GetOptions().Level),
			PromoteLevel:   int32(v.GetOptions().PromoteLevel),
			Star:           int32(v.GetOptions().Star),
			NormalSpellId:  v.GetOptions().NormalSpellId,
			SpecialSpellId: v.GetOptions().SpecialSpellId,
			RageSpellId:    v.GetOptions().RageSpellId,
			Friendship:     v.GetOptions().Friendship,
			FashionId:      v.GetOptions().FashionId,
		}
		reply.Heros = append(reply.Heros, h)
	}
	acct.SendProtoMessage(reply)
	return nil
}

func (m *MsgHandler) handleQueryHeros(ctx context.Context, acct *player.Account, p *transport.Message) error {
	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleQueryHeros.AccountExecute failed: %w", err)
	}

	list := pl.HeroManager().GetHeroList()
	reply := &pbGlobal.S2C_HeroList{}
	for _, v := range list {
		h := &pbGlobal.Hero{
			Id:             v.GetOptions().Id,
			TypeId:         int32(v.GetOptions().TypeId),
			Exp:            v.GetOptions().Exp,
			Level:          int32(v.GetOptions().Level),
			PromoteLevel:   int32(v.GetOptions().PromoteLevel),
			Star:           int32(v.GetOptions().Star),
			NormalSpellId:  v.GetOptions().NormalSpellId,
			SpecialSpellId: v.GetOptions().SpecialSpellId,
			RageSpellId:    v.GetOptions().RageSpellId,
			Friendship:     v.GetOptions().Friendship,
			FashionId:      v.GetOptions().FashionId,
		}
		reply.Heros = append(reply.Heros, h)
	}
	acct.SendProtoMessage(reply)
	return nil
}

func (m *MsgHandler) handleQueryHeroAtt(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_QueryHeroAtt)
	if !ok {
		return errors.New("handelQueryHeroAtt failed: recv message body error")
	}

	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleQueryHeroAtt failed: %w", err)
	}

	h := pl.HeroManager().GetHero(msg.HeroId)
	if h == nil {
		return fmt.Errorf("handleQueryHeroAtt failed: cannot find hero<%d>", msg.HeroId)
	}

	h.GetAttManager().CalcAtt()
	pl.HeroManager().SendHeroAtt(h)
	return nil
}

func (m *MsgHandler) handleHeroLevelup(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_HeroLevelup)
	if !ok {
		return errors.New("handelHeroLevelup failed: recv message body error")
	}

	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleHeroLevelup failed: %w", err)
	}

	return pl.HeroManager().HeroLevelup(msg.GetHeroId(), msg.GetStuffItems())
}

func (m *MsgHandler) handleHeroPromote(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_HeroPromote)
	if !ok {
		return errors.New("handleHeroPromote failed: recv message body error")
	}

	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleHeroPromote failed: %w", err)
	}

	return pl.HeroManager().HeroPromote(msg.GetHeroId())
}
