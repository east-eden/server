package combat

import (
	"context"

	pbCombat "github.com/yokaiio/yokai_server/proto/combat"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

type RpcHandler struct {
	c         *Combat
	combatSrv pbCombat.CombatService
	gameSrv   pbGame.GameService
}

func NewRpcHandler(c *Combat) *RpcHandler {
	h := &RpcHandler{
		c: c,
		combatSrv: pbCombat.NewCombatService(
			"",
			c.mi.srv.Client(),
		),

		gameSrv: pbGame.NewGameService(
			"",
			c.mi.srv.Client(),
		),
	}

	pbCombat.RegisterCombatServiceHandler(c.mi.srv.Server(), h)

	return h
}

/////////////////////////////////////////////
// rpc call
/////////////////////////////////////////////

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
func (h *RpcHandler) StartStageCombat(ctx context.Context, req *pbCombat.StartStageCombatReq, rsp *pbCombat.StartStageCombatReply) error {
	h.c.sm.CreateScene(req.SceneType)
	la := h.g.am.GetLiteAccount(req.Id)
	if la == nil {
		rsp.Info = nil
		return nil
	}

	rsp.Info = &pbAccount.LiteAccount{
		Id:    la.ID,
		Name:  la.Name,
		Level: la.Level,
	}

	return nil
}
