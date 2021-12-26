package comment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"e.coding.net/mmstudio/blade/server/define"
	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
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
func (h *RpcHandler) CallKickCommentTopicData(topic define.CommentTopic, nodeId int32) (*pbComment.KickCommentTopicDataRs, error) {
	if !topic.Valid() {
		return nil, ErrInvalidComment
	}

	if nodeId == int32(h.m.ID) {
		return nil, errors.New("same comment node id")
	}

	req := &pbComment.KickCommentTopicDataRq{
		Topic:              topic.ToPB(),
		CommentTopicNodeId: nodeId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()

	return h.commentSrv.KickCommentTopicData(
		ctx,
		req,
		client.WithSelectOption(
			utils.SpecificIDSelector(
				fmt.Sprintf("comment-%d", nodeId),
			),
		),
	)
}

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
// 查询单个话题评论
func (h *RpcHandler) QueryCommentTopic(
	ctx context.Context,
	req *pbComment.QueryCommentTopicRq,
	rsp *pbComment.QueryCommentTopicRs,
) error {
	var topic define.CommentTopic
	topic.FromPB(req.GetTopic())
	metadatas, err := h.m.manager.QueryCommentTopic(ctx, topic)
	if utils.ErrCheck(err, "QueryCommentTopic failed when RpcHandler.QueryCommentTopic") {
		return err
	}

	rsp.Metadatas = make([]*pbGlobal.CommentMetadata, 0, len(metadatas))
	for _, v := range metadatas {
		rsp.Metadatas = append(rsp.Metadatas, v.ToPB())
	}
	return err
}

// 查询一定数量的单个话题评论
func (h *RpcHandler) QueryCommentTopicRange(
	ctx context.Context,
	req *pbComment.QueryCommentTopicRangeRq,
	rsp *pbComment.QueryCommentTopicRangeRs,
) error {
	var topic define.CommentTopic
	topic.FromPB(req.GetTopic())
	metadatas, err := h.m.manager.QueryCommentTopicRange(ctx, topic, req.GetStart(), req.GetEnd())
	rsp.Metadatas = make([]*pbGlobal.CommentMetadata, 0, len(metadatas))
	for _, metadata := range metadatas {
		rsp.Metadatas = append(rsp.Metadatas, metadata.ToPB())
	}
	return err
}

// 设置评论赞数
func (h *RpcHandler) ModCommentThumbs(
	ctx context.Context,
	req *pbComment.ModCommentThumbsRq,
	rsp *pbComment.ModCommentThumbsRs,
) error {
	var topic define.CommentTopic
	topic.FromPB(req.GetTopic())
	err := h.m.manager.ModCommentThumbs(ctx, topic, req.GetCommentId(), req.GetModThumbs())
	rsp.Error = err.Error()
	return err
}

// 踢出评论cache
func (h *RpcHandler) KickCommentTopicData(
	ctx context.Context,
	req *pbComment.KickCommentTopicDataRq,
	rsp *pbComment.KickCommentTopicDataRs,
) error {
	var topic define.CommentTopic
	topic.FromPB(req.GetTopic())
	err := h.m.manager.KickCommentTopicData(topic, req.GetCommentTopicNodeId())
	rsp.Error = err.Error()
	return err
}
