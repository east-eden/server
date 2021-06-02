package game

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/services/game/player"
	"bitbucket.org/funplus/server/services/game/prom"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/transport"
	"bitbucket.org/funplus/server/utils"
	"bitbucket.org/funplus/server/utils/cache"
	"github.com/golang/groupcache/lru"
	"github.com/hellodudu/task"
	log "github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"

	// "github.com/sasha-s/go-deadlock"
	"github.com/urfave/cli/v2"
)

var (
	maxPlayerInfoLruCache = 10000            // max number of lite player, expire non used PlayerInfo
	UserCacheExpire       = 10 * time.Minute // user cache缓存10分钟
	AccountCacheExpire    = 10 * time.Minute // 账号cache缓存10分钟

	ErrAccountHasNoPlayer = errors.New("account has no player")
	ErrAccountNotFound    = errors.New("account not found")
)

type AccountManagerFace interface {
}

type AccountManager struct {
	cacheAccounts *cache.Cache
	cacheUsers    *cache.Cache
	mapSocks      map[transport.Socket]int64 // socket->accountId

	g  *Game
	wg utils.WaitGroupWrapper

	accountConnectMax int

	userPool        sync.Pool
	playerPool      sync.Pool
	accountPool     sync.Pool
	playerInfoPool  sync.Pool
	playerInfoCache *lru.Cache

	sync.RWMutex
}

func NewAccountManager(ctx *cli.Context, g *Game) *AccountManager {
	am := &AccountManager{
		g:                 g,
		cacheUsers:        cache.New(UserCacheExpire, UserCacheExpire),
		cacheAccounts:     cache.New(AccountCacheExpire, AccountCacheExpire),
		mapSocks:          make(map[transport.Socket]int64),
		accountConnectMax: ctx.Int("account_connect_max"),
		playerInfoCache:   lru.New(maxPlayerInfoLruCache),
	}

	// user pool
	am.userPool.New = NewUser

	// user cache evicted
	am.cacheUsers.OnEvicted(func(k, v interface{}) {
		log.Info().Interface("key", k).Interface("value", v).Msg("user cache evicted")
		am.userPool.Put(v)
	})

	// 账号缓存删除时处理
	am.cacheAccounts.OnEvicted(func(k, v interface{}) {
		acct := v.(*player.Account)

		am.Lock()
		delete(am.mapSocks, acct.GetSock())
		am.Unlock()

		acct.Stop()
		am.playerPool.Put(acct.GetPlayer())
		am.accountPool.Put(v)
		log.Info().Interface("key", k).Msg("account cache evicted")
	})

	am.playerPool.New = player.NewPlayer
	am.accountPool.New = player.NewAccount
	am.playerInfoPool.New = player.NewPlayerInfo
	am.playerInfoCache.OnEvicted = am.OnPlayerInfoEvicted

	// add store info
	store.GetStore().AddStoreInfo(define.StoreType_User, "user", "_id")
	store.GetStore().AddStoreInfo(define.StoreType_Account, "account", "_id")
	store.GetStore().AddStoreInfo(define.StoreType_Player, "player", "_id")
	store.GetStore().AddStoreInfo(define.StoreType_Item, "player_item", "_id")
	store.GetStore().AddStoreInfo(define.StoreType_Hero, "player_hero", "_id")
	store.GetStore().AddStoreInfo(define.StoreType_Token, "player_token", "_id")
	store.GetStore().AddStoreInfo(define.StoreType_Fragment, "player_fragment", "_id")
	store.GetStore().AddStoreInfo(define.StoreType_Collection, "player_collection", "_id")

	// migrate user table
	if err := store.GetStore().MigrateDbTable("user", "account_id", "player_id"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection user failed")
	}

	// migrate account table
	if err := store.GetStore().MigrateDbTable("account", "user_id"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection account failed")
	}

	// migrate player table
	if err := store.GetStore().MigrateDbTable("player", "account_id"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection player failed")
	}

	// migrate item table
	if err := store.GetStore().MigrateDbTable("player_item", "owner_id"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection player_item failed")
	}

	// migrate hero table
	if err := store.GetStore().MigrateDbTable("player_hero", "owner_id"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection player_hero failed")
	}

	// migrate hero table
	if err := store.GetStore().MigrateDbTable("player_token", "owner_id"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection player_token failed")
	}

	// migrate fragment table
	if err := store.GetStore().MigrateDbTable("player_fragment", "owner_id"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection player_fragment failed")
	}

	// migrate collection table
	if err := store.GetStore().MigrateDbTable("player_collection", "type_id", "owner_id"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection player_collection failed")
	}

	log.Info().Msg("AccountManager init ok ...")
	return am
}

