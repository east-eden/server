package client

import (
	"context"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/transport"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

func (cmd *Commander) CmdCrystalLevelup(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_CrystalLevelup",
		Body: &pbGlobal.C2S_CrystalLevelup{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdCrystalLevelup command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_CrystalAttUpdate"
}

func (cmd *Commander) CmdPutonCrystal(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_PutonCrystal",
		Body: &pbGlobal.C2S_PutonCrystal{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdPutonCrystal command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_CrystalAttUpdate"
}

func (cmd *Commander) CmdTakeoffCrystal(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_TakeoffCrystal",
		Body: &pbGlobal.C2S_TakeoffCrystal{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdTakeoffCrystal command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_CrystalAttUpdate"
}
