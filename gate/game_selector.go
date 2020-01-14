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
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
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
	defaultGameID int16
	gameMetadatas map[int16]Metadata  // all game's metadata
	sectionGames  map[int16]([]int16) // map[section_id]game_ids
	syncTimer     *time.Timer

	ctx    context.Context
	cancel context.CancelFunc
	g      *Gate

	coll *mongo.Collection
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

	gs.ctx, gs.cancel = context.WithCancel(c)

	gs.migrate()

	gs.cacheUsers = utils.NewCacheLoader(
		gs.ctx,
		gs.coll,
		"_id",
		10000,
		NewUserInfo,
		nil,
	)

	return gs
}

func (gs *GameSelector) migrate() {
	gs.coll = gs.g.ds.Database().Collection("users")

	// check index
	idx := gs.coll.Indexes()

	opts := options.ListIndexes().SetMaxTime(2 * time.Second)
	cursor, err := idx.List(context.Background(), opts)
	if err != nil {
		log.Fatal(err)
	}

	accountIndex := false
	playerIndex := false
	for cursor.Next(context.Background()) {
		var result bson.M
		cursor.Decode(&result)
		if result["name"] == "account_id" {
			accountIndex = true
		}

		if result["name"] == "player_id" {
			playerIndex = true
		}
	}

	// create index
	if !accountIndex {
		_, err := gs.coll.Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys:    bsonx.Doc{{"account_id", bsonx.Int32(1)}},
				Options: options.Index().SetName("account_id"),
			},
		)

		if err != nil {
			logger.Warn("collection users create index account_id failed:", err)
		}
	}

	if !playerIndex {
		_, err := gs.coll.Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys:    bsonx.Doc{{"player_id", bsonx.Int32(1)}},
				Options: options.Index().SetName("player_id"),
			},
		)

		if err != nil {
			logger.Warn("collection users create index player_id failed:", err)
		}
	}
}

func (gs *GameSelector) peekGameBySection(section int16) int16 {
	gs.RLock()
	defer gs.RUnlock()

	ids, ok := gs.sectionGames[section]
	if !ok {
		return -1
	}

	return ids[rand.Intn(len(ids))]
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
	gs.cacheUsers.Store(newUser)

	return newUser
}

func (gs *GameSelector) getMetadata(id int16) Metadata {
	gs.RLock()
	defer gs.RUnlock()
	return gs.gameMetadatas[id]
}

func (gs *GameSelector) save(u *UserInfo) {
	filter := bson.D{{"_id", u.UserID}}
	update := bson.D{{"$set", u}}
	res, err := gs.coll.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	logger.Info("user save result:", res, err)
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
		if value, ok := metadata["game_id"]; ok {
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
		if mt, ok := gs.gameMetadatas[gs.peekGameBySection(gameID/10)]; ok {
			gs.RUnlock()
			return userInfo, mt
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

func (gs *GameSelector) Run() error {
	for {
		select {
		case <-gs.ctx.Done():
			logger.Print("game selector context done!")
			return nil
		case <-gs.syncTimer.C:
			gs.syncDefaultGame()
		}
	}

	return nil
}
