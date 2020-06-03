package gate

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/store/memory"
	"github.com/yokaiio/yokai_server/utils"
)

var userExpireTime time.Duration = 30 * time.Minute
var defaultGameIDSyncTimer time.Duration = 10 * time.Second

type UserInfo struct {
	UserID      int64       `bson:"_id" json:"_id" redis:"_id"`
	AccountID   int64       `bson:"account_id" json:"account_id" redis:"account_id"`
	GameID      int16       `bson:"game_id" json:"game_id" redis:"game_id"`
	PlayerID    int64       `bson:"player_id" json:"player_id" redis:"player_id"`
	PlayerName  string      `bson:"player_name" json:"player_name" redis:"player_name"`
	PlayerLevel int32       `bson:"player_level" json:"player_level" redis:"player_level"`
	Expire      *time.Timer `bson:"-" json:"-" redis:"-"`
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

	// init users memory
	if err := g.store.AddMemExpire(c, memory.MemExpireType_Users, NewUserInfo); err != nil {
		logger.Warning("store add memory expire failed:", err)
	}

	// migrate users table
	if err := g.store.MigrateDbTable("users", "account_id", "player_id"); err != nil {
		logger.Warning("migrate collection user failed:", err)
	}

	return gs
}

func (gs *GameSelector) newUserInfo(info *UserInfo, userId int64) {
	// create new user
	accountID, err := utils.NextID(define.SnowFlake_Account)
	if err != nil {
		logger.Warn("new user nextid error:", err)
		return
	}

	// default game id
	gs.RLock()
	gameID := gs.defaultGameID
	gs.RUnlock()

	if gameID == -1 {
		logger.Warn("cannot find default game_id")
		return
	}

	info.UserID = userId
	info.AccountID = accountID
	info.GameID = gameID

	gs.g.store.SaveObject(memory.MemExpireType_Users, info)
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
	userId, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		logger.Warn("invalid user_id when call SelectGame:", err)
		return nil, Metadata{}
	}

	// old user
	obj, err := gs.g.store.LoadObject(memory.MemExpireType_Users, "_id", userId)
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
		user := obj.(*UserInfo)
		if user == nil {
			return user, Metadata{}
		}

		// create new user info
		gs.newUserInfo(user, userId)
		return user, gs.getMetadata(user.GameID)
	}
}

func (gs *GameSelector) UpdateUserInfo(info *UserInfo) {
	gs.g.store.SaveObject(memory.MemExpireType_Users, info)
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
