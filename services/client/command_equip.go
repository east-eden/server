package client

import (
	"context"

	pbCommon "bitbucket.org/funplus/server/proto/global/common"
	"bitbucket.org/funplus/server/transport"
	"bitbucket.org/funplus/server/utils"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

func (cmd *Commander) initEquipCommands() {
	cmd.registerCommandPage(&CommandPage{PageID: Cmd_Page_Equip, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// 返回上页
	cmd.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Equip, GotoPageID: Cmd_Page_Main, Cb: nil})

	// 2穿装备
	cmd.registerCommand(&Command{Text: "穿装备", PageID: Cmd_Page_Equip, GotoPageID: -1, InputText: "请输入英雄ID和物品ID:", DefaultInput: "1,1", Cb: cmd.CmdHeroPutonEquip})

	// 3脱装备
	cmd.registerCommand(&Command{Text: "脱装备", PageID: Cmd_Page_Equip, GotoPageID: -1, InputText: "请输入英雄ID和装备位置索引:", DefaultInput: "1,0", Cb: cmd.CmdHeroTakeoffEquip})

	// 4装备升级
	cmd.registerCommand(&Command{Text: "装备升级", PageID: Cmd_Page_Equip, GotoPageID: -1, InputText: "请输入装备ID:", Cb: cmd.CmdEquipLevelup})
}

func (cmd *Commander) CmdHeroPutonEquip(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_PutonEquip",
		Body: &pbCommon.C2S_PutonEquip{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdHeroPutonEquip command")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_HeroInfo"
}

func (cmd *Commander) CmdHeroTakeoffEquip(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_TakeoffEquip",
		Body: &pbCommon.C2S_TakeoffEquip{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdHeroTakeoffEquip command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_HeroInfo"
}

func (cmd *Commander) CmdEquipLevelup(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_EquipLevelup",
		Body: &pbCommon.C2S_EquipLevelup{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if !utils.ErrCheck(err, "CmdEquipoLevelup failed") {
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_EquipUpdate"
}
