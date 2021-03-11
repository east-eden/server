package client

import (
	"context"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/transport"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

func (cmd *Commander) CmdQueryHeros(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_QueryHeros",
		Body: &pbGlobal.C2S_QueryHeros{},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_HeroList"
}

func (cmd *Commander) CmdDelHero(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_DelHero",
		Body: &pbGlobal.C2S_DelHero{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdDelHero command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_HeroList"
}

func (cmd *Commander) CmdQueryHeroAtt(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		Name: "C2S_QueryHeroAtt",
		Body: &pbGlobal.C2S_QueryHeroAtt{},
	}

	err := reflectIntoMsg(msg.Body.(proto.Message), result)
	if err != nil {
		log.Error().Err(err).Msg("CmdQueryHeroAtt command failed")
		return false, ""
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_HeroAttUpdate"
}
