package client

import (
	"context"

	"bitbucket.org/east-eden/server/excel/auto"
	pbGlobal "bitbucket.org/east-eden/server/proto/global"
	"bitbucket.org/east-eden/server/transport"
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
	registerFn(&pbGlobal.S2C_QueryPlayerInfo{}, h.OnS2C_QueryPlayerInfo)
	registerFn(&pbGlobal.S2C_ExpUpdate{}, h.OnS2C_ExpUpdate)
	registerFn(&pbGlobal.S2C_SyncPlayerInfo{}, h.OnS2C_SyncPlayerInfo)
	registerFn(&pbGlobal.S2C_PublicSyncPlayerInfo{}, h.OnS2C_PublicSyncPlayerInfo)

	registerFn(&pbGlobal.S2C_HeroList{}, h.OnS2C_HeroList)
	registerFn(&pbGlobal.S2C_HeroInfo{}, h.OnS2C_HeroInfo)
	registerFn(&pbGlobal.S2C_HeroAttUpdate{}, h.OnS2C_HeroAttUpdate)

	registerFn(&pbGlobal.S2C_ItemList{}, h.OnS2C_ItemList)
	registerFn(&pbGlobal.S2C_DelItem{}, h.OnS2C_DelItem)
	registerFn(&pbGlobal.S2C_ItemAdd{}, h.OnS2C_ItemAdd)
	registerFn(&pbGlobal.S2C_ItemUpdate{}, h.OnS2C_ItemUpdate)

	registerFn(&pbGlobal.S2C_TokenList{}, h.OnS2C_TokenList)

	registerFn(&pbGlobal.S2C_TalentList{}, h.OnS2C_TalentList)

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
		Int64("角色id", m.GetInfo().GetAccountId()).
		Str("角色名字", m.GetInfo().GetName()).
		Int64("角色经验", m.GetInfo().GetExp()).
		Int32("角色等级", m.GetInfo().GetLevel()).
		Msg("角色创建成功")

	return nil
}

func (h *MsgHandler) OnS2C_QueryPlayerInfo(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_QueryPlayerInfo)
	if m.Info == nil {
		log.Info().Msg("该账号下还没有角色，请先创建一个角色")
		return nil
	}

	log.Info().
		Int64("角色id", m.Info.Id).
		Str("角色名字", m.Info.Name).
		Int64("角色经验", m.Info.Exp).
		Int32("角色等级", m.Info.Level).
		Msg("角色信息")

	return nil
}

func (h *MsgHandler) OnS2C_ExpUpdate(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_ExpUpdate)

	log.Info().
		Int64("当前经验", m.Exp).
		Int32("当前等级", m.Level).
		Msg("角色信息")

	return nil
}

func (h *MsgHandler) OnS2C_SyncPlayerInfo(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	log.Info().Msg("rpc同步玩家信息成功")
	return nil
}

func (h *MsgHandler) OnS2C_PublicSyncPlayerInfo(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	log.Info().Msg("MQ同步玩家信息成功")
	return nil
}

func (h *MsgHandler) OnS2C_HeroList(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_HeroList)

	if len(m.Heros) == 0 {
		log.Info().Msg("未拥有任何英雄，请先添加一个")
		return nil
	}

	log.Info().Msg("拥有英雄：")
	for k, v := range m.Heros {
		entry, ok := auto.GetHeroEntry(v.TypeId)
		if !ok {
			continue
		}

		event := log.Info()
		event.Int64("id", v.Id).
			Int32("type_id", v.TypeId).
			Int64("经验", v.Exp).
			Int32("等级", v.Level).
			Str("名字", entry.Name).
			Msgf("英雄%d", k+1)
	}

	return nil
}

func (h *MsgHandler) OnS2C_HeroInfo(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_HeroInfo)

	entry, _ := auto.GetHeroEntry(m.Info.TypeId)
	log.Info().
		Int64("id", m.Info.Id).
		Int32("TypeID", m.Info.TypeId).
		Int64("经验", m.Info.Exp).
		Int32("等级", m.Info.Level).
		Str("名字", entry.Name).
		Msg("英雄信息")

	return nil
}

func (h *MsgHandler) OnS2C_HeroAttUpdate(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	//m := msg.Body.(*pbGame.S2C_HeroAttUpdate)

	log.Info().Msg("英雄属性更新")
	//logger.WithFields(logger.Fields{
	//"id":     m.Info.Id,
	//"TypeID": m.Info.TypeId,
	//"经验":     m.Info.Exp,
	//"等级":     m.Info.Level,
	//"名字":     entry.Name,
	//}).Info("英雄属性更新：")
	return nil
}

func (h *MsgHandler) OnS2C_ItemList(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_ItemList)

	if len(m.Items) == 0 {
		log.Info().Msg("未拥有任何英雄，请先添加一个")
		return nil
	}

	log.Info().Msg("拥有物品：")
	for k, v := range m.Items {
		entry, ok := auto.GetItemEntry(v.TypeId)
		if !ok {
			continue
		}

		event := log.Info()
		event.Int64("id", v.Id).
			Int32("type_id", v.TypeId).
			Str("name", entry.Name).
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
		Int64("equip_obj", m.Item.EquipObjId).
		Msg("添加了新物品")

	return nil
}

func (h *MsgHandler) OnS2C_ItemUpdate(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_ItemUpdate)
	log.Info().
		Int64("item_id", m.Item.Id).
		Int32("type_id", m.Item.TypeId).
		Int32("item_num", m.Item.Num).
		Int64("equip_obj", m.Item.EquipObjId).
		Msg("物品更新")

	return nil
}

func (h *MsgHandler) OnS2C_TokenList(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_TokenList)

	log.Info().Msg("拥有代币：")
	for k, v := range m.Tokens {
		entry, ok := auto.GetTokenEntry(v.Type)
		if !ok {
			continue
		}

		event := log.Info()
		event.Int32("type", v.Type).
			Int64("value", v.Value).
			Int64("max_hold", v.MaxHold).
			Str("name", entry.Name).
			Msgf("代币%d", k+1)
	}

	return nil
}

func (h *MsgHandler) OnS2C_TalentList(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_TalentList)

	log.Info().Msg("已点击天赋:")
	for k, v := range m.Talents {
		entry, ok := auto.GetTalentEntry(v.Id)
		if !ok {
			continue
		}

		event := log.Info()
		event.Int32("talent_id", v.Id).
			Str("名字", entry.Name).
			Str("描述", entry.Desc).
			Msgf("天赋%d", k+1)
	}

	return nil
}

func (h *MsgHandler) OnS2C_StartStageCombat(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGlobal.S2C_StartStageCombat)

	log.Info().Interface("result", m).Msg("战斗返回结果")
	return nil
}
