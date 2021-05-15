package global

import (
	"bitbucket.org/funplus/server/define"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	pbMail "bitbucket.org/funplus/server/proto/server/mail"
	"bitbucket.org/funplus/server/utils"
)

// 发送爬塔结算邮件
func (g *GlobalController) SendTowerSettleRewardMail(receiverId int64, attachments *define.MailAttachments) (int64, error) {
	req := &pbMail.CreateMailRq{
		ReceiverId:  receiverId,
		Type:        pbGlobal.MailType_System,
		SenderName:  "系统",
		Title:       "爬塔每日结算奖励",
		Content:     "这是爬塔每日结算奖励，请查收",
		Attachments: attachments.GenAttachmentsPB(),
	}

	rsp, err := g.rpcCaller.CallCreateMail(req)
	if !utils.ErrCheck(err, "CallCreateMail failed when GlobalController.SendTowerSettleRewardMail", receiverId, attachments) {
		return -1, err
	}

	return rsp.MailId, nil
}
