package game

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/services/game/player"
	"e.coding.net/mmstudio/blade/server/services/game/prom"
	"e.coding.net/mmstudio/blade/server/store"
	"e.coding.net/mmstudio/blade/server/transport"
	"e.coding.net/mmstudio/blade/server/utils"
	"e.coding.net/mmstudio/blade/server/utils/cache"
	"github.com/hellodudu/task"
	log "github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"

	// "github.com/sasha-s/go-deadlock"
	"github.com/urfave/cli/v2"
)

var (
	CacheCleanupInterval  = 1 * time.Minute  // cache cleanup interval
	UserCacheExpire       = 10 * time.Minute // user cache缓存10分钟
	AccountCacheExpire    = 1 * time.Minute  // 账号cache缓存10分钟
	PlayerInfoCacheExpire = time.Hour        // 玩家简易信息cache缓存1小时

	ErrAccountHasNoPlayer = errors.New("account has no player")
	ErrAccountNotFound    = errors.New("account not found")
	ErrPlayerInfoNotFound = errors.New("player info not found")
	ErrPlayerLoadFailed   = errors.New("player load failed")
)

type AccountManagerFace interface {
}

type AccountManager struct {
	cacheAccounts    *cache.Cache
	cacheUsers       *cache.Cache
	cachePlayerInfos *cache.Cache
	mapSocks         map[transport.Socket]int64 // socket->accountId

	g  *Game
	wg utils.WaitGroupWrapper

	accountConnectMax int

	userPool       sync.Pool
	playerPool     sync.Pool
	accountPool    sync.Pool
	playerInfoPool sync.Pool

	sync.RWMutex
}

func NewAccountManager(ctx *cli.Context, g *Game) *AccountManager {
	am := &AccountManager{
		g:                 g,
		cacheAccounts:     cache.New(AccountCacheExpire, CacheCleanupInterval),
		cacheUsers:        cache.New(UserCacheExpire, CacheCleanupInterval),
		cachePlayerInfos:  cache.New(PlayerInfoCacheExpire, CacheCleanupInterval),
		mapSocks:          make(map[transport.Socket]int64),
		accountConnectMax: ctx.Int("account_connect_max"),
	}

	// heart beat timeout
	player.AccountTaskTimeout = ctx.Duration("heart_beat_timeout")

	// user pool
	am.userPool.New = NewUser

	// user cache evicted
	am.cacheUsers.OnEvicted(func(k, v interface{}) {
		log.Info().Interface("key", k).Interface("value", v).Msg("user cache evicted")
		am.userPool.Put(v)
	})

	// player info cache evicted
	am.cachePlayerInfos.OnEvicted(func(k, v interface{}) {
		log.Info().Interface("key", k).Interface("value", v).Msg("player info cache evicted")
		am.playerInfoPool.Put(v)
	})

	// 账号缓存删除时处理
	am.cacheAccounts.OnEvicted(func(k, v interface{}) {
		acct := v.(*player.Account)

		event := log.Info().Caller().Int64("account_id", acct.Id)
		if acct.GetSock() != nil {
			event = event.Str("sock_local", acct.GetSock().Local()).
				Str("sock_remote", acct.GetSock().Remote())
		}
		event.Msg("account cache evicted")

		acct.StopTask()
		if acct.GetPlayer() != nil {
			acct.GetPlayer().Destroy()
			am.playerPool.Put(acct.GetPlayer())
		}
		am.accountPool.Put(acct)
	})

	am.playerPool.New = player.NewPlayer
	am.accountPool.New = player.NewAccount
	am.playerInfoPool.New = player.NewPlayerInfo

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

		pl := am.playerPool.Get().(*player.Player)
		pl.Init(ids[0])
		pl.SetAccount(acct)
		err := store.GetStore().FindOne(context.Background(), define.StoreType_Player, ids[0], pl)
		if !utils.ErrCheck(err, "load player object failed", ids[0]) {
			am.playerPool.Put(pl)
			return err
		}

		// 加载玩家其他数据
		err = pl.AfterLoad()
		if !utils.ErrCheck(err, "player.AfterLoad failed", ids[0]) {
			am.playerPool.Put(pl)
			return fmt.Errorf("%w: %s", ErrPlayerLoadFailed, err.Error())
		}

		acct.SetPlayer(pl)
		return nil
	}

	// 加载玩家
	return load(acct)
}

