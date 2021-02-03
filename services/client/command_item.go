package client

import (
	"context"

	pbGlobal "bitbucket.org/east-eden/server/proto/global"
	"bitbucket.org/east-eden/server/transport"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

func (cmd *Commander) CmdQueryItems(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_QueryItems",
		Body: &pbGlobal.C2M_QueryItems{},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "M2C_ItemList"
}

func (cmd *Commander) CmdAddItem(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_AddItem",
		Body: &pbGlobal.C2M_AddItem{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdAddItem command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "M2C_ItemUpdate,M2C_ItemAdd"
}

func (cmd *Commander) CmdDelItem(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_DelItem",
		Body: &pbGlobal.C2M_DelItem{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdDelItem command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "M2C_DelItem"
}

func (cmd *Commander) CmdUseItem(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_UseItem",
		Body: &pbGlobal.C2M_UseItem{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdUseItem command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "M2C_DelItem,M2C_ItemUpdate"
}

func (cmd *Commander) CmdHeroPutonEquip(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_PutonEquip",
		Body: &pbGlobal.C2M_PutonEquip{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdHeroPutonEquip command")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "M2C_HeroInfo"
}

func (cmd *Commander) CmdHeroTakeoffEquip(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_TakeoffEquip",
		Body: &pbGlobal.C2M_TakeoffEquip{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdHeroTakeoffEquip command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "M2C_HeroInfo"
}

func (cmd *Commander) CmdQueryTalents(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_QueryTalents",
		Body: &pbGlobal.C2M_QueryTalents{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdQueryTalents command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "M2C_TalentList"
}

func (cmd *Commander) CmdAddTalent(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_AddTalent",
		Body: &pbGlobal.C2M_AddTalent{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdAddTalent command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "M2C_TalentList"
}
