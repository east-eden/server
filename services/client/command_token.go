package client

import (
	"context"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/transport"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

func (cmd *Commander) CmdQueryTokens(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2S_QueryTokens",
		Body: &pbGlobal.C2S_QueryTokens{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdQueryTokens command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_TokenList"
}
