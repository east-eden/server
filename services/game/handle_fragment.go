package game

import (
	"context"
	"errors"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/player"
)

func (m *MsgRegister) handleHeroFragmentsCompose(ctx context.Context, p ...any) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_HeroFragmentsCompose)
	if !ok {
		return errors.New("handleHeroFragmentsCompose failed: recv message body error")
	}
	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	h := pl.HeroManager().GetHeroByTypeId(msg.FragId)
	if h != nil {
		return errors.New("hero already existed, no need to compose")
	}

	return pl.FragmentManager().HeroCompose(msg.FragId)
}

func (m *MsgRegister) handleCollectionFragmentsCompose(ctx context.Context, p ...any) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_CollectionFragmentsCompose)
	if !ok {
		return errors.New("handleCollectionFragmentsCompose failed: recv message body error")
	}
	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	return pl.FragmentManager().CollectionCompose(msg.FragId)
}