// todo get by userId???
func (am *AccountManager) getUser(userId int64) (*User, error) {
	u, ok := am.cacheUsers.Get(userId)
	if ok {
		return u.(*User), nil
	}

	u = am.userPool.Get()
	err := store.GetStore().FindOne(context.Background(), define.StoreType_User, userId, u)
	if err == nil {
		am.cacheUsers.Set(userId, u, UserCacheExpire)
		return u.(*User), nil
	}

	if errors.Is(err, store.ErrNoResult) {
		accountId, err := utils.NextID(define.SnowFlake_Account)
		if err != nil {
			am.userPool.Put(u)
			return nil, err
		}

		user := u.(*User)
		user.UserID = userId
		user.AccountID = accountId

		err = store.GetStore().UpdateOne(context.Background(), define.StoreType_User, user.UserID, user, true)
		if !utils.ErrCheck(err, "UpdateOne failed when AccountManager.getUser", user) {
			am.userPool.Put(user)
			return nil, err
		}

		am.cacheUsers.Set(userId, user, UserCacheExpire)
		return user, nil
	}

	am.userPool.Put(u)
	return nil, err
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
		defer utils.CaptureException()
		exitFunc(am.Run(ctx))
	})

	return <-exitCh
}

func (am *AccountManager) Exit() {
	am.wg.Wait()
	log.Info().Msg("account manager exit...")
}

func (am *AccountManager) handleLoadPlayer(ctx context.Context, p ...interface{}) error {
	acct := p[0].(*player.Account)

	load := func(acct *player.Account) error {
		ids := acct.GetPlayerIDs()
		if len(ids) < 1 {
			return ErrAccountHasNoPlayer
		}

		p := am.playerPool.Get().(*player.Player)
		p.Init(ids[0])
		p.SetAccount(acct)
		err := store.GetStore().FindOne(context.Background(), define.StoreType_Player, ids[0], p)
		if !utils.ErrCheck(err, "load player object failed", ids[0]) {
			am.playerPool.Put(p)
			return err
		}

		// 加载玩家其他数据
		err = p.AfterLoad()
		if !utils.ErrCheck(err, "player.AfterLoad failed", ids[0]) {
			am.playerPool.Put(p)
			return err
		}

		acct.SetPlayer(p)
		return nil
	}

	// 加载玩家
	err := load(acct)

	return err
}

// 踢掉account对象
func (am *AccountManager) KickAccount(ctx context.Context, acctId int64, gameId int32) error {
	if acctId == -1 {
		return nil
	}

	// 踢掉本服account
	if int16(gameId) == am.g.ID {

		acct := am.GetAccountById(acctId)
		if acct == nil {
			return nil
		}

		acct.Stop()
		store.GetStore().Flush()
		return nil

	} else {
		// game节点不存在的话不用发送rpc
		nodeId := fmt.Sprintf("game-%d", gameId)
		srvs, err := am.g.mi.srv.Server().Options().Registry.GetService("game")
		if err != nil {
			return nil
		}

		hit := false
		for _, srv := range srvs {
			for _, node := range srv.Nodes {
				if node.Id == nodeId {
					hit = true
					break
				}
			}
		}

		if !hit {
			return nil
		}

		// 发送rpc踢掉其他服account
		rs, err := am.g.rpcHandler.CallKickAccountOffline(acctId, gameId)
		if !utils.ErrCheck(err, "kick account offline failed", acctId, gameId, rs) {
			return err
		}

		// rpc调用成功
		if rs.GetAccountId() == acctId {
			return nil
		}

		return errors.New("kick account invalid error")
	}
}

