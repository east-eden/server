package game

import (
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"github.com/yokaiio/yokai_server/transport"
)

type MsgHandler struct {
	g *Game
	r transport.Register
}

func NewMsgHandler(g *Game) *MsgHandler {
	m := &MsgHandler{
		g: g,
		r: transport.NewTransportRegister(),
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
	m.r.RegisterProtobufMessage(&pbAccount.C2M_AccountLogon{}, m.handleAccountLogon)
	m.r.RegisterProtobufMessage(&pbAccount.C2M_HeartBeat{}, m.handleHeartBeat)
	m.r.RegisterProtobufMessage(&pbAccount.MC_AccountConnected{}, m.handleAccountConnected)
	m.r.RegisterProtobufMessage(&pbAccount.C2M_AccountDisconnect{}, m.handleAccountDisconnect)

	// player
	m.r.RegisterProtobufMessage(&pbGame.C2M_QueryPlayerInfo{}, m.handleQueryPlayerInfo)
	m.r.RegisterProtobufMessage(&pbGame.C2M_CreatePlayer{}, m.handleCreatePlayer)
	m.r.RegisterProtobufMessage(&pbGame.MC_SelectPlayer{}, m.handleSelectPlayer)
	m.r.RegisterProtobufMessage(&pbGame.C2M_ChangeExp{}, m.handleChangeExp)
	m.r.RegisterProtobufMessage(&pbGame.C2M_ChangeLevel{}, m.handleChangeLevel)

	// heros
	m.r.RegisterProtobufMessage(&pbGame.C2M_AddHero{}, m.handleAddHero)
	m.r.RegisterProtobufMessage(&pbGame.C2M_DelHero{}, m.handleDelHero)
	m.r.RegisterProtobufMessage(&pbGame.C2M_QueryHeros{}, m.handleQueryHeros)
	//m.r.RegisterMessage("yokai_game.MC_HeroAddExp", &pbGame.MC_HeroAddExp{}, m.handleHeroAddExp)
	//m.r.RegisterMessage("yokai_game.MC_HeroAddLevel", &pbGame.MC_HeroAddLevel{}, m.handleHeroAddLevel)

	// items & equips
	m.r.RegisterProtobufMessage(&pbGame.C2M_AddItem{}, m.handleAddItem)
	m.r.RegisterProtobufMessage(&pbGame.C2M_DelItem{}, m.handleDelItem)
	m.r.RegisterProtobufMessage(&pbGame.C2M_UseItem{}, m.handleUseItem)
	m.r.RegisterProtobufMessage(&pbGame.C2M_QueryItems{}, m.handleQueryItems)

	m.r.RegisterProtobufMessage(&pbGame.C2M_PutonEquip{}, m.handlePutonEquip)
	m.r.RegisterProtobufMessage(&pbGame.C2M_TakeoffEquip{}, m.handleTakeoffEquip)

	// tokens
	m.r.RegisterProtobufMessage(&pbGame.C2M_AddToken{}, m.handleAddToken)
	m.r.RegisterProtobufMessage(&pbGame.C2M_QueryTokens{}, m.handleQueryTokens)

	// talent
	m.r.RegisterProtobufMessage(&pbGame.MC_AddTalent{}, m.handleAddTalent)
	m.r.RegisterProtobufMessage(&pbGame.MC_QueryTalents{}, m.handleQueryTalents)

	// rune
	m.r.RegisterProtobufMessage(&pbGame.C2M_AddRune{}, m.handleAddRune)
	m.r.RegisterProtobufMessage(&pbGame.C2M_DelRune{}, m.handleDelRune)
	m.r.RegisterProtobufMessage(&pbGame.C2M_QueryRunes{}, m.handleQueryRunes)
	m.r.RegisterProtobufMessage(&pbGame.C2M_PutonRune{}, m.handlePutonRune)
	m.r.RegisterProtobufMessage(&pbGame.C2M_TakeoffRune{}, m.handleTakeoffRune)

	// scene
	m.r.RegisterProtobufMessage(&pbGame.C2M_StartStageCombat{}, m.handleStartStageCombat)
}
