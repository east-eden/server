package game

import (
	pbGlobal "bitbucket.org/east-eden/server/proto/global"
	pbGame "bitbucket.org/east-eden/server/proto/server/game"
	"bitbucket.org/east-eden/server/transport"
	"bitbucket.org/east-eden/server/transport/codec"
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/rs/zerolog/log"
)

type MsgHandler struct {
	g             *Game
	r             transport.Register
	timeHistogram *prometheus.HistogramVec
}

func NewMsgHandler(g *Game) *MsgHandler {
	m := &MsgHandler{
		g: g,
		r: transport.NewTransportRegister(),
		timeHistogram: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "account",
				Name:      "handle_latency",
				Help:      "account goroutine handle latency",
			},
			[]string{"method"},
		),
	}

	m.registerAllMessage()
	return m
}

type MC_AccountTest struct {
	AccountId int64  `protobuf:"varint,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (msg *MC_AccountTest) GetName() string {
	return "MC_AccountTest"
}

func (m *MsgHandler) registerAllMessage() {
	registerJsonFn := func(c codec.JsonCodec, mf transport.MessageFunc) {
		err := m.r.RegisterJsonMessage(c, mf)
		if err != nil {
			log.Fatal().
				Err(err).
				Str("name", c.GetName()).
				Msg("register message failed")
		}
	}

	// json
	registerJsonFn(&MC_AccountTest{}, m.handleAccountTest)

	registerPbFn := func(p proto.Message, mf transport.MessageFunc) {
		err := m.r.RegisterProtobufMessage(p, mf)
		if err != nil {
			log.Fatal().
				Err(err).
				Str("name", string(proto.MessageReflect(p).Descriptor().Name())).
				Msg("register message failed")
		}
	}

	// protobuf

	// account
	registerPbFn(&pbGlobal.C2M_Ping{}, m.handleAccountPing)
	registerPbFn(&pbGlobal.C2M_AccountLogon{}, m.handleAccountLogon)
	registerPbFn(&pbGlobal.C2M_HeartBeat{}, m.handleHeartBeat)
	registerPbFn(&pbGlobal.C2M_AccountDisconnect{}, m.handleAccountDisconnect)

	// player
	registerPbFn(&pbGame.C2M_QueryPlayerInfo{}, m.handleQueryPlayerInfo)
	registerPbFn(&pbGame.C2M_CreatePlayer{}, m.handleCreatePlayer)
	registerPbFn(&pbGame.MC_SelectPlayer{}, m.handleSelectPlayer)
	registerPbFn(&pbGame.C2M_ChangeExp{}, m.handleChangeExp)
	registerPbFn(&pbGame.C2M_ChangeLevel{}, m.handleChangeLevel)
	registerPbFn(&pbGame.C2M_SyncPlayerInfo{}, m.handleSyncPlayerInfo)
	registerPbFn(&pbGame.C2M_PublicSyncPlayerInfo{}, m.handlePublicSyncPlayerInfo)

	// heros
	registerPbFn(&pbGame.C2M_AddHero{}, m.handleAddHero)
	registerPbFn(&pbGame.C2M_DelHero{}, m.handleDelHero)
	registerPbFn(&pbGame.C2M_QueryHeros{}, m.handleQueryHeros)
	//m.r.RegisterMessage("game.MC_HeroAddExp", &pbGame.MC_HeroAddExp{}, m.handleHeroAddExp)
	//m.r.RegisterMessage("game.MC_HeroAddLevel", &pbGame.MC_HeroAddLevel{}, m.handleHeroAddLevel)

	// items & equips
	registerPbFn(&pbGame.C2M_AddItem{}, m.handleAddItem)
	registerPbFn(&pbGame.C2M_DelItem{}, m.handleDelItem)
	registerPbFn(&pbGame.C2M_UseItem{}, m.handleUseItem)
	registerPbFn(&pbGame.C2M_QueryItems{}, m.handleQueryItems)

	registerPbFn(&pbGame.C2M_PutonEquip{}, m.handlePutonEquip)
	registerPbFn(&pbGame.C2M_TakeoffEquip{}, m.handleTakeoffEquip)

	// tokens
	registerPbFn(&pbGame.C2M_AddToken{}, m.handleAddToken)
	registerPbFn(&pbGame.C2M_QueryTokens{}, m.handleQueryTokens)

	// talent
	registerPbFn(&pbGame.C2M_AddTalent{}, m.handleAddTalent)
	registerPbFn(&pbGame.C2M_QueryTalents{}, m.handleQueryTalents)

	// rune
	registerPbFn(&pbGame.C2M_AddRune{}, m.handleAddRune)
	registerPbFn(&pbGame.C2M_DelRune{}, m.handleDelRune)
	registerPbFn(&pbGame.C2M_QueryRunes{}, m.handleQueryRunes)
	registerPbFn(&pbGame.C2M_PutonRune{}, m.handlePutonRune)
	registerPbFn(&pbGame.C2M_TakeoffRune{}, m.handleTakeoffRune)

	// scene
	registerPbFn(&pbGame.C2M_StartStageCombat{}, m.handleStartStageCombat)
}
