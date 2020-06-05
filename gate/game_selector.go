package gate

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/store"
	"github.com/yokaiio/yokai_server/utils"
)

var defaultGameIDSyncTimer time.Duration = 10 * time.Second

type Metadata map[string]string

type GameSelector struct {
	userPool      sync.Pool
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
		defaultGameID: -1,
		gameMetadatas: make(map[int16]Metadata),
		sectionGames:  make(map[int16]([]int16)),
		syncTimer:     time.NewTimer(defaultGameIDSyncTimer),
	}

	// user pool new function
	gs.userPool.New = NewUserInfo

	// init users memory
	if err := g.store.AddMemExpire(c, store.ExpireType_User, &gs.userPool, userExpireTime); err != nil {
		logger.Warning("store add memory expire failed:", err)
	}

	// migrate users table
	if err := g.store.MigrateDbTable("users", "account_id", "player_id"); err != nil {
		logger.Warning("migrate collection user failed:", err)
	}

	return gs
}

func (gs *GameSelector) getMetadata(id int16) Metadata {
	gs.RLock()
	defer gs.RUnlock()
	return gs.gameMetadatas[id]
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

func (gs *GameSelector) SelectGame(userID string, userName string) (*UserInfo, Metadata) {
	userId, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		logger.Warn("invalid user_id when call SelectGame:", err)
		return nil, Metadata{}
	}

	// old user
	obj, err := gs.g.store.LoadObject(store.ExpireType_User, "_id", userId)
	if err == nil {
		userInfo := obj.(*UserInfo)
		gameID := userInfo.GameID

		// first find in game's gameMetadatas
		gs.RLock()
		if mt, ok := gs.gameMetadatas[gameID]; ok {
			gs.RUnlock()
			return userInfo, mt
		}

		// previous game node offline, peek another game node in same section
		if ids, ok := gs.sectionGames[gameID/10]; ok {
			if mt, ok := gs.gameMetadatas[ids[rand.Intn(len(ids))]]; ok {
				gs.RUnlock()
				return userInfo, mt
			}
		}

		gs.RUnlock()
		return userInfo, Metadata{}
	} else {

		// new user
		accountID, err := utils.NextID(define.SnowFlake_Account)
		if err != nil {
			logger.Warn("new user nextid error:", err)
			return nil, Metadata{}
		}

		// default game id
		gameID := atomic.LoadInt32(&gs.defaultGameID)

		if gameID == -1 {
			logger.Warn("cannot find default game_id")
			return nil, Metadata{}
		}

		user := gs.userPool.Get().(*UserInfo)
		user.UserID = userId
		user.AccountID = accountID
		user.GameID = int16(gameID)
		gs.g.store.SaveObject(store.ExpireType_User, user)

		return user, gs.getMetadata(user.GameID)
	}
}

func (gs *GameSelector) UpdateUserInfo(info *UserInfo) {
	gs.g.store.SaveObject(store.ExpireType_User, info)
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
