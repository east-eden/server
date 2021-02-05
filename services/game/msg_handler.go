package game

import (
	pbGlobal "bitbucket.org/east-eden/server/proto/global"
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
	registerPbFn(&pbGlobal.C2S_WaitResponseMessage{}, m.handleWaitResponseMessage)
	registerPbFn(&pbGlobal.C2S_Ping{}, m.handleAccountPing)
	registerPbFn(&pbGlobal.C2S_AccountLogon{}, m.handleAccountLogon)
	registerPbFn(&pbGlobal.C2S_HeartBeat{}, m.handleHeartBeat)
	registerPbFn(&pbGlobal.C2S_AccountDisconnect{}, m.handleAccountDisconnect)

	// player
	registerPbFn(&pbGlobal.C2S_QueryPlayerInfo{}, m.handleQueryPlayerInfo)
	registerPbFn(&pbGlobal.C2S_CreatePlayer{}, m.handleCreatePlayer)
	registerPbFn(&pbGlobal.C2S_ChangeExp{}, m.handleChangeExp)
	registerPbFn(&pbGlobal.C2S_ChangeLevel{}, m.handleChangeLevel)
	registerPbFn(&pbGlobal.C2S_SyncPlayerInfo{}, m.handleSyncPlayerInfo)
	registerPbFn(&pbGlobal.C2S_PublicSyncPlayerInfo{}, m.handlePublicSyncPlayerInfo)

	// heros
	registerPbFn(&pbGlobal.C2S_AddHero{}, m.handleAddHero)
	registerPbFn(&pbGlobal.C2S_DelHero{}, m.handleDelHero)
	registerPbFn(&pbGlobal.C2S_QueryHeros{}, m.handleQueryHeros)
	//m.r.RegisterMessage("game.MC_HeroAddExp", &pbGame.MC_HeroAddExp{}, m.handleHeroAddExp)
	//m.r.RegisterMessage("game.MC_HeroAddLevel", &pbGame.MC_HeroAddLevel{}, m.handleHeroAddLevel)

	// items & equips
	registerPbFn(&pbGlobal.C2S_AddItem{}, m.handleAddItem)
	registerPbFn(&pbGlobal.C2S_DelItem{}, m.handleDelItem)
	registerPbFn(&pbGlobal.C2S_UseItem{}, m.handleUseItem)
	registerPbFn(&pbGlobal.C2S_QueryItems{}, m.handleQueryItems)

	registerPbFn(&pbGlobal.C2S_PutonEquip{}, m.handlePutonEquip)
	registerPbFn(&pbGlobal.C2S_TakeoffEquip{}, m.handleTakeoffEquip)

	// tokens
	registerPbFn(&pbGlobal.C2S_AddToken{}, m.handleAddToken)
	registerPbFn(&pbGlobal.C2S_QueryTokens{}, m.handleQueryTokens)

	// talent
	registerPbFn(&pbGlobal.C2S_AddTalent{}, m.handleAddTalent)
	registerPbFn(&pbGlobal.C2S_QueryTalents{}, m.handleQueryTalents)

	// rune
	registerPbFn(&pbGlobal.C2S_AddRune{}, m.handleAddRune)
	registerPbFn(&pbGlobal.C2S_DelRune{}, m.handleDelRune)
	registerPbFn(&pbGlobal.C2S_QueryRunes{}, m.handleQueryRunes)
	registerPbFn(&pbGlobal.C2S_PutonRune{}, m.handlePutonRune)
	registerPbFn(&pbGlobal.C2S_TakeoffRune{}, m.handleTakeoffRune)

	// scene
	registerPbFn(&pbGlobal.C2S_StartStageCombat{}, m.handleStartStageCombat)
}
