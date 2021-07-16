package game

import (
	"context"
	"errors"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/player"
)

func (m *MsgRegister) handleDelHero(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_DelHero)
	if !ok {
		return errors.New("handelDelHero failed: recv message body error")
	}
	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	pl.HeroManager().DelHero(msg.Id)
	return nil
}

func (m *MsgRegister) handleHeroLevelup(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_HeroLevelup)
	if !ok {
		return errors.New("handelHeroLevelup failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	return pl.HeroManager().HeroLevelup(msg.GetHeroId(), msg.GetStuffItemTypeId(), msg.GetUseNum())
}

func (m *MsgRegister) handleHeroPromote(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_HeroPromote)
	if !ok {
		return errors.New("handleHeroPromote failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	return pl.HeroManager().HeroPromote(msg.GetHeroId())
}

func (m *MsgRegister) handleHeroStarup(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_HeroStarup)
	if !ok {
		return errors.New("handleHeroStarup failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	return pl.HeroManager().HeroStarup(msg.GetHeroId())
}

func (m *MsgRegister) handleHeroTalentChoose(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_HeroTalentChoose)
	if !ok {
		return errors.New("handleHeroTalentChoose failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	return pl.HeroManager().HeroTalentChoose(msg.GetHeroId(), msg.GetTalentId())
}
