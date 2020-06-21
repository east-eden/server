package gate

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/groupcache/lru"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/define"
	pbGate "github.com/yokaiio/yokai_server/proto/gate"
	"github.com/yokaiio/yokai_server/store"
	"github.com/yokaiio/yokai_server/utils"
)

var (
	defaultGameIDSyncTimer = 10 * time.Second
	maxUserLruCache        = 5000
)

type Metadata map[string]string

type GameSelector struct {
	userPool      sync.Pool
	userCache     *lru.Cache
	defaultGameID int32
	gameMetadatas map[int16]Metadata  // all game's metadata
	sectionGames  map[int16]([]int16) // map[section_id]game_ids
	syncTimer     *time.Timer

	wg utils.WaitGroupWrapper
	g  *Gate

	sync.RWMutex
}

func NewGameSelector(g *Gate, c *cli.Context) *GameSelector {
	gs := &GameSelector{
		g:             g,
		userCache:     lru.New(maxUserLruCache),
		defaultGameID: -1,
		gameMetadatas: make(map[int16]Metadata),
		sectionGames:  make(map[int16]([]int16)),
		syncTimer:     time.NewTimer(defaultGameIDSyncTimer),
	}

	// user pool new function
	gs.userPool.New = NewUserInfo

	// user cache evicted function
	gs.userCache.OnEvicted = gs.OnUserEvicted

	// add user store info
	store.GetStore().AddStoreInfo(define.StoreType_User, "user", "_id", "")

	// migrate users table
	if err := store.GetStore().MigrateDbTable("user", "account_id", "player_id"); err != nil {
		logger.Warning("migrate collection user failed:", err)
	}

	return gs
}

// user evicted callback
func (gs *GameSelector) OnUserEvicted(key lru.Key, value interface{}) {
	logger.WithFields(logger.Fields{
		"key":   key,
		"value": value,
	}).Info("user info evicted callback")

	gs.userPool.Put(value)
}

func (gs *GameSelector) syncDefaultGame() {
	defaultGameID := gs.g.mi.GetDefaultGameID()
	gameMetadatas := gs.g.mi.GetServiceMetadatas("yokai_game")

	gs.Lock()
	defer gs.Unlock()

	gs.sectionGames = make(map[int16]([]int16))
	atomic.StoreInt32(&gs.defaultGameID, int32(defaultGameID))
	gs.syncTimer.Reset(defaultGameIDSyncTimer)

	gs.gameMetadatas = make(map[int16]Metadata)
	for _, metadata := range gameMetadatas {
		if value, ok := metadata["gameId"]; ok {
			gameID, err := strconv.ParseInt(value, 10, 16)
			if err != nil {
				logger.Warn("convert game_id to int16 failed when call syncDefaultGame:", err)
				continue
			}

			gs.gameMetadatas[int16(gameID)] = metadata
		}
	}

	for gameID := range gs.gameMetadatas {
		sectionID := int16(gameID / 10)
		ids, ok := gs.sectionGames[sectionID]
		if !ok {
			gs.sectionGames[sectionID] = make([]int16, 0)
			gs.sectionGames[sectionID] = append(gs.sectionGames[sectionID], int16(gameID))
		} else {
			hit := false
			for _, v := range ids {
				if v == int16(gameID) {
					hit = true
					break
				}
			}

			if !hit {
				gs.sectionGames[sectionID] = append(gs.sectionGames[sectionID], int16(gameID))
			}
		}
	}
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

	gameID := atomic.LoadInt32(&gs.defaultGameID)
	if gameID == -1 {
		return nil, errors.New("cannot find default game_id")
	}

	user := gs.userPool.Get().(*UserInfo)
	user.UserID = userId
	user.AccountID = accountId
	user.GameID = int16(gameID)

	// add to lru cache
	gs.userCache.Add(user.UserID, user)

	// save to cache and database
	store.GetStore().SaveObject(define.StoreType_User, user)

	return user, nil
}

func (gs *GameSelector) SelectGame(userID string, userName string) (*UserInfo, Metadata) {
	userId, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		logger.Warn("invalid user_id when call SelectGame:", err)
		return nil, Metadata{}
	}

	userInfo, errUser := gs.loadUserInfo(userId)
	if errUser != nil {
		return userInfo, Metadata{}
	}

	// first find in game's gameMetadatas
	if mt, ok := gs.gameMetadatas[userInfo.GameID]; ok {
		return userInfo, mt
	}

	// previous game node offline, peek another game node in same section
	if ids, ok := gs.sectionGames[userInfo.GameID/10]; ok {
		if mt, ok := gs.gameMetadatas[ids[rand.Intn(len(ids))]]; ok {
			return userInfo, mt
		}
	}

	return userInfo, Metadata{}
}

func (gs *GameSelector) UpdateUserInfo(req *pbGate.UpdateUserInfoRequest) error {
	user, err := gs.getUserInfo(req.Info.UserId)
	if err != nil {
		return err
	}

	user.UserID = req.Info.UserId
	user.AccountID = req.Info.AccountId
	user.GameID = int16(req.Info.GameId)
	user.PlayerID = req.Info.PlayerId
	user.PlayerName = req.Info.PlayerName
	user.PlayerLevel = req.Info.PlayerLevel
	store.GetStore().SaveObject(define.StoreType_User, user)
	return nil
}

func (gs *GameSelector) Main(ctx context.Context) error {
	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal("Game Run() error:", err)
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
	for {
		select {
		case <-ctx.Done():
			logger.Print("game selector context done!")
			return nil
		case <-gs.syncTimer.C:
			gs.syncDefaultGame()
			//logger.WithFields(logger.Fields{
			//"metadata":     gs.gameMetadatas,
			//"section_game": gs.sectionGames,
			//}).Info("sync default game result")
		}
	}

	return nil
}

func (gs *GameSelector) Exit(ctx context.Context) {
	gs.wg.Wait()
	logger.Info("game selector exit...")
}
