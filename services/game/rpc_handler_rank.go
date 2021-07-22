package game

import (
	"context"

	pbRank "e.coding.net/mmstudio/blade/server/proto/server/rank"
	"github.com/spf13/cast"
)

/////////////////////////////////////////////
// rpc call
/////////////////////////////////////////////

// 设置排行积分
func (h *RpcHandler) CallSetRankScore(req *pbRank.SetRankScoreRq) (*pbRank.SetRankScoreRs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()
	return h.rankSrv.SetRankScore(
		ctx,
		req,
		h.consistentHashCallOption(cast.ToString(req.GetRankId())),
		h.retries(3),
	)
}

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
