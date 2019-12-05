package game

import (
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/define"
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

	// items
	m.r.RegisterMessage("yokai_game.MC_AddItem", &pbGame.MC_AddItem{}, m.handleAddItem)
	m.r.RegisterMessage("yokai_game.MC_DelItem", &pbGame.MC_DelItem{}, m.handleDelItem)
	m.r.RegisterMessage("yokai_game.MC_QueryItems", &pbGame.MC_QueryItems{}, m.handleQueryItems)

	// tokens
	m.r.RegisterMessage("yokai_game.MC_AddToken", &pbGame.MC_AddToken{}, m.handleAddToken)
	m.r.RegisterMessage("yokai_game.MC_QueryTokens", &pbGame.MC_QueryTokens{}, m.handleQueryTokens)

	// json
	m.r.RegisterMessage("MC_ClientTest", &MC_ClientTest{}, m.handleClientTest)

	/* m.regProtoHandle("ultimate_service_game.MWU_RequestPlayerInfo", m.handleRequestPlayerInfo)*/
	//m.regProtoHandle("ultimate_service_game.MWU_RequestGuildInfo", m.handleRequestGuildInfo)
	//m.regProtoHandle("ultimate_service_game.MWU_PlayUltimateRecord", m.handlePlayUltimateRecord)
	//m.regProtoHandle("ultimate_service_game.MWU_RequestUltimatePlayer", m.handleRequestUltimatePlayer)
	//m.regProtoHandle("ultimate_service_game.MWU_ViewFormation", m.handleViewFormation)
	//m.regProtoHandle("ultimate_service_game.MWU_AddInvite", m.handleAddInvite)
	//m.regProtoHandle("ultimate_service_game.MWU_CheckInviteResult", m.handleCheckInviteResult)
	//m.regProtoHandle("ultimate_service_game.MWU_InviteRecharge", m.handleInviteRecharge)
	//m.regProtoHandle("ultimate_service_game.MWU_ReplacePlayerInfo", m.handleReplacePlayerInfo)
	//m.regProtoHandle("ultimate_service_game.MWU_ReplaceGuildInfo", m.handleReplaceGuildInfo)
	//m.regProtoHandle("ultimate_service_arena.MWU_ArenaMatching", m.handleArenaMatching)
	//m.regProtoHandle("ultimate_service_arena.MWU_ArenaAddRecord", m.handleArenaAddRecord)
	//m.regProtoHandle("ultimate_service_arena.MWU_ArenaBattleResult", m.handleArenaBattleResult)
	//m.regProtoHandle("ultimate_service_arena.MWU_RequestArenaRank", m.handleRequestArenaRank)
	/*m.regProtoHandle("ultimate_service_arena.MWU_ArenaChampionOnline", m.handleArenaChampionOnline)*/

}

func (m *MsgHandler) handleClientLogon(sock transport.Socket, p *transport.Message) {
	msg, ok := p.Body.(*pbClient.MC_ClientLogon)
	if !ok {
		logger.Warn("Cannot assert value to message")
		return
	}

	client := m.g.cm.GetClientBySock(sock)
	if client != nil {
		logger.Warn("client had logon:", sock)
		return
	}

	client, err := m.g.cm.ClientLogon(msg.ClientId, msg.ClientName, sock)
	if err != nil {
		logger.WithFields(logger.Fields{
			"id":   msg.ClientId,
			"name": msg.ClientName,
			"sock": sock,
		}).Warn("add client failed")
		return
	}

	reply := &pbClient.MS_ClientLogon{}
	client.SendProtoMessage(reply)
}

func (m *MsgHandler) handleHeartBeat(sock transport.Socket, p *transport.Message) {
	if client := m.g.cm.GetClientBySock(sock); client != nil {
		if t := int32(time.Now().Unix()); t == -1 {
			logger.Warn("Heart beat get time err")
			return
		}

		client.HeartBeat()
	}
}

func (m *MsgHandler) handleClientConnected(sock transport.Socket, p *transport.Message) {
	if client := m.g.cm.GetClientBySock(sock); client != nil {
		clientID := p.Body.(*pbClient.MC_ClientConnected).ClientId
		logger.WithFields(logger.Fields{
			"client_id": clientID,
		}).Info("client connected")

		// todo after connected
	}
}

