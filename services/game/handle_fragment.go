package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/player"
)

func (m *MsgRegister) handleQueryFragments(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	_, ok := p[1].(*pbGlobal.C2S_QueryFragments)
	if !ok {
		return errors.New("handleQueryFragments failed: recv message body error")
	}
	pl, err := m.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleQueryFragments.AccountExecute failed: %w", err)
	}

	reply := &pbGlobal.S2C_FragmentsList{}
	reply.Frags = pl.FragmentManager().GetFragmentList()
	acct.SendProtoMessage(reply)
	return nil
}

func (m *MsgRegister) handleFragmentsCompose(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_FragmentsCompose)
	if !ok {
		return errors.New("handleFragmentsCompose failed: recv message body error")
	}
	pl, err := m.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleFragmentsCompose.AccountExecute failed: %w", err)
	}

	return pl.FragmentManager().Compose(msg.FragId)
}