// kick all cache
func (am *AccountManager) KickAllCache() {
	am.cacheUsers.DeleteAll()
	am.cacheAccounts.DeleteAll()
	am.cachePlayerInfos.DeleteAll()
	store.GetStore().Flush()
}

// 踢掉account对象
func (am *AccountManager) KickAccount(ctx context.Context, acctId int64, gameId int32) error {
	if acctId == -1 {
		return nil
	}

	// 踢掉本服account
	if int16(gameId) == am.g.ID {
		am.cacheAccounts.Delete(acctId)
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

func (am *AccountManager) newAccount(ctx context.Context, userId int64, accountId int64, accountName string, sock transport.Socket) (*player.Account, error) {
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
		acct.SaveAccount()

	} else {
		// 更新game节点
		if acct.GameId != am.g.ID {
			acct.SaveGameNode(am.g.ID)
		}
	}

	// prometheus ops
	// prom.OpsOnlineAccountGauge.Set(float64(am.cacheAccounts.ItemCount()))
	prom.OpsLogonAccountCounter.Inc()

	return acct, nil
}

func (am *AccountManager) startAccountTask(ctx context.Context, sock transport.Socket, acct *player.Account, start task.StartFn) {
	// account init task
	startFn := func() {
		// 增加连接
		am.Lock()
		am.mapSocks[sock] = acct.GetId()
		am.Unlock()
		acct.SetSock(sock)
		start()
	}
	stopFn := func() {
		// 删除连接
		am.Lock()
		delete(am.mapSocks, acct.GetSock())
		am.Unlock()
	}
	acct.InitTask(startFn, stopFn)

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

		log.Info().Caller().Int64("account_id", acct.GetId()).Str("remote_sock", sock.Remote()).Msg("account run new task")
		errAcct := acct.TaskRun(ctx)
		utils.ErrPrint(errAcct, "account run failed", acct.GetId())

		// 被踢下线、连接超时、登陆失败，都立即删除缓存
		if errors.Is(errAcct, player.ErrAccountKicked) || errors.Is(errAcct, task.ErrTimeout) {
			am.cacheAccounts.Delete(acct.GetId())
			return
		}
	})
}

func (am *AccountManager) Logon(ctx context.Context, userId int64, newSock transport.Socket) error {
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
		prevSock := acct.GetSock()

		// connect with new socket
		if prevSock != newSock && prevSock != nil {
			log.Info().
				Caller().
				Int64("account_id", acct.Id).
				Str("prev_socket_local", prevSock.Local()).
				Str("prev_socket_remote", prevSock.Remote()).
				Str("new_sock_local", newSock.Local()).
				Str("new_sock_remote", newSock.Remote()).
				Msg("logon with new socket replacing prev socket")

			acct.StopTask()
		}

		// run new task
		am.startAccountTask(ctx, newSock, acct, func() {
			acct.LogonSucceed()
		})

		log.Info().
			Caller().
			Int64("account_id", acct.Id).
			Str("new_sock_local", newSock.Local()).
			Str("new_sock_remote", newSock.Remote()).
			Msg("logon with task is not running")

	} else {
		// cache not exist, add a new account with socket
		acct, err := am.newAccount(ctx, userId, user.AccountID, user.PlayerName, newSock)
		if !utils.ErrCheck(err, "addNewAccount failed when AccountManager.Logon", userId, user.AccountID) {
			return err
		}

		am.cacheAccounts.Set(acct.GetId(), acct, AccountCacheExpire)

		// account run
		am.startAccountTask(ctx, newSock, acct, func() {
			err := am.handleLoadPlayer(ctx, acct)

			// 加载玩家成功或者账号下没有玩家
			if err == nil || errors.Is(err, ErrAccountHasNoPlayer) {
				acct.LogonSucceed()
				return
			}

			// 加载失败
			am.cacheAccounts.Delete(acct.GetId())
		})

		log.Info().
			Int64("user_id", acct.UserId).
			Int64("account_id", acct.Id).
			Str("name", acct.GetName()).
			Str("socket_remote", newSock.Remote()).
			Msg("add account success")
	}

	return nil
}

