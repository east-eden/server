package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	logger "github.com/sirupsen/logrus"
)

func BotCmdAccountLogon(c *TcpClient, userID int64, userName string) bool {
	logger.WithFields(logger.Fields{
		"user_id":    userID,
		"user_name":  userName,
		"tcp_client": c,
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
		logger.Warn("json marshal failed when call CmdAccountLogon:", err)
		return false
	}

	resp, err := httpPost(c, header, body)
	if err != nil {
		logger.Warn("http post failed when call CmdAccountLogon:", err)
		return false
	}

	var metadata map[string]string
	if err := json.Unmarshal(resp, &metadata); err != nil {
		logger.Warn("json unmarshal failed when call CmdAccountLogon:", err)
		return false
	}

	if len(metadata["public_addr"]) == 0 {
		logger.Warn("invalid game_addr")
		return false
	}

	respUserID, err := strconv.ParseInt(metadata["user_id"], 10, 64)
	if err != nil {
		logger.Warn("parser_int user_id failed:", err)
		return false
	}

	respAccountID, err := strconv.ParseInt(metadata["account_id"], 10, 64)
	if err != nil {
		logger.Warn("parser_int account_id failed:", err)
		return false
	}

	c.SetTcpAddress(metadata["public_addr"])
	c.SetUserInfo(respUserID, respAccountID, metadata["user_name"])
	c.Connect()
	return true
}
