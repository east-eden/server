package client

import (
	"context"

	pbGlobal "bitbucket.org/east-eden/server/proto/global"
	"bitbucket.org/east-eden/server/transport"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

func (cmd *Commander) CmdQueryTokens(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_QueryTokens",
		Body: &pbGlobal.C2M_QueryTokens{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdQueryTokens command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "M2C_TokenList"
}

func (cmd *Commander) CmdAddToken(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_AddToken",
		Body: &pbGlobal.C2M_AddToken{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdAddToken command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "M2C_TokenList"
}
