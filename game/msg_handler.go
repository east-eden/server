package game

import (
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/transport"
	pbClient "github.com/yokaiio/yokai_server/proto/client"
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
	m.r.RegisterMessage("yokai_client.MC_QueryPlayerInfo", &pbClient.MC_QueryPlayerInfo{}, m.handleQueryPlayerInfo)
	m.r.RegisterMessage("yokai_client.MC_CreatePlayer", &pbClient.MC_CreatePlayer{}, m.handleCreatePlayer)
	m.r.RegisterMessage("yokai_client.MC_SelectPlayer", &pbClient.MC_SelectPlayer{}, m.handleSelectPlayer)
	m.r.RegisterMessage("yokai_client.MC_ChangeExp", &pbClient.MC_ChangeExp{}, m.handleChangeExp)
	m.r.RegisterMessage("yokai_client.MC_ChangeLevel", &pbClient.MC_ChangeLevel{}, m.handleChangeLevel)

	// heros
	m.r.RegisterMessage("yokai_client.MC_AddHero", &pbClient.MC_AddHero{}, m.handleAddHero)
	m.r.RegisterMessage("yokai_client.MC_DelHero", &pbClient.MC_DelHero{}, m.handleDelHero)
	m.r.RegisterMessage("yokai_client.MC_QueryHeros", &pbClient.MC_QueryHeros{}, m.handleQueryHeros)

	// items
	m.r.RegisterMessage("yokai_client.MC_AddItem", &pbClient.MC_AddItem{}, m.handleAddItem)
	m.r.RegisterMessage("yokai_client.MC_DelItem", &pbClient.MC_DelItem{}, m.handleDelItem)
	m.r.RegisterMessage("yokai_client.MC_QueryItems", &pbClient.MC_QueryItems{}, m.handleQueryItems)

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

func (m *MsgHandler) handleQueryPlayerInfo(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("query player info failed")
		return
	}

	reply := &pbClient.MS_QueryPlayerInfo{}
	pl := cli.Player()
	if pl == nil {
		logger.WithFields(logger.Fields{
			"client_id": cli.ID(),
		}).Warn("client has no player")
		reply.Info = nil
	} else {
		reply.Info = &pbClient.PlayerInfo{
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

func (m *MsgHandler) handleCreatePlayer(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("create player failed")
		return
	}

	msg, ok := p.Body.(*pbClient.MC_CreatePlayer)
	if !ok {
		logger.Warn("create player failed, recv message body error")
		return
	}

	pl, err := m.g.cm.CreatePlayer(cli, msg.Name)
	reply := &pbClient.MS_CreatePlayer{
		ErrorCode: 0,
	}

	if err != nil {
		reply.ErrorCode = -1
	}

	if pl != nil {
		reply.Info = &pbClient.PlayerInfo{
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

	msg, ok := p.Body.(*pbClient.MC_SelectPlayer)
	if !ok {
		logger.Warn("Select player failed, recv message body error")
		return
	}

	pl, err := m.g.cm.SelectPlayer(cli, msg.Id)
	reply := &pbClient.MS_CreatePlayer{
		ErrorCode: 0,
	}

	if err != nil {
		reply.ErrorCode = -1
	}

	if pl != nil {
		reply.Info = &pbClient.PlayerInfo{
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

	msg, ok := p.Body.(*pbClient.MC_ChangeExp)
	if !ok {
		logger.Warn("change exp failed, recv message body error")
		return
	}

	cli.Player().ChangeExp(msg.AddExp)

	// sync player info
	pl := cli.Player()
	reply := &pbClient.MS_QueryPlayerInfo{
		Info: &pbClient.PlayerInfo{
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

	msg, ok := p.Body.(*pbClient.MC_ChangeLevel)
	if !ok {
		logger.Warn("change level failed, recv message body error")
		return
	}

	cli.Player().ChangeLevel(msg.AddLevel)

	// sync player info
	pl := cli.Player()
	reply := &pbClient.MS_QueryPlayerInfo{
		Info: &pbClient.PlayerInfo{
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

	msg, ok := p.Body.(*pbClient.MC_AddHero)
	if !ok {
		logger.Warn("Add Hero failed, recv message body error")
		return
	}

	cli.Player().HeroManager().AddHero(msg.TypeId)
	list := cli.Player().HeroManager().GetHeroList()
	reply := &pbClient.MS_HeroList{Heros: make([]*pbClient.Hero, len(list))}
	for _, v := range list {
		h := &pbClient.Hero{
			Id:     v.GetID(),
			TypeId: v.GetTypeID(),
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

	msg, ok := p.Body.(*pbClient.MC_DelHero)
	if !ok {
		logger.Warn("Delete Hero failed, recv message body error")
		return
	}

	cli.Player().HeroManager().DelHero(msg.Id)
	list := cli.Player().HeroManager().GetHeroList()
	reply := &pbClient.MS_HeroList{Heros: make([]*pbClient.Hero, len(list))}
	for _, v := range list {
		h := &pbClient.Hero{
			Id:     v.GetID(),
			TypeId: v.GetTypeID(),
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
	reply := &pbClient.MS_HeroList{Heros: make([]*pbClient.Hero, len(list))}
	for _, v := range list {
		h := &pbClient.Hero{
			Id:     v.GetID(),
			TypeId: v.GetTypeID(),
		}
		reply.Heros = append(reply.Heros, h)
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

	msg, ok := p.Body.(*pbClient.MC_AddItem)
	if !ok {
		logger.Warn("Add Item failed, recv message body error")
		return
	}

	cli.Player().ItemManager().AddItem(msg.TypeId)
	list := cli.Player().ItemManager().GetItemList()
	reply := &pbClient.MS_ItemList{Items: make([]*pbClient.Item, len(list))}
	for _, v := range list {
		i := &pbClient.Item{
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

	msg, ok := p.Body.(*pbClient.MC_DelItem)
	if !ok {
		logger.Warn("Delete Item failed, recv message body error")
		return
	}

	cli.Player().ItemManager().DelItem(msg.Id)
	list := cli.Player().ItemManager().GetItemList()
	reply := &pbClient.MS_ItemList{Items: make([]*pbClient.Item, len(list))}
	for _, v := range list {
		i := &pbClient.Item{
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
	reply := &pbClient.MS_ItemList{Items: make([]*pbClient.Item, len(list))}
	for _, v := range list {
		i := &pbClient.Item{
			Id:     v.GetID(),
			TypeId: v.GetTypeID(),
		}
		reply.Items = append(reply.Items, i)
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
