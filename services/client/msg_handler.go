package client

import (
	"context"

	"bitbucket.org/funplus/server/excel/auto"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/transport"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type MsgHandler struct {
	c *Client
	r transport.Register
}

func NewMsgHandler(c *Client, ctx *cli.Context) *MsgHandler {
	h := &MsgHandler{
		c: c,
		r: transport.NewTransportRegister(),
	}

	h.registerMessage()

	return h
}

func (h *MsgHandler) registerMessage() {
	registerFn := func(m proto.Message, mf transport.MessageFunc) {
		err := h.r.RegisterProtobufMessage(m, mf)
		if err != nil {
			log.Fatal().
				Err(err).
				Str("name", string(proto.MessageReflect(m).Descriptor().Name())).
				Msg("register message failed")
		}
	}

	registerFn(&pbGlobal.S2C_Pong{}, h.OnS2C_Pong)
	registerFn(&pbGlobal.S2C_AccountLogon{}, h.OnS2C_AccountLogon)
	registerFn(&pbGlobal.S2C_HeartBeat{}, h.OnS2C_HeartBeat)
	registerFn(&pbGlobal.S2C_WaitResponseMessage{}, h.OnS2C_WaitResponseMessage)

	registerFn(&pbGlobal.S2C_CreatePlayer{}, h.OnS2C_CreatePlayer)
	registerFn(&pbGlobal.S2C_PlayerInitInfo{}, h.OnS2C_PlayerInitInfo)
	registerFn(&pbGlobal.S2C_ExpUpdate{}, h.OnS2C_ExpUpdate)
	registerFn(&pbGlobal.S2C_VipUpdate{}, h.OnS2C_VipUpdate)

	registerFn(&pbGlobal.S2C_HeroList{}, h.OnS2C_HeroList)
	registerFn(&pbGlobal.S2C_HeroInfo{}, h.OnS2C_HeroInfo)
	registerFn(&pbGlobal.S2C_HeroAttUpdate{}, h.OnS2C_HeroAttUpdate)

	registerFn(&pbGlobal.S2C_FragmentsList{}, h.OnS2C_FragmentsList)
	registerFn(&pbGlobal.S2C_FragmentsUpdate{}, h.OnS2C_FragmentsUpdate)

	registerFn(&pbGlobal.S2C_ItemList{}, h.OnS2C_ItemList)
	registerFn(&pbGlobal.S2C_DelItem{}, h.OnS2C_DelItem)
	registerFn(&pbGlobal.S2C_ItemAdd{}, h.OnS2C_ItemAdd)
	registerFn(&pbGlobal.S2C_ItemUpdate{}, h.OnS2C_ItemUpdate)
	registerFn(&pbGlobal.S2C_EquipUpdate{}, h.OnS2C_EquipUpdate)
	registerFn(&pbGlobal.S2C_TestCrystalRandomReport{}, h.OnS2C_TestCrystalRandomReport)

	registerFn(&pbGlobal.S2C_TokenList{}, h.OnS2C_TokenList)
	registerFn(&pbGlobal.S2C_TokenUpdate{}, h.OnS2C_TokenUpdate)

	registerFn(&pbGlobal.S2C_StartStageCombat{}, h.OnS2C_StartStageCombat)
}

func (h *MsgHandler) OnS2C_Pong(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	return nil
}

func (h *MsgHandler) OnS2C_AccountLogon(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_AccountLogon)

	log.Info().
		Str("local", sock.Local()).
		Int64("user_id", m.UserId).
		Int64("account_id", m.AccountId).
		Int64("player_id", m.PlayerId).
		Str("player_name", m.PlayerName).
		Int32("player_level", m.PlayerLevel).Msg("账号登录成功")

	return nil
}

func (h *MsgHandler) OnS2C_HeartBeat(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	return nil
}

func (h *MsgHandler) OnS2C_WaitResponseMessage(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_WaitResponseMessage)
	log.Info().Int32("msg_id", m.MsgId).Int32("err_code", m.ErrCode).Msg("收到解除锁屏消息")
	return nil
}

func (h *MsgHandler) OnS2C_CreatePlayer(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_CreatePlayer)
	m.GetInfo().GetAccountId()
	log.Info().
		Int64("角色id", m.GetInfo().GetId()).
		Str("角色名字", m.GetInfo().GetName()).
		Int32("角色经验", m.GetInfo().GetExp()).
		Int32("角色等级", m.GetInfo().GetLevel()).
		Msg("角色创建成功")

	return nil
}

