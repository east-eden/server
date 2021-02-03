package client

import (
	"context"

	pbGlobal "bitbucket.org/east-eden/server/proto/global"
	"bitbucket.org/east-eden/server/transport"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

func (cmd *Commander) CmdQueryHeros(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_QueryHeros",
		Body: &pbGlobal.C2M_QueryHeros{},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "M2C_HeroList"
}

func (cmd *Commander) CmdAddHero(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_AddHero",
		Body: &pbGlobal.C2M_AddHero{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdAddHero command failed")
		return false, ""
	}

	log.Info().Interface("body", msg.Body).Send()
	cmd.c.transport.SendMessage(msg)
	return true, "M2C_HeroList"
}

func (cmd *Commander) CmdDelHero(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_DelHero",
		Body: &pbGlobal.C2M_DelHero{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdDelHero command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "M2C_HeroList"
}
