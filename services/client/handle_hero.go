package client

import (
	"context"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/transport"
	log "github.com/rs/zerolog/log"
)

func (h *MsgHandler) OnS2C_HeroList(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_HeroList)

	if len(m.Heros) == 0 {
		log.Info().Msg("未拥有任何英雄，请先添加一个")
		return nil
	}

	log.Info().Msg("拥有英雄：")
	for k, v := range m.Heros {
		_, ok := auto.GetHeroEntry(v.TypeId)
		if !ok {
			continue
		}

		event := log.Info()
		event.Int64("id", v.Id).
			Int32("type_id", v.TypeId).
			Int32("经验", v.Exp).
			Int32("等级", v.Level).
			Msgf("英雄%d", k+1)
	}

	return nil
}

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
