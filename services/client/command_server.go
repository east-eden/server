package client

import (
	"context"
	"hash/crc32"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/utils"
	log "github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func (cmd *Commander) initServerCommands() {
	// page server connection options
	cmd.registerCommandPage(&CommandPage{PageID: Cmd_Page_Server, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// 返回上页
	cmd.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Server, GotoPageID: Cmd_Page_Main, Cb: nil})

	// 1登录
	cmd.registerCommand(&Command{Text: "登录", PageID: Cmd_Page_Server, GotoPageID: -1, InputText: "请输入登录user ID", DefaultInput: "1", Cb: cmd.CmdAccountLogon})

	// websocket连接登录
	cmd.registerCommand(&Command{Text: "websocket登录", PageID: Cmd_Page_Server, GotoPageID: -1, InputText: "请输入登录user ID和名字，以逗号分隔", DefaultInput: "1,dudu", Cb: cmd.CmdWebSocketAccountLogon})

	// 2发送心跳
	cmd.registerCommand(&Command{Text: "发送心跳", PageID: Cmd_Page_Server, GotoPageID: -1, Cb: cmd.CmdSendHeartBeat})

	// 3发送ClientMessage
	cmd.registerCommand(&Command{Text: "发送等待服务器返回消息", PageID: Cmd_Page_Server, GotoPageID: -1, Cb: cmd.CmdWaitResponseMessage})

	// 4客户端断开连接
	cmd.registerCommand(&Command{Text: "客户端断开连接", PageID: Cmd_Page_Server, GotoPageID: -1, Cb: cmd.CmdCliAccountDisconnect})

	// 5服务器断开连接
	cmd.registerCommand(&Command{Text: "服务器断开连接", PageID: Cmd_Page_Server, GotoPageID: -1, Cb: cmd.CmdServerAccountDisconnect})

	// 6重发登陆请求
	cmd.registerCommand(&Command{Text: "重发登陆请求", PageID: Cmd_Page_Server, GotoPageID: -1, Cb: cmd.CmdResendLogon})
}

func (cmd *Commander) CmdAccountLogon(ctx context.Context, result []string) (bool, string) {
	// transfer gate
	var gateInfo GateInfo
	gateInfo.UserID = result[0]
	gateInfo.PublicTcpAddr = cmd.c.GateAddr

	log.Info().Interface("info", gateInfo).Msg("metadata unmarshaled result")

	if len(gateInfo.PublicTcpAddr) == 0 {
		log.Warn().Msg("invalid gate public tcp address")
		return false, ""
	}

	cmd.c.transport.SetGateInfo(&gateInfo)
	cmd.c.transport.SetProtocol("tcp")
	if err := cmd.c.transport.StartConnect(ctx); err != nil {
		log.Warn().Err(err).Msg("tcp connect failed")
	}

	return true, "S2C_AccountLogon"
}

func (cmd *Commander) CmdWebSocketAccountLogon(ctx context.Context, result []string) (bool, string) {
	// transfer gate
	var gateInfo GateInfo
	gateInfo.UserID = result[0]
	gateInfo.PublicTcpAddr = cmd.c.GateAddr

	log.Info().Interface("info", gateInfo).Msg("metadata unmarshaled result")

	if len(gateInfo.PublicWsAddr) == 0 {
		log.Warn().Msg("invalid gate public tcp address")
		return false, ""
	}

	cmd.c.transport.SetGateInfo(&gateInfo)
	cmd.c.transport.SetProtocol("ws")
	if err := cmd.c.transport.StartConnect(ctx); err != nil {
		log.Warn().Err(err).Msg("ws connect failed")
	}
	return true, "S2C_AccountLogon"
}

func (cmd *Commander) CmdSendHeartBeat(ctx context.Context, result []string) (bool, string) {
	msg := &pbGlobal.C2S_HeartBeat{}

	cmd.c.transport.SendMessage(msg)

	return false, ""
}

func (cmd *Commander) CmdWaitResponseMessage(ctx context.Context, result []string) (bool, string) {
	// inner message
	innerMsg := pbGlobal.C2S_Ping{}
	data, err := proto.Marshal(&innerMsg)
	utils.ErrPrint(err, "marshal proto message failed")

	// send wait response message
	msg := &pbGlobal.C2S_WaitResponseMessage{
		MsgId:        1001,
		InnerMsgCrc:  crc32.ChecksumIEEE([]byte("C2S_Ping")),
		InnerMsgData: data,
	}

	cmd.c.transport.SendMessage(msg)

	return false, ""
}

func (cmd *Commander) CmdCliAccountDisconnect(ctx context.Context, result []string) (bool, string) {
	cmd.c.transport.StartDisconnect()
	return false, ""
}

func (cmd *Commander) CmdServerAccountDisconnect(ctx context.Context, result []string) (bool, string) {
	msg := &pbGlobal.C2S_AccountDisconnect{}

	cmd.c.transport.SendMessage(msg)

	return false, ""
}

func (cmd *Commander) CmdResendLogon(ctx context.Context, result []string) (bool, string) {
	cmd.c.transport.sendLogon()
	return false, ""
}
