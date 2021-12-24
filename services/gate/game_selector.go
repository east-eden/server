package gate

import (
	"context"
	"hash/crc32"
	"sync"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	"github.com/golang/groupcache/lru"
	log "github.com/rs/zerolog/log"
	"github.com/spf13/cast"
	"github.com/urfave/cli/v2"
	"stathat.com/c/consistent"
)

var (
	maxUserLruCache = 5000
	maxGameNode     = 128 // max game node number, used in constent hash
)

type Metadata map[string]string

type GameSelector struct {
	userPool      sync.Pool
	userCache     *lru.Cache
	gameMetadatas map[int16]Metadata // all game's metadata

	wg utils.WaitGroupWrapper
	g  *Gate

	sync.RWMutex

	consistent *consistent.Consistent
}

func NewGameSelector(c *cli.Context, g *Gate) *GameSelector {
	gs := &GameSelector{
		g:             g,
		userCache:     lru.New(maxUserLruCache),
		gameMetadatas: make(map[int16]Metadata),
		consistent:    consistent.New(),
	}

	// constent hash node number
	gs.consistent.NumberOfReplicas = maxGameNode

	// user pool new function
	gs.userPool.New = NewUserInfo

	// user cache evicted function
	gs.userCache.OnEvicted = gs.OnUserEvicted

	// add user store info
	store.GetStore().AddStoreInfo(define.StoreType_User, "user", "_id")

	// migrate users table
	if err := store.GetStore().MigrateDbTable("user", "account_id", "player_id"); err != nil {
		log.Warn().
			Err(err).
			Msg("migrate collection user failed")
	}

	return gs
}

// user evicted callback
func (gs *GameSelector) OnUserEvicted(key lru.Key, value interface{}) {
	log.Info().
		Interface("key", key).
		Interface("value", value).
		Msg("user info evicted callback")

	gs.userPool.Put(value)
}

func (gs *GameSelector) getUserInfo(userId int64) (*UserInfo, error) {
	gs.RLock()
	defer gs.RUnlock()

	// find in lru cache
	obj, ok := gs.userCache.Get(userId)
	if ok {
		return obj.(*UserInfo), nil
	}

	// find in store
	obj = gs.userPool.Get()
	err := store.GetStore().FindOne(context.Background(), define.StoreType_User, userId, obj)
	if err == nil {
		return obj.(*UserInfo), nil
	}

	gs.userPool.Put(obj)
	return nil, err
}

func (gs *GameSelector) loadUserInfo(userId int64) (*UserInfo, error) {
	// get old user
	if user, err := gs.getUserInfo(userId); err == nil {
		return user, nil
	}

	gs.Lock()
	defer gs.Unlock()

	// create new user
	accountId, err := utils.NextID(define.SnowFlake_Account)
	if err != nil {
		return nil, err
	}

	user := gs.userPool.Get().(*UserInfo)
	user.UserID = userId
	user.AccountID = accountId

	// add to lru cache
	gs.userCache.Add(user.UserID, user)

	// save to cache and database
	if err := store.GetStore().UpdateOne(context.Background(), define.StoreType_User, user.UserID, user, true); err != nil {
		return user, err
	}

	return user, nil
}

func (gs *GameSelector) SelectGame(userID string) (*UserInfo, Metadata) {
	// todo userId 暂时为userID(string)的crc32
	userId := crc32.ChecksumIEEE([]byte(userID))
	userInfo, errUser := gs.loadUserInfo(int64(userId))
	if errUser != nil {
		return userInfo, Metadata{}
	}

	// every time select calls, consistent hash will be refreshed
	next, err := gs.g.mi.srv.Client().Options().Selector.Select("game", utils.ConsistentHashSelector(gs.consistent, cast.ToString(userId)))
	if !utils.ErrCheck(err, "select game failed") {
		return nil, Metadata{}
	}

	node, err := next()
	if !utils.ErrCheck(err, "get next node failed") {
		return nil, Metadata{}
	}

	log.Info().Interface("node", node).Msg("select game node success")
	return userInfo, node.Metadata
}

func (gs *GameSelector) Main(ctx context.Context) error {
	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("GameSelector Main() error")
			}
			exitCh <- err
		})
	}

	gs.wg.Wrap(func() {
		defer utils.CaptureException()
		exitFunc(gs.Run(ctx))
	})

	return <-exitCh
}

func (gs *GameSelector) Run(ctx context.Context) error {
	<-ctx.Done()
	log.Info().Msg("game selector context done...")
	return nil
}

func (gs *GameSelector) Exit(ctx context.Context) {
	gs.wg.Wait()
	log.Info().Msg("game selector exit...")
}
