package game

import (
	"context"

	pbMail "github.com/east-eden/server/proto/server/mail"
	"github.com/spf13/cast"
)

/////////////////////////////////////////////
// rpc call
/////////////////////////////////////////////

// 创建玩家邮件
func (h *RpcHandler) CallCreateMail(req *pbMail.CreateMailRq) (*pbMail.CreateMailRs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()
	return h.mailSrv.CreateMail(
		ctx,
		req,
		h.consistentHashCallOption(cast.ToString(req.GetReceiverId())),
	)
}

// 查询玩家邮件
func (h *RpcHandler) CallQueryPlayerMails(req *pbMail.QueryPlayerMailsRq) (*pbMail.QueryPlayerMailsRs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()
	return h.mailSrv.QueryPlayerMails(
		ctx,
		req,
		h.consistentHashCallOption(cast.ToString(req.GetOwnerId())),
		h.retries(3),
	)
}

// 读取邮件
func (h *RpcHandler) CallReadMail(req *pbMail.ReadMailRq) (*pbMail.ReadMailRs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()
	return h.mailSrv.ReadMail(
		ctx,
		req,
		h.consistentHashCallOption(cast.ToString(req.GetOwnerId())),
	)
}

// 获取附件
func (h *RpcHandler) CallGainAttachments(req *pbMail.GainAttachmentsRq) (*pbMail.GainAttachmentsRs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()
	return h.mailSrv.GainAttachments(
		ctx,
		req,
		h.consistentHashCallOption(cast.ToString(req.GetOwnerId())),
	)
}

// 删除邮件
func (h *RpcHandler) CallDelMail(req *pbMail.DelMailRq) (*pbMail.DelMailRs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()
	return h.mailSrv.DelMail(
		ctx,
		req,
		h.consistentHashCallOption(cast.ToString(req.GetOwnerId())),
	)
}

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