func (am *AccountManager) GetAccountIdBySock(sock transport.Socket) (int64, bool) {
	am.RLock()
	defer am.RUnlock()

	id, ok := am.mapSocks[sock]
	return id, ok
}

func (am *AccountManager) GetAccountById(acctId int64) *player.Account {
	acct, ok := am.cacheAccounts.Get(acctId)
	if ok {
		return acct.(*player.Account)
	}

	return nil
}

func (am *AccountManager) GetPlayerInfoById(playerId int64) *player.PlayerInfo {
	c, ok := am.cachePlayerInfos.Get(playerId)
	if ok {
		return c.(*player.PlayerInfo)
	}

	info := am.playerInfoPool.Get().(*player.PlayerInfo)
	info.Init()
	err := store.GetStore().FindOne(context.Background(), define.StoreType_Player, playerId, info)
	if err == nil {
		am.cachePlayerInfos.Set(playerId, info, PlayerInfoCacheExpire)
		return info
	}

	am.playerInfoPool.Put(info)
	return nil
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

// add handler to account's execute channel, will be dealed by account's run goroutine
func (am *AccountManager) AddAccountTask(ctx context.Context, acctId int64, fn task.TaskHandler, p ...interface{}) error {
	acct := am.GetAccountById(acctId)

	if acct == nil {
		return fmt.Errorf("AddAccountTask err:%w, account_id:%d", ErrAccountNotFound, acctId)
	}

	acct.AddTask(ctx, fn, p...)
	return nil
}

func (am *AccountManager) AddPlayerTask(ctx context.Context, playerId int64, fn task.TaskHandler, p ...interface{}) error {
	info := am.GetPlayerInfoById(playerId)
	if info == nil {
		return fmt.Errorf("error:%w, player_id:%d", ErrPlayerInfoNotFound, playerId)
	}

	return am.AddAccountTask(ctx, info.AccountID, fn, p...)
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
		return store.GetStore().UpdateOne(context.Background(), define.StoreType_Player, p.ID, p)
	})

	// errHandle(func() error {
	// 	return store.GetStore().UpdateOne(context.Background(), define.StoreType_Hero, p.ID, p.HeroManager())
	// })

	// errHandle(func() error {
	// 	return store.GetStore().UpdateOne(context.Background(), define.StoreType_Item, p.ID, p.ItemManager())
	// })

	errHandle(func() error {
		return store.GetStore().UpdateOne(context.Background(), define.StoreType_Token, p.ID, p.TokenManager())
	})

	errHandle(func() error {
		return store.GetStore().UpdateOne(context.Background(), define.StoreType_Fragment, p.ID, p.FragmentManager())
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

func (am *AccountManager) Broadcast(msg proto.Message) {
	am.cacheAccounts.Range(func(v interface{}) bool {
		acct := v.(*cache.Item).Object.(*player.Account)
		acct.AddTask(context.Background(), func(c context.Context, p ...interface{}) error {
			a := p[0].(*player.Account)
			message := p[1].(proto.Message)
			a.SendProtoMessage(message)
			return nil
		}, msg)
		return true
	})
}

func (am *AccountManager) Run(ctx context.Context) error {
	<-ctx.Done()
	log.Info().Msg("world session context done...")
	return nil
}
