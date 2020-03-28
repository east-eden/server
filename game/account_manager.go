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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type AccountManager struct {
	mapAccount map[int64]*player.Account
	mapSocks   map[transport.Socket]*player.Account

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
		mapAccount:        make(map[int64]*player.Account),
		mapSocks:          make(map[transport.Socket]*player.Account),
		accountConnectMax: ctx.Int("account_connect_max"),
	}

	am.ctx, am.cancel = context.WithCancel(ctx)
	am.coll = game.ds.Database().Collection(am.TableName())
	am.cacheLiteAccount = utils.NewCacheLoader(
		ctx,
		am.coll,
		"_id",
		player.NewLiteAccount,
		nil,
	)

	am.migrate()

	logger.Info("AccountManager Init OK ...")

	return am
}

func (am *AccountManager) migrate() {

	// check index
	idx := am.coll.Indexes()

	opts := options.ListIndexes().SetMaxTime(2 * time.Second)
	cursor, err := idx.List(context.Background(), opts)
	if err != nil {
		log.Fatal(err)
	}

	indexExist := false
	for cursor.Next(context.Background()) {
		var result bson.M
		cursor.Decode(&result)
		if result["name"] == "user_id" {
			indexExist = true
			break
		}
	}

	// create index
	if !indexExist {
		_, err := am.coll.Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys:    bsonx.Doc{{"user_id", bsonx.Int32(1)}},
				Options: options.Index().SetName("user_id"),
			},
		)

		if err != nil {
			logger.Warn("collection account create index user_id failed:", err)
		}
	}
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

func (am *AccountManager) save(acct *player.Account) {
	filter := bson.D{{"_id", acct.GetID()}}
	update := bson.D{{"$set", acct}}
	if _, err := am.coll.UpdateOne(am.ctx, filter, update, options.Update().SetUpsert(true)); err != nil {
		logger.Warn("account manager save failed:", err)
	}
}

func (am *AccountManager) addAccount(userID int64, accountID int64, accountName string, sock transport.Socket) (*player.Account, error) {
	if accountID == -1 {
		return nil, errors.New("add account id invalid!")
	}

	if len(am.mapSocks) >= am.accountConnectMax {
		return nil, fmt.Errorf("Reach game server's max account connect num")
	}

	var account *player.Account
	obj := am.cacheLiteAccount.Load(accountID)
	if obj == nil {
		// create new account
		la := player.NewLiteAccount().(*player.LiteAccount)
		la.ID = accountID
		la.UserID = userID
		la.GameID = am.g.ID
		la.Name = accountName

		account = player.NewAccount(am.ctx, la, sock)
		am.save(account)

	} else {
		// exist account logon
		la := obj.(*player.LiteAccount)
		account = player.NewAccount(am.ctx, la, sock)

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

	logger.Info(fmt.Sprintf("add account <user_id:%d, account_id:%d, name:%s, sock:%v> success!", account.UserID, account.ID, account.GetName(), account.GetSock()))

	// account main
	am.waitGroup.Wrap(func() {
		err := account.Main()
		if err != nil {
			logger.Info("account Main() return err:", err)
		}

		am.Lock()
		delete(am.mapAccount, account.GetID())
		delete(am.mapSocks, account.GetSock())
		am.Unlock()

		account.Exit()
	})

	// cache store
	am.cacheLiteAccount.Store(account.LiteAccount)

	return account, nil
}

func (am *AccountManager) AccountLogon(userID int64, accountID int64, accountName string, sock transport.Socket) (*player.Account, error) {
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
	return am.addAccount(userID, accountID, accountName, sock)
}

func (am *AccountManager) GetLiteAccount(acctID int64) *player.LiteAccount {
	cacheObj := am.cacheLiteAccount.Load(acctID)
	if cacheObj != nil {
		return cacheObj.(*player.LiteAccount)
	}

	return nil
}

func (am *AccountManager) GetAccountByID(acctID int64) *player.Account {
	am.RLock()
	account, ok := am.mapAccount[acctID]
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

	account.Cancel()

	logger.WithFields(logger.Fields{
		"id":     account.GetID(),
		"reason": reason,
	}).Warn("Account disconnected!")
}

func (am *AccountManager) CreatePlayer(c *player.Account, name string) (*player.Player, error) {
	// only can create one player
	if c.GetPlayer() != nil {
		return nil, fmt.Errorf("only can create one player")
	}

	if len(c.GetPlayerIDs()) > 0 {
		return nil, fmt.Errorf("already create one player before")
	}

	p, err := am.g.pm.CreatePlayer(c, name)
	if err != nil {
		return nil, err
	}

	c.SetPlayer(p)
	c.Name = name
	c.Level = p.GetLevel()
	c.AddPlayerID(p.GetID())
	am.save(c)

	// update account info
	am.g.rpcHandler.CallUpdateUserInfo(c)

	return p, err
}

func (am *AccountManager) SelectPlayer(c *player.Account, id int64) (*player.Player, error) {
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
