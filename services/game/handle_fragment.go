package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "bitbucket.org/east-eden/server/proto/global"
	"bitbucket.org/east-eden/server/services/game/player"
	"bitbucket.org/east-eden/server/transport"
)

func (m *MsgHandler) handleQueryFragments(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	_, ok := p.Body.(*pbGlobal.C2S_QueryFragments)
	if !ok {
		return errors.New("handleQueryFragments failed: recv message body error")
	}

	m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		pl, err := m.g.am.GetPlayerByAccount(acct)
		if err != nil {
			return fmt.Errorf("handleQueryFragments.AccountExecute failed: %w", err)
		}

		reply := &pbGlobal.S2C_FragmentsList{}
		reply.Frags = pl.FragmentManager().GetFragmentList()
		acct.SendProtoMessage(reply)
		return nil
	})

	return nil
}

func (m *MsgHandler) handleFragmentsCompose(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_FragmentsCompose)
	if !ok {
		return errors.New("handleFragmentsCompose failed: recv message body error")
	}

	m.g.am.AccountExecute(sock, func(acct *player.Account) error {
		pl, err := m.g.am.GetPlayerByAccount(acct)
		if err != nil {
			return fmt.Errorf("handleFragmentsCompose.AccountExecute failed: %w", err)
		}

		return pl.FragmentManager().Compose(msg.FragId)
	})

	return nil
}
