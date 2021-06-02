package client

import (
	"context"
	"hash/crc32"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/transport"
	"bitbucket.org/funplus/server/utils"
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
}

func (cmd *Commander) CmdAccountLogon(ctx context.Context, result []string) (bool, string) {
	var req struct {
		UserID string `json:"userId"`
	}

	req.UserID = result[0]

	var gameInfo GameInfo
	gameInfo.UserID = req.UserID
	gameInfo.PublicTcpAddr = "127.0.0.1:8989"

	log.Info().Interface("info", gameInfo).Msg("metadata unmarshaled result")

	if len(gameInfo.PublicTcpAddr) == 0 {
		log.Warn().Msg("invalid game public tcp address")
		return false, ""
	}

	cmd.c.transport.SetGameInfo(&gameInfo)
	cmd.c.transport.SetProtocol("tcp")
	if err := cmd.c.transport.StartConnect(ctx); err != nil {
		log.Warn().Err(err).Msg("tcp connect failed")
	}

	return true, "S2C_AccountLogon"
}

func (cmd *Commander) CmdWebSocketAccountLogon(ctx context.Context, result []string) (bool, string) {
	var req struct {
		UserID string `json:"userId"`
	}

	req.UserID = result[0]

	var gameInfo GameInfo
	gameInfo.UserID = req.UserID
	gameInfo.PublicTcpAddr = "127.0.0.1:8989"

	log.Info().Interface("info", gameInfo).Msg("metadata unmarshaled result")

	if len(gameInfo.PublicWsAddr) == 0 {
		log.Warn().Msg("invalid game public tcp address")
		return false, ""
	}

	cmd.c.transport.SetGameInfo(&gameInfo)
	cmd.c.transport.SetProtocol("ws")
	if err := cmd.c.transport.StartConnect(ctx); err != nil {
		log.Warn().Err(err).Msg("ws connect failed")
	}
	return true, "S2C_AccountLogon"
}

func (cmd *Commander) CmdSendHeartBeat(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_HeartBeat",
		Body: &pbGlobal.C2S_HeartBeat{},
	}

	cmd.c.transport.SendMessage(msg)

	return false, ""
}

func (cmd *Commander) CmdWaitResponseMessage(ctx context.Context, result []string) (bool, string) {
	// inner message
	innerMsg := pbGlobal.C2S_Ping{}
	data, err := proto.Marshal(&innerMsg)
	utils.ErrPrint(err, "marshal proto message failed")

	// send wait response message
	msg := &transport.Message{
		Name: "C2S_WaitResponseMessage",
		Body: &pbGlobal.C2S_WaitResponseMessage{
			MsgId:        1001,
			InnerMsgCrc:  crc32.ChecksumIEEE([]byte("C2S_Ping")),
			InnerMsgData: data,
		},
	}

	cmd.c.transport.SendMessage(msg)

	return false, ""
}

func (cmd *Commander) CmdCliAccountDisconnect(ctx context.Context, result []string) (bool, string) {
	cmd.c.transport.StartDisconnect()
	return false, ""
}

func (cmd *Commander) CmdServerAccountDisconnect(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_AccountDisconnect",
		Body: &pbGlobal.C2S_AccountDisconnect{},
	}

	cmd.c.transport.SendMessage(msg)

	return false, ""
}
