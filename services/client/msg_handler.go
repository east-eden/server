package client

import (
	"context"

	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/entries"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"github.com/yokaiio/yokai_server/transport"
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

	h.r.RegisterProtobufMessage(&pbAccount.M2C_AccountLogon{}, h.OnM2C_AccountLogon)
	h.r.RegisterProtobufMessage(&pbAccount.M2C_HeartBeat{}, h.OnM2C_HeartBeat)

	h.r.RegisterProtobufMessage(&pbGame.M2C_CreatePlayer{}, h.OnM2C_CreatePlayer)
	h.r.RegisterProtobufMessage(&pbGame.MS_SelectPlayer{}, h.OnMS_SelectPlayer)
	h.r.RegisterProtobufMessage(&pbGame.M2C_QueryPlayerInfo{}, h.OnM2C_QueryPlayerInfo)
	h.r.RegisterProtobufMessage(&pbGame.M2C_ExpUpdate{}, h.OnM2C_ExpUpdate)
	h.r.RegisterProtobufMessage(&pbGame.M2C_SyncPlayerInfo{}, h.OnM2C_SyncPlayerInfo)
	h.r.RegisterProtobufMessage(&pbGame.M2C_PublicSyncPlayerInfo{}, h.OnM2C_PublicSyncPlayerInfo)

	h.r.RegisterProtobufMessage(&pbGame.M2C_HeroList{}, h.OnM2C_HeroList)
	h.r.RegisterProtobufMessage(&pbGame.M2C_HeroInfo{}, h.OnM2C_HeroInfo)
	h.r.RegisterProtobufMessage(&pbGame.M2C_HeroAttUpdate{}, h.OnM2C_HeroAttUpdate)

	h.r.RegisterProtobufMessage(&pbGame.M2C_ItemList{}, h.OnM2C_ItemList)
	h.r.RegisterProtobufMessage(&pbGame.M2C_DelItem{}, h.OnM2C_DelItem)
	h.r.RegisterProtobufMessage(&pbGame.M2C_ItemAdd{}, h.OnM2C_ItemAdd)
	h.r.RegisterProtobufMessage(&pbGame.M2C_ItemUpdate{}, h.OnM2C_ItemUpdate)

	h.r.RegisterProtobufMessage(&pbGame.M2C_TokenList{}, h.OnM2C_TokenList)

	h.r.RegisterProtobufMessage(&pbGame.M2C_TalentList{}, h.OnM2C_TalentList)

	h.r.RegisterProtobufMessage(&pbGame.M2C_StartStageCombat{}, h.OnM2C_StartStageCombat)
}

func (h *MsgHandler) OnM2C_AccountLogon(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbAccount.M2C_AccountLogon)

	log.Info().
		Str("local", sock.Local()).
		Int64("user_id", m.UserId).
		Int64("account_id", m.AccountId).
		Int64("player_id", m.PlayerId).
		Str("player_name", m.PlayerName).
		Int32("player_level", m.PlayerLevel).Msg("账号登录成功")

	return nil
}

func (h *MsgHandler) OnM2C_HeartBeat(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	return nil
}

func (h *MsgHandler) OnM2C_CreatePlayer(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_CreatePlayer)
	if m.Error == 0 {
		log.Info().
			Int64("角色id", m.Info.LiteInfo.Id).
			Str("角色名字", m.Info.LiteInfo.Name).
			Int64("角色经验", m.Info.LiteInfo.Exp).
			Int32("角色等级", m.Info.LiteInfo.Level).
			Int32("角色拥有英雄数量", m.Info.HeroNums).
			Int32("角色拥有物品数量", m.Info.ItemNums).
			Msg("角色创建成功")
	} else {
		log.Info().Int32("error_code", m.Error).Msg("角色创建失败")
	}

	return nil
}

func (h *MsgHandler) OnMS_SelectPlayer(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.MS_SelectPlayer)
	if m.ErrorCode == 0 {
		log.Info().
			Int64("角色id", m.Info.LiteInfo.Id).
			Str("角色名字", m.Info.LiteInfo.Name).
			Int64("角色经验", m.Info.LiteInfo.Exp).
			Int32("角色等级", m.Info.LiteInfo.Level).
			Int32("角色拥有英雄数量", m.Info.HeroNums).
			Int32("角色拥有物品数量", m.Info.ItemNums).
			Msg("使用此角色")
	} else {
		log.Info().Int32("error_code", m.ErrorCode).Msg("选择角色失败")
	}

	return nil
}

