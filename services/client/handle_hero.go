package client

import (
	"context"

	"e.coding.net/mmstudio/blade/server/define"
	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
	"e.coding.net/mmstudio/blade/server/transport"
	log "github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func (h *MsgHandler) OnS2C_HeroInfo(ctx context.Context, sock transport.Socket, msg proto.Message) error {
	m := msg.(*pbGlobal.S2C_HeroInfo)

	log.Info().
		Int64("id", m.Info.Id).
		Int32("TypeID", m.Info.TypeId).
		Int32("经验", m.Info.Exp).
		Int32("等级", m.Info.Level).
		Msg("英雄信息")

	return nil
}

func (h *MsgHandler) OnS2C_HeroAttUpdate(ctx context.Context, sock transport.Socket, msg proto.Message) error {
	m := msg.(*pbGlobal.S2C_HeroAttUpdate)

	log.Info().Msg("英雄属性更新")
	attValues := m.GetAtts()
	for _, att := range attValues {
		log.Info().Float32(define.AttNames[att.AttType], att.AttValue).Send()
	}
	return nil
}