func (m *MsgHandler) handleClientDisconnect(sock transport.Socket, p *transport.Message) {
	m.g.cm.DisconnectClient(sock, "client disconnect initiativly")
}

func (m *MsgHandler) handleQueryPlayerInfos(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("query player info failed")
		return
	}

	playerList := m.g.pm.GetPlayersByClientID(cli.ID())
	reply := &pbGame.MS_QueryPlayerInfos{
		Infos: make([]*pbGame.PlayerInfo, 0),
	}

	for _, v := range playerList {
		info := &pbGame.PlayerInfo{
			Id:       v.GetID(),
			Name:     v.GetName(),
			Exp:      v.GetExp(),
			Level:    v.GetLevel(),
			HeroNums: int32(v.HeroManager().GetHeroNums()),
			ItemNums: int32(v.ItemManager().GetItemNums()),
		}

		reply.Infos = append(reply.Infos, info)
	}

	cli.SendProtoMessage(reply)
}

func (m *MsgHandler) handleCreatePlayer(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("create player failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_CreatePlayer)
	if !ok {
		logger.Warn("create player failed, recv message body error")
		return
	}

	pl, err := m.g.cm.CreatePlayer(cli, msg.Name)
	reply := &pbGame.MS_CreatePlayer{
		ErrorCode: 0,
	}

	if err != nil {
		reply.ErrorCode = -1
	}

	if pl != nil {
		reply.Info = &pbGame.PlayerInfo{
			Id:       pl.GetID(),
			Name:     pl.GetName(),
			Exp:      pl.GetExp(),
			Level:    pl.GetLevel(),
			HeroNums: int32(pl.HeroManager().GetHeroNums()),
			ItemNums: int32(pl.ItemManager().GetItemNums()),
		}
	}

	cli.SendProtoMessage(reply)
}

func (m *MsgHandler) handleSelectPlayer(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("select player failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_SelectPlayer)
	if !ok {
		logger.Warn("Select player failed, recv message body error")
		return
	}

	pl, err := m.g.cm.SelectPlayer(cli, msg.Id)
	reply := &pbGame.MS_SelectPlayer{
		ErrorCode: 0,
	}

	if err != nil {
		reply.ErrorCode = -1
	}

	if pl != nil {
		reply.Info = &pbGame.PlayerInfo{
			Id:       pl.GetID(),
			Name:     pl.GetName(),
			Exp:      pl.GetExp(),
			Level:    pl.GetLevel(),
			HeroNums: int32(pl.HeroManager().GetHeroNums()),
			ItemNums: int32(pl.ItemManager().GetItemNums()),
		}
	}

	cli.SendProtoMessage(reply)
}

func (m *MsgHandler) handleChangeExp(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("change exp failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_ChangeExp)
	if !ok {
		logger.Warn("change exp failed, recv message body error")
		return
	}

	cli.Player().ChangeExp(msg.AddExp)

	// sync player info
	pl := cli.Player()
	reply := &pbGame.MS_QueryPlayerInfo{
		Info: &pbGame.PlayerInfo{
			Id:       pl.GetID(),
			Name:     pl.GetName(),
			Exp:      pl.GetExp(),
			Level:    pl.GetLevel(),
			HeroNums: int32(pl.HeroManager().GetHeroNums()),
			ItemNums: int32(pl.ItemManager().GetItemNums()),
		},
	}

	cli.SendProtoMessage(reply)
}

func (m *MsgHandler) handleChangeLevel(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("change level failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_ChangeLevel)
	if !ok {
		logger.Warn("change level failed, recv message body error")
		return
	}

	cli.Player().ChangeLevel(msg.AddLevel)

	// sync player info
	pl := cli.Player()
	reply := &pbGame.MS_QueryPlayerInfo{
		Info: &pbGame.PlayerInfo{
			Id:       pl.GetID(),
			Name:     pl.GetName(),
			Exp:      pl.GetExp(),
			Level:    pl.GetLevel(),
			HeroNums: int32(pl.HeroManager().GetHeroNums()),
			ItemNums: int32(pl.ItemManager().GetItemNums()),
		},
	}

	cli.SendProtoMessage(reply)
}

func (m *MsgHandler) handleAddHero(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("add hero failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_AddHero)
	if !ok {
		logger.Warn("Add Hero failed, recv message body error")
		return
	}

	cli.Player().HeroManager().AddHero(msg.TypeId)
	list := cli.Player().HeroManager().GetHeroList()
	reply := &pbGame.MS_HeroList{Heros: make([]*pbGame.Hero, 0)}
	for _, v := range list {
		h := &pbGame.Hero{
			Id:     v.GetID(),
			TypeId: v.GetTypeID(),
			Exp:    v.GetExp(),
			Level:  v.GetLevel(),
		}
		reply.Heros = append(reply.Heros, h)
	}
	cli.SendProtoMessage(reply)
}

func (m *MsgHandler) handleDelHero(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("delete hero failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_DelHero)
	if !ok {
		logger.Warn("Delete Hero failed, recv message body error")
		return
	}

	cli.Player().HeroManager().DelHero(msg.Id)
	list := cli.Player().HeroManager().GetHeroList()
	reply := &pbGame.MS_HeroList{Heros: make([]*pbGame.Hero, 0)}
	for _, v := range list {
		h := &pbGame.Hero{
			Id:     v.GetID(),
			TypeId: v.GetTypeID(),
			Exp:    v.GetExp(),
			Level:  v.GetLevel(),
		}
		reply.Heros = append(reply.Heros, h)
	}
	cli.SendProtoMessage(reply)
}

func (m *MsgHandler) handleQueryHeros(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("query heros failed")
		return
	}

	list := cli.Player().HeroManager().GetHeroList()
	reply := &pbGame.MS_HeroList{Heros: make([]*pbGame.Hero, 0)}
	for _, v := range list {
		h := &pbGame.Hero{
			Id:     v.GetID(),
			TypeId: v.GetTypeID(),
			Exp:    v.GetExp(),
			Level:  v.GetLevel(),
		}
		reply.Heros = append(reply.Heros, h)
	}
	cli.SendProtoMessage(reply)
}

func (m *MsgHandler) handleHeroAddExp(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("hero add exp failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_HeroAddExp)
	if !ok {
		logger.Warn("hero add exp failed, recv message body error")
		return
	}

	if cli.Player() == nil {
		logger.Warn("client has no player", cli.ID())
		return
	}

	cli.Player().HeroManager().HeroAddExp(msg.HeroId, msg.Exp)
	hero := cli.Player().HeroManager().GetHero(msg.HeroId)
	if hero == nil {
		logger.Warn("get hero by id error:", msg.HeroId)
		return
	}

	reply := &pbGame.MS_HeroInfo{
		Info: &pbGame.Hero{
			Id:     hero.GetID(),
			TypeId: hero.GetTypeID(),
			Exp:    hero.GetExp(),
			Level:  hero.GetLevel(),
		},
	}

	cli.SendProtoMessage(reply)
}

func (m *MsgHandler) handleHeroAddLevel(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("hero add level failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_HeroAddLevel)
	if !ok {
		logger.Warn("hero add level failed, recv message body error")
		return
	}

	if cli.Player() == nil {
		logger.Warn("client has no player", cli.ID())
		return
	}

	cli.Player().HeroManager().HeroAddLevel(msg.HeroId, msg.Level)
	hero := cli.Player().HeroManager().GetHero(msg.HeroId)
	if hero == nil {
		logger.Warn("get hero by id error:", msg.HeroId)
		return
	}

	reply := &pbGame.MS_HeroInfo{
		Info: &pbGame.Hero{
			Id:     hero.GetID(),
			TypeId: hero.GetTypeID(),
			Exp:    hero.GetExp(),
			Level:  hero.GetLevel(),
		},
	}

	cli.SendProtoMessage(reply)
}

func (m *MsgHandler) handleAddItem(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("add item failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_AddItem)
	if !ok {
		logger.Warn("Add Item failed, recv message body error")
		return
	}

	cli.Player().ItemManager().AddItem(msg.TypeId)
	list := cli.Player().ItemManager().GetItemList()
	reply := &pbGame.MS_ItemList{Items: make([]*pbGame.Item, 0)}
	for _, v := range list {
		i := &pbGame.Item{
			Id:     v.GetID(),
			TypeId: v.GetTypeID(),
		}
		reply.Items = append(reply.Items, i)
	}
	cli.SendProtoMessage(reply)
}

func (m *MsgHandler) handleDelItem(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("delete item failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_DelItem)
	if !ok {
		logger.Warn("Delete Item failed, recv message body error")
		return
	}

	cli.Player().ItemManager().DelItem(msg.Id)
	list := cli.Player().ItemManager().GetItemList()
	reply := &pbGame.MS_ItemList{Items: make([]*pbGame.Item, 0)}
	for _, v := range list {
		i := &pbGame.Item{
			Id:     v.GetID(),
			TypeId: v.GetTypeID(),
		}
		reply.Items = append(reply.Items, i)
	}
	cli.SendProtoMessage(reply)
}

func (m *MsgHandler) handleQueryItems(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("query items failed")
		return
	}

	list := cli.Player().ItemManager().GetItemList()
	reply := &pbGame.MS_ItemList{Items: make([]*pbGame.Item, 0)}
	for _, v := range list {
		i := &pbGame.Item{
			Id:     v.GetID(),
			TypeId: v.GetTypeID(),
		}
		reply.Items = append(reply.Items, i)
	}
	cli.SendProtoMessage(reply)
}

func (m *MsgHandler) handleAddToken(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("add token failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_AddToken)
	if !ok {
		logger.Warn("Add Item failed, recv message body error")
		return
	}

	err := cli.Player().TokenManager().TokenInc(msg.Type, msg.Value)
	if err != nil {
		logger.Warn("token inc failed:", err)
	}

	cli.Player().TokenManager().Save()

	reply := &pbGame.MS_TokenList{Tokens: make([]*pbGame.Token, 0)}
	for n := 0; n < define.Token_End; n++ {
		v, err := cli.Player().TokenManager().GetToken(int32(n))
		if err != nil {
			logger.Warn("token get value failed:", err)
			return
		}

		t := &pbGame.Token{
			Type:    v.ID,
			Value:   v.Value,
			MaxHold: v.MaxHold,
		}
		reply.Tokens = append(reply.Tokens, t)
	}
	cli.SendProtoMessage(reply)
}

func (m *MsgHandler) handleQueryTokens(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("query tokens failed")
		return
	}

	reply := &pbGame.MS_TokenList{Tokens: make([]*pbGame.Token, 0)}
	for n := 0; n < define.Token_End; n++ {
		v, err := cli.Player().TokenManager().GetToken(int32(n))
		if err != nil {
			logger.Warn("token get value failed:", err)
			return
		}

		t := &pbGame.Token{
			Type:    v.ID,
			Value:   v.Value,
			MaxHold: v.MaxHold,
		}
		reply.Tokens = append(reply.Tokens, t)
	}
	cli.SendProtoMessage(reply)
}

func (m *MsgHandler) handleClientTest(sock transport.Socket, p *transport.Message) {
	logger.WithFields(logger.Fields{
		"type": p.Type,
		"name": p.Name,
		"body": p.Body,
	}).Info("recv client test")
}

/*func (m *MsgHandler) handleRequestPlayerInfo(con iface.ITCPConn, p proto.Message) {*/
//if world := m.wm.GetWorldByCon(con); world != nil {
//msg, ok := p.(*pbGame.MWU_RequestPlayerInfo)
//if !ok {
//logger.WithFields(logger.Fields{
//"msg_name": proto.MessageName(p),
//}).Warn("parsing message name error")
//return
//}

//m.gm.AddPlayerInfoList(msg.Info)
//}
//}

//func (m *MsgHandler) handleRequestGuildInfo(con iface.ITCPConn, p proto.Message) {
//if world := m.wm.GetWorldByCon(con); world != nil {
//msg, ok := p.(*pbGame.MWU_RequestGuildInfo)
//if !ok {
//logger.WithFields(logger.Fields{
//"msg_name": proto.MessageName(p),
//}).Warn("parsing message name error")
//return
//}

//m.gm.AddGuildInfoList(msg.Info)
//}
//}

//func (m *MsgHandler) handlePlayUltimateRecord(con iface.ITCPConn, p proto.Message) {
//if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
//msg, ok := p.(*pbGame.MWU_PlayUltimateRecord)
//if !ok {
//logger.WithFields(logger.Fields{
//"msg_name": proto.MessageName(p),
//}).Warn("parsing message name error")
//return
//}

//dstWorld := m.wm.GetWorldByID(msg.DstServerId)
//if dstWorld == nil {
//return
//}

//msgSend := &pbGame.MUW_PlayUltimateRecord{
//SrcPlayerId: msg.SrcPlayerId,
//SrcServerId: msg.SrcServerId,
//RecordId:    msg.RecordId,
//DstServerId: msg.DstServerId,
//}
//dstWorld.SendProtoMessage(msgSend)
//}
//}

//func (m *MsgHandler) handleRequestUltimatePlayer(con iface.ITCPConn, p proto.Message) {
//if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
//msg, ok := p.(*pbGame.MWU_RequestUltimatePlayer)
//if !ok {
//logger.WithFields(logger.Fields{
//"msg_name": proto.MessageName(p),
//}).Warn("parsing message name error")
//return
//}

//dstInfo, err := m.gm.GetPlayerInfoByID(msg.DstPlayerId)
//dstWorld := m.wm.GetWorldByID(msg.DstServerId)
//if err != nil {
//return
//}

//if int32(msg.DstServerId) == -1 {
//dstWorld = m.wm.GetWorldByID(dstInfo.ServerId)
//}

//if dstWorld == nil {
//return
//}

//msgSend := &pbGame.MUW_RequestUltimatePlayer{
//SrcPlayerId: msg.SrcPlayerId,
//SrcServerId: msg.SrcServerId,
//DstPlayerId: msg.DstPlayerId,
//DstServerId: dstWorld.GetID(),
//}
//dstWorld.SendProtoMessage(msgSend)
//}
//}

//func (m *MsgHandler) handleViewFormation(con iface.ITCPConn, p proto.Message) {
//if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
//msg, ok := p.(*pbGame.MWU_ViewFormation)
//if !ok {
//logger.WithFields(logger.Fields{
//"msg_name": proto.MessageName(p),
//}).Warn("parsing message name error")
//return
//}

//dstInfo, err := m.gm.GetPlayerInfoByID(msg.DstPlayerId)
//dstWorld := m.wm.GetWorldByID(msg.DstServerId)
//if err != nil {
//return
//}

//if int32(msg.DstServerId) == -1 {
//dstWorld = m.wm.GetWorldByID(dstInfo.ServerId)
//}

//if dstWorld == nil {
//return
//}

//msgSend := &pbGame.MUW_ViewFormation{
//SrcPlayerId: msg.SrcPlayerId,
//SrcServerId: msg.SrcServerId,
//DstPlayerId: msg.DstPlayerId,
//DstServerId: dstWorld.GetID(),
//}
//dstWorld.SendProtoMessage(msgSend)
//}
//}

/////////////////////////////////
//// arena battle
////////////////////////////////
//func (m *MsgHandler) handleArenaMatching(con iface.ITCPConn, p proto.Message) {
//if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
//msg, ok := p.(*pbArena.MWU_ArenaMatching)
//if !ok {
//logger.WithFields(logger.Fields{
//"msg_name": proto.MessageName(p),
//}).Warn("parsing message name error")
//return
//}

//m.gm.ArenaMatching(msg.PlayerId)
//}
//}

//func (m *MsgHandler) handleArenaAddRecord(con iface.ITCPConn, p proto.Message) {
//if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
//msg, ok := p.(*pbArena.MWU_ArenaAddRecord)
//if !ok {
//logger.WithFields(logger.Fields{
//"msg_name": proto.MessageName(p),
//}).Warn("parsing message name error")
//return
//}

//m.gm.ArenaAddRecord(msg.Record)
//}
//}

//func (m *MsgHandler) handleArenaBattleResult(con iface.ITCPConn, p proto.Message) {
//if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
//msg, ok := p.(*pbArena.MWU_ArenaBattleResult)
//if !ok {
//logger.WithFields(logger.Fields{
//"msg_name": proto.MessageName(p),
//}).Warn("parsing message name error")
//return
//}

//m.gm.ArenaBattleResult(msg.AttackId, msg.TargetId, msg.AttackWin)
//}
//}

//func (m *MsgHandler) handleReplacePlayerInfo(con iface.ITCPConn, p proto.Message) {
//if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
//msg, ok := p.(*pbGame.MWU_ReplacePlayerInfo)
//if !ok {
//logger.WithFields(logger.Fields{
//"msg_name": proto.MessageName(p),
//}).Warn("parsing message name error")
//return
//}

//m.gm.AddPlayerInfo(msg.Info)
//}
//}

//func (m *MsgHandler) handleReplaceGuildInfo(con iface.ITCPConn, p proto.Message) {
//if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
//msg, ok := p.(*pbGame.MWU_ReplaceGuildInfo)
//if !ok {
//logger.WithFields(logger.Fields{
//"msg_name": proto.MessageName(p),
//}).Warn("parsing message name error")
//return
//}

//m.gm.AddGuildInfo(msg.Info)
//}
//}

//func (m *MsgHandler) handleRequestArenaRank(con iface.ITCPConn, p proto.Message) {
//if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
//msg, ok := p.(*pbArena.MWU_RequestArenaRank)
//if !ok {
//logger.WithFields(logger.Fields{
//"msg_name": proto.MessageName(p),
//}).Warn("parsing message name error")
//return
//}

//m.gm.ArenaGetRank(msg.PlayerId, msg.Page)
//}
//}

//func (m *MsgHandler) handleAddInvite(con iface.ITCPConn, p proto.Message) {
//if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
//msg, ok := p.(*pbGame.MWU_AddInvite)
//if !ok {
//logger.WithFields(logger.Fields{
//"msg_name": proto.MessageName(p),
//}).Warn("parsing message name error")
//return
//}

//ret := m.gm.Invite().AddInvite(msg.NewbieId, msg.InviterId)
//if ret != 0 {
//msgRet := &pbGame.MUW_AddInviteResult{
//NewbieId:  msg.NewbieId,
//InviterId: msg.InviterId,
//ErrorCode: ret,
//}

//srcWorld.SendProtoMessage(msgRet)
//}
//}
//}

//func (m *MsgHandler) handleCheckInviteResult(con iface.ITCPConn, p proto.Message) {
//if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
//msg, ok := p.(*pbGame.MWU_CheckInviteResult)
//if !ok {
//logger.WithFields(logger.Fields{
//"msg_name": proto.MessageName(p),
//}).Warn("parsing message name error")
//return
//}

//m.gm.Invite().CheckInviteResult(msg.NewbieId, msg.InviterId, msg.ErrorCode)
//}
//}

//func (m *MsgHandler) handleInviteRecharge(con iface.ITCPConn, p proto.Message) {
//if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
//msg, ok := p.(*pbGame.MWU_InviteRecharge)
//if !ok {
//logger.WithFields(logger.Fields{
//"msg_name": proto.MessageName(p),
//}).Warn("parsing message name error")
//return
//}

//m.gm.Invite().InviteRecharge(msg.NewbieId, msg.NewbieName, msg.InviterId, msg.DiamondGift)
//}
//}

//func (m *MsgHandler) handleArenaChampionOnline(con iface.ITCPConn, p proto.Message) {
//if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
//msg, ok := p.(*pbArena.MWU_ArenaChampionOnline)
//if !ok {
//logger.WithFields(logger.Fields{
//"msg_name": proto.MessageName(p),
//}).Warn("parsing message name error")
//return
//}

//msgSend := &pbArena.MUW_ArenaChampionOnline{
//PlayerId:   msg.PlayerId,
//PlayerName: msg.PlayerName,
//ServerName: msg.ServerName,
//}

//m.wm.BroadCast(msgSend)
//}
/*}*/
