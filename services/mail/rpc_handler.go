package mail

import (
	"context"
	"errors"
	"time"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	pbGame "bitbucket.org/funplus/server/proto/server/game"
	pbMail "bitbucket.org/funplus/server/proto/server/mail"
	"bitbucket.org/funplus/server/utils"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	ErrInvalidGlobalConfig = errors.New("invalid global config")
)

type RpcHandler struct {
	m       *Mail
	gameSrv pbGame.GameService
}

func NewRpcHandler(cli *cli.Context, m *Mail) *RpcHandler {
	h := &RpcHandler{
		m: m,
		gameSrv: pbGame.NewGameService(
			"game",
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

/////////////////////////////////////////////
// rpc receive
/////////////////////////////////////////////

// 创建系统邮件
func (h *RpcHandler) CreateSystemMail(
	ctx context.Context,
	req *pbMail.CreateSystemMailRq,
	rsp *pbMail.CreateMailRs,
) error {

	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		return ErrInvalidGlobalConfig
	}

	newMail := &define.Mail{}
	mailId, err := utils.NextID(define.SnowFlake_Mail)
	if !utils.ErrCheck(err, "NextID failed when RpcHandler.CreateSystemMail", req) {
		return err
	}

	newMail.Id = mailId
	newMail.Type = int32(req.GetType())
	newMail.Date = int32(time.Now().Unix())
	newMail.ExpireDate = int32(time.Now().Add(time.Duration(globalConfig.MailExpireTime) * time.Second).Unix())
	newMail.Title = req.Title
	newMail.Content = req.Content

	err = h.m.manager.CreateMail(req.ReceiverId, newMail)
	if !utils.ErrCheck(err, "CreateMail failed when RpcHandler.CreateSystemMail", req) {
		return err
	}

	return nil
}

// 创建玩家邮件
func (h *RpcHandler) CreatePlayerMail(
	ctx context.Context,
	req *pbMail.CreatePlayerMailRq,
	rsp *pbMail.CreateMailRs,
) error {

	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		return ErrInvalidGlobalConfig
	}

	newMail := &define.Mail{}
	mailId, err := utils.NextID(define.SnowFlake_Mail)
	if !utils.ErrCheck(err, "NextID failed when RpcHandler.CreatePlayerMail", req) {
		return err
	}

	newMail.Id = mailId
	newMail.SenderId = req.GetSenderId()
	newMail.Type = int32(req.GetType())
	newMail.Date = int32(time.Now().Unix())
	newMail.ExpireDate = int32(time.Now().Add(time.Duration(globalConfig.MailExpireTime) * time.Second).Unix())
	newMail.SenderName = req.GetSenderName()
	newMail.Title = req.Title
	newMail.Content = req.Content

	err = h.m.manager.CreateMail(req.ReceiverId, newMail)
	if !utils.ErrCheck(err, "CreateMail failed when RpcHandler.CreatePlayerMail", req) {
		return err
	}

	return nil
}
