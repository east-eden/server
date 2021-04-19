package combat

import (
	"context"
	"errors"

	"bitbucket.org/funplus/server/excel/auto"
	pbCombat "bitbucket.org/funplus/server/proto/server/combat"
	pbGame "bitbucket.org/funplus/server/proto/server/game"
	"bitbucket.org/funplus/server/services/combat/scene"
	"bitbucket.org/funplus/server/utils"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidStage     = errors.New("invalid stage id")
	ErrInvalidScene     = errors.New("invalid scene id")
	ErrInvalidUnitGroup = errors.New("invalid unit group entry")
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

	_ = pbCombat.RegisterCombatServiceHandler(c.mi.srv.Server(), h)

	return h
}

/////////////////////////////////////////////
// rpc call
/////////////////////////////////////////////

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
func (h *RpcHandler) StageCombat(ctx context.Context, req *pbCombat.StageCombatRq, rsp *pbCombat.StageCombatRs) error {
	log.Info().Interface("request", req).Msg("recv rpc call StageCombat")

	stageEntry, ok := auto.GetStageEntry(req.GetStageId())
	if !ok {
		return ErrInvalidStage
	}

	sceneEntry, ok := auto.GetSceneEntry(stageEntry.SceneId)
	if !ok {
		return ErrInvalidScene
	}

	unitGroupEntry, ok := auto.GetUnitGroupEntry(stageEntry.UnitGroupId)
	if !ok {
		return ErrInvalidUnitGroup
	}

	sc, err := h.c.sm.CreateScene(
		ctx,
		scene.WithSceneAttackId(req.AttackId),
		scene.WithSceneAttackUnitList(req.AttackEntityList),
		scene.WithSceneEntry(sceneEntry),
		scene.WithSceneUnitGroupEntry(unitGroupEntry),
	)

	if !utils.ErrCheck(err, "CreateScene failed when RpcHandler.StageCombat", req.GetStageId(), req.GetAttackId()) {
		return err
	}

	rsp.Win = sc.GetResult()
	// todo 关卡条件

	return nil
}
