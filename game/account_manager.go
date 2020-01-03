package game

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/game/player"
	"github.com/yokaiio/yokai_server/internal/transport"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type AccountManager struct {
	mapAccount sync.Map
	mapSocks   sync.Map
	g          *Game
	waitGroup  utils.WaitGroupWrapper
	ctx        context.Context
	cancel     context.CancelFunc

	accountConnectMax int
	accountTimeout    time.Duration
}

func NewAccountManager(game *Game, ctx *cli.Context) *AccountManager {
	am := &AccountManager{
		g:                 game,
		accountConnectMax: ctx.Int("account_connect_max"),
		accountTimeout:    ctx.Duration("account_timeout"),
	}

	am.ctx, am.cancel = context.WithCancel(ctx)

	logger.Info("AccountManager Init OK ...")

	return am
}

func (am *AccountManager) Main() error {
	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal("AccountManager Main() error:", err)
			}
			exitCh <- err
		})
	}

	am.waitGroup.Wrap(func() {
		exitFunc(am.Run())
	})

	return <-exitCh
}

func (am *AccountManager) Exit() {
	logger.Info("AccountManager context done...")
	am.cancel()
	am.waitGroup.Wait()
}

func (am *AccountManager) addAccount(accID int64, name string, sock transport.Socket) (*Account, error) {
	if accID == -1 {
		return nil, errors.New("add account id invalid!")
	}

	var numSocks uint32
	am.mapSocks.Range(func(_, _ interface{}) bool {
		numSocks++
		return true
	})

	if numSocks >= uint32(am.accountConnectMax) {
		return nil, fmt.Errorf("Reach game server's max account connect num")
	}

	// new account
	info := &AccountInfo{
		ID:   accID,
		Name: name,
		sock: sock,
		p:    nil,
	}

	// rand peek one player from account
	mapPlayer := am.g.pm.GetPlayers(accID)
	if len(mapPlayer) > 0 {
		for _, p := range mapPlayer {
			info.p = p
			break
		}
	}

	account := NewAccount(am, info)
	am.mapAccount.Store(account.ID(), account)
	am.mapSocks.Store(sock, account)
	logger.Info(fmt.Sprintf("add account <id:%d, name:%s, sock:%v> success!", account.ID(), account.Name(), account.Sock()))

	// account main
	am.waitGroup.Wrap(func() {
		err := account.Main()
		if err != nil {
			logger.Info("account Main() return err:", err)
		}
		am.mapSocks.Delete(account.info.sock)
		am.mapAccount.Delete(account.ID())
		account.Exit()

	})

	return account, nil
}

func (am *AccountManager) AccountLogon(accID int64, name string, sock transport.Socket) (*Account, error) {
	account, ok := am.mapAccount.Load(accID)
	if ok {
		// return exist connection sock
		ac := account.(*Account)
		if ac.info.sock == sock {
			return ac, nil
		}

		// disconnect last account sock
		am.DisconnectAccount(ac, "AddAccount")
	}

	return am.addAccount(accID, name, sock)
}

func (am *AccountManager) GetAccountByID(id int64) *Account {
	v, ok := am.mapAccount.Load(id)
	if !ok {
		return nil
	}

	return v.(*Account)
}

func (am *AccountManager) GetAccountBySock(sock transport.Socket) *Account {
	v, ok := am.mapSocks.Load(sock)
	if !ok {
		return nil
	}

	return v.(*Account)
}

func (am *AccountManager) GetAllAccounts() []*Account {
	ret := make([]*Account, 0)
	am.mapAccount.Range(func(k, v interface{}) bool {
		a := v.(*Account)
		ret = append(ret, a)
		return true
	})

	return ret
}

func (am *AccountManager) DisconnectAccount(ac *Account, reason string) {
	if ac == nil {
		return
	}

	am.DisconnectAccountBySock(ac.info.sock, reason)
}

func (am *AccountManager) DisconnectAccountByID(id int64, reason string) {
	am.DisconnectAccount(am.GetAccountByID(id), reason)
}

func (am *AccountManager) DisconnectAccountBySock(sock transport.Socket, reason string) {
	v, ok := am.mapSocks.Load(sock)
	if !ok {
		return
	}

	account, ok := v.(*Account)
	if !ok {
		return
	}

	logger.WithFields(logger.Fields{
		"id":     account.ID(),
		"reason": reason,
	}).Warn("Account disconnected!")

	account.cancel()
}

func (am *AccountManager) CreatePlayer(c *Account, name string) (*player.Player, error) {
	// only can create one player
	if c.info.p != nil {
		return nil, fmt.Errorf("only can create one player")
	}

	mapPlayer := am.g.pm.GetPlayers(c.ID())
	if len(mapPlayer) > 0 {
		return nil, fmt.Errorf("already create one player before")
	}

	p, err := am.g.pm.CreatePlayer(c.ID(), name)
	if p != nil {
		c.info.p = p
	}

	return p, err
}

func (am *AccountManager) SelectPlayer(c *Account, id int64) (*player.Player, error) {
	playerList := am.g.pm.GetPlayers(c.ID())
	for _, v := range playerList {
		if v.GetID() == id {
			c.info.p = v
			return v, nil
		}
	}

	return nil, fmt.Errorf("select player with wrong id<%d>", id)
}

func (am *AccountManager) BroadCast(msg proto.Message) {
	am.mapAccount.Range(func(_, v interface{}) bool {
		if account, ok := v.(*Account); ok {
			account.SendProtoMessage(msg)
		}
		return true
	})
}

func (am *AccountManager) Run() error {
	for {
		select {
		case <-am.ctx.Done():
			logger.Print("world session context done!")
			return nil
		}
	}

	return nil
}
