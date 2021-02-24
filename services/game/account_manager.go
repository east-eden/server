package game

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/services/game/player"
	"github.com/east-eden/server/services/game/prom"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/transport"
	"github.com/east-eden/server/utils"
	"github.com/east-eden/server/utils/cache"
	"github.com/golang/groupcache/lru"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"

	// "github.com/sasha-s/go-deadlock"
	"github.com/urfave/cli/v2"
)

var (
	maxPlayerInfoLruCache = 10000            // max number of lite player, expire non used PlayerInfo
	maxAccountSlowHandler = 100              // max account execute channel number
	AccountCacheExpire    = 10 * time.Minute // 账号cache缓存10分钟
)

type AccountManager struct {
	cacheAccounts *cache.Cache
	mapSocks      map[transport.Socket]int64 // socket->accountId

	g  *Game
	wg utils.WaitGroupWrapper

	accountConnectMax int

	playerPool      sync.Pool
	accountPool     sync.Pool
	playerInfoPool  sync.Pool
	playerInfoCache *lru.Cache

	sync.RWMutex
}

func NewAccountManager(ctx *cli.Context, g *Game) *AccountManager {
	am := &AccountManager{
		g:                 g,
		cacheAccounts:     cache.New(AccountCacheExpire, AccountCacheExpire),
		mapSocks:          make(map[transport.Socket]int64),
		accountConnectMax: ctx.Int("account_connect_max"),
		playerInfoCache:   lru.New(maxPlayerInfoLruCache),
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
	am.playerInfoPool.New = player.NewPlayerInfo
	am.playerInfoCache.OnEvicted = am.OnPlayerInfoEvicted

	// add store info
	store.GetStore().AddStoreInfo(define.StoreType_Account, "account", "_id", "")
	store.GetStore().AddStoreInfo(define.StoreType_Player, "player", "_id", "")
	store.GetStore().AddStoreInfo(define.StoreType_PlayerInfo, "player", "_id", "")
	store.GetStore().AddStoreInfo(define.StoreType_Item, "item", "_id", "owner_id")
	store.GetStore().AddStoreInfo(define.StoreType_Hero, "hero", "_id", "owner_id")
	store.GetStore().AddStoreInfo(define.StoreType_Rune, "rune", "_id", "owner_id")
	store.GetStore().AddStoreInfo(define.StoreType_Token, "token", "_id", "owner_id")
	store.GetStore().AddStoreInfo(define.StoreType_Blade, "blade", "_id", "owner_id")
	store.GetStore().AddStoreInfo(define.StoreType_Fragment, "fragment", "_id", "owner_id")

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

	// migrate fragment table
	if err := store.GetStore().MigrateDbTable("fragment", "owner_id"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection fragment failed")
	}

	log.Info().Msg("AccountManager init ok ...")
	return am
}

func (am *AccountManager) OnPlayerInfoEvicted(key lru.Key, value interface{}) {
	am.playerInfoPool.Put(value)
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
	acct.SlowHandler = make(chan *player.AccountSlowHandler, maxAccountSlowHandler)

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
		if err := store.GetStore().SaveObject(define.StoreType_Account, acct.ID, acct); err != nil {
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
func (am *AccountManager) AccountSlowHandle(sock transport.Socket, handler *player.AccountSlowHandler) {
	id := am.GetAccountIdBySock(sock)
	acct := am.GetAccountById(id)

	if acct == nil {
		log.Warn().Int64("account_id", id).Msg("AccountExecute failed: cannot find account by id")
		return
	}

	acct.SlowHandler <- handler
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
	p.Init()
	p.AccountID = acct.ID
	p.SetAccount(acct)
	p.SetID(id)
	p.SetName(name)

	// save handle
	errHandle := func(f func() error) {
		if err != nil {
			return
		}

		err = f()
	}
	errHandle(func() error {
		return store.GetStore().SaveObject(define.StoreType_Player, p.ID, p)
	})

	errHandle(func() error {
		return store.GetStore().SaveObject(define.StoreType_Token, p.ID, p.TokenManager())
	})

	errHandle(func() error {
		return store.GetStore().SaveObject(define.StoreType_Hero, p.ID, p.HeroManager())
	})

	errHandle(func() error {
		return store.GetStore().SaveObject(define.StoreType_Item, p.ID, p.ItemManager())
	})

	errHandle(func() error {
		return store.GetStore().SaveObject(define.StoreType_Rune, p.ID, p.RuneManager())
	})

	errHandle(func() error {
		return store.GetStore().SaveObject(define.StoreType_Blade, p.ID, p.BladeManager())
	})

	errHandle(func() error {
		return store.GetStore().SaveObject(define.StoreType_Fragment, p.ID, p.FragmentManager())
	})

	// 保存失败处理
	if pass := utils.ErrCheck(err, "save player failed when CreatePlayer", id, name); !pass {
		am.playerPool.Put(p)
		return nil, err
	}

	acct.SetPlayer(p)
	acct.Name = name
	acct.Level = p.GetLevel()
	acct.AddPlayerID(p.GetID())
	if err := store.GetStore().SaveObject(define.StoreType_Account, acct.ID, acct); err != nil {
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
	p.Init()
	err := store.GetStore().LoadObject(define.StoreType_Player, ids[0], p)
	if err != nil {
		return nil, fmt.Errorf("AccountManager.GetPlayerByAccount failed: %w", err)
	}

	p.AfterLoad()

	acct.SetPlayer(p)
	return p, nil
}

func (am *AccountManager) GetPlayerInfo(playerId int64) (player.PlayerInfo, error) {
	am.RLock()
	defer am.RUnlock()

	if lp, ok := am.playerInfoCache.Get(playerId); ok {
		return *(lp.(*player.PlayerInfo)), nil
	}

	lp := am.playerInfoPool.Get().(*player.PlayerInfo)
	err := store.GetStore().LoadObject(define.StoreType_PlayerInfo, playerId, lp)
	if err == nil {
		am.playerInfoCache.Add(lp.ID, lp)
		return *lp, nil
	}

	am.playerInfoPool.Put(lp)
	return *(player.NewPlayerInfo().(*player.PlayerInfo)), err
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

		acct.SlowHandler <- &player.AccountSlowHandler{
			F: func(ctx context.Context, a *player.Account, p *transport.Message) error {
				a.SendProtoMessage(msg)
				return nil
			},
		}
	}
}

func (am *AccountManager) Run(ctx context.Context) error {
	<-ctx.Done()
	log.Info().Msg("world session context done...")
	return nil
}
