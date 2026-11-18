package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/player"
)

func (m *MsgRegister) handleDelItem(ctx context.Context, p ...any) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_DelItem)
	if !ok {
		return errors.New("handleDelItem failed: recv message body error")
	}
	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	if _, err := pl.ItemManager().GetItem(msg.Id); err != nil {
		return fmt.Errorf("handleDelItem.AccountExecute failed: %w", err)
	}

	// delete item
	return pl.ItemManager().DeleteItem(msg.Id)
}

func (m *MsgRegister) handleUseItem(ctx context.Context, p ...any) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_UseItem)
	if !ok {
		return errors.New("handleUseItem failed: recv message body error")
	}
	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	if err := pl.ItemManager().UseItem(msg.ItemId); err != nil {
		return fmt.Errorf("handleUseItem.AccountExecute failed: %w", err)
	}

	return nil
}
