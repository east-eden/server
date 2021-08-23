package comment

import (
	"errors"
	"time"

	pbComment "e.coding.net/mmstudio/blade/server/proto/server/comment"
	pbGame "e.coding.net/mmstudio/blade/server/proto/server/game"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/asim/go-micro/v3/client"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	ErrInvalidGlobalConfig = errors.New("invalid global config")
)

var (
	DefaultRpcTimeout = 5 * time.Second // 默认rpc超时时间
)

type RpcHandler struct {
	m          *Comment
	commentSrv pbComment.CommentService
	gameSrv    pbGame.GameService
}

func NewRpcHandler(cli *cli.Context, m *Comment) *RpcHandler {
	h := &RpcHandler{
		m: m,
		commentSrv: pbComment.NewCommentService(
			"comment",
			m.mi.srv.Client(),
		),
		gameSrv: pbGame.NewGameService(
			"game",
			m.mi.srv.Client(),
		),
	}

	err := pbComment.RegisterCommentServiceHandler(m.mi.srv.Server(), h)
	if err != nil {
		log.Fatal().Err(err).Msg("RegisterCommentServiceHandler failed")
	}

	return h
}

// 一致性哈希
func (h *RpcHandler) consistentHashCallOption(key string) client.CallOption {
	return client.WithSelectOption(
		utils.ConsistentHashSelector(h.m.cons, key),
	)
}

// 重试次数
func (h *RpcHandler) retries(times int) client.CallOption {
	return client.WithRetries(times)
}

/////////////////////////////////////////////
// rpc call
/////////////////////////////////////////////

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
