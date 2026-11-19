package game

import (
	"context"
	"errors"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/player"
)

// handleArenaClaim 领取竞技场日榜/周榜奖励
func (m *MsgRegister) handleArenaClaim(ctx context.Context, p ...any) error {
	acct := p[0].(*player.Account)
	msg, ok := p[1].(*pbGlobal.C2S_ArenaClaim)
	if !ok {
		return errors.New("handleArenaClaim failed: recv message body error")
	}

	pl := acct.GetPlayer()
	if pl == nil {
		return ErrPlayerNotFound
	}

	mgr := pl.ArenaManager()

	var err error
	if msg.GetWeekly() {
		err = mgr.ClaimWeekly(msg.GetRank())
	} else {
		err = mgr.ClaimDaily(msg.GetRank())
	}

	if err != nil {
		return err
	}

	st := mgr.State
	pl.SendProtoMessage(&pbGlobal.S2C_ArenaInfo{
		Info: &pbGlobal.ArenaInfo{
			Score:         st.Score,
			Rank:          st.Rank,
			WeekWins:      st.WeekWins,
			DailyClaimed:  st.DailyClaimed,
			WeeklyClaimed: st.WeeklyClaimed,
		},
	})

	return nil
}
