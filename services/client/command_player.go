package client

import (
	"context"

	pbGlobal "bitbucket.org/east-eden/server/proto/global"
	"bitbucket.org/east-eden/server/transport"
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

func (cmd *Commander) CmdChangeExp(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_ChangeExp",
		Body: &pbGlobal.C2S_ChangeExp{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdChangeExp command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_ExpUpdate"
}

func (cmd *Commander) CmdChangeLevel(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_ChangeLevel",
		Body: &pbGlobal.C2S_ChangeLevel{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdChangeLevel command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_ExpUpdate"
}

func (cmd *Commander) CmdSyncPlayerInfo(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_SyncPlayerInfo",
		Body: &pbGlobal.C2S_SyncPlayerInfo{},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_SyncPlayerInfo"
}

func (cmd *Commander) CmdPublicSyncPlayerInfo(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_PublicSyncPlayerInfo",
		Body: &pbGlobal.C2S_PublicSyncPlayerInfo{},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_PublicSyncPlayerInfo"
}
