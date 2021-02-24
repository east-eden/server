package game

import (
	"context"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/player"
	"github.com/east-eden/server/transport"
	"github.com/east-eden/server/transport/codec"
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

	// normal protobuf handler
	registerPBHandler := func(p proto.Message, mf transport.MessageFunc) {
		err := m.r.RegisterProtobufMessage(p, mf)
		if err != nil {
			log.Fatal().
				Err(err).
				Str("name", string(proto.MessageReflect(p).Descriptor().Name())).
				Msg("register message failed")
		}
	}

	// account protobuf handler
	registerPBAccountHandler := func(p proto.Message, handle player.SlowHandleFunc) {
		mf := func(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
			m.g.am.AccountSlowHandle(sock, &player.AccountSlowHandler{
				F: handle,
				M: msg,
			})
			return nil
		}

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
	registerPBHandler(&pbGlobal.C2S_WaitResponseMessage{}, m.handleWaitResponseMessage)
	registerPBHandler(&pbGlobal.C2S_Ping{}, m.handleAccountPing)
	registerPBHandler(&pbGlobal.C2S_AccountLogon{}, m.handleAccountLogon)
	registerPBHandler(&pbGlobal.C2S_HeartBeat{}, m.handleHeartBeat)
	registerPBHandler(&pbGlobal.C2S_AccountDisconnect{}, m.handleAccountDisconnect)

	// player
	registerPBAccountHandler(&pbGlobal.C2S_QueryPlayerInfo{}, m.handleQueryPlayerInfo)
	registerPBAccountHandler(&pbGlobal.C2S_CreatePlayer{}, m.handleCreatePlayer)
	registerPBAccountHandler(&pbGlobal.C2S_ChangeExp{}, m.handleChangeExp)
	registerPBAccountHandler(&pbGlobal.C2S_ChangeLevel{}, m.handleChangeLevel)
	registerPBAccountHandler(&pbGlobal.C2S_SyncPlayerInfo{}, m.handleSyncPlayerInfo)
	registerPBAccountHandler(&pbGlobal.C2S_PublicSyncPlayerInfo{}, m.handlePublicSyncPlayerInfo)

	// heros
	registerPBAccountHandler(&pbGlobal.C2S_AddHero{}, m.handleAddHero)
	registerPBAccountHandler(&pbGlobal.C2S_DelHero{}, m.handleDelHero)
	registerPBAccountHandler(&pbGlobal.C2S_QueryHeros{}, m.handleQueryHeros)

	// fragment
	registerPBAccountHandler(&pbGlobal.C2S_QueryFragments{}, m.handleQueryFragments)
	registerPBAccountHandler(&pbGlobal.C2S_FragmentsCompose{}, m.handleFragmentsCompose)

	// items & equips
	registerPBAccountHandler(&pbGlobal.C2S_AddItem{}, m.handleAddItem)
	registerPBAccountHandler(&pbGlobal.C2S_DelItem{}, m.handleDelItem)
	registerPBAccountHandler(&pbGlobal.C2S_UseItem{}, m.handleUseItem)
	registerPBAccountHandler(&pbGlobal.C2S_QueryItems{}, m.handleQueryItems)
	registerPBAccountHandler(&pbGlobal.C2S_EquipLevelup{}, m.handleEquipLevelup)

	registerPBAccountHandler(&pbGlobal.C2S_PutonEquip{}, m.handlePutonEquip)
	registerPBAccountHandler(&pbGlobal.C2S_TakeoffEquip{}, m.handleTakeoffEquip)

	// tokens
	registerPBAccountHandler(&pbGlobal.C2S_AddToken{}, m.handleAddToken)
	registerPBAccountHandler(&pbGlobal.C2S_QueryTokens{}, m.handleQueryTokens)

	// talent
	registerPBAccountHandler(&pbGlobal.C2S_AddTalent{}, m.handleAddTalent)
	registerPBAccountHandler(&pbGlobal.C2S_QueryTalents{}, m.handleQueryTalents)

	// rune
	registerPBAccountHandler(&pbGlobal.C2S_AddRune{}, m.handleAddRune)
	registerPBAccountHandler(&pbGlobal.C2S_DelRune{}, m.handleDelRune)
	registerPBAccountHandler(&pbGlobal.C2S_QueryRunes{}, m.handleQueryRunes)
	registerPBAccountHandler(&pbGlobal.C2S_PutonRune{}, m.handlePutonRune)
	registerPBAccountHandler(&pbGlobal.C2S_TakeoffRune{}, m.handleTakeoffRune)

	// scene
	registerPBAccountHandler(&pbGlobal.C2S_StartStageCombat{}, m.handleStartStageCombat)
}
