package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/player"
)

func (m *MsgRegister) handleQueryHeroFragments(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	_, ok := p[1].(*pbGlobal.C2S_QueryHeroFragments)
	if !ok {
		return errors.New("handleQueryHeroFragments failed: recv message body error")
	}
	pl, err := m.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("GetPlayerByAccount failed: %w", err)
	}

	reply := &pbGlobal.S2C_HeroFragmentsList{}
	reply.Frags = pl.FragmentManager().HeroFragmentManager.GetFragmentList()
	acct.SendProtoMessage(reply)
	return nil
}

func (m *MsgRegister) handleQueryCollectionFragments(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	_, ok := p[1].(*pbGlobal.C2S_QueryCollectionFragments)
	if !ok {
		return errors.New("handleQueryCollectionFragments failed: recv message body error")
	}
	pl, err := m.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("GetPlayerByAccount failed: %w", err)
	}

	reply := &pbGlobal.S2C_CollectionFragmentsList{}
	reply.Frags = pl.FragmentManager().CollectionFragmentManager.GetFragmentList()
	acct.SendProtoMessage(reply)
	return nil
}

func (m *MsgRegister) handleHeroFragmentsCompose(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_HeroFragmentsCompose)
	if !ok {
		return errors.New("handleHeroFragmentsCompose failed: recv message body error")
	}
	pl, err := m.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("GetPlayerByAccount failed: %w", err)
	}

	h := pl.HeroManager().GetHeroByTypeId(msg.FragId)
	if h != nil {
		return errors.New("hero already existed, no need to compose")
	}

	return pl.FragmentManager().HeroCompose(msg.FragId)
}

func (m *MsgRegister) handleCollectionFragmentsCompose(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_CollectionFragmentsCompose)
	if !ok {
		return errors.New("handleCollectionFragmentsCompose failed: recv message body error")
	}
	pl, err := m.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("GetPlayerByAccount failed: %w", err)
	}

	return pl.FragmentManager().CollectionCompose(msg.FragId)
}
