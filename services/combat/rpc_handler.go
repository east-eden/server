package combat

import (
	"context"
	"errors"

	"e.coding.net/mmstudio/blade/server/excel/auto"
	pbCombat "e.coding.net/mmstudio/blade/server/proto/server/combat"
	pbGame "e.coding.net/mmstudio/blade/server/proto/server/game"
	"e.coding.net/mmstudio/blade/server/services/combat/scene"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidStage      = errors.New("invalid stage id")
	ErrInvalidScene      = errors.New("invalid scene id")
	ErrInvalidBattleWave = errors.New("invalid battle wave entry")
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

	sceneEntry, ok := auto.GetSceneEntry(1)
	if !ok {
		return ErrInvalidScene
	}

	battleWaveEntries := make([]*auto.BattleWaveEntry, 0, len(stageEntry.WaveID))
	for _, id := range stageEntry.WaveID {
		entry, ok := auto.GetBattleWaveEntry(id)
		if !ok {
			log.Error().Caller().
				Int32("stage_id", req.GetStageId()).
				Int32("wave_id", id).
				Msg("can not find BattleWaveEntry")
			return ErrInvalidBattleWave
		}

		battleWaveEntries = append(battleWaveEntries, entry)
	}

	sc, err := h.c.sm.CreateScene(
		ctx,
		scene.WithSceneAttackId(req.AttackId),
		scene.WithSceneAttackUnitList(req.AttackEntityList),
		scene.WithSceneEntry(sceneEntry),
		scene.WithSceneBattleWaveEntries(battleWaveEntries...),
	)

	if !utils.ErrCheck(err, "CreateScene failed when RpcHandler.StageCombat", req.GetStageId(), req.GetAttackId()) {
		return err
	}

	rsp.Win = sc.GetResult()
	// todo 关卡条件

	return nil
}
