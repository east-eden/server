package game

import (
	"context"
	"strconv"

	pbMail "bitbucket.org/funplus/server/proto/server/mail"
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
		consistentHashCallOption(h.g.cons, strconv.Itoa(int(req.GetReceiverId()))),
	)
}

// 查询玩家邮件
func (h *RpcHandler) CallQueryPlayerMails(req *pbMail.QueryPlayerMailsRq) (*pbMail.QueryPlayerMailsRs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()
	return h.mailSrv.QueryPlayerMails(
		ctx,
		req,
		consistentHashCallOption(h.g.cons, strconv.Itoa(int(req.GetOwnerId()))),
	)
}

// 读取邮件
func (h *RpcHandler) CallReadMail(req *pbMail.ReadMailRq) (*pbMail.ReadMailRs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()
	return h.mailSrv.ReadMail(
		ctx,
		req,
		consistentHashCallOption(h.g.cons, strconv.Itoa(int(req.GetOwnerId()))),
	)
}

// 获取附件
func (h *RpcHandler) CallGainAttachments(req *pbMail.GainAttachmentsRq) (*pbMail.GainAttachmentsRs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()
	return h.mailSrv.GainAttachments(
		ctx,
		req,
		consistentHashCallOption(h.g.cons, strconv.Itoa(int(req.GetOwnerId()))),
	)
}

// 删除邮件
func (h *RpcHandler) CallDelMail(req *pbMail.DelMailRq) (*pbMail.DelMailRs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()
	return h.mailSrv.DelMail(
		ctx,
		req,
		consistentHashCallOption(h.g.cons, strconv.Itoa(int(req.GetOwnerId()))),
	)
}

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////
