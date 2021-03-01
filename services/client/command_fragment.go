package client

import (
	"context"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/transport"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

func (cmd *Commander) CmdQueryFragments(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_QueryFragments",
		Body: &pbGlobal.C2S_QueryFragments{},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_FragmentsList"
}

func (cmd *Commander) CmdFragmentsCompose(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_FragmentsCompose",
		Body: &pbGlobal.C2S_FragmentsCompose{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdFragmentsCompose command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_FragmentsUpdate,S2C_HeroInfo"
}
