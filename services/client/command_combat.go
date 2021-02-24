package client

import (
	"context"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/transport"
)

func (cmd *Commander) CmdStartStageCombat(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2S_StartStageCombat",
		Body: &pbGlobal.C2S_StartStageCombat{RpcId: 1},
	}

	cmd.c.transport.SendMessage(msg)
	return true, "S2C_StartStageCombat"
}
