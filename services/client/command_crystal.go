package client

import (
	"context"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/transport"
	log "github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

// 注册选项
func (cmd *Commander) initCrystalCommands() {
	cmd.registerCommandPage(&CommandPage{PageID: Cmd_Page_Crystal, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// 返回上页
	cmd.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Crystal, GotoPageID: Cmd_Page_Main, Cb: nil})

	// 晶石升级
	cmd.registerCommand(&Command{Text: "晶石升级", PageID: Cmd_Page_Crystal, GotoPageID: -1, InputText: "请输入晶石ID:", Cb: cmd.CmdCrystalLevelup})

	// 装备晶石
	cmd.registerCommand(&Command{Text: "装备晶石", PageID: Cmd_Page_Crystal, GotoPageID: -1, InputText: "请输入英雄ID和晶石ID:", DefaultInput: "1,1", Cb: cmd.CmdPutonCrystal})

	// 卸下晶石
	cmd.registerCommand(&Command{Text: "卸下晶石", PageID: Cmd_Page_Crystal, GotoPageID: -1, InputText: "请输入英雄ID和位置:", DefaultInput: "1,0", Cb: cmd.CmdTakeoffCrystal})

	// 批量测试晶石属性
	cmd.registerCommand(&Command{Text: "批量测试晶石属性", PageID: Cmd_Page_Crystal, GotoPageID: -1, InputText: "输入批量生成晶石数量:", DefaultInput: "10000", Cb: cmd.CmdCrystalBulkRandom})
}

func (cmd *Commander) CmdCrystalLevelup(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_CrystalLevelup",
		Body: &pbGlobal.C2S_CrystalLevelup{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdCrystalLevelup command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_CrystalAttUpdate"
}

func (cmd *Commander) CmdPutonCrystal(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_PutonCrystal",
		Body: &pbGlobal.C2S_PutonCrystal{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdPutonCrystal command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_CrystalAttUpdate"
}

func (cmd *Commander) CmdTakeoffCrystal(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_TakeoffCrystal",
		Body: &pbGlobal.C2S_TakeoffCrystal{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdTakeoffCrystal command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_CrystalAttUpdate"
}

func (cmd *Commander) CmdCrystalBulkRandom(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_TestCrystalRandom",
		Body: &pbGlobal.C2S_TestCrystalRandom{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdTakeoffCrystal command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_TestCrystalReport"
}
