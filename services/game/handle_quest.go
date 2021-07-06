package game

import (
	"context"
	"errors"

	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
	"e.coding.net/mmstudio/blade/server/services/game/player"
)

func (m *MsgRegister) handlePlayerQuestReward(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_PlayerQuestReward)
	if !ok {
		return errors.New("handlePlayerQuestReward failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	return pl.QuestManager.QuestReward(msg.GetId())
}
