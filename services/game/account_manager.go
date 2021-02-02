package game

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/services/game/player"
	"bitbucket.org/east-eden/server/services/game/prom"
	"bitbucket.org/east-eden/server/store"
	"bitbucket.org/east-eden/server/transport"
	"bitbucket.org/east-eden/server/utils"
	"bitbucket.org/east-eden/server/utils/cache"
	"github.com/golang/groupcache/lru"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"

	// "github.com/sasha-s/go-deadlock"
	"github.com/urfave/cli/v2"
)

var (
	maxLitePlayerLruCache    = 10000            // max number of lite player, expire non used LitePlayer
	maxAccountExecuteChannel = 100              // max account execute channel number
	AccountCacheExpire       = 10 * time.Minute // 账号cache缓存10分钟
)

type AccountManager struct {
	cacheAccounts *cache.Cache
	mapSocks      map[transport.Socket]int64 // socket->accountId

	g  *Game
	wg utils.WaitGroupWrapper

	accountConnectMax int

	playerPool      sync.Pool
	accountPool     sync.Pool
	litePlayerPool  sync.Pool
	litePlayerCache *lru.Cache

	sync.RWMutex
}

func NewAccountManager(ctx *cli.Context, g *Game) *AccountManager {
	am := &AccountManager{
		g:                 g,
		cacheAccounts:     cache.New(AccountCacheExpire, AccountCacheExpire),
		mapSocks:          make(map[transport.Socket]int64),
		accountConnectMax: ctx.Int("account_connect_max"),
		litePlayerCache:   lru.New(maxLitePlayerLruCache),
	}

	// 账号缓存删除时处理
	am.cacheAccounts.OnEvicted(func(k, v interface{}) {
		acct := v.(*player.Account)

		am.Lock()
		delete(am.mapSocks, acct.GetSock())
		am.Unlock()

		acct.Close()
		am.playerPool.Put(acct.GetPlayer())
		am.accountPool.Put(v)
		log.Info().Interface("key", k).Msg("account cache evicted")
	})

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
		log.Fatal().Err(err).Msg("migrate collection account failed")
	}

	// migrate player table
	if err := store.GetStore().MigrateDbTable("player", "account_id"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection player failed")
	}

	// migrate item table
	if err := store.GetStore().MigrateDbTable("item", "owner_id"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection item failed")
	}

	// migrate hero table
	if err := store.GetStore().MigrateDbTable("hero", "owner_id"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection hero failed")
	}

	// migrate hero table
	if err := store.GetStore().MigrateDbTable("rune", "owner_id"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection rune failed")
	}

	// migrate hero table
	if err := store.GetStore().MigrateDbTable("token", "owner_id"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection token failed")
	}

	// migrate blade table
	if err := store.GetStore().MigrateDbTable("blade", "owner_id"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection blade failed")
	}

	log.Info().Msg("AccountManager init ok ...")
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
				log.Fatal().Err(err).Msg("AccountManager Main() failed")
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
	log.Info().Msg("account manager exit...")
}

// func (am *AccountManager) onSocketEvicted(sock transport.Socket) {
// 	am.Lock()
// 	delete(am.mapSocks, sock)
// 	am.Unlock()

// 	// prometheus ops
// 	prom.OpsOnlineAccountGauge.Set(float64(am.cacheAccounts.ItemCount()))
// }

func (am *AccountManager) addAccount(ctx context.Context, userId int64, accountId int64, accountName string, sock transport.Socket) error {
	if accountId == -1 {
		return errors.New("AccountManager.addAccount failed: account id invalid!")
	}

	am.RLock()
	socksNum := len(am.mapSocks)
	am.RUnlock()
	if socksNum >= am.accountConnectMax {
		return errors.New("AccountManager.addAccount failed: Reach game server's max account connect num")
	}

	acct := am.accountPool.Get().(*player.Account)
	acct.DelayHandler = make(chan player.DelayHandleFunc, maxAccountExecuteChannel)

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
			log.Warn().
				Int64("account_id", accountId).
				Int64("user_id", userId).
				Err(err).
				Msg("save account failed")
		}

	}

	// add account to manager
	am.Lock()
	am.cacheAccounts.Set(acct.GetID(), acct, AccountCacheExpire)
	am.mapSocks[sock] = acct.GetID()
	am.Unlock()

	acct.SetSock(sock)

	// peek one player from account
	p, err := am.g.am.GetPlayerByAccount(acct)
	if err == nil {
		p.SetAccount(acct)
	}

	log.Info().
		Int64("user_id", acct.UserId).
		Int64("account_id", acct.ID).
		Str("name", acct.GetName()).
		Str("socket_remote", acct.GetSock().Remote()).
		Msg("add account success")

	// account run
	am.wg.Wrap(func() {
		err := acct.Run(ctx)
		if !utils.ErrCheck(err, "account run failed", acct.GetID()) {
			am.cacheAccounts.Delete(acct.GetID())
		}
	})

	// prometheus ops
	prom.OpsOnlineAccountGauge.Set(float64(am.cacheAccounts.ItemCount()))
	prom.OpsLogonAccountCounter.Inc()

	return nil
}

