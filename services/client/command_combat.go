package client

import (
	"context"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/transport"
)

func (cmd *Commander) initCombatCommands() {
	cmd.registerCommandPage(&CommandPage{PageID: Cmd_Page_Combat, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// 返回上页
	cmd.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Combat, GotoPageID: Cmd_Page_Main, Cb: nil})

	// 1关卡战斗
	cmd.registerCommand(&Command{Text: "普通关卡战斗", PageID: Cmd_Page_Combat, GotoPageID: -1, Cb: cmd.CmdStartStageCombat})
}

func (cmd *Commander) CmdStartStageCombat(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2S_StartStageCombat",
		Body: &pbGlobal.C2S_StartStageCombat{RpcId: 1},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_StartStageCombat"
}
