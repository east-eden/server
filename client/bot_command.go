package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	logger "github.com/sirupsen/logrus"
)

type BotCommand struct {
	ai     *BotAI
	ctx    context.Context
	cancel context.CancelFunc
}

func NewBotCommand(ctx context.Context, ai *BotAI) *BotCommand {
	bc := &BotCommand{ai: ai}
	bc.ctx, bc.cancel = context.WithCancel(ctx)
	return bc
}

func (bc *BotCommand) BotCmdAccountLogon(userID int64, userName string) error {
	logger.WithFields(logger.Fields{
		"user_id":    userID,
		"user_name":  userName,
		"tcp_client": bc.ai.tcpCli,
	}).Info("call BotCmdAccountLogon")

	header := map[string]string{
		"Content-Type": "application/json",
	}

	var req struct {
		UserID   string `json:"user_id"`
		UserName string `json:"user_name"`
	}

	req.UserID = fmt.Sprintf("%d", userID)
	req.UserName = userName

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("json marshal failed when call CmdAccountLogon:", err)
	}

	resp, err := httpPost(bc.ai.tcpCli, header, body)
	if err != nil {
		return fmt.Errorf("http post failed when call CmdAccountLogon:", err)
	}

	var metadata map[string]string
	if err := json.Unmarshal(resp, &metadata); err != nil {
		return fmt.Errorf("json unmarshal failed when call CmdAccountLogon:", err)
	}

	if len(metadata["public_addr"]) == 0 {
		return fmt.Errorf("invalid game_addr")
	}

	respUserID, err := strconv.ParseInt(metadata["user_id"], 10, 64)
	if err != nil {
		return fmt.Errorf("parser_int user_id failed:", err)
	}

	respAccountID, err := strconv.ParseInt(metadata["account_id"], 10, 64)
	if err != nil {
		return fmt.Errorf("parser_int account_id failed:", err)
	}

	bc.ai.tcpCli.SetTcpAddress(metadata["public_addr"])
	bc.ai.tcpCli.SetUserInfo(respUserID, respAccountID, metadata["user_name"])
	return bc.ai.tcpCli.Connect()
}

func (bc *BotCommand) BotCmdCreatePlayer() error {
	return nil
}
