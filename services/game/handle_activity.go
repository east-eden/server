package game

import (
	"context"
	"errors"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/player"
)

// handleActivitySignIn 活动签到
func (m *MsgRegister) handleActivitySignIn(ctx context.Context, p ...any) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_ActivitySignIn)
	if !ok {
		return errors.New("handleActivitySignIn failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	return pl.ActivityManager().SignIn(msg.GetActivityId(), msg.GetDay())
}

// handleActivityClaimTask 领取活动任务奖励
func (m *MsgRegister) handleActivityClaimTask(ctx context.Context, p ...any) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_ActivityClaimTask)
	if !ok {
		return errors.New("handleActivityClaimTask failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	return pl.ActivityManager().ClaimTask(msg.GetActivityId(), msg.GetTaskId())
}
