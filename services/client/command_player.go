package client

import (
	"context"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/transport"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

func (cmd *Commander) CmdQueryPlayerInfo(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_QueryPlayerInfo",
		Body: &pbGlobal.C2S_QueryPlayerInfo{},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_QueryPlayerInfo"
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
