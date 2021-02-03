package client

import (
	"context"

	pbGame "bitbucket.org/east-eden/server/proto/server/game"
	"bitbucket.org/east-eden/server/transport"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

func (cmd *Commander) CmdQueryPlayerInfo(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_QueryPlayerInfo",
		Body: &pbGame.C2M_QueryPlayerInfo{},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "M2C_QueryPlayerInfo"
}

func (cmd *Commander) CmdCreatePlayer(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_CreatePlayer",
		Body: &pbGame.C2M_CreatePlayer{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdCreatePlayer command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "M2C_CreatePlayer"
}

func (cmd *Commander) CmdChangeExp(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_ChangeExp",
		Body: &pbGame.C2M_ChangeExp{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdChangeExp command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "M2C_ExpUpdate"
}

func (cmd *Commander) CmdChangeLevel(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_ChangeLevel",
		Body: &pbGame.C2M_ChangeLevel{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdChangeLevel command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "M2C_ExpUpdate"
}

func (cmd *Commander) CmdSyncPlayerInfo(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_SyncPlayerInfo",
		Body: &pbGame.C2M_SyncPlayerInfo{},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "M2C_SyncPlayerInfo"
}

func (cmd *Commander) CmdPublicSyncPlayerInfo(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_PublicSyncPlayerInfo",
		Body: &pbGame.C2M_PublicSyncPlayerInfo{},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "M2C_PublicSyncPlayerInfo"
}
