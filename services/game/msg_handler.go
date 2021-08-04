package game

import (
	"context"
	"errors"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/player"
	"github.com/east-eden/server/transport"
	"github.com/east-eden/server/transport/codec"
	"github.com/hellodudu/task"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

var (
	ErrPlayerNotFound = errors.New("player not found")
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
				Str("name", string(p.ProtoReflect().Descriptor().Name())).
				Msg("register message failed")
		}
	}

	// account protobuf handler
	registerPBAccountHandler := func(p proto.Message, handle task.TaskHandler) {
		mf := func(ctx context.Context, sock transport.Socket, msg *transport.Message) error {

			// wrap heartbeat
			wrappedHandle := func(ctx context.Context, p ...interface{}) error {
				acct := p[0].(*player.Account)
				acct.HeartBeat()
				return handle(ctx, p...)
			}

			accountId, ok := m.am.GetAccountIdBySock(sock)
			if !ok {
				return ErrAccountNotFound
			}

			return m.am.AddAccountTask(
				ctx,
				accountId,
				wrappedHandle,
				msg.Body.(proto.Message),
			)
		}

		err := m.r.RegisterProtobufMessage(p, mf)
		if err != nil {
			log.Fatal().
				Err(err).
				Str("name", string(p.ProtoReflect().Descriptor().Name())).
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
	registerPBAccountHandler(&pbGlobal.C2S_GuidePass{}, m.handleGuidePass)
	registerPBAccountHandler(&pbGlobal.C2S_SaveBattleArray{}, m.handleSaveBattleArray)

	// heros
	registerPBAccountHandler(&pbGlobal.C2S_DelHero{}, m.handleDelHero)
	registerPBAccountHandler(&pbGlobal.C2S_HeroLevelup{}, m.handleHeroLevelup)
	registerPBAccountHandler(&pbGlobal.C2S_HeroPromote{}, m.handleHeroPromote)
	registerPBAccountHandler(&pbGlobal.C2S_HeroStarup{}, m.handleHeroStarup)
	registerPBAccountHandler(&pbGlobal.C2S_HeroTalentChoose{}, m.handleHeroTalentChoose)

	// fragment
	registerPBAccountHandler(&pbGlobal.C2S_HeroFragmentsCompose{}, m.handleHeroFragmentsCompose)
	registerPBAccountHandler(&pbGlobal.C2S_CollectionFragmentsCompose{}, m.handleCollectionFragmentsCompose)

	// items & equips & crystals
	registerPBAccountHandler(&pbGlobal.C2S_DelItem{}, m.handleDelItem)
	registerPBAccountHandler(&pbGlobal.C2S_UseItem{}, m.handleUseItem)

	registerPBAccountHandler(&pbGlobal.C2S_EquipLevelup{}, m.handleEquipLevelup)
	registerPBAccountHandler(&pbGlobal.C2S_PutonEquip{}, m.handlePutonEquip)
	registerPBAccountHandler(&pbGlobal.C2S_TakeoffEquip{}, m.handleTakeoffEquip)
	registerPBAccountHandler(&pbGlobal.C2S_EquipPromote{}, m.handleEquipPromote)
	registerPBAccountHandler(&pbGlobal.C2S_EquipStarup{}, m.handleEquipStarup)

	registerPBAccountHandler(&pbGlobal.C2S_PutonCrystal{}, m.handlePutonCrystal)
	registerPBAccountHandler(&pbGlobal.C2S_TakeoffCrystal{}, m.handleTakeoffCrystal)
	registerPBAccountHandler(&pbGlobal.C2S_CrystalLevelup{}, m.handleCrystalLevelup)
	registerPBAccountHandler(&pbGlobal.C2S_TestCrystalRandom{}, m.handleTestCrystalRandom)

	// collections
	registerPBAccountHandler(&pbGlobal.C2S_CollectionActive{}, m.handleCollectionActive)
	registerPBAccountHandler(&pbGlobal.C2S_CollectionStarup{}, m.handleCollectionStarup)
	registerPBAccountHandler(&pbGlobal.C2S_CollectionWakeup{}, m.handleCollectionWakeup)

	// tokens

	// chapter & stage
	registerPBAccountHandler(&pbGlobal.C2S_StageSweep{}, m.handleStageSweep)
	registerPBAccountHandler(&pbGlobal.C2S_StageChallenge{}, m.handleStageChallenge)
	registerPBAccountHandler(&pbGlobal.C2S_ChapterReward{}, m.handleChapterReward)

	// tower
	registerPBAccountHandler(&pbGlobal.C2S_TowerChallenge{}, m.handleTowerChallenge)

	// scene

	// quest
	registerPBAccountHandler(&pbGlobal.C2S_PlayerQuestReward{}, m.handlePlayerQuestReward)

	// rank
	registerPBAccountHandler(&pbGlobal.C2S_QueryRank{}, m.handleQueryRank)
}
