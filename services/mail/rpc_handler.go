package mail

import (
	"context"
	"errors"
	"fmt"
	"time"

	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/excel/auto"
	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
	pbMail "e.coding.net/mmstudio/blade/server/proto/server/mail"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/micro/go-micro/v2/client"
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
	m       *Mail
	mailSrv pbMail.MailService
}

func NewRpcHandler(cli *cli.Context, m *Mail) *RpcHandler {
	h := &RpcHandler{
		m: m,
		mailSrv: pbMail.NewMailService(
			"mail",
			m.mi.srv.Client(),
		),
	}

	err := pbMail.RegisterMailServiceHandler(m.mi.srv.Server(), h)
	if err != nil {
		log.Fatal().Err(err).Msg("RegisterMailServiceHandler failed")
	}

	return h
}

/////////////////////////////////////////////
// rpc call
/////////////////////////////////////////////
func (h *RpcHandler) CallKickMailBox(ownerId int64, nodeId int32) (*pbMail.KickMailBoxRs, error) {
	if ownerId == -1 {
		return nil, errors.New("invalid mail box owner id")
	}

	if nodeId == int32(h.m.ID) {
		return nil, errors.New("same mail node id")
	}

	req := &pbMail.KickMailBoxRq{
		OwnerId:    ownerId,
		MailNodeId: nodeId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()

	return h.mailSrv.KickMailBox(
		ctx,
		req,
		client.WithSelectOption(
			utils.SpecificIDSelector(
				fmt.Sprintf("mail-%d", nodeId),
			),
		),
	)
}

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////

// 创建邮件
func (h *RpcHandler) CreateMail(
	ctx context.Context,
	req *pbMail.CreateMailRq,
	rsp *pbMail.CreateMailRs,
) error {

	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		return ErrInvalidGlobalConfig
	}

	newMail := &define.Mail{}
	newMail.Init()
	mailId, err := utils.NextID(define.SnowFlake_Mail)
	if !utils.ErrCheck(err, "NextID failed when RpcHandler.CreateSystemMail", req) {
		return err
	}

	newMail.Id = mailId
	newMail.OwnerId = req.GetReceiverId()
	newMail.SenderId = req.GetSenderId()
	newMail.Type = int32(req.GetType())
	newMail.Date = int32(time.Now().Unix())
	newMail.ExpireDate = int32(time.Now().Add(time.Duration(globalConfig.MailExpireTime) * time.Second).Unix())
	newMail.SenderName = req.GetSenderName()
	newMail.Title = req.GetTitle()
	newMail.Content = req.GetContent()

	attachments := req.GetAttachments()
	if attachments != nil {
		newMail.Attachments = make([]*define.LootData, 0, len(attachments))
		for _, attachment := range attachments {
			newMail.Attachments = append(newMail.Attachments, &define.LootData{
				LootType: int32(attachment.Type),
				LootMisc: attachment.Misc,
				LootNum:  attachment.Num,
			})
		}
	}

	err = h.m.manager.CreateMail(ctx, req.GetReceiverId(), newMail)
	rsp.NewMail = newMail.ToPB()
	if !utils.ErrCheck(err, "CreateMail failed when RpcHandler.CreateSystemMail", req) {
		return err
	}

	return nil
}

// 查询玩家邮件 (不应该被频繁调用)
func (h *RpcHandler) QueryPlayerMails(
	ctx context.Context,
	req *pbMail.QueryPlayerMailsRq,
	rsp *pbMail.QueryPlayerMailsRs,
) error {
	mails, err := h.m.manager.QueryPlayerMails(ctx, req.GetOwnerId())
	if !utils.ErrCheck(err, "QueryPlayerMails failed when RpcHandler.QueryPlayerMails", req) {
		return err
	}

	rsp.Mails = make([]*pbGlobal.Mail, 0, len(mails))
	for _, mail := range mails {
		rsp.Mails = append(rsp.Mails, mail.ToPB())
	}
	return nil
}

// 读取邮件
func (h *RpcHandler) ReadMail(
	ctx context.Context,
	req *pbMail.ReadMailRq,
	rsp *pbMail.ReadMailRs,
) error {
	err := h.m.manager.ReadMail(ctx, req.GetOwnerId(), req.GetMailId())
	if !utils.ErrCheck(err, "ReadMail failed when RpcHandler.ReadMail", req) {
		return err
	}

	rsp.MailId = req.GetMailId()
	rsp.Status = pbGlobal.MailStatus_Readed
	return nil
}

// 获取附件
func (h *RpcHandler) GainAttachments(
	ctx context.Context,
	req *pbMail.GainAttachmentsRq,
	rsp *pbMail.GainAttachmentsRs,
) error {
	err := h.m.manager.GainAttachments(ctx, req.GetOwnerId(), req.GetMailId())
	if !utils.ErrCheck(err, "GainAttachments failed when RpcHandler.GainAttachments", req) {
		return err
	}

	rsp.MailId = req.GetMailId()
	return nil
}

// 删除邮件
func (h *RpcHandler) DelMail(
	ctx context.Context,
	req *pbMail.DelMailRq,
	rsp *pbMail.DelMailRs,
) error {
	err := h.m.manager.DelMail(ctx, req.GetOwnerId(), req.GetMailId())
	if !utils.ErrCheck(err, "DelMail failed when RpcHandler.DelMail", req) {
		return err
	}

	rsp.MailId = req.GetMailId()
	return nil
}

// 踢出邮件cache
func (h *RpcHandler) KickMailBox(
	ctx context.Context,
	req *pbMail.KickMailBoxRq,
	rsp *pbMail.KickMailBoxRs,
) error {
	return h.m.manager.KickMailBox(req.GetOwnerId(), req.GetMailNodeId())
}