func (h *MsgHandler) OnM2C_QueryPlayerInfo(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_QueryPlayerInfo)
	if m.Info == nil {
		log.Info().Msg("该账号下还没有角色，请先创建一个角色")
		return nil
	}

	log.Info().
		Int64("角色id", m.Info.LiteInfo.Id).
		Str("角色名字", m.Info.LiteInfo.Name).
		Int64("角色经验", m.Info.LiteInfo.Exp).
		Int32("角色等级", m.Info.LiteInfo.Level).
		Int32("角色拥有英雄数量", m.Info.HeroNums).
		Int32("角色拥有物品数量", m.Info.ItemNums).
		Msg("角色信息")

	return nil
}

func (h *MsgHandler) OnM2C_ExpUpdate(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_ExpUpdate)

	log.Info().
		Int64("当前经验", m.Exp).
		Int32("当前等级", m.Level).
		Msg("角色信息")

	return nil
}

func (h *MsgHandler) OnM2C_SyncPlayerInfo(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	log.Info().Msg("rpc同步玩家信息成功")
	return nil
}

func (h *MsgHandler) OnM2C_PublicSyncPlayerInfo(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	log.Info().Msg("MQ同步玩家信息成功")
	return nil
}

func (h *MsgHandler) OnM2C_HeroList(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_HeroList)

	if len(m.Heros) == 0 {
		log.Info().Msg("未拥有任何英雄，请先添加一个")
		return nil
	}

	log.Info().Msg("拥有英雄：")
	for k, v := range m.Heros {
		entry := entries.GetHeroEntry(v.TypeId)
		if entry == nil {
			continue
		}

		event := log.Info()
		event.Int64("id", v.Id).
			Int32("type_id", v.TypeId).
			Int64("经验", v.Exp).
			Int32("等级", v.Level).
			Str("名字".entry.Name).
			Msgf("英雄%d", k+1)
	}

	return nil
}

func (h *MsgHandler) OnM2C_HeroInfo(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_HeroInfo)

	entry := entries.GetHeroEntry(m.Info.TypeId)
	log.Info().
		Int64("id", m.Info.Id).
		Int32("TypeID", m.Info.TypeId).
		Int64("经验", m.Info.Exp).
		Int32("等级", m.Info.Level).
		Str("名字", entry.Name).
		Msg("英雄信息")

	return nil
}

func (h *MsgHandler) OnM2C_HeroAttUpdate(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	//m := msg.Body.(*pbGame.M2C_HeroAttUpdate)

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

func (h *MsgHandler) OnM2C_ItemList(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_ItemList)

	if len(m.Items) == 0 {
		log.Info().Msg("未拥有任何英雄，请先添加一个")
		return nil
	}

	log.Info().Msg("拥有物品：")
	for k, v := range m.Items {
		entry := entries.GetItemEntry(v.TypeId)
		if entry == nil {
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

func (h *MsgHandler) OnM2C_DelItem(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_DelItem)
	log.Info().Int64("item_id", m.ItemId).Msg("物品已删除")

	return nil
}

func (h *MsgHandler) OnM2C_ItemAdd(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_ItemAdd)
	log.Info().
		Int64("item_id", m.Item.Id).
		Int32("type_id", m.Item.TypeId).
		Int32("item_num", m.Item.Num).
		Int64("equip_obj", m.Item.EquipObjId).
		Msg("添加了新物品")

	return nil
}

func (h *MsgHandler) OnM2C_ItemUpdate(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_ItemUpdate)
	log.Info().
		Int64("item_id", m.Item.Id).
		Int32("type_id", m.Item.TypeId).
		Int32("item_num", m.Item.Num).
		Int64("equip_obj", m.Item.EquipObjId).
		Msg("物品更新")

	return nil
}

func (h *MsgHandler) OnM2C_TokenList(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_TokenList)

	log.Info().Msg("拥有代币：")
	for k, v := range m.Tokens {
		entry := entries.GetTokenEntry(v.Type)
		if entry == nil {
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

func (h *MsgHandler) OnM2C_TalentList(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_TalentList)

	log.Info().Msg("已点击天赋:")
	for k, v := range m.Talents {
		entry := entries.GetTalentEntry(v.Id)
		if entry == nil {
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

func (h *MsgHandler) OnM2C_StartStageCombat(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_StartStageCombat)

	log.Info().Interface("result", m).Msg("战斗返回结果")
	return nil
}
