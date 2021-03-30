package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/player"
	"github.com/east-eden/server/transport"
)

func (m *MsgRegister) handleQueryFragments(ctx context.Context, acct *player.Account, p *transport.Message) error {
	_, ok := p.Body.(*pbGlobal.C2S_QueryFragments)
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

func (m *MsgRegister) handleFragmentsCompose(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_FragmentsCompose)
	if !ok {
		return errors.New("handleFragmentsCompose failed: recv message body error")
	}
	pl, err := m.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleFragmentsCompose.AccountExecute failed: %w", err)
	}

	return pl.FragmentManager().Compose(msg.FragId)
}
