package game

import (
	"context"

	pbBattle "github.com/yokaiio/yokai_server/proto/battle"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

type RpcHandler struct {
	g         *Game
	battleSrv pbBattle.BattleService
}

func NewRpcHandler(g *Game) *RpcHandler {
	h := &RpcHandler{
		g: g,
		battleSrv: pbBattle.NewBattleService(
			"yokai_battle",
			g.mi.srv.Client(),
		),
	}

	pbGame.RegisterGameServiceHandler(g.mi.srv.Server(), h)

	return h
}

/////////////////////////////////////////////
// rpc call
/////////////////////////////////////////////
func (h *RpcHandler) GetBattleStatus() (*pbBattle.GetBattleStatusReply, error) {
	req := &pbBattle.GetBattleStatusRequest{}
	return h.battleSrv.GetBattleStatus(h.g.ctx, req)
}

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
func (h *RpcHandler) GetClientByID(ctx context.Context, req *pbGame.GetClientByIDRequest, rsp *pbGame.GetClientByIDReply) error {
	client := h.g.cm.GetClientByID(req.Id)
	rsp.Info.Id = client.GetID()
	rsp.Info.Name = client.GetName()
	return nil
}
