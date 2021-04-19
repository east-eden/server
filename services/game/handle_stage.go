package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/player"
)

func (m *MsgRegister) handleStageChallenge(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_StageChallenge)
	if !ok {
		return errors.New("handleStageChallenge failed: recv message body error")
	}

	pl, err := m.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleStageChallenge.AccountExecute failed: %w", err)
	}

	return pl.ChapterStageManager.StageChallenge(msg.StageId)
}

func (m *MsgRegister) handleStageSweep(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_StageSweep)
	if !ok {
		return errors.New("handleStageSweep failed: recv message body error")
	}

	pl, err := m.am.GetPlayerByAccount(acct)
	if err != nil {
		return fmt.Errorf("handleStageSweep.GetPlayerByAccount failed: %w", err)
	}

	return pl.ChapterStageManager.StageSweep(msg.StageId, msg.Times)
}
