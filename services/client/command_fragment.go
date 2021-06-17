package client

import (
	"context"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/transport"
	log "github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func (cmd *Commander) initFragmentCommands() {
	cmd.registerCommandPage(&CommandPage{PageID: Cmd_Page_Fragment, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// 返回上页
	cmd.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Fragment, GotoPageID: Cmd_Page_Main, Cb: nil})

	// 2英雄碎片合成
	cmd.registerCommand(&Command{Text: "英雄碎片合成", PageID: Cmd_Page_Fragment, GotoPageID: -1, InputText: "请输入碎片ID:", DefaultInput: "1", Cb: cmd.CmdHeroFragmentsCompose})

	// 4收集品碎片合成
	cmd.registerCommand(&Command{Text: "收集品碎片合成", PageID: Cmd_Page_Fragment, GotoPageID: -1, InputText: "请输入碎片ID:", DefaultInput: "1", Cb: cmd.CmdCollectionFragmentsCompose})
}

func (cmd *Commander) CmdHeroFragmentsCompose(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_HeroFragmentsCompose",
		Body: &pbGlobal.C2S_HeroFragmentsCompose{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdHeroFragmentsCompose command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_HeroFragmentsUpdate,S2C_HeroInfo"
}

func (cmd *Commander) CmdCollectionFragmentsCompose(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_CollectionFragmentsCompose",
		Body: &pbGlobal.C2S_CollectionFragmentsCompose{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdCollectionFragmentsCompose command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_CollectionFragmentsUpdate"
}
