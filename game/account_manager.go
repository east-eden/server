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
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/transport"
	"github.com/yokaiio/yokai_server/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountManager struct {
	mapLiteAccount sync.Map
	mapAccount     map[int64]*Account
	mapSocks       map[transport.Socket]*Account

	g         *Game
	waitGroup utils.WaitGroupWrapper
	ctx       context.Context
	cancel    context.CancelFunc

	accountConnectMax int
	accountTimeout    time.Duration
	chExpire          chan int64

	coll *mongo.Collection
	sync.RWMutex
}

func NewAccountManager(game *Game, ctx *cli.Context) *AccountManager {
	am := &AccountManager{
		g:                 game,
		mapAccount:        make(map[int64]*Account),
		mapSocks:          make(map[transport.Socket]*Account),
		accountConnectMax: ctx.Int("account_connect_max"),
		accountTimeout:    ctx.Duration("account_timeout"),
		chExpire:          make(chan int64, define.Account_ExpireChanNum),
	}

	am.ctx, am.cancel = context.WithCancel(ctx)
	am.coll = game.ds.Database().Collection(am.TableName())

	logger.Info("AccountManager Init OK ...")

	return am
}

func (am *AccountManager) TableName() string {
	return "player"
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

	if len(am.mapSocks) >= am.accountConnectMax {
		return nil, fmt.Errorf("Reach game server's max account connect num")
	}

	// new account
	info := &LiteAccount{
		ID:   accID,
		Name: name,
	}

	// rand peek one player from account
	var player *player.Player
	mapPlayer := am.g.pm.GetPlayers(accID)
	if len(mapPlayer) > 0 {
		for _, p := range mapPlayer {
			player = p
			break
		}
	}

	account := NewAccount(am, info, sock, player)
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

	// lite account
	am.mapLiteAccount.Store(account.GetID(), account.LiteAccount)
	am.beginTimeExpire(account.LiteAccount)

	return account, nil
}

func (am *AccountManager) beginTimeExpire(la *LiteAccount) {
	// memcache time expired
	go func() {
		for {
			select {
			case <-am.ctx.Done():
				return

			case <-la.GetExpire().C:
				am.chExpire <- la.ID
			}
		}
	}()
}

// load one account, first hit in memory, then find in database
func (am *AccountManager) loadDBLiteAccount(acctID int64) *LiteAccount {
	res := am.coll.FindOne(am.ctx, bson.D{{"account_id", acctID}})
	if res.Err() == nil {
		la := NewLiteAccount(-1)
		res.Decode(la)

		am.mapLiteAccount.Store(acctID, la)
		am.beginTimeExpire(la)
		return la
	}

	return nil
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
	v, ok := am.mapLiteAccount.Load(acctID)
	if ok {
		v.(*LiteAccount).ResetExpire()
		return v.(*LiteAccount)
	}

	return am.loadDBLiteAccount(acctID)
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

	mapPlayer := am.g.pm.GetPlayers(c.GetID())
	if len(mapPlayer) > 0 {
		return nil, fmt.Errorf("already create one player before")
	}

	p, err := am.g.pm.CreatePlayer(c.GetID(), name)
	if err != nil {
		return nil, err
	}

	c.p = p
	return p, err
}

func (am *AccountManager) SelectPlayer(c *Account, id int64) (*player.Player, error) {
	playerList := am.g.pm.GetPlayers(c.GetID())
	for _, v := range playerList {
		if v.GetID() == id {
			c.p = v
			return v, nil
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

		// memcache time expired
		case id := <-am.chExpire:
			am.mapLiteAccount.Delete(id)
		}
	}

	return nil
}
