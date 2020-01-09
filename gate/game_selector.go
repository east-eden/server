package gate

import (
	"context"
	"sync"

	"github.com/grafana/grafana/pkg/cmd/grafana-cli/logger"
	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type UserInfo struct {
	UserID      int64  `bson:"_id"`
	AccountID   int64  `bson:"account_id"`
	PlayerID    int64  `bson:"player_id"`
	PlayerName  string `bson:"player_name"`
	PlayerLevel int32  `bson:"player_level"`
}

type GameSelector struct {
	mapUsers map[int64]*UserInfo

	ctx    context.Context
	cancel context.CancelFunc
	g      *Gate

	coll *mongo.Collection
	sync.RWMutex
}

func NewGameSelector(g *Gate, c *cli.Context) *HttpServer {
	gs := &GameSelector{
		g: g,
	}

	gs.ctx, gs.cancel = context.WithCancel(c)
	return gs
}

func (gs *GameSelector) migrate() {
	m.coll = gs.g.ds.Database().Collection("users")

	// create index
	_, err := m.coll.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bsonx.Doc{{"account_id", bsonx.Int32(1)}},
		},
	)

	if err != nil {
		logger.Warn("player manager create index failed:", err)
	}

	//player.Migrate(ds)
}
