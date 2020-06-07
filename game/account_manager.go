package game

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/golang/protobuf/proto"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/game/player"
	"github.com/yokaiio/yokai_server/store"
	"github.com/yokaiio/yokai_server/transport"
	"github.com/yokaiio/yokai_server/utils"
)

type AccountManager struct {
	mapAccount map[int64]*player.Account
	mapSocks   map[transport.Socket]*player.Account

	g  *Game
	wg utils.WaitGroupWrapper

	accountConnectMax int

	accountPool     sync.Pool
	liteAccountPool sync.Pool

	sync.RWMutex
}

func NewAccountManager(g *Game, ctx *cli.Context) *AccountManager {
	am := &AccountManager{
		g:                 g,
		mapAccount:        make(map[int64]*player.Account),
		mapSocks:          make(map[transport.Socket]*player.Account),
		accountConnectMax: ctx.Int("account_connect_max"),
	}

	am.accountPool.New = player.NewAccount
	am.liteAccountPool.New = player.NewLiteAccount

	// init account memory
	if err := g.store.AddMemExpire(ctx, store.StoreType_Account, &am.accountPool, player.Account_MemExpire); err != nil {
		logger.Warning("store add account memory expire failed:", err)
	}

	// init account memory
	if err := g.store.AddMemExpire(ctx, store.StoreType_LiteAccount, &am.liteAccountPool, player.Account_MemExpire); err != nil {
		logger.Warning("store add lite account memory expire failed:", err)
	}

	// migrate users table
	if err := g.store.MigrateDbTable("account", "user_id"); err != nil {
		logger.Warning("migrate collection account failed:", err)
	}

	logger.Info("AccountManager Init OK ...")
	return am
}

func (am *AccountManager) TableName() string {
	return "account"
}

func (am *AccountManager) Main(ctx context.Context) error {
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

	am.wg.Wrap(func() {
		exitFunc(am.Run(ctx))
	})

	return <-exitCh
}

func (am *AccountManager) Exit() {
	am.wg.Wait()
	logger.Info("account manager exit...")
}

func (am *AccountManager) addAccount(ctx context.Context, userId int64, accountId int64, accountName string, sock transport.Socket) (*player.Account, error) {
	if accountId == -1 {
		return nil, errors.New("add account id invalid!")
	}

	if len(am.mapSocks) >= am.accountConnectMax {
		return nil, fmt.Errorf("Reach game server's max account connect num")
	}

	acct := am.accountPool.Get().(*player.Account)
	err := am.g.store.LoadObjectFromCacheAndDB(store.StoreType_Account, "_id", accountId, acct)
	if err != nil {
		// store cannot load account, create a new account
		acct.ID = accountId
		acct.UserId = userId
		acct.GameId = am.g.ID
		acct.Name = accountName

		// save object
		if err := am.g.store.SaveObjectToCacheAndDB(store.StoreType_Account, acct); err != nil {
			logger.WithFields(logger.Fields{
				"account_id": accountId,
				"user_id":    userId,
			}).Warn("save account failed")
		}

	}

	acct.SetSock(sock)

	// peek one player from account
	if p := am.g.pm.GetPlayerByAccount(acct); p != nil {
		p.SetAccount(acct)
	}

	am.Lock()
	am.mapAccount[acct.GetID()] = acct
	am.mapSocks[sock] = acct
	am.Unlock()

	logger.WithFields(logger.Fields{
		"user_id":    acct.UserId,
		"account_id": acct.ID,
		"name":       acct.GetName(),
		"socket":     acct.GetSock(),
	}).Info("add account success")

	// account main
	am.wg.Wrap(func() {
		err := acct.Main(ctx)
		if err != nil {
			logger.Info("account Main() return err:", err)
		}

		am.Lock()
		delete(am.mapAccount, acct.GetID())
		delete(am.mapSocks, acct.GetSock())
		am.Unlock()

		acct.Exit()
		am.accountPool.Put(acct)
	})

	return acct, nil
}

