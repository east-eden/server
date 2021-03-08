package client

import (
	"context"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/transport"
	"bitbucket.org/funplus/server/utils"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

func (cmd *Commander) CmdQueryItems(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_QueryItems",
		Body: &pbGlobal.C2S_QueryItems{},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_ItemList"
}

func (cmd *Commander) CmdAddItem(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_AddItem",
		Body: &pbGlobal.C2S_AddItem{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdAddItem command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_ItemUpdate,S2C_ItemAdd"
}

func (cmd *Commander) CmdDelItem(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_DelItem",
		Body: &pbGlobal.C2S_DelItem{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdDelItem command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_DelItem"
}

func (cmd *Commander) CmdUseItem(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_UseItem",
		Body: &pbGlobal.C2S_UseItem{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdUseItem command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_DelItem,S2C_ItemUpdate"
}

func (cmd *Commander) CmdHeroPutonEquip(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_PutonEquip",
		Body: &pbGlobal.C2S_PutonEquip{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdHeroPutonEquip command")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_HeroInfo"
}

func (cmd *Commander) CmdHeroTakeoffEquip(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_TakeoffEquip",
		Body: &pbGlobal.C2S_TakeoffEquip{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdHeroTakeoffEquip command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_HeroInfo"
}

func (cmd *Commander) CmdEquipLevelup(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_EquipLevelup",
		Body: &pbGlobal.C2S_EquipLevelup{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if pass := utils.ErrCheck(err, "CmdEquipoLevelup failed"); !pass {
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_EquipUpdate"
}
