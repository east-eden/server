package client

import (
	"context"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/transport"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

func (cmd *Commander) initRoleCommands() {
	cmd.registerCommandPage(&CommandPage{PageID: Cmd_Page_Role, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// 返回上页
	cmd.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Role, GotoPageID: Cmd_Page_Main, Cb: nil})

	// 创建角色
	cmd.registerCommand(&Command{Text: "创建角色", PageID: Cmd_Page_Role, GotoPageID: -1, InputText: "请输入rpcid和角色名字:", DefaultInput: "加百列", Cb: cmd.CmdCreatePlayer})

	// 关卡扫荡
	cmd.registerCommand(&Command{Text: "扫荡关卡", PageID: Cmd_Page_Role, GotoPageID: -1, InputText: "请输入关卡id:", DefaultInput: "1", Cb: cmd.CmdStageSweep})

	// gm命令
	cmd.registerCommand(&Command{Text: "gm命令", PageID: Cmd_Page_Role, GotoPageID: -1, InputText: "请输入gm命令", DefaultInput: "gm player exp 100", Cb: cmd.CmdGmCmd})
}

func (cmd *Commander) CmdCreatePlayer(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_CreatePlayer",
		Body: &pbGlobal.C2S_CreatePlayer{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdCreatePlayer command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_CreatePlayer"
}

func (cmd *Commander) CmdStageSweep(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_StageSweep",
		Body: &pbGlobal.C2S_StageSweep{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdStageSweep command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_StageUpdate"
}

func (cmd *Commander) CmdGmCmd(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_GmCmd",
		Body: &pbGlobal.C2S_GmCmd{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdGmCmd command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return false, ""
}
