package client

import (
	"context"

	pbCommon "bitbucket.org/funplus/server/proto/global/common"
	"bitbucket.org/funplus/server/transport"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

func (cmd *Commander) initFragmentCommands() {
	cmd.registerCommandPage(&CommandPage{PageID: Cmd_Page_Fragment, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// 返回上页
	cmd.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Fragment, GotoPageID: Cmd_Page_Main, Cb: nil})

	// 1请求碎片信息
	cmd.registerCommand(&Command{Text: "请求碎片信息", PageID: Cmd_Page_Fragment, GotoPageID: -1, Cb: cmd.CmdQueryFragments})

	// 2碎片合成
	cmd.registerCommand(&Command{Text: "碎片合成", PageID: Cmd_Page_Fragment, GotoPageID: -1, InputText: "请输入碎片ID:", DefaultInput: "1", Cb: cmd.CmdFragmentsCompose})
}

func (cmd *Commander) CmdQueryFragments(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_QueryFragments",
		Body: &pbCommon.C2S_QueryFragments{},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_FragmentsList"
}

func (cmd *Commander) CmdFragmentsCompose(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_FragmentsCompose",
		Body: &pbCommon.C2S_FragmentsCompose{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdFragmentsCompose command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_FragmentsUpdate,S2C_HeroInfo"
}
