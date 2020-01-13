package gate

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"gopkg.in/mgo.v2/bson"
)

var userExpireTime time.Duration = 30 * time.Minute
var defaultGameIDSyncTimer time.Duration = 10 * time.Second

type UserInfo struct {
	UserID      int64       `bson:"_id"`
	AccountID   int64       `bson:"account_id"`
	GameID      uint16      `bson:"game_id"`
	PlayerID    int64       `bson:"player_id"`
	PlayerName  string      `bson:"player_name"`
	PlayerLevel int32       `bson:"player_level"`
	Expire      *time.Timer `bson:"-"`
}

func (u *UserInfo) GetObjID() int64 {
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

func NewUserInfo() *UserInfo {
	return &UserInfo{
		UserID:      -1,
		AccountID:   -1,
		GameID:      -1,
		PlayerID:    -1,
		PlayerName:  "",
		PlayerLevel: 1,
		Expire:      time.NewTimer(userExpireTime),
	}
}

type Metadata map[string]string

type GameSelector struct {
	cacheUsers    *utils.CacheLoader
	defaultGameID uint16
	gameMetadatas map[uint16]Metadata   // all game's metadata
	sectionGames  map[uint16]([]uint16) // map[section_id]game_ids
	syncTimer     *time.Timer

	ctx    context.Context
	cancel context.CancelFunc
	g      *Gate

	coll *mongo.Collection
	sync.RWMutex
}

func NewGameSelector(g *Gate, c *cli.Context) *HttpServer {
	gs := &GameSelector{
		g:             g,
		defaultGameID: uint16(-1),
		gameMetadatas: make(map[uint16]Metadata),
		sectionGames:  make(map[uint16]([]uint16)),
		syncTimer:     time.NewTimer(defaultGameIDSyncTimer),
	}

	gs.ctx, gs.cancel = context.WithCancel(c)

	cacheUsers = utils.NewCacheLoader(
		gs.ctx,
		coll,
		"_id",
		10000,
		NewUserInfo,
		nil,
	)

	return gs
}

func (gs *GameSelector) migrate() {
	m.coll = gs.g.ds.Database().Collection("users")

	// check index
	idx := m.coll.Indexes()

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
		_, err := coll.Indexes().CreateOne(
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
		_, err := coll.Indexes().CreateOne(
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

func (gs *GameSelector) peekGameBySection(section uint16) uint16 {
	gs.RLock()
	defer gs.RUnlock()

	ids, ok := gs.sectionGames[section]
	if !ok {
		return uint16(-1)
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

	newUser := NewUserInfo()
	newUser.UserID = userID
	newUser.AccountID = accountID
	newUser.GameID = gameID
	gs.save(newUser)
	gs.cacheUsers.Store(newUser)

	return newUser
}

func (gs *GameSelector) getMetadata(id uint16) Metadata {
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
	gameMetadatas := gs.g.mi.GetServicegameMetadatas()

	gs.Lock()
	defer gs.Unlock()

	gs.defaultGameID = defaultGameID
	gs.gameMetadatas = gameMetadatas
	gs.syncTimer.Reset(defaultGameIDSyncTimer)
	for gameID, data := range gs.gameMetadatas {
		ids, ok := gs.sectionGames[gameID/10]
		if !ok {
			gs.sectionGames[gameID/10] = make([]uint16, 0)
		}
	}
}

func (gs *GameSelector) SelectGame(userID int64) Metadata {
	// old user
	if obj := gs.cacheUsers.Load(userID); obj != nil {
		gameID := obj.(*UserInfo).GameID

		// first find in game's gameMetadatas
		gs.RLock()
		if mt, ok := gs.gameMetadatas[gameID]; ok {
			gs.RUnlock()
			return mt
		}

		// previous game node offline, peek another game node in same section
		if mt, ok := gs.gameMetadatas[gs.peekGameBySection(gameID/10)]; ok {
			gs.RUnlock()
			return mt
		}

		gs.RUnlock()
		return Metadata{}
	}

	// create new user
	user := gs.newUser(userID)
	if user == nil {
		return Metadata{}
	}

	return gs.getMetadata(user.GameID)
}

func (gs *GameSelector) Run() error {
	for {
		select {
		case <-gs.ctx.Done():
			logger.Print("game selector context done!")
			return nil
		case <-gs.syncTimer.C:
			gs.syncDefaultGame()
			gs.syncTimer.Reset(defaultGameIDSyncTimer)
		}
	}

	return nil
}
