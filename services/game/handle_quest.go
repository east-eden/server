package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/player"
)

func (m *MsgRegister) handlePlayerQuestReward(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_PlayerQuestReward)
	if !ok {
		return errors.New("handlePlayerQuestReward failed: recv message body error")
	}

	pl, err := m.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handlePlayerQuestReward.GetPlayerByAccount failed: %w", err)
	}

	return pl.QuestManager().QuestReward(msg.GetId())
}
