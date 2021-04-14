package game

import (
	"context"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/transport"
	"bitbucket.org/funplus/server/transport/codec"
	"github.com/golang/protobuf/proto"
	"github.com/hellodudu/task"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/rs/zerolog/log"
)

type MsgRegister struct {
	am            *AccountManager
	rpcHandler    *RpcHandler
	pubSub        *PubSub
	r             transport.Register
	timeHistogram *prometheus.HistogramVec
}

func NewMsgRegister(am *AccountManager, rpcHandler *RpcHandler, pubSub *PubSub) *MsgRegister {
	m := &MsgRegister{
		am:         am,
		rpcHandler: rpcHandler,
		pubSub:     pubSub,
		r:          transport.NewTransportRegister(),
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

func (m *MsgRegister) registerAllMessage() {
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
	registerPBAccountHandler := func(p proto.Message, handle task.TaskHandler) {
		mf := func(ctx context.Context, sock transport.Socket, msg *transport.Message) error {
			return m.am.AddAccountTask(
				ctx,
				m.am.GetAccountIdBySock(sock),
				handle,
				msg.Body.(proto.Message),
			)
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
	registerPBAccountHandler(&pbGlobal.C2S_CreatePlayer{}, m.handleCreatePlayer)
	registerPBAccountHandler(&pbGlobal.C2S_GmCmd{}, m.handleGmCmd)
	registerPBAccountHandler(&pbGlobal.C2S_WithdrawStrengthen{}, m.handleWithdrawStrengthen)
	registerPBAccountHandler(&pbGlobal.C2S_BuyStrengthen{}, m.handleBuyStrengthen)

	// heros
	registerPBAccountHandler(&pbGlobal.C2S_DelHero{}, m.handleDelHero)
	registerPBAccountHandler(&pbGlobal.C2S_QueryHeros{}, m.handleQueryHeros)
	registerPBAccountHandler(&pbGlobal.C2S_QueryHeroAtt{}, m.handleQueryHeroAtt)
	registerPBAccountHandler(&pbGlobal.C2S_HeroLevelup{}, m.handleHeroLevelup)
	registerPBAccountHandler(&pbGlobal.C2S_HeroPromote{}, m.handleHeroPromote)

	// fragment
	registerPBAccountHandler(&pbGlobal.C2S_QueryFragments{}, m.handleQueryFragments)
	registerPBAccountHandler(&pbGlobal.C2S_FragmentsCompose{}, m.handleFragmentsCompose)

	// items & equips & crystals
	registerPBAccountHandler(&pbGlobal.C2S_DelItem{}, m.handleDelItem)
	registerPBAccountHandler(&pbGlobal.C2S_UseItem{}, m.handleUseItem)
	registerPBAccountHandler(&pbGlobal.C2S_QueryItems{}, m.handleQueryItems)
	registerPBAccountHandler(&pbGlobal.C2S_EquipLevelup{}, m.handleEquipLevelup)

	registerPBAccountHandler(&pbGlobal.C2S_PutonEquip{}, m.handlePutonEquip)
	registerPBAccountHandler(&pbGlobal.C2S_TakeoffEquip{}, m.handleTakeoffEquip)
	registerPBAccountHandler(&pbGlobal.C2S_EquipPromote{}, m.handleEquipPromote)

	registerPBAccountHandler(&pbGlobal.C2S_PutonCrystal{}, m.handlePutonCrystal)
	registerPBAccountHandler(&pbGlobal.C2S_TakeoffCrystal{}, m.handleTakeoffCrystal)
	registerPBAccountHandler(&pbGlobal.C2S_CrystalLevelup{}, m.handleCrystalLevelup)
	registerPBAccountHandler(&pbGlobal.C2S_TestCrystalRandom{}, m.handleTestCrystalRandom)

	// tokens
	registerPBAccountHandler(&pbGlobal.C2S_QueryTokens{}, m.handleQueryTokens)

	// chapter & stage
	registerPBAccountHandler(&pbGlobal.C2S_StageSweep{}, m.handleStageSweep)

	// scene
	registerPBAccountHandler(&pbGlobal.C2S_StartStageCombat{}, m.handleStartStageCombat)
}
