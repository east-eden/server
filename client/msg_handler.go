package client

import (
	"context"
	"fmt"

	logger "github.com/sirupsen/logrus"
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

	logger.WithFields(logger.Fields{
		"local":        sock.Local(),
		"user_id":      m.UserId,
		"account_id":   m.AccountId,
		"player_id":    m.PlayerId,
		"player_name":  m.PlayerName,
		"player_level": m.PlayerLevel,
	}).Info("帐号登录成功")

	return nil
}

func (h *MsgHandler) OnM2C_HeartBeat(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	return nil
}

func (h *MsgHandler) OnM2C_CreatePlayer(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_CreatePlayer)
	if m.Error == 0 {
		logger.WithFields(logger.Fields{
			"角色id":     m.Info.LiteInfo.Id,
			"角色名字":     m.Info.LiteInfo.Name,
			"角色经验":     m.Info.LiteInfo.Exp,
			"角色等级":     m.Info.LiteInfo.Level,
			"角色拥有英雄数量": m.Info.HeroNums,
			"角色拥有物品数量": m.Info.ItemNums,
		}).Info("角色创建成功：")
	} else {
		logger.Info("角色创建失败，error_code=", m.Error)
	}

	return nil
}

func (h *MsgHandler) OnMS_SelectPlayer(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.MS_SelectPlayer)
	if m.ErrorCode == 0 {
		logger.WithFields(logger.Fields{
			"角色id":     m.Info.LiteInfo.Id,
			"角色名字":     m.Info.LiteInfo.Name,
			"角色经验":     m.Info.LiteInfo.Exp,
			"角色等级":     m.Info.LiteInfo.Level,
			"角色拥有英雄数量": m.Info.HeroNums,
			"角色拥有物品数量": m.Info.ItemNums,
		}).Info("使用此角色：")
	} else {
		logger.Info("选择角色失败，error_code=", m.ErrorCode)
	}

	return nil
}

func (h *MsgHandler) OnM2C_QueryPlayerInfo(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_QueryPlayerInfo)
	if m.Info == nil {
		logger.Info("该账号下还没有角色，请先创建一个角色")
		return nil
	}

	logger.WithFields(logger.Fields{
		"角色id":     m.Info.LiteInfo.Id,
		"角色名字":     m.Info.LiteInfo.Name,
		"角色经验":     m.Info.LiteInfo.Exp,
		"角色等级":     m.Info.LiteInfo.Level,
		"角色拥有英雄数量": m.Info.HeroNums,
		"角色拥有物品数量": m.Info.ItemNums,
	}).Info("角色信息：")

	return nil
}

func (h *MsgHandler) OnM2C_ExpUpdate(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_ExpUpdate)

	logger.WithFields(logger.Fields{
		"当前经验": m.Exp,
		"当前等级": m.Level,
	}).Info("角色信息：")

	return nil
}

func (h *MsgHandler) OnM2C_HeroList(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_HeroList)
	fields := logger.Fields{}

	if len(m.Heros) == 0 {
		logger.Info("未拥有任何英雄，请先添加一个")
		return nil
	}

	logger.Info("拥有英雄：")
	for k, v := range m.Heros {
		fields["id"] = v.Id
		fields["TypeID"] = v.TypeId
		fields["经验"] = v.Exp
		fields["等级"] = v.Level

		entry := entries.GetHeroEntry(v.TypeId)
		if entry != nil {
			fields["名字"] = entry.Name
		}

		logger.WithFields(fields).Info(fmt.Sprintf("英雄%d", k+1))
	}

	return nil
}

func (h *MsgHandler) OnM2C_HeroInfo(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_HeroInfo)

	entry := entries.GetHeroEntry(m.Info.TypeId)
	logger.WithFields(logger.Fields{
		"id":     m.Info.Id,
		"TypeID": m.Info.TypeId,
		"经验":     m.Info.Exp,
		"等级":     m.Info.Level,
		"名字":     entry.Name,
	}).Info("英雄信息：")

	return nil
}

func (h *MsgHandler) OnM2C_HeroAttUpdate(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	//m := msg.Body.(*pbGame.M2C_HeroAttUpdate)

	logger.Info("英雄属性更新")
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
	fields := logger.Fields{}

	if len(m.Items) == 0 {
		logger.Info("未拥有任何英雄，请先添加一个")
		return nil
	}

	logger.Info("拥有物品：")
	for k, v := range m.Items {
		fields["id"] = v.Id
		fields["type_id"] = v.TypeId

		entry := entries.GetItemEntry(v.TypeId)
		if entry != nil {
			fields["name"] = entry.Name
		}
		logger.WithFields(fields).Info(fmt.Sprintf("物品%d", k+1))
	}

	return nil
}

func (h *MsgHandler) OnM2C_DelItem(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_DelItem)
	logger.Info("物品已删除：", m.ItemId)

	return nil
}

func (h *MsgHandler) OnM2C_ItemAdd(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_ItemAdd)
	logger.WithFields(logger.Fields{
		"item_id":   m.Item.Id,
		"type_id":   m.Item.TypeId,
		"item_num":  m.Item.Num,
		"equip_obj": m.Item.EquipObjId,
	}).Info("添加了新物品")

	return nil
}

func (h *MsgHandler) OnM2C_ItemUpdate(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_ItemUpdate)
	logger.WithFields(logger.Fields{
		"item_id":   m.Item.Id,
		"type_id":   m.Item.TypeId,
		"item_num":  m.Item.Num,
		"equip_obj": m.Item.EquipObjId,
	}).Info("物品更新")

	return nil
}

func (h *MsgHandler) OnM2C_TokenList(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_TokenList)
	fields := logger.Fields{}

	logger.Info("拥有代币：")
	for k, v := range m.Tokens {
		fields["type"] = v.Type
		fields["value"] = v.Value
		fields["max_hold"] = v.MaxHold

		entry := entries.GetTokenEntry(v.Type)
		if entry != nil {
			fields["name"] = entry.Name
		}
		logger.WithFields(fields).Info(fmt.Sprintf("代币%d", k+1))
	}

	return nil
}

func (h *MsgHandler) OnM2C_TalentList(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_TalentList)
	fields := logger.Fields{}

	logger.Info("已点击天赋：")
	for k, v := range m.Talents {
		fields["id"] = v.Id

		entry := entries.GetTalentEntry(v.Id)
		if entry != nil {
			fields["名字"] = entry.Name
			fields["描述"] = entry.Desc
		}

		logger.WithFields(fields).Info(fmt.Sprintf("天赋%d", k+1))
	}

	return nil
}

func (h *MsgHandler) OnM2C_StartStageCombat(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
	m := msg.Body.(*pbGame.M2C_StartStageCombat)

	logger.Info("战斗返回结果:", m)
	return nil
}