func (am *AccountManager) addNewAccount(ctx context.Context, userId int64, accountId int64, accountName string, sock transport.Socket) (*player.Account, error) {
	// check max connections
	am.RLock()
	socksNum := len(am.mapSocks)
	am.RUnlock()
	if socksNum >= am.accountConnectMax {
		return nil, errors.New("AccountManager.addAccount failed: Reach game server's max account connect num")
	}

	// init new account
	acct := am.accountPool.Get().(*player.Account)
	acct.Init()
	acct.SetRpcCaller(am.g.rpcHandler)

	// load account info from store
	err := store.GetStore().FindOne(context.Background(), define.StoreType_Account, accountId, acct)
	if err != nil && !errors.Is(err, store.ErrNoResult) {
		return nil, fmt.Errorf("AccountManager.addAccount failed: %w", err)
	}

	// 如果account的上次登陆game节点不是此节点，则发rpc提掉上一个登陆节点的account
	if acct.GameId != -1 && acct.GameId != am.g.ID {
		err := am.KickAccount(ctx, acct.Id, int32(acct.GameId))
		if !utils.ErrCheck(err, "kick account failed", acct.Id, acct.GameId, am.g.ID) {
			return nil, err
		}
	}

	if errors.Is(err, store.ErrNoResult) {
		// 账号首次登陆
		acct.Id = accountId
		acct.UserId = userId
		acct.GameId = am.g.ID
		acct.Name = accountName

		// save account
		err := store.GetStore().UpdateOne(context.Background(), define.StoreType_Account, acct.Id, acct, true)
		utils.ErrPrint(err, "UpdateOne failed when AccountManager.addAccount", accountId, userId)
	} else {
		// 更新account节点id
		acct.GameId = am.g.ID
		fields := map[string]interface{}{
			"game_id": acct.GameId,
		}

		err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Account, acct.Id, fields, true)
		_ = utils.ErrCheck(err, "UpdateFields failed when AccountManager.addAccount", acct.Id, acct.GameId)
	}

	// add account to manager
	am.Lock()
	am.mapSocks[sock] = acct.GetId()
	am.Unlock()

	acct.SetSock(sock)
	am.cacheAccounts.Set(acct.GetId(), acct, AccountCacheExpire)

	log.Info().
		Int64("user_id", acct.UserId).
		Int64("account_id", acct.Id).
		Str("name", acct.GetName()).
		Str("socket_remote", acct.GetSock().Remote()).
		Msg("add account success")

	// prometheus ops
	// prom.OpsOnlineAccountGauge.Set(float64(am.cacheAccounts.ItemCount()))
	prom.OpsLogonAccountCounter.Inc()

	return acct, nil
}

func (am *AccountManager) runAccountTask(ctx context.Context, acct *player.Account, startFns ...task.StartFn) {
	// account init task
	acct.InitTask(startFns...)

	// account task run
	acct.ResetTimeout()

	am.wg.Wrap(func() {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				log.Error().Caller().Msgf("catch exception:%v, panic recovered with stack:%s", err, stack)

				// 立即删除缓存
				am.cacheAccounts.Delete(acct.GetId())
			}
		}()

		errAcct := acct.TaskRun(ctx)
		utils.ErrPrint(errAcct, "account run failed", acct.GetId())

		// 记录下线时间
		acct.LastLogoffTime = int32(time.Now().Unix())
		fields := map[string]interface{}{
			"last_logoff_time": acct.LastLogoffTime,
		}
		err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Account, acct.Id, fields, true)
		utils.ErrPrint(err, "account save last_logoff_time failed", acct.Id, acct.LastLogoffTime)

		// 被踢下线或者连接超时，立即删除缓存
		if errors.Is(errAcct, player.ErrAccountKicked) || errors.Is(errAcct, task.ErrTimeout) {
			am.cacheAccounts.Delete(acct.GetId())
			return
		}
	})
}

