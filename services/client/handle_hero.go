package client

import (
	"context"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/transport"
	log "github.com/rs/zerolog/log"
)

// 属性名
var AttNames = [define.Att_End]string{
	"攻击力",
	"护甲",
	"总伤害加成",
	"暴击值",
	"暴击倍数加成",
	"治疗强度",
	"真实伤害",
	"战场移动速度",
	"时间槽速度",
	"技能效果命中",
	"技能效果抵抗",
	"生命值上限",
	"当前生命值",
	"蓝量上限",
	"当前蓝量",
	"mp恢复值",
	"怒气值",
	"物理系伤害加成",
	"地系伤害加成",
	"水系伤害加成",
	"火系伤害加成",
	"风系伤害加成",
	"时系伤害加成",
	"空系伤害加成",
	"幻系伤害加成",
	"光系伤害加成",
	"暗系伤害加成",
	"物理系伤害抗性",
	"地系伤害抗性",
	"水系伤害抗性",
	"火系伤害抗性",
	"风系伤害抗性",
	"时系伤害抗性",
	"空系伤害抗性",
	"幻系伤害抗性",
	"光系伤害抗性",
	"暗系伤害抗性",
}

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
		log.Info().Int32(AttNames[n], attValues[n]).Send()
	}
	return nil
}
