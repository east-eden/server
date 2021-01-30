package client

import (
	"context"

	pbGame "bitbucket.org/east-eden/server/proto/game"
	"bitbucket.org/east-eden/server/transport"
)

func (cmd *Commander) CmdStartStageCombat(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_StartStageCombat",
		Body: &pbGame.C2M_StartStageCombat{RpcId: 1},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "M2C_StartStageCombat"
}
