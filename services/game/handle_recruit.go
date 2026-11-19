package game

import (
	"context"
	"errors"
	"time"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/player"
	"github.com/east-eden/server/services/game/recruit"
)

// handleRecruit 招募：免费或高级
func (m *MsgRegister) handleRecruit(ctx context.Context, p ...any) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_Recruit)
	if !ok {
		return errors.New("handleRecruit failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	mgr := pl.RecruitManager()

	var (
		rewards []recruit.Reward
		err     error
	)

	switch msg.GetRecruitType() {
	case recruit.TypeFree:
		rewards, err = mgr.Free(time.Now())
	case recruit.TypeSenior:
		rewards, err = mgr.Senior(msg.GetTimes())
	default:
		return errors.New("handleRecruit failed: invalid recruit type")
	}

	if err != nil {
		return err
	}

	result := &pbGlobal.S2C_RecruitResult{SeniorTimes: mgr.State.SeniorTimes}
	for _, r := range rewards {
		result.Rewards = append(result.Rewards, &pbGlobal.RecruitReward{ItemId: r.ItemId, Num: r.Num})
	}

	pl.SendProtoMessage(result)

	st := mgr.State
	pl.SendProtoMessage(&pbGlobal.S2C_RecruitInfo{
		FreeUsed:       st.FreeUsed,
		SeniorTimes:    st.SeniorTimes,
		BigGuarantee:   st.BigGuarantee,
		SmallGuarantee: st.SmallGuarantee,
	})

	return nil
}
