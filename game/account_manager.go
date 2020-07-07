package game

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/golang/groupcache/lru"
	"github.com/golang/protobuf/proto"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/game/player"
	"github.com/yokaiio/yokai_server/store"
	"github.com/yokaiio/yokai_server/transport"
	"github.com/yokaiio/yokai_server/utils"
)

var (
	ErrAccountDisconnect   = errors.New("account disconnect") // handleSocket got this error will disconnect account
	maxLitePlayerLruCache  = 10000                            // max number of lite player, expire non used LitePlayer
	maxAccountLaterHandler = 100                              // max account later handler function number
)

type AccountManager struct {
	mapAccounts       map[int64]*player.Account
	mapSocks          map[transport.Socket]*player.Account
	mapAccountHandler map[int64]chan func(*player.Account)

	g  *Game
	wg utils.WaitGroupWrapper

	accountConnectMax int

	playerPool      sync.Pool
	accountPool     sync.Pool
	litePlayerPool  sync.Pool
	litePlayerCache *lru.Cache

	sync.RWMutex
}

func NewAccountManager(g *Game, ctx *cli.Context) *AccountManager {
	am := &AccountManager{
		g:                 g,
		mapAccounts:       make(map[int64]*player.Account),
		mapSocks:          make(map[transport.Socket]*player.Account),
		mapAccountHandler: make(map[int64]chan func(*player.Account)),
		accountConnectMax: ctx.Int("account_connect_max"),
		litePlayerCache:   lru.New(maxLitePlayerLruCache),
	}

	am.playerPool.New = player.NewPlayer
	am.accountPool.New = player.NewAccount
	am.litePlayerPool.New = player.NewLitePlayer
	am.litePlayerCache.OnEvicted = am.OnLitePlayerEvicted

	// add store info
	store.GetStore().AddStoreInfo(define.StoreType_Account, "account", "_id", "")
	store.GetStore().AddStoreInfo(define.StoreType_Player, "player", "_id", "")
	store.GetStore().AddStoreInfo(define.StoreType_LitePlayer, "player", "_id", "")
	store.GetStore().AddStoreInfo(define.StoreType_Item, "item", "_id", "owner_id")
	store.GetStore().AddStoreInfo(define.StoreType_Hero, "hero", "_id", "owner_id")
	store.GetStore().AddStoreInfo(define.StoreType_Rune, "rune", "_id", "owner_id")
	store.GetStore().AddStoreInfo(define.StoreType_Token, "token", "_id", "owner_id")
	store.GetStore().AddStoreInfo(define.StoreType_Blade, "blade", "_id", "owner_id")

	// migrate users table
	if err := store.GetStore().MigrateDbTable("account", "user_id"); err != nil {
		logger.Warning("migrate collection account failed:", err)
	}

	// migrate player table
	if err := store.GetStore().MigrateDbTable("player", "account_id"); err != nil {
		logger.Warning("migrate collection player failed:", err)
	}

	// migrate item table
	if err := store.GetStore().MigrateDbTable("item", "owner_id"); err != nil {
		logger.Warning("migrate collection item failed:", err)
	}

	// migrate hero table
	if err := store.GetStore().MigrateDbTable("hero", "owner_id"); err != nil {
		logger.Warning("migrate collection hero failed:", err)
	}

	// migrate hero table
	if err := store.GetStore().MigrateDbTable("rune", "owner_id"); err != nil {
		logger.Warning("migrate collection rune failed:", err)
	}

	// migrate hero table
	if err := store.GetStore().MigrateDbTable("token", "owner_id"); err != nil {
		logger.Warning("migrate collection token failed:", err)
	}

	// migrate blade table
	if err := store.GetStore().MigrateDbTable("blade", "owner_id"); err != nil {
		logger.Warning("migrate collection blade failed:", err)
	}

	logger.Info("AccountManager Init OK ...")
	return am
}

func (am *AccountManager) OnLitePlayerEvicted(key lru.Key, value interface{}) {
	am.litePlayerPool.Put(value)
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

func (am *AccountManager) onSocketEvicted(sock transport.Socket) {
	am.Lock()
	defer am.Unlock()

	acct, ok := am.mapSocks[sock]
	if ok {
		delete(am.mapAccounts, acct.GetID())
	}

	// close account message channel
	if chHandler, ok := am.mapAccountHandler[acct.GetID()]; ok {
		delete(am.mapAccountHandler, acct.GetID())
		close(chHandler)
	}

	delete(am.mapSocks, sock)
	am.playerPool.Put(acct.GetPlayer())
	am.accountPool.Put(acct)
}

func (am *AccountManager) addAccount(ctx context.Context, userId int64, accountId int64, accountName string, sock transport.Socket) error {
	if accountId == -1 {
		return errors.New("AccountManager.addAccount failed: account id invalid!")
	}

	if len(am.mapSocks) >= am.accountConnectMax {
		return errors.New("AccountManager.addAccount failed: Reach game server's max account connect num")
	}

	acct := am.accountPool.Get().(*player.Account)
	err := store.GetStore().LoadObject(define.StoreType_Account, accountId, acct)
	if err != nil && !errors.Is(err, store.ErrNoResult) {
		return fmt.Errorf("AccountManager.addAccount failed: %w", err)
	}

	if errors.Is(err, store.ErrNoResult) {
		// store cannot load account, create a new account
		acct.ID = accountId
		acct.UserId = userId
		acct.GameId = am.g.ID
		acct.Name = accountName

		// save object
		if err := store.GetStore().SaveObject(define.StoreType_Account, acct); err != nil {
			logger.WithFields(logger.Fields{
				"account_id": accountId,
				"user_id":    userId,
			}).Warn("save account failed")
		}

	}

	// add account to manager
	chHandler := make(chan func(*player.Account), maxAccountLaterHandler)
	am.Lock()
	am.mapAccounts[acct.GetID()] = acct
	am.mapSocks[sock] = acct
	am.mapAccountHandler[acct.GetID()] = chHandler
	am.Unlock()

	acct.SetSock(sock)
	acct.SetLaterHandlerChannel(chHandler)
	sock.AddEvictedHandle(am.onSocketEvicted)

	// peek one player from account
	if p := am.g.am.GetPlayerByAccount(acct); p != nil {
		p.SetAccount(acct)
	}

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

	})

	return nil
}

