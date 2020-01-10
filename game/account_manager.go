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
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/transport"
	"github.com/yokaiio/yokai_server/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AccountManager struct {
	mapAccount map[int64]*Account
	mapSocks   map[transport.Socket]*Account

	g         *Game
	waitGroup utils.WaitGroupWrapper
	ctx       context.Context
	cancel    context.CancelFunc

	accountConnectMax int
	cacheLiteAccount  *utils.CacheLoader

	coll *mongo.Collection
	sync.RWMutex
}

func NewAccountManager(game *Game, ctx *cli.Context) *AccountManager {
	am := &AccountManager{
		g:                 game,
		mapAccount:        make(map[int64]*Account),
		mapSocks:          make(map[transport.Socket]*Account),
		accountConnectMax: ctx.Int("account_connect_max"),
	}

	am.ctx, am.cancel = context.WithCancel(ctx)
	am.coll = game.ds.Database().Collection(am.TableName())
	am.cacheLiteAccount = utils.NewCacheLoader(
		ctx,
		am.coll,
		"_id",
		define.Account_ExpireChanNum,
		NewLiteAccount,
		nil,
	)

	logger.Info("AccountManager Init OK ...")

	return am
}

func (am *AccountManager) TableName() string {
	return "account"
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

func (am *AccountManager) save(acct *Account) {
	filter := bson.D{{"_id", acct.GetID()}}
	update := bson.D{{"$set", acct}}
	if _, err := am.coll.UpdateOne(am.ctx, filter, update, options.Update().SetUpsert(true)); err != nil {
		logger.Warn("account manager save failed:", err)
	}
}

func (am *AccountManager) addAccount(accID int64, name string, sock transport.Socket) (*Account, error) {
	if accID == -1 {
		return nil, errors.New("add account id invalid!")
	}

	if len(am.mapSocks) >= am.accountConnectMax {
		return nil, fmt.Errorf("Reach game server's max account connect num")
	}

	var account *Account
	obj := am.cacheLiteAccount.Load(accID)
	if obj == nil {
		// create new account
		la := NewLiteAccount().(*LiteAccount)
		la.SetID(accID)
		la.SetName(name)

		account = NewAccount(am.ctx, la, sock)
		am.save(account)

	} else {
		// exist account logon
		la := obj.(*LiteAccount)
		account = NewAccount(am.ctx, la, sock)

		// peek one player from account
		listPlayerID := account.GetPlayerIDs()
		if len(listPlayerID) > 0 {
			account.SetPlayer(am.g.pm.GetPlayer(listPlayerID[0]))
		}
	}

	am.Lock()
	am.mapAccount[account.GetID()] = account
	am.mapSocks[sock] = account
	am.Unlock()

	logger.Info(fmt.Sprintf("add account <id:%d, name:%s, sock:%v> success!", account.GetID(), account.GetName(), account.GetSock()))

	// account main
	am.waitGroup.Wrap(func() {
		err := account.Main()
		if err != nil {
			logger.Info("account Main() return err:", err)
		}

		am.Lock()
		delete(am.mapAccount, account.GetID())
		delete(am.mapSocks, account.sock)
		am.Unlock()

		account.Exit()
	})

	// cache store
	am.cacheLiteAccount.Store(account.LiteAccount)

	return account, nil
}

func (am *AccountManager) AccountLogon(accID int64, name string, sock transport.Socket) (*Account, error) {
	am.RLock()
	account, acctOK := am.mapAccount[accID]
	am.RUnlock()

	// if reconnect with same socket, then do nothing
	if acctOK && account.sock == sock {
		return account, nil
	}

	// if reconnect with another socket, replace socket in account
	if acctOK {
		am.Lock()
		if account.sock != nil {
			delete(am.mapSocks, account.sock)
			account.sock.Close()
		}

		am.mapSocks[sock] = account
		account.sock = sock
		am.Unlock()

		return account, nil
	}

	// add a new account with socket
	return am.addAccount(accID, name, sock)
}

func (am *AccountManager) GetLiteAccount(acctID int64) *LiteAccount {
	cacheObj := am.cacheLiteAccount.Load(acctID)
	if cacheObj != nil {
		return cacheObj.(*LiteAccount)
	}

	return nil
}

func (am *AccountManager) GetAccountByID(acctID int64) *Account {
	am.RLock()
	account, ok := am.mapAccount[acctID]
	am.RUnlock()

	if !ok {
		return nil
	}

	return account
}

func (am *AccountManager) GetAccountBySock(sock transport.Socket) *Account {
	am.RLock()
	account, ok := am.mapSocks[sock]
	am.RUnlock()

	if !ok {
		return nil
	}

	return account
}

func (am *AccountManager) GetAllAccounts() []*Account {
	ret := make([]*Account, 0)

	am.RLock()
	for _, account := range am.mapAccount {
		ret = append(ret, account)
	}
	am.RUnlock()

	return ret
}

func (am *AccountManager) DisconnectAccount(ac *Account, reason string) {
	if ac == nil {
		return
	}

	am.DisconnectAccountBySock(ac.sock, reason)
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

	account.cancel()

	logger.WithFields(logger.Fields{
		"id":     account.GetID(),
		"reason": reason,
	}).Warn("Account disconnected!")
}

func (am *AccountManager) CreatePlayer(c *Account, name string) (*player.Player, error) {
	// only can create one player
	if c.p != nil {
		return nil, fmt.Errorf("only can create one player")
	}

	if len(c.GetPlayerIDs()) > 0 {
		return nil, fmt.Errorf("already create one player before")
	}

	p, err := am.g.pm.CreatePlayer(c.GetID(), name)
	if err != nil {
		return nil, err
	}

	c.SetPlayer(p)
	c.AddPlayerID(p.GetID())
	am.save(c)

	return p, err
}

func (am *AccountManager) SelectPlayer(c *Account, id int64) (*player.Player, error) {
	playerIDs := c.GetPlayerIDs()
	for _, v := range playerIDs {
		if p := am.g.pm.GetPlayer(v); p != nil && v == id {
			c.SetPlayer(p)
			return p, nil
		}
	}

	return nil, fmt.Errorf("select player with wrong id<%d>", id)
}

func (am *AccountManager) BroadCast(msg proto.Message) {
	accounts := am.GetAllAccounts()
	for _, account := range accounts {
		account.SendProtoMessage(msg)
	}
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
