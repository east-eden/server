package client

import (
	"context"

	pbGlobal "github.com/east-eden/server/proto/global"
	log "github.com/rs/zerolog/log"
)

func (cmd *Commander) initItemCommands() {
	cmd.registerCommandPage(&CommandPage{PageID: Cmd_Page_Item, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// 返回上页
	cmd.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Item, GotoPageID: Cmd_Page_Main, Cb: nil})

	// 3删除物品
	cmd.registerCommand(&Command{Text: "删除物品", PageID: Cmd_Page_Item, GotoPageID: -1, InputText: "请输入要删除的物品ID:", DefaultInput: "1", Cb: cmd.CmdDelItem})

	// 4使用物品
	cmd.registerCommand(&Command{Text: "使用物品", PageID: Cmd_Page_Item, GotoPageID: -1, InputText: "请输入要使用的物品ID:", Cb: cmd.CmdUseItem})
}

func (cmd *Commander) CmdDelItem(ctx context.Context, result []string) (bool, string) {
	msg := &pbGlobal.C2S_DelItem{}

	err := reflectIntoMsg(msg, result)
	if err != nil {
		log.Error().Err(err).Msg("CmdDelItem command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_DelItem"
}

func (cmd *Commander) CmdUseItem(ctx context.Context, result []string) (bool, string) {
	msg := &pbGlobal.C2S_UseItem{}

	err := reflectIntoMsg(msg, result)
	if err != nil {
		log.Error().Err(err).Msg("CmdUseItem command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_DelItem,S2C_ItemUpdate"
}
