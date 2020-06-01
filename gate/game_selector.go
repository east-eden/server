package gate

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/utils"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var userExpireTime time.Duration = 30 * time.Minute
var defaultGameIDSyncTimer time.Duration = 10 * time.Second

type UserInfo struct {
	UserID      int64       `bson:"_id"`
	AccountID   int64       `bson:"account_id"`
	GameID      int16       `bson:"game_id"`
	PlayerID    int64       `bson:"player_id"`
	PlayerName  string      `bson:"player_name"`
	PlayerLevel int32       `bson:"player_level"`
	Expire      *time.Timer `bson:"-"`
}

func (u *UserInfo) ToJson() []byte {
	data, err := json.Marshal(u)
	if err != nil {
		return []byte("")
	}

	return data
}

func (u *UserInfo) TableName() string {
	return "users"
}

func (u *UserInfo) GetObjID() interface{} {
	return u.UserID
}

func (u *UserInfo) GetExpire() *time.Timer {
	return u.Expire
}

func (u *UserInfo) ResetExpire() {
	u.Expire.Reset(userExpireTime)
}

func (u *UserInfo) StopExpire() {
	u.Expire.Stop()
}

func NewUserInfo() interface{} {
	return &UserInfo{
		UserID:      -1,
		AccountID:   -1,
		GameID:      int16(-1),
		PlayerID:    -1,
		PlayerName:  "",
		PlayerLevel: 1,
		Expire:      time.NewTimer(userExpireTime),
	}
}

type Metadata map[string]string

type GameSelector struct {
	cacheUsers    *utils.CacheLoader
	cacheCancel   context.CancelFunc
	defaultGameID int16
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

	if err := g.store.MigrateCollection("users", "account_id", "player_id"); err != nil {
		logger.Warning("migrate collection user failed:", err)
	}

	gs.cacheUsers = utils.NewCacheLoader(
		g.store.GetCollection("users"),
		"_id",
		NewUserInfo,
		nil,
	)

	return gs
}

func (gs *GameSelector) newUser(userID int64) *UserInfo {
	// create new user
	accountID, err := utils.NextID(define.SnowFlake_Account)
	if err != nil {
		logger.Warn("new user nextid error:", err)
		return nil
	}

	// default game id
	gs.RLock()
	gameID := gs.defaultGameID
	gs.RUnlock()

	if gameID == -1 {
		logger.Warn("cannot find default game_id")
		return nil
	}

	newUser := NewUserInfo().(*UserInfo)
	newUser.UserID = userID
	newUser.AccountID = accountID
	newUser.GameID = gameID
	gs.save(newUser)

	return newUser
}

func (gs *GameSelector) getMetadata(id int16) Metadata {
	gs.RLock()
	defer gs.RUnlock()
	return gs.gameMetadatas[id]
}

func (gs *GameSelector) save(u *UserInfo) {
	// memory cache store
	gs.cacheUsers.Store(u)

	// store to cache
	gs.g.store.CacheDoAsync("SET", func(reply interface{}, err error) {
		if err != nil {
			logger.WithFields(logger.Fields{
				"user_info": u,
				"error":     err,
			}).Error("store user info to cache failed")
		}
	}, u.UserID, u.ToJson())

	// store to database
	filter := bson.D{{"_id", u.UserID}}
	update := bson.D{{"$set", u}}
	opts := options.Update().SetUpsert(true)
	timeout, _ := context.WithTimeout(context.Background(), time.Second*5)
	_, err := gs.g.store.CollationUpdate(timeout, u.TableName(), filter, update, opts)
	if err != nil {
		logger.Warning("collation update failed:", err)
	}
}

func (gs *GameSelector) syncDefaultGame() {
	defaultGameID := gs.g.mi.GetDefaultGameID()
	gameMetadatas := gs.g.mi.GetServiceMetadatas("yokai_game")

	gs.Lock()
	defer gs.Unlock()

	gs.sectionGames = make(map[int16]([]int16))
	gs.defaultGameID = defaultGameID
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
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		logger.Warn("invalid user_id when call SelectGame:", err)
		return nil, Metadata{}
	}

	// old user
	if obj := gs.cacheUsers.Load(id); obj != nil {
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
	}

	// create new user
	user := gs.newUser(id)
	if user == nil {
		return user, Metadata{}
	}

	return user, gs.getMetadata(user.GameID)
}

func (gs *GameSelector) UpdateUserInfo(info *UserInfo) {
	if obj := gs.cacheUsers.Load(info.UserID); obj == nil {
		gs.cacheUsers.Store(info)
	}

	gs.save(info)
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

	// cache loader
	var cacheCtx context.Context
	cacheCtx, gs.cacheCancel = context.WithCancel(ctx)
	gs.wg.Wrap(func() {
		gs.cacheUsers.Run(cacheCtx)
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
	gs.cacheCancel()
	gs.wg.Wait()
	logger.Info("game selector exit...")
}