func (am *AccountManager) AccountLogon(ctx context.Context, userID int64, accountID int64, accountName string, sock transport.Socket) (*player.Account, error) {
	am.RLock()
	account, acctOK := am.mapAccount[accountID]
	am.RUnlock()

	// if reconnect with same socket, then do nothing
	if acctOK && account.GetSock() == sock {
		return account, nil
	}

	// if reconnect with another socket, replace socket in account
	if acctOK {
		am.Lock()
		if account.GetSock() != nil {
			delete(am.mapSocks, account.GetSock())
			account.GetSock().Close()
		}

		am.mapSocks[sock] = account
		account.SetSock(sock)
		am.Unlock()

		return account, nil
	}

	// add a new account with socket
	return am.addAccount(ctx, userID, accountID, accountName, sock)
}

func (am *AccountManager) GetLiteAccount(acctId int64) (*player.LiteAccount, error) {
	x, err := am.g.store.LoadObject(store.StoreType_LiteAccount, "_id", acctId)
	if err != nil {
		return nil, err
	}

	return x.(*player.LiteAccount), nil
}

func (am *AccountManager) GetAccountByID(acctId int64) *player.Account {
	am.RLock()
	account, ok := am.mapAccount[acctId]
	am.RUnlock()

	if !ok {
		return nil
	}

	return account
}

func (am *AccountManager) GetAccountBySock(sock transport.Socket) *player.Account {
	am.RLock()
	account, ok := am.mapSocks[sock]
	am.RUnlock()

	if !ok {
		return nil
	}

	return account
}

func (am *AccountManager) GetAllAccounts() []*player.Account {
	ret := make([]*player.Account, 0)

	am.RLock()
	for _, account := range am.mapAccount {
		ret = append(ret, account)
	}
	am.RUnlock()

	return ret
}

func (am *AccountManager) DisconnectAccount(ac *player.Account, reason string) {
	if ac == nil {
		return
	}

	am.DisconnectAccountBySock(ac.GetSock(), reason)
}

func (am *AccountManager) DisconnectAccountByID(id int64, reason string) {
	am.DisconnectAccount(am.GetAccountByID(id), reason)
}

func (am *AccountManager) DisconnectAccountBySock(sock transport.Socket, reason string) {
	am.RLock()
	account, ok := am.mapSocks[sock]
	am.RUnlock()

	if !ok {
		return
	}

	sock.Close()

	logger.WithFields(logger.Fields{
		"id":     account.GetID(),
		"reason": reason,
	}).Warn("Account disconnected!")
}

func (am *AccountManager) CreatePlayer(acct *player.Account, name string) (*player.Player, error) {
	// only can create one player
	if am.g.pm.GetPlayerByAccount(acct) != nil {
		return nil, fmt.Errorf("only can create one player")
	}

	p, err := am.g.pm.CreatePlayer(acct, name)
	if err != nil {
		return nil, err
	}

	acct.SetPlayer(p)
	acct.Name = name
	acct.Level = p.GetLevel()
	acct.AddPlayerID(p.GetID())
	if err := am.g.store.SaveObjectToCacheAndDB(store.StoreType_Account, acct); err != nil {
		logger.WithFields(logger.Fields{
			"account_id": acct.ID,
			"user_id":    acct.UserId,
		}).Warn("save account failed")
	}

	// update account info
	am.g.rpcHandler.CallUpdateUserInfo(acct)

	return p, err
}

func (am *AccountManager) SelectPlayer(acct *player.Account, id int64) (*player.Player, error) {
	if pl := am.g.pm.GetPlayerByAccount(acct); pl != nil {
		return pl, nil
	}

	return nil, fmt.Errorf("select player with wrong id<%d>", id)
}

func (am *AccountManager) BroadCast(msg proto.Message) {
	accounts := am.GetAllAccounts()
	for _, account := range accounts {
		acct := account
		m := msg
		acct.PushAsyncHandler(func() {
			acct.SendProtoMessage(m)
		})
	}
}

func (am *AccountManager) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			logger.Print("world session context done!")
			return nil

		}
	}

	return nil
}
