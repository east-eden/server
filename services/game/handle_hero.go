package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/player"
	"github.com/east-eden/server/transport"
)

func (m *MsgHandler) handleAddHero(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_AddHero)
	if !ok {
		return errors.New("handleAddHero failed: recv message body error")
	}
	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleAddHero.AccountExecute failed: %w", err)
	}

	_ = pl.HeroManager().AddHeroByTypeID(msg.TypeId)
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

//func (m *MsgHandler) handleHeroAddExp(sock transport.Socket, p *transport.Message) {
//cli := m.g.cm.GetClientBySock(sock)
//if cli == nil {
//logger.WithFields(logger.Fields{
//"client_id":   cli.ID(),
//"client_name": cli.Name(),
//}).Warn("hero add exp failed")
//return
//}

//msg, ok := p.Body.(*pbGame.MC_HeroAddExp)
//if !ok {
//logger.Warn("hero add exp failed, recv message body error")
//return
//}

//if cli.Player() == nil {
//logger.Warn("client has no player", cli.ID())
//return
//}

//cli.Player().HeroManager().HeroAddExp(msg.HeroId, msg.Exp)
//hero := cli.Player().HeroManager().GetHero(msg.HeroId)
//if hero == nil {
//logger.Warn("get hero by id error:", msg.HeroId)
//return
//}

//reply := &pbGame.MS_HeroInfo{
//Info: &pbGame.Hero{
//Id:     hero.GetID(),
//TypeId: hero.GetTypeID(),
//Exp:    hero.GetExp(),
//Level:  hero.GetLevel(),
//},
//}

//cli.SendProtoMessage(reply)
//}

//func (m *MsgHandler) handleHeroAddLevel(sock transport.Socket, p *transport.Message) {
//cli := m.g.cm.GetClientBySock(sock)
//if cli == nil {
//logger.WithFields(logger.Fields{
//"client_id":   cli.ID(),
//"client_name": cli.Name(),
//}).Warn("hero add level failed")
//return
//}

//msg, ok := p.Body.(*pbGame.MC_HeroAddLevel)
//if !ok {
//logger.Warn("hero add level failed, recv message body error")
//return
//}

//if cli.Player() == nil {
//logger.Warn("client has no player", cli.ID())
//return
//}

//cli.Player().HeroManager().HeroAddLevel(msg.HeroId, msg.Level)
//hero := cli.Player().HeroManager().GetHero(msg.HeroId)
//if hero == nil {
//logger.Warn("get hero by id error:", msg.HeroId)
//return
//}

//reply := &pbGame.MS_HeroInfo{
//Info: &pbGame.Hero{
//Id:     hero.GetID(),
//TypeId: hero.GetTypeID(),
//Exp:    hero.GetExp(),
//Level:  hero.GetLevel(),
//},
//}

//cli.SendProtoMessage(reply)
//}
