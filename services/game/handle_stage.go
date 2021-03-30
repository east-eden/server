package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/player"
	"github.com/east-eden/server/transport"
)

func (m *MsgRegister) handleStageSweep(ctx context.Context, acct *player.Account, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_StageSweep)
	if !ok {
		return errors.New("handleStageSweep failed: recv message body error")
	}

	pl, err := m.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleStageSweep.GetPlayerByAccount failed: %w", err)
	}

	return pl.ChapterStageManager.StageSweep(msg.StageId, msg.Times)
}
