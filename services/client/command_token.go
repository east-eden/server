package client

import (
	"context"

	pbCommon "bitbucket.org/funplus/server/proto/global/common"
	"bitbucket.org/funplus/server/transport"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

func (cmd *Commander) initTokenCommands() {
	cmd.registerCommandPage(&CommandPage{PageID: Cmd_Page_Token, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// 返回上页
	cmd.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Token, GotoPageID: Cmd_Page_Main, Cb: nil})

	// 1查询代币信息
	cmd.registerCommand(&Command{Text: "查询代币信息", PageID: Cmd_Page_Token, GotoPageID: -1, Cb: cmd.CmdQueryTokens})
}

func (cmd *Commander) CmdQueryTokens(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2S_QueryTokens",
		Body: &pbCommon.C2S_QueryTokens{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdQueryTokens command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_TokenList"
}
