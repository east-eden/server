package game

import (
	"context"
	"errors"
	"fmt"
	"time"

	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
	pbCombat "e.coding.net/mmstudio/blade/server/proto/server/combat"
	pbGame "e.coding.net/mmstudio/blade/server/proto/server/game"
	pbGate "e.coding.net/mmstudio/blade/server/proto/server/gate"
	pbMail "e.coding.net/mmstudio/blade/server/proto/server/mail"
	pbRank "e.coding.net/mmstudio/blade/server/proto/server/rank"
	"e.coding.net/mmstudio/blade/server/services/game/player"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/asim/go-micro/v3/client"
	log "github.com/rs/zerolog/log"
	"github.com/spf13/cast"
)

var (
	DefaultRpcTimeout          = 5 * time.Second // 默认rpc超时时间
	ErrRpcCannotFindPlayerInfo = errors.New("cannot find player info")
)

type RpcHandler struct {
	g         *Game
	gateSrv   pbGate.GateService
	gameSrv   pbGame.GameService
	combatSrv pbCombat.CombatService
	mailSrv   pbMail.MailService
	rankSrv   pbRank.RankService
}

func NewRpcHandler(g *Game) *RpcHandler {
	h := &RpcHandler{
		g: g,
		gateSrv: pbGate.NewGateService(
			"gate",
			g.mi.srv.Client(),
		),

		gameSrv: pbGame.NewGameService(
			"game",
			g.mi.srv.Client(),
		),

		combatSrv: pbCombat.NewCombatService(
			"combat",
			g.mi.srv.Client(),
		),

		mailSrv: pbMail.NewMailService(
			"mail",
			g.mi.srv.Client(),
		),

		rankSrv: pbRank.NewRankService(
			"rank",
			g.mi.srv.Client(),
		),
	}

	err := pbGame.RegisterGameServiceHandler(g.mi.srv.Server(), h)
	if err != nil {
		log.Fatal().Err(err).Msg("RegisterGameServiceHandler failed")
	}

	return h
}

// 一致性哈希
func (h *RpcHandler) consistentHashCallOption(key string) client.CallOption {
	return client.WithSelectOption(
		utils.ConsistentHashSelector(h.g.cons, key),
	)
}

// 重试次数
func (h *RpcHandler) retries(times int) client.CallOption {
	return client.WithRetries(times)
}

/////////////////////////////////////////////
// rpc call
/////////////////////////////////////////////
func (h *RpcHandler) CallGetGateStatus(ctx context.Context) (*pbGate.GetGateStatusReply, error) {
	req := &pbGate.GateEmptyMessage{}
	return h.gateSrv.GetGateStatus(ctx, req)
}

func (h *RpcHandler) CallGetRemotePlayerInfo(playerID int64) (*pbGame.GetRemotePlayerInfoRs, error) {
	req := &pbGame.GetRemotePlayerInfoRq{Id: playerID}
	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()
	return h.gameSrv.GetRemotePlayerInfo(
		ctx,
		req,
		h.consistentHashCallOption(cast.ToString(playerID)),
	)
}

func (h *RpcHandler) CallSyncPlayerInfo(userId int64, info *player.PlayerInfo) (*pbGate.SyncPlayerInfoReply, error) {
	req := &pbGate.SyncPlayerInfoRequest{
		UserId: userId,
		Info: &pbGlobal.PlayerInfo{
			Id:        info.ID,
			AccountId: info.AccountID,
			Name:      info.Name,
			Exp:       info.Exp,
			Level:     info.Level,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()
	return h.gateSrv.SyncPlayerInfo(
		ctx,
		req,
		h.consistentHashCallOption(cast.ToString(info.ID)),
	)
}

// 踢account下线
func (h *RpcHandler) CallKickAccountOffline(accountId int64, gameId int32) (*pbGame.KickAccountOfflineRs, error) {
	if accountId == -1 {
		return nil, errors.New("invalid account id")
	}

	if gameId == int32(h.g.ID) {
		return nil, errors.New("same game id")
	}

	req := &pbGame.KickAccountOfflineRq{
		AccountId: accountId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()

	return h.gameSrv.KickAccountOffline(
		ctx,
		req,
		client.WithSelectOption(
			utils.SpecificIDSelector(
				fmt.Sprintf("game-%d", gameId),
			),
		),
	)
}

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
func (h *RpcHandler) GetRemotePlayerInfo(ctx context.Context, req *pbGame.GetRemotePlayerInfoRq, rsp *pbGame.GetRemotePlayerInfoRs) error {
	info := h.g.am.GetPlayerInfoById(req.Id)
	if info == nil {
		return ErrRpcCannotFindPlayerInfo
	}

	rsp.Info = info.ToPB()
	return nil
}

func (h *RpcHandler) UpdatePlayerExp(ctx context.Context, req *pbGame.UpdatePlayerExpRequest, rsp *pbGame.UpdatePlayerExpReply) error {
	//lp, err := h.g.pm.getLitePlayer(req.Id)
	//if err != nil {
	//return err
	//}

	//rsp.Id = lp.GetID()
	//rsp.Exp = lp.GetExp()
	return nil
}

func (h *RpcHandler) KickAccountOffline(ctx context.Context, req *pbGame.KickAccountOfflineRq, rsp *pbGame.KickAccountOfflineRs) error {
	return h.g.am.KickAccount(ctx, req.GetAccountId(), req.GetGameId())
}
