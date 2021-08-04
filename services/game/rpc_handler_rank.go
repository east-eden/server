package game

import (
	"context"
	"errors"

	"github.com/east-eden/server/excel/auto"
	pbRank "github.com/east-eden/server/proto/server/rank"
	"github.com/spf13/cast"
)

var (
	ErrRpcInvalidRankId = errors.New("invalid rank id")
)

/////////////////////////////////////////////
// rpc call
/////////////////////////////////////////////

// 请求排行
func (h *RpcHandler) CallQueryRankByObjId(req *pbRank.QueryRankByObjIdRq) (*pbRank.QueryRankByObjIdRs, error) {
	rankEntry, ok := auto.GetRankEntry(req.GetRankId())
	if !ok {
		return nil, ErrRpcInvalidRankId
	}

	var consistentKey string
	if rankEntry.Local {
		consistentKey = cast.ToString(h.g.ID)
	} else {
		consistentKey = cast.ToString(req.GetRankId())
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()
	return h.rankSrv.QueryRankByObjId(
		ctx,
		req,
		h.consistentHashCallOption(consistentKey),
		h.retries(3),
	)
}

// 设置排行积分
func (h *RpcHandler) CallSetRankScore(req *pbRank.SetRankScoreRq) (*pbRank.SetRankScoreRs, error) {
	rankEntry, ok := auto.GetRankEntry(req.GetRankId())
	if !ok {
		return nil, ErrRpcInvalidRankId
	}

	var consistentKey string
	if rankEntry.Local {
		consistentKey = cast.ToString(h.g.ID)
	} else {
		consistentKey = cast.ToString(req.GetRankId())
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()
	return h.rankSrv.SetRankScore(
		ctx,
		req,
		h.consistentHashCallOption(consistentKey),
		h.retries(3),
	)
}

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
