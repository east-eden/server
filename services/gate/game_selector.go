package gate

import (
	"context"
	"strconv"
	"sync"

	"github.com/east-eden/server/define"
	pbGate "github.com/east-eden/server/proto/gate"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	"github.com/golang/groupcache/lru"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"stathat.com/c/consistent"
)

var (
	maxUserLruCache = 5000
	maxGameNode     = 200 // max game node number, used in constent hash
)

type Metadata map[string]string

type GameSelector struct {
	userPool      sync.Pool
	userCache     *lru.Cache
	gameMetadatas map[int16]Metadata  // all game's metadata
	sectionGames  map[int16]([]int16) // map[section_id]game_ids

	wg utils.WaitGroupWrapper
	g  *Gate

	sync.RWMutex

	consistent *consistent.Consistent
}

func NewGameSelector(g *Gate, c *cli.Context) *GameSelector {
	gs := &GameSelector{
		g:             g,
		userCache:     lru.New(maxUserLruCache),
		gameMetadatas: make(map[int16]Metadata),
		sectionGames:  make(map[int16]([]int16)),
		consistent:    consistent.New(),
	}

	// constent hash node number
	gs.consistent.NumberOfReplicas = maxGameNode

	// user pool new function
	gs.userPool.New = NewUserInfo

	// user cache evicted function
	gs.userCache.OnEvicted = gs.OnUserEvicted

	// add user store info
	store.GetStore().AddStoreInfo(define.StoreType_User, "user", "_id", "")

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

// func (gs *GameSelector) SyncDefaultGame() {
// 	defaultGameID, _ := gs.g.mi.SelectGameEntry()
// 	gameMetadatas := gs.g.mi.GetServiceMetadatas("game")

// 	gs.Lock()
// 	defer gs.Unlock()

// 	gs.sectionGames = make(map[int16]([]int16))
// 	atomic.StoreInt32(&gs.defaultGameID, int32(defaultGameID))

// 	gs.gameMetadatas = make(map[int16]Metadata)
// 	for _, metadata := range gameMetadatas {
// 		if value, ok := metadata["gameId"]; ok {
// 			gameID, err := strconv.ParseInt(value, 10, 16)
// 			if err != nil {
// 				log.Warn().
// 					Err(err).
// 					Msg("convert game_id to int16 failed when call syncDefaultGame")
// 				continue
// 			}

// 			gs.gameMetadatas[int16(gameID)] = metadata
// 		}
// 	}

// 	for gameID := range gs.gameMetadatas {
// 		sectionID := int16(gameID / 10)
// 		ids, ok := gs.sectionGames[sectionID]
// 		if !ok {
// 			gs.sectionGames[sectionID] = make([]int16, 0)
// 			gs.sectionGames[sectionID] = append(gs.sectionGames[sectionID], int16(gameID))
// 		} else {
// 			hit := false
// 			for _, v := range ids {
// 				if v == int16(gameID) {
// 					hit = true
// 					break
// 				}
// 			}

// 			if !hit {
// 				gs.sectionGames[sectionID] = append(gs.sectionGames[sectionID], int16(gameID))
// 			}
// 		}
// 	}
// }

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
	err := store.GetStore().LoadObject(define.StoreType_User, userId, obj.(store.StoreObjector))
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
	if err := store.GetStore().SaveObject(define.StoreType_User, user); err != nil {
		return user, err
	}

	return user, nil
}

func (gs *GameSelector) SelectGame(userID string, userName string) (*UserInfo, Metadata) {
	userId, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("invalid user_id when call SelectGame")
		return nil, Metadata{}
	}

	userInfo, errUser := gs.loadUserInfo(userId)
	if errUser != nil {
		return userInfo, Metadata{}
	}

	// every time select calls, consistent hash will be refreshed
	next, err := gs.g.mi.srv.Client().Options().Selector.Select("game", utils.ConsistentHashSelector(gs.consistent, strconv.Itoa(int(userId))))
	if err != nil {
		log.Warn().Err(err).Msg("select game failed")
		return nil, Metadata{}
	}

	node, err := next()
	if err != nil {
		log.Warn().Err(err).Msg("get next node failed")
		return nil, Metadata{}
	}

	log.Info().Interface("node", node).Msg("select game node success")
	return userInfo, node.Metadata
}

func (gs *GameSelector) UpdateUserInfo(req *pbGate.UpdateUserInfoRequest) error {
	user, err := gs.getUserInfo(req.Info.UserId)
	if err != nil {
		return err
	}

	user.UserID = req.Info.UserId
	user.AccountID = req.Info.AccountId
	user.PlayerID = req.Info.PlayerId
	user.PlayerName = req.Info.PlayerName
	user.PlayerLevel = req.Info.PlayerLevel
	return store.GetStore().SaveObject(define.StoreType_User, user)
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