func (h *MsgHandler) OnS2C_PlayerInitInfo(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_PlayerInitInfo)

	log.Info().
		Interface("角色信息", m.GetInfo()).
		Interface("英雄数据", m.GetHeros()).
		Interface("物品数据", m.GetItems()).
		Interface("装备数据", m.GetEquips()).
		Interface("晶石数据", m.GetCrystals()).
		Interface("碎片数据", m.GetFrags()).
		Interface("章节数据", m.GetChapters()).
		Interface("关卡数据", m.GetStages()).
		Msg("角色上线数据同步")

	return nil
}

func (h *MsgHandler) OnS2C_ExpUpdate(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_ExpUpdate)

	log.Info().
		Int32("当前经验", m.Exp).
		Int32("当前等级", m.Level).
		Msg("角色信息")

	return nil
}

func (h *MsgHandler) OnS2C_VipUpdate(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_VipUpdate)

	log.Info().
		Int32("当前vip经验", m.GetVipExp()).
		Int32("当前vip等级", m.GetVipLevel()).
		Msg("角色信息")

	return nil
}

func (h *MsgHandler) OnS2C_FragmentsList(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_FragmentsList)
	event := log.Info()
	for _, frag := range m.Frags {
		event.Interface("英雄碎片", frag)
	}
	event.Msg("英雄碎片信息")

	return nil
}

func (h *MsgHandler) OnS2C_FragmentsUpdate(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_FragmentsUpdate)
	event := log.Info()
	for _, frag := range m.Frags {
		event.Interface("英雄碎片", frag)
	}
	event.Msg("英雄碎片更新")

	return nil
}

func (h *MsgHandler) OnS2C_ItemList(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_ItemList)

	if len(m.Items) == 0 {
		log.Info().Msg("未拥有任何物品，请先添加一个")
		return nil
	}

	log.Info().Msg("拥有物品：")
	for k, v := range m.Items {
		_, ok := auto.GetItemEntry(v.TypeId)
		if !ok {
			continue
		}

		event := log.Info()
		event.Int64("id", v.Id).
			Int32("type_id", v.TypeId).
			Msgf("物品%d", k+1)
	}

	return nil
}

func (h *MsgHandler) OnS2C_DelItem(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_DelItem)
	log.Info().Int64("item_id", m.ItemId).Msg("物品已删除")

	return nil
}

func (h *MsgHandler) OnS2C_ItemAdd(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_ItemAdd)
	log.Info().
		Int64("item_id", m.Item.Id).
		Int32("type_id", m.Item.TypeId).
		Int32("item_num", m.Item.Num).
		Msg("添加了新物品")

	return nil
}

func (h *MsgHandler) OnS2C_ItemUpdate(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_ItemUpdate)
	log.Info().
		Int64("item_id", m.Item.Id).
		Int32("type_id", m.Item.TypeId).
		Int32("item_num", m.Item.Num).
		Msg("物品更新")

	return nil
}

func (h *MsgHandler) OnS2C_EquipUpdate(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_EquipUpdate)
	log.Info().
		Int64("equip_id", m.EquipId).
		Int32("level", m.EquipData.Level).
		Int32("exp", m.EquipData.Exp).
		Int32("promote", m.EquipData.Promote).
		Bool("lock", m.EquipData.Lock).
		Int64("equip_obj_id", m.EquipData.EquipObj).
		Msg("装备更新")

	return nil
}

func (h *MsgHandler) OnS2C_TestCrystalRandomReport(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_TestCrystalRandomReport)
	for _, report := range m.Report {
		log.Info().Str("report", report).Send()
	}

	return nil
}

func (h *MsgHandler) OnS2C_TokenList(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_TokenList)

	log.Info().Msg("拥有代币：")
	for k, v := range m.Tokens {
		entry, ok := auto.GetTokenEntry(int32(k))
		if !ok {
			continue
		}

		event := log.Info()
		event.Int("type", k).
			Int32("value", v).
			Int32("max_hold", entry.MaxHold).
			Msgf("代币%d", k+1)
	}

	return nil
}

func (h *MsgHandler) OnS2C_TokenUpdate(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_TokenUpdate)

	log.Info().Int32("token_type", m.Type).Int32("token_value", m.Value).Msg("代币更新")
	return nil
}

func (h *MsgHandler) OnS2C_StartStageCombat(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_StartStageCombat)

	log.Info().Interface("result", m).Msg("战斗返回结果")
	return nil
}