func (am *AccountManager) AccountLogon(ctx context.Context, userID int64, accountID int64, accountName string, sock transport.Socket) error {
	k, ok := am.cacheAccounts.Get(accountID)

	// if reconnect with same socket, then do nothing
	if ok && k.(*player.Account).GetSock() == sock {
		return nil
	}

	// if reconnect with another socket, replace socket in account
	if ok {
		acct := k.(*player.Account)
		if acct.GetSock() != nil {
			acct.GetSock().Close()
		}

		am.Lock()
		am.mapSocks[sock] = acct.GetID()
		am.Unlock()

		acct.SetSock(sock)
	}

	// add a new account with socket
	return am.addAccount(ctx, userID, accountID, accountName, sock)
}

func (am *AccountManager) GetAccountIdBySock(sock transport.Socket) int64 {
	am.RLock()
	defer am.RUnlock()

	return am.mapSocks[sock]
}

func (am *AccountManager) GetAccountById(acctId int64) *player.Account {
	acct, ok := am.cacheAccounts.Get(acctId)
	if ok {
		return acct.(*player.Account)
	}

	return nil
}

// add handler to account's execute channel, will be dealed by account's run goroutine
func (am *AccountManager) AccountExecute(sock transport.Socket, handler player.DelayHandleFunc) {
	id := am.GetAccountIdBySock(sock)
	acct := am.GetAccountById(id)

	if acct == nil {
		log.Warn().Int64("account_id", id).Msg("AccountExecute failed: cannot find account by id")
	}

	acct.DelayHandler <- handler
}

func (am *AccountManager) CreatePlayer(acct *player.Account, name string) (*player.Player, error) {
	// only can create one player
	if pl, _ := am.GetPlayerByAccount(acct); pl != nil {
		return nil, player.ErrCreateMoreThanOnePlayer
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
		log.Error().
			Int64("player_id", id).
			Str("player_name", name).
			Err(err).
			Msg("save player failed")
	}

	acct.SetPlayer(p)
	acct.Name = name
	acct.Level = p.GetLevel()
	acct.AddPlayerID(p.GetID())
	if err := store.GetStore().SaveObject(define.StoreType_Account, acct); err != nil {
		log.Warn().
			Int64("account_id", acct.ID).
			Int64("user_id", acct.UserId).
			Err(err).
			Msg("save account failed")
	}

	// update account info
	if _, err := am.g.rpcHandler.CallUpdateUserInfo(acct); err != nil {
		log.Warn().Err(err).Msg("CallUpdateUserInfo failed")
		return p, err
	}

	return p, err
}

func (am *AccountManager) GetPlayerByAccount(acct *player.Account) (*player.Player, error) {
	if acct == nil {
		return nil, errors.New("invalid account")
	}

	ids := acct.GetPlayerIDs()
	if len(ids) < 1 {
		return nil, errors.New("there was no player in this account")
	}

	if p := acct.GetPlayer(); p != nil {
		return p, nil
	}

	// todo load multiple players
	p := am.playerPool.Get().(*player.Player)
	err := store.GetStore().LoadObject(define.StoreType_Player, ids[0], p)
	if err != nil {
		return nil, fmt.Errorf("AccountManager.GetPlayerByAccount failed: %w", err)
	}

	acct.SetPlayer(p)
	return p, nil
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
	if pl, _ := am.g.am.GetPlayerByAccount(acct); pl != nil {
		return pl, nil
	}

	return nil, fmt.Errorf("select player with wrong id<%d>", id)
}

func (am *AccountManager) BroadCast(msg proto.Message) {
	items := am.cacheAccounts.Items()
	for _, v := range items {
		acct := v.Object.(*player.Account)
		acct.DelayHandler <- func(a *player.Account) error {
			a.SendProtoMessage(msg)
			return nil
		}
	}
}

func (am *AccountManager) Run(ctx context.Context) error {
	<-ctx.Done()
	log.Info().Msg("world session context done...")
	return nil
}
