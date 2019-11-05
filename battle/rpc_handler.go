package battle

import (
	"context"

	pbBattle "github.com/yokaiio/yokai_server/proto/battle"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

type RpcHandler struct {
	b       *Battle
	gameSrv pbGame.GameService
}

func NewRpcHandler(b *Battle) *RpcHandler {
	h := &RpcHandler{
		b: b,
		gameSrv: pbGame.NewGameService(
			"yokai_game",
			b.mi.srv.Client(),
		),
	}

	pbBattle.RegisterBattleServiceHandler(b.mi.srv.Server(), h)

	return h
}

/////////////////////////////////////////////
// rpc call
/////////////////////////////////////////////
func (h *RpcHandler) GetClientByID(id int64) (*pbGame.GetClientByIDReply, error) {
	req := &pbGame.GetClientByIDRequest{Id: id}
	return h.gameSrv.GetClientByID(h.b.ctx, req)
}

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
func (h *RpcHandler) GetBattleStatus(ctx context.Context, req *pbBattle.GetBattleStatusRequest, rsp *pbBattle.GetBattleStatusReply) error {
	rsp.Status.BattleId = int32(h.b.opts.BattleID)
	rsp.Status.Health = 2
	return nil
}
