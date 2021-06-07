package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
	"e.coding.net/mmstudio/blade/server/services/game/player"
)

func (m *MsgRegister) handleTowerChallenge(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_TowerChallenge)
	if !ok {
		return errors.New("handleTowerChallenge failed: recv message body error")
	}

	pl, err := m.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("GetPlayerByAccount failed: %w", err)
	}

	return pl.TowerManager.Challenge(msg.TowerType, msg.TowerFloor, msg.BattleArray)
}
