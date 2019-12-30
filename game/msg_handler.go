package game

import (
	"github.com/yokaiio/yokai_server/internal/transport"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

type MsgHandler struct {
	g *Game
	r transport.Register
}

func NewMsgHandler(g *Game) *MsgHandler {
	m := &MsgHandler{
		g: g,
		r: transport.DefaultRegister,
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

	// json
	m.r.RegisterJsonMessage(&MC_AccountTest{}, m.handleAccountTest)

	// account
	m.r.RegisterProtobufMessage(&pbAccount.MC_AccountLogon{}, m.handleAccountLogon)
	m.r.RegisterProtobufMessage(&pbAccount.MC_AccountLogon{}, m.handleAccountLogon)
	m.r.RegisterProtobufMessage(&pbAccount.MC_HeartBeat{}, m.handleHeartBeat)
	m.r.RegisterProtobufMessage(&pbAccount.MC_AccountConnected{}, m.handleAccountConnected)
	m.r.RegisterProtobufMessage(&pbAccount.MC_AccountDisconnect{}, m.handleAccountDisconnect)

	// player
	m.r.RegisterProtobufMessage(&pbGame.MC_QueryPlayerInfos{}, m.handleQueryPlayerInfos)
	m.r.RegisterProtobufMessage(&pbGame.MC_CreatePlayer{}, m.handleCreatePlayer)
	m.r.RegisterProtobufMessage(&pbGame.MC_SelectPlayer{}, m.handleSelectPlayer)
	m.r.RegisterProtobufMessage(&pbGame.MC_ChangeExp{}, m.handleChangeExp)
	m.r.RegisterProtobufMessage(&pbGame.MC_ChangeLevel{}, m.handleChangeLevel)

	// heros
	m.r.RegisterProtobufMessage(&pbGame.MC_AddHero{}, m.handleAddHero)
	m.r.RegisterProtobufMessage(&pbGame.MC_DelHero{}, m.handleDelHero)
	m.r.RegisterProtobufMessage(&pbGame.MC_QueryHeros{}, m.handleQueryHeros)
	//m.r.RegisterMessage("yokai_game.MC_HeroAddExp", &pbGame.MC_HeroAddExp{}, m.handleHeroAddExp)
	//m.r.RegisterMessage("yokai_game.MC_HeroAddLevel", &pbGame.MC_HeroAddLevel{}, m.handleHeroAddLevel)

	// items & equips
	m.r.RegisterProtobufMessage(&pbGame.MC_AddItem{}, m.handleAddItem)
	m.r.RegisterProtobufMessage(&pbGame.MC_DelItem{}, m.handleDelItem)
	m.r.RegisterProtobufMessage(&pbGame.MC_QueryItems{}, m.handleQueryItems)
	m.r.RegisterProtobufMessage(&pbGame.MC_QueryHeroEquips{}, m.handleQueryHeroEquips)

	m.r.RegisterProtobufMessage(&pbGame.MC_PutonEquip{}, m.handlePutonEquip)
	m.r.RegisterProtobufMessage(&pbGame.MC_TakeoffEquip{}, m.handleTakeoffEquip)

	// tokens
	m.r.RegisterProtobufMessage(&pbGame.MC_AddToken{}, m.handleAddToken)
	m.r.RegisterProtobufMessage(&pbGame.MC_QueryTokens{}, m.handleQueryTokens)

	// talent
	m.r.RegisterProtobufMessage(&pbGame.MC_AddTalent{}, m.handleAddTalent)
	m.r.RegisterProtobufMessage(&pbGame.MC_QueryTalents{}, m.handleQueryTalents)
}
