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
	sc, err := h.c.sm.CreateScene(req.SceneId, req.SceneType, req.AttackId, req.DefenceId, req.AttackUnitList, req.DefenceUnitList)
	if err != nil {
		logger.WithFields(logger.Fields{
			"scene_type":  req.SceneType,
			"attack_id":   req.AttackId,
			"defenece_id": req.DefenceId,
		}).Warn("CreateScene failed")
		return nil
	}

	rsp.SceneId = sc.GetID()
	rsp.AttackId = req.AttackId
	rsp.DefenceId = req.defenceId
	rsp.Result = sc.GetResult()

	return nil
}
