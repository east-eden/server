package client

import (
	"context"
	"encoding/json"

	pbAccount "bitbucket.org/east-eden/server/proto/account"
	"bitbucket.org/east-eden/server/transport"
	log "github.com/rs/zerolog/log"
)

func (cmd *Commander) CmdAccountLogon(ctx context.Context, result []string) (bool, string) {
	header := map[string]string{
		"Content-Type": "application/json",
	}

	var req struct {
		UserID   string `json:"userId"`
		UserName string `json:"userName"`
	}

	req.UserID = result[0]
	req.UserName = result[1]

	body, err := json.Marshal(req)
	if err != nil {
		log.Warn().Err(err).Msg("json marshal failed when call CmdAccountLogon")
		return false, ""
	}

	resp, err := httpPost(cmd.c.transport.GetGateEndPoints(), header, body)
	if err != nil {
		log.Warn().Err(err).Msg("http post failed when call CmdAccountLogon")
		return false, ""
	}

	var gameInfo GameInfo
	if err := json.Unmarshal(resp, &gameInfo); err != nil {
		log.Warn().Err(err).Msg("json unmarshal failed when call CmdAccountLogon")
		return false, ""
	}

	log.Info().Interface("info", gameInfo).Msg("metadata unmarshaled result")

	if len(gameInfo.PublicTcpAddr) == 0 {
		log.Warn().Msg("invalid game public tcp address")
		return false, ""
	}

	cmd.c.transport.SetGameInfo(&gameInfo)
	cmd.c.transport.SetProtocol("tcp")
	if err := cmd.c.transport.StartConnect(ctx); err != nil {
		log.Warn().Err(err).Msg("tcp connect failed")
	}

	return true, "M2C_AccountLogon"
}

func (cmd *Commander) CmdWebSocketAccountLogon(ctx context.Context, result []string) (bool, string) {
	header := map[string]string{
		"Content-Type": "application/json",
	}

	var req struct {
		UserID   string `json:"userId"`
		UserName string `json:"userName"`
	}

	req.UserID = result[0]
	req.UserName = result[1]

	body, err := json.Marshal(req)
	if err != nil {
		log.Warn().Err(err).Msg("json marshal failed when call CmdWebSocketAccountLogon")
		return false, ""
	}

	resp, err := httpPost(cmd.c.transport.GetGateEndPoints(), header, body)
	if err != nil {
		log.Warn().Err(err).Msg("http post failed when call CmdAccountLogon")
		return false, ""
	}

	var gameInfo GameInfo
	if err := json.Unmarshal(resp, &gameInfo); err != nil {
		log.Warn().Err(err).Msg("json unmarshal failed when call CmdAccountLogon")
		return false, ""
	}

	log.Info().Interface("info", gameInfo).Msg("metadata unmarshaled result")

	if len(gameInfo.PublicWsAddr) == 0 {
		log.Warn().Msg("invalid game public tcp address")
		return false, ""
	}

	cmd.c.transport.SetGameInfo(&gameInfo)
	cmd.c.transport.SetProtocol("ws")
	if err := cmd.c.transport.StartConnect(ctx); err != nil {
		log.Warn().Err(err).Msg("ws connect failed")
	}
	return true, "M2C_AccountLogon"
}

func (cmd *Commander) CmdSendHeartBeat(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_HeartBeat",
		Body: &pbAccount.C2M_HeartBeat{},
	}

	cmd.c.transport.SendMessage(msg)

	return false, ""
}

func (cmd *Commander) CmdCliAccountDisconnect(ctx context.Context, result []string) (bool, string) {
	cmd.c.transport.StartDisconnect()
	return false, ""
}

func (cmd *Commander) CmdServerAccountDisconnect(ctx context.Context, result []string) (bool, string) {
	msg := &transport.Message{
		// Type: transport.BodyProtobuf,
		Name: "C2M_AccountDisconnect",
		Body: &pbAccount.C2M_AccountDisconnect{},
	}

	cmd.c.transport.SendMessage(msg)

	return false, ""
}
