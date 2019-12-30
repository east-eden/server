package battle

import (
	"context"

	"github.com/urfave/cli/v2"
	pbBattle "github.com/yokaiio/yokai_server/proto/battle"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

type RpcHandler struct {
	b       *Battle
	gameSrv pbGame.GameService
}

func NewRpcHandler(b *Battle, ucli *cli.Context) *RpcHandler {
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
func (h *RpcHandler) GetAccountByID(id int64) (*pbGame.GetAccountByIDReply, error) {
	req := &pbGame.GetAccountByIDRequest{Id: id}
	return h.gameSrv.GetAccountByID(h.b.ctx, req)
}

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
func (h *RpcHandler) GetBattleStatus(ctx context.Context, req *pbBattle.GetBattleStatusRequest, rsp *pbBattle.GetBattleStatusReply) error {
	rsp.Status = &pbBattle.BattleStatus{BattleId: int32(h.b.ID), Health: 2}
	return nil
}
