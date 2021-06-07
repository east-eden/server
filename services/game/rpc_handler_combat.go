package game

import (
	"context"

	pbCombat "e.coding.net/mmstudio/blade/server/proto/server/combat"
	"github.com/spf13/cast"
)

/////////////////////////////////////////////
// rpc call
/////////////////////////////////////////////

// 关卡战斗
func (h *RpcHandler) CallStageCombat(req *pbCombat.StageCombatRq) (*pbCombat.StageCombatRs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()
	return h.combatSrv.StageCombat(
		ctx,
		req,
		h.consistentHashCallOption(cast.ToString(req.GetAttackId())),
	)
}
