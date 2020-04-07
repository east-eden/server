package client

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/global"
	"github.com/yokaiio/yokai_server/internal/transport"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
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
		UserID   string `json:"userId"`
		UserName string `json:"userName"`
	}

	req.UserID = fmt.Sprintf("%d", userID)
	req.UserName = userName

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("json marshal failed when call CmdAccountLogon:%v", err)
	}

	resp, err := httpPost(bc.ai.tcpCli, header, body)
	if err != nil {
		return fmt.Errorf("http post failed when call CmdAccountLogon:%v", err)
	}

	var metadata map[string]string
	if err := json.Unmarshal(resp, &metadata); err != nil {
		return fmt.Errorf("json unmarshal failed when call CmdAccountLogon:%v", err)
	}

	if len(metadata["publicAddr"]) == 0 {
		return fmt.Errorf("invalid game_addr")
	}

	respUserID, err := strconv.ParseInt(metadata["userId"], 10, 64)
	if err != nil {
		return fmt.Errorf("parser_int user_id failed:%v", err)
	}

	respAccountID, err := strconv.ParseInt(metadata["accountId"], 10, 64)
	if err != nil {
		return fmt.Errorf("parser_int account_id failed:%v", err)
	}

	bc.ai.tcpCli.SetTcpAddress(metadata["publicAddr"])
	bc.ai.tcpCli.SetUserInfo(respUserID, respAccountID, metadata["userName"])
	return bc.ai.tcpCli.Connect()
}

func (bc *BotCommand) BotCmdCreatePlayer() error {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_CreatePlayer",
		Body: &pbGame.C2M_CreatePlayer{Name: bc.ai.userName},
	}

	bc.ai.tcpCli.SendMessage(msg)

	return nil
}

func (bc *BotCommand) BotCmdQueryPlayerInfo() error {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_QueryPlayerInfo",
		Body: &pbGame.C2M_QueryPlayerInfo{},
	}

	bc.ai.tcpCli.SendMessage(msg)
	return nil
}

func (bc *BotCommand) BotCmdChangeExp() error {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_ChangeExp",
		Body: &pbGame.C2M_ChangeExp{AddExp: rand.Int63n(10)},
	}

	bc.ai.tcpCli.SendMessage(msg)
	return nil
}

func (bc *BotCommand) BotCmdChangeLevel() error {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_ChangeLevel",
		Body: &pbGame.C2M_ChangeLevel{AddLevel: rand.Int31n(5)},
	}

	bc.ai.tcpCli.SendMessage(msg)
	return nil
}

func (bc *BotCommand) BotCmdAddHero() error {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_AddHero",
		Body: &pbGame.C2M_AddHero{TypeId: int32(rand.Intn(len(global.DefaultEntries.HeroEntries)) + 1)},
	}

	bc.ai.tcpCli.SendMessage(msg)
	return nil
}

func (bc *BotCommand) BotCmdQueryHeros() error {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_QueryHeros",
		Body: &pbGame.C2M_QueryHeros{},
	}

	bc.ai.tcpCli.SendMessage(msg)
	return nil
}

func (bc *BotCommand) BotCmdAddItem() error {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_AddItem",
		Body: &pbGame.C2M_AddItem{TypeId: int32(rand.Intn(len(global.DefaultEntries.ItemEntries)) + 1)},
	}

	bc.ai.tcpCli.SendMessage(msg)
	return nil
}

func (bc *BotCommand) BotCmdQueryItems() error {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.C2M_QueryItems",
		Body: &pbGame.C2M_QueryItems{},
	}

	bc.ai.tcpCli.SendMessage(msg)
	return nil
}

func (bc *BotCommand) BotCmdAddToken() error {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_game.MC_AddToken",
		Body: &pbGame.MC_AddToken{
			Type:  int32(rand.Intn(len(global.DefaultEntries.TokenEntries))),
			Value: rand.Int31n(10),
		},
	}

	bc.ai.tcpCli.SendMessage(msg)
	return nil
}
