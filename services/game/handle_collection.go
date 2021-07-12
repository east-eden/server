package game

import (
	"context"
	"errors"

	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
	"e.coding.net/mmstudio/blade/server/services/game/player"
)

func (m *MsgRegister) handleCollectionActive(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_CollectionActive)
	if !ok {
		return errors.New("handleCollectionActive failed: recv message body error")
	}
	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	return pl.CollectionManager().CollectionActive(msg.TypeId)
}

func (m *MsgRegister) handleCollectionStarup(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_CollectionStarup)
	if !ok {
		return errors.New("handleCollectionStarup failed: recv message body error")
	}
	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	return pl.CollectionManager().CollectionStarup(msg.TypeId)
}

func (m *MsgRegister) handleCollectionWakeup(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_CollectionWakeup)
	if !ok {
		return errors.New("handleCollectionWakeup failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	return pl.CollectionManager().CollectionWakeup(msg.TypeId)
}