func (am *AccountManager) Logon(ctx context.Context, userId int64, sock transport.Socket) error {
	// if accountId == -1 {
	// 	return errors.New("AccountManager.addAccount failed: account id invalid!")
	// }

	user, err := am.getUser(userId)
	if !utils.ErrCheck(err, "getUser failed when AccountManager.Logon", userId) {
		return err
	}

	c, ok := am.cacheAccounts.Get(user.AccountID)

	if ok {
		// cache exist
		acct := c.(*player.Account)

		// connect with new socket
		if acct.GetSock() != sock {
			if acct.GetSock() != nil {
				acct.GetSock().Close()
			}

			am.Lock()
			am.mapSocks[sock] = acct.GetId()
			am.Unlock()

			acct.SetSock(sock)
		}

		// if task is running, return
		if acct.IsTaskRunning() {
			_ = am.AddAccountTask(
				ctx,
				acct.GetId(),
				func(ctx context.Context, p ...interface{}) error {
					acct := p[0].(*player.Account)
					acct.LogonSucceed()
					return nil
				},
				nil,
			)
			return nil
		}

		// account run
		am.runAccountTask(ctx, acct, func() {
			acct.LogonSucceed()
		})

	} else {
		// cache not exist, add a new account with socket
		acct, err := am.addNewAccount(ctx, userId, user.AccountID, user.PlayerName, sock)
		if !utils.ErrCheck(err, "addNewAccount failed when AccountManager.Logon", userId, user.AccountID) {
			return err
		}

		// account run
		am.runAccountTask(ctx, acct, func() {
			err := am.handleLoadPlayer(ctx, acct)

			// 加载玩家成功或者账号下没有玩家
			if err == nil || errors.Is(err, ErrAccountHasNoPlayer) {
				acct.LogonSucceed()
				return
			}

			// 加载失败
			acct.Stop()
		})
	}

	return nil
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
func (am *AccountManager) AddAccountTask(ctx context.Context, acctId int64, fn task.TaskHandler, m proto.Message) error {
	acct := am.GetAccountById(acctId)

	if acct == nil {
		return ErrAccountNotFound
	}

	return acct.AddTask(ctx, fn, m)
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
	p.Init(id)
	p.AccountID = acct.Id
	p.SetAccount(acct)
	p.SetName(name)

	// save handle
	errHandle := func(f func() error) {
		if err != nil {
			return
		}

		err = f()
	}
	errHandle(func() error {
		return store.GetStore().UpdateOne(context.Background(), define.StoreType_Player, p.ID, p, true)
	})

	// errHandle(func() error {
	// 	return store.GetStore().UpdateOne(context.Background(), define.StoreType_Hero, p.ID, p.HeroManager())
	// })

	// errHandle(func() error {
	// 	return store.GetStore().UpdateOne(context.Background(), define.StoreType_Item, p.ID, p.ItemManager())
	// })

	errHandle(func() error {
		return store.GetStore().UpdateOne(context.Background(), define.StoreType_Token, p.ID, p.TokenManager(), true)
	})

	errHandle(func() error {
		return store.GetStore().UpdateOne(context.Background(), define.StoreType_Fragment, p.ID, p.FragmentManager(), true)
	})

	// 保存失败处理
	if !utils.ErrCheck(err, "save player failed when CreatePlayer", id, name) {
		am.playerPool.Put(p)
		return nil, err
	}

	acct.SetPlayer(p)
	acct.Name = name
	acct.Level = p.GetLevel()
	acct.AddPlayerID(p.GetId())
	if err := store.GetStore().UpdateOne(context.Background(), define.StoreType_Account, acct.Id, acct, true); err != nil {
		log.Warn().
			Int64("account_id", acct.Id).
			Int64("user_id", acct.UserId).
			Err(err).
			Msg("save account failed")
	}

	// 第一次上线处理
	p.OnFirstLogon()

	// 同步玩家初始信息
	p.SendInitInfo()

	return p, err
}

func (am *AccountManager) GetPlayerByAccount(acct *player.Account) (*player.Player, error) {
	if acct == nil {
		return nil, errors.New("invalid account")
	}

	if p := acct.GetPlayer(); p != nil {
		return p, nil
	}

	return nil, errors.New("invalid player")
}

func (am *AccountManager) GetPlayerInfo(playerId int64) (player.PlayerInfo, error) {
	am.RLock()
	defer am.RUnlock()

	if lp, ok := am.playerInfoCache.Get(playerId); ok {
		return *(lp.(*player.PlayerInfo)), nil
	}

	lp := am.playerInfoPool.Get().(*player.PlayerInfo)
	lp.Init()
	err := store.GetStore().FindOne(context.Background(), define.StoreType_Player, playerId, lp)
	if err == nil {
		am.playerInfoCache.Add(lp.ID, lp)
		return *lp, nil
	}

	am.playerInfoPool.Put(lp)
	return player.PlayerInfo{}, err
}

func (am *AccountManager) BroadCast(msg proto.Message) {
	items := am.cacheAccounts.Items()
	for _, v := range items {
		acct := v.Object.(*player.Account)
		_ = acct.AddTask(context.Background(), func(c context.Context, p ...interface{}) error {
			a := p[0].(*player.Account)
			message := p[1].(proto.Message)
			a.SendProtoMessage(message)
			return nil
		}, msg)
	}
}

func (am *AccountManager) Run(ctx context.Context) error {
	<-ctx.Done()
	log.Info().Msg("world session context done...")
	return nil
}
