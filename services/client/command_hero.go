package client

import (
	"context"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/transport"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

func (cmd *Commander) initHeroCommands() {
	cmd.registerCommandPage(&CommandPage{PageID: Cmd_Page_Hero, ParentPageID: Cmd_Page_Main, Cmds: make([]*Command, 0)})

	// 返回上页
	cmd.registerCommand(&Command{Text: "返回上页", PageID: Cmd_Page_Hero, GotoPageID: Cmd_Page_Main, Cb: nil})

	// 1查询英雄信息
	cmd.registerCommand(&Command{Text: "查询英雄信息", PageID: Cmd_Page_Hero, GotoPageID: -1, Cb: cmd.CmdQueryHeros})

	// 3删除英雄
	cmd.registerCommand(&Command{Text: "删除英雄", PageID: Cmd_Page_Hero, GotoPageID: -1, InputText: "请输入要删除的英雄ID:", DefaultInput: "1", Cb: cmd.CmdDelHero})

	// 4查询英雄属性
	cmd.registerCommand(&Command{Text: "查询英雄属性", PageID: Cmd_Page_Hero, GotoPageID: -1, InputText: "请输入要查询的英雄ID:", DefaultInput: "1", Cb: cmd.CmdQueryHeroAtt})
}

func (cmd *Commander) CmdQueryHeros(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_QueryHeros",
		Body: &pbGlobal.C2S_QueryHeros{},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_HeroList"
}

func (cmd *Commander) CmdDelHero(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_DelHero",
		Body: &pbGlobal.C2S_DelHero{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdDelHero command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_HeroList"
}

func (cmd *Commander) CmdQueryHeroAtt(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_QueryHeroAtt",
		Body: &pbGlobal.C2S_QueryHeroAtt{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdQueryHeroAtt command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_HeroAttUpdate"
}
