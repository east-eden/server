package client

import (
	"context"

	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/internal/utils"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
)

type BotAI struct {
	userID      int64
	userName    string
	accountID   int64
	playerID    int64
	playerName  string
	playerLevel int32

	bc     *BotCommand
	ctx    context.Context
	cancel context.CancelFunc

	tcpCli    *TcpClient
	waitGroup utils.WaitGroupWrapper
}

func NewBotAI(ctx *cli.Context, userID int64, userName string) *BotAI {
	ai := &BotAI{
		userID:   userID,
		userName: userName,
	}

	ai.ctx, ai.cancel = context.WithCancel(ctx)
	ai.bc = NewBotCommand(ai.ctx, ai)
	ai.tcpCli = NewTcpClient(ctx, ai)

	return ai
}

func (ai *BotAI) Run() error {

	// first logon
	if err := ai.bc.BotCmdAccountLogon(ai.userID, ai.userName); err != nil {
		return err
	}

	// random message
	return ai.bc.BotCmdCreatePlayer()
}

func (ai *BotAI) Exit() {
	ai.cancel()
	ai.waitGroup.Wait()
}

func (ai *BotAI) accountLogon(m *pbAccount.MS_AccountLogon) {
	ai.accountID = m.AccountId
	ai.playerID = m.PlayerId
	ai.playerName = m.PlayerName
	ai.playerLevel = m.PlayerLevel
}
