package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/player"
	"bitbucket.org/funplus/server/transport"
)

func (m *MsgHandler) handleQueryPlayerInfo(ctx context.Context, acct *player.Account, p *transport.Message) error {
	reply := &pbGlobal.S2C_QueryPlayerInfo{Error: 0}
	if pl, err := m.g.am.GetPlayerByAccount(acct); err == nil {
		reply.Info = &pbGlobal.PlayerInfo{
			Id:        pl.GetID(),
			AccountId: pl.GetAccountID(),
			Name:      pl.GetName(),
			Exp:       pl.GetExp(),
			Level:     pl.GetLevel(),
		}
	}

	acct.SendProtoMessage(reply)
	return nil
}

func (m *MsgHandler) handleCreatePlayer(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_CreatePlayer)
	if !ok {
		return errors.New("handleCreatePlayer failed: recv message body error")
	}

	reply := &pbGlobal.S2C_CreatePlayer{}

	pl, err := m.g.am.CreatePlayer(acct, msg.Name)
	if err != nil {
		acct.SendProtoMessage(reply)
		return fmt.Errorf("handleCreatePlayer.AccountExecute failed: %w", err)
	}

	reply.Info = &pbGlobal.PlayerInfo{
		Id:        pl.GetID(),
		AccountId: pl.GetAccountID(),
		Name:      pl.GetName(),
		Exp:       pl.GetExp(),
		Level:     pl.GetLevel(),
	}

	acct.SendProtoMessage(reply)
	return nil
}

func (m *MsgHandler) handleGmCmd(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_GmCmd)
	if !ok {
		return errors.New("handleGmCmd failed: recv message body error")
	}

	pl, err := m.g.am.GetPlayerByAccount(acct)
	if err != nil {
		return err
	}

	return player.GmCmd(pl, msg.Cmd)
}