func (am *AccountManager) AccountLogon(ctx context.Context, userID int64, accountID int64, accountName string, sock transport.Socket) error {
	am.RLock()
	account, acctOK := am.mapAccounts[accountID]
	am.RUnlock()

	// if reconnect with same socket, then do nothing
	if acctOK && account.GetSock() == sock {
		return nil
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

		return nil
	}

	// add a new account with socket
	return am.addAccount(ctx, userID, accountID, accountName, sock)
}

func (am *AccountManager) GetAccountByID(acctId int64) *player.Account {
	am.RLock()
	defer am.RUnlock()

	account, ok := am.mapAccounts[acctId]

	if !ok {
		return nil
	}

	return account
}

func (am *AccountManager) GetAccountBySock(sock transport.Socket) *player.Account {
	am.RLock()
	defer am.RUnlock()

	account, ok := am.mapSocks[sock]

	if !ok {
		return nil
	}

	return account
}

func (am *AccountManager) GetAllAccounts() []*player.Account {
	ret := make([]*player.Account, 0)

	am.RLock()
	for _, account := range am.mapAccounts {
		ret = append(ret, account)
	}
	am.RUnlock()

	return ret
}

// add handler to account's handler channel, will be dealed by account's run goroutine
func (am *AccountManager) AccountLaterHandle(sock transport.Socket, handle func(*player.Account)) error {
	am.RLock()
	defer am.RUnlock()

	acct, ok := am.mapSocks[sock]
	if !ok {
		return fmt.Errorf("AccountManager.AccountLaterHandle failed: can't find account by socket")
	}

	if chFunc, ok := am.mapAccountHandler[acct.GetID()]; ok {
		chFunc <- handle
	}

	return nil
}

func (am *AccountManager) CreatePlayer(acct *player.Account, name string) (*player.Player, error) {
	// only can create one player
	if am.GetPlayerByAccount(acct) != nil {
		return nil, errors.New("AccountManager.CreatePlayer failed: only can create one player")
	}

	id, err := utils.NextID(define.SnowFlake_Player)
	if err != nil {
		return nil, err
	}

	p := am.playerPool.Get().(*player.Player)
	p.AccountID = acct.ID
	p.SetAccount(acct)
	p.SetID(id)
	p.SetName(name)
	if err := store.GetStore().SaveObject(define.StoreType_Player, p); err != nil {
		logger.WithFields(logger.Fields{
			"player_id":   id,
			"player_name": name,
		}).Error("save player failed")
	}

	acct.SetPlayer(p)
	acct.Name = name
	acct.Level = p.GetLevel()
	acct.AddPlayerID(p.GetID())
	if err := store.GetStore().SaveObject(define.StoreType_Account, acct); err != nil {
		logger.WithFields(logger.Fields{
			"account_id": acct.ID,
			"user_id":    acct.UserId,
		}).Warn("save account failed")
	}

	// update account info
	am.g.rpcHandler.CallUpdateUserInfo(acct)

	return p, err
}

func (am *AccountManager) GetPlayerByAccount(acct *player.Account) *player.Player {
	if acct == nil {
		return nil
	}

	ids := acct.GetPlayerIDs()
	if len(ids) < 1 {
		return nil
	}

	if p := acct.GetPlayer(); p != nil {
		return p
	}

	// todo load multiple players
	p := am.playerPool.Get().(*player.Player)
	err := store.GetStore().LoadObject(define.StoreType_Player, ids[0], p)
	if err != nil {
		return nil
	}

	acct.SetPlayer(p)
	return p
}

func (am *AccountManager) GetLitePlayer(playerId int64) (player.LitePlayer, error) {
	am.RLock()
	defer am.RUnlock()

	if lp, ok := am.litePlayerCache.Get(playerId); ok {
		return *(lp.(*player.LitePlayer)), nil
	}

	lp := am.litePlayerPool.Get().(*player.LitePlayer)
	err := store.GetStore().LoadObject(define.StoreType_LitePlayer, playerId, lp)
	if err == nil {
		am.litePlayerCache.Add(lp.ID, lp)
		return *lp, nil
	}

	am.litePlayerPool.Put(lp)
	return *(player.NewLitePlayer().(*player.LitePlayer)), err
}

// todo omitempty
func (am *AccountManager) SelectPlayer(acct *player.Account, id int64) (*player.Player, error) {
	if pl := am.g.am.GetPlayerByAccount(acct); pl != nil {
		return pl, nil
	}

	return nil, fmt.Errorf("select player with wrong id<%d>", id)
}

func (am *AccountManager) BroadCast(msg proto.Message) {
	am.RLock()
	defer am.RUnlock()

	fn := func(acct *player.Account) {
		acct.SendProtoMessage(msg)
	}

	for _, acct := range am.mapAccounts {
		if chFunc, ok := am.mapAccountHandler[acct.GetID()]; ok {
			chFunc <- fn
		}
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
