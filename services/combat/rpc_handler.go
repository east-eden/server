package combat

import (
	"context"

	"github.com/rs/zerolog/log"
	pbCombat "github.com/yokaiio/yokai_server/proto/combat"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"github.com/yokaiio/yokai_server/services/combat/scene"
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
	log.Info().Interface("request", req).Msg("recv rpc call StartStageCombat")

	sc, err := h.c.sm.CreateScene(
		ctx,
		req.SceneId,
		req.SceneType,
		scene.WithSceneAttackId(req.AttackId),
		scene.WithSceneDefenceId(req.DefenceId),
		scene.WithSceneAttackUnitList(req.AttackUnitList),
		scene.WithSceneDefenceUnitList(req.DefenceUnitList),
	)

	if err != nil {
		log.Warn().
			Int32("scene_type", req.SceneType).
			Int64("attack_id", req.AttackId).
			Int64("defence_id", req.DefenceId).
			Msg("CreateScene failed")
		return nil
	}

	rsp.SceneId = sc.GetID()
	rsp.AttackId = req.AttackId
	rsp.DefenceId = req.DefenceId
	rsp.Result = sc.GetResult()

	return nil
}
