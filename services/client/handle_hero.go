package client

import (
	"context"

	"github.com/east-eden/server/define"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/transport"
	log "github.com/rs/zerolog/log"
)

func (h *MsgHandler) OnS2C_HeroInfo(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_HeroInfo)

	log.Info().
		Int64("id", m.Info.Id).
		Int32("TypeID", m.Info.TypeId).
		Int32("经验", m.Info.Exp).
		Int32("等级", m.Info.Level).
		Msg("英雄信息")

	return nil
}

func (h *MsgHandler) OnS2C_HeroAttUpdate(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_HeroAttUpdate)

	log.Info().Msg("英雄属性更新")
	attValues := m.GetAttValue()
	for n := range attValues {
		log.Info().Int32(define.AttNames[n], attValues[n]).Send()
	}
	return nil
}
