package game

import (
	"github.com/yokaiio/yokai_server/internal/transport"
	pbClient "github.com/yokaiio/yokai_server/proto/client"
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

type MC_ClientTest struct {
	ClientId int64  `protobuf:"varint,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	Name     string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (m *MsgHandler) registerAllMessage() {

	// json
	m.r.RegisterMessage("MC_ClientTest", &MC_ClientTest{}, m.handleClientTest)

	// client
	m.r.RegisterMessage("yokai_client.MC_ClientLogon", &pbClient.MC_ClientLogon{}, m.handleClientLogon)
	m.r.RegisterMessage("yokai_client.MC_HeartBeat", &pbClient.MC_HeartBeat{}, m.handleHeartBeat)
	m.r.RegisterMessage("yokai_client.MC_ClientConnected", &pbClient.MC_ClientConnected{}, m.handleClientConnected)
	m.r.RegisterMessage("yokai_client.MC_ClientDisconnect", &pbClient.MC_ClientDisconnect{}, m.handleClientDisconnect)

	// player
	m.r.RegisterMessage("yokai_game.MC_QueryPlayerInfos", &pbGame.MC_QueryPlayerInfos{}, m.handleQueryPlayerInfos)
	m.r.RegisterMessage("yokai_game.MC_CreatePlayer", &pbGame.MC_CreatePlayer{}, m.handleCreatePlayer)
	m.r.RegisterMessage("yokai_game.MC_SelectPlayer", &pbGame.MC_SelectPlayer{}, m.handleSelectPlayer)
	m.r.RegisterMessage("yokai_game.MC_ChangeExp", &pbGame.MC_ChangeExp{}, m.handleChangeExp)
	m.r.RegisterMessage("yokai_game.MC_ChangeLevel", &pbGame.MC_ChangeLevel{}, m.handleChangeLevel)

	// heros
	m.r.RegisterMessage("yokai_game.MC_AddHero", &pbGame.MC_AddHero{}, m.handleAddHero)
	m.r.RegisterMessage("yokai_game.MC_DelHero", &pbGame.MC_DelHero{}, m.handleDelHero)
	m.r.RegisterMessage("yokai_game.MC_QueryHeros", &pbGame.MC_QueryHeros{}, m.handleQueryHeros)
	m.r.RegisterMessage("yokai_game.MC_HeroAddExp", &pbGame.MC_HeroAddExp{}, m.handleHeroAddExp)
	m.r.RegisterMessage("yokai_game.MC_HeroAddLevel", &pbGame.MC_HeroAddLevel{}, m.handleHeroAddLevel)

	// items & equips
	m.r.RegisterMessage("yokai_game.MC_AddItem", &pbGame.MC_AddItem{}, m.handleAddItem)
	m.r.RegisterMessage("yokai_game.MC_DelItem", &pbGame.MC_DelItem{}, m.handleDelItem)
	m.r.RegisterMessage("yokai_game.MC_QueryItems", &pbGame.MC_QueryItems{}, m.handleQueryItems)

	m.r.RegisterMessage("yokai_game.MC_PutonEquip", &pbGame.MC_PutonEquip{}, m.handlePutonEquip)
	m.r.RegisterMessage("yokai_game.MC_TakeoffEquip", &pbGame.MC_TakeoffEquip{}, m.handleTakeoffEquip)

	// tokens
	m.r.RegisterMessage("yokai_game.MC_AddToken", &pbGame.MC_AddToken{}, m.handleAddToken)
	m.r.RegisterMessage("yokai_game.MC_QueryTokens", &pbGame.MC_QueryTokens{}, m.handleQueryTokens)

	// talent
	m.r.RegisterMessage("yokai_game.MC_AddTalent", &pbGame.MC_AddTalent{}, m.handleAddTalent)
	m.r.RegisterMessage("yokai_game.MC_QueryTalents", &pbGame.MC_QueryTalents{}, m.handleQueryTalents)
}
