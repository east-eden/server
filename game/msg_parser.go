package game

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/transport"
	pbClient "github.com/yokaiio/yokai_server/proto/client"
)

// ProtoHandler handle function
type ProtoHandler func(transport.Socket, *transport.Message)

type MsgParser struct {
	protoHandler map[uint32]ProtoHandler
	g            *Game
}

func NewMsgParser(g *Game) *MsgParser {
	m := &MsgParser{
		protoHandler: make(map[uint32]ProtoHandler),
		g:            g,
	}

	m.registerAllMessage()
	return m
}

func (m *MsgParser) registerAllMessage() {
	m.regProtoHandle("yokai_client.MC_ClientLogon", m.handleClientLogon)
	m.regProtoHandle("yokai_client.MC_HeartBeat", m.handleHeartBeat)
	m.regProtoHandle("yokai_client.MC_ClientConnected", m.handleClientConnected)

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

func (m *MsgParser) getRegProtoHandle(id uint32) (ProtoHandler, error) {
	v, ok := m.protoHandler[id]
	if ok {
		return v, nil
	}

	return nil, errors.New("cannot find proto type registered in msg_handle!")
}

func (m *MsgParser) regProtoHandle(name string, fn ProtoHandler) {
	id := crc32.ChecksumIEEE([]byte(name))
	if v, ok := m.protoHandler[id]; ok {
		logger.WithFields(logger.Fields{
			"id":   id,
			"type": v,
		}).Warn("register proto msg id existed")
		return
	}

	m.protoHandler[id] = fn
}

// decode binarys to proto message
func (m *MsgParser) decodeToProto(data []byte) (proto.Message, error) {

	// discard top 8 bytes of message size and message crc id
	byProto := data[8:]

	// get next 2 bytes of message name length
	protoNameLen := binary.LittleEndian.Uint16(byProto[:2])

	if uint16(len(byProto)) < 2+protoNameLen {
		return nil, fmt.Errorf("recv proto msg length < 2+protoNameLen:" + string(byProto))
	}

	// get proto name
	protoTypeName := string(byProto[2 : 2+protoNameLen])
	pType := proto.MessageType(protoTypeName)
	if pType == nil {
		return nil, fmt.Errorf("invalid message<%s>, won't deal with it", protoTypeName)
	}

	// get proto data
	protoData := byProto[2+protoNameLen:]

	// prepare proto struct to be unmarshaled in
	newProto, ok := reflect.New(pType.Elem()).Interface().(proto.Message)
	if !ok {
		return nil, fmt.Errorf("invalid message<%s>, won't deal with it", protoTypeName)
	}

	// unmarshal
	if err := proto.Unmarshal(protoData, newProto); err != nil {
		logger.WithFields(logger.Fields{
			"proto": newProto,
			"error": err,
		}).Warn("Failed to parse proto msg")
		return nil, fmt.Errorf("invalid message<%s>, won't deal with it", protoTypeName)
	}

	return newProto, nil
}

/*
msg Example:
	Header: map[string]string{
		"Content-Type": "application/json",
		"Name": "yokai_client.MC_ClientLogon",
	}
	Body: protoBuf byte
*/
func (m *MsgParser) ParserMessage(sock transport.Socket, msg *transport.Message) {

	// protobuf
	if msg.Type == transport.BodyProtobuf {
		m.parserProtobufBody(sock, msg)
	}

	// json
	if msg.Type == transport.BodyJson {
		m.parserJsonBody(sock, msg)
	}
}

func (m *MsgParser) parserProtobufBody(sock transport.Socket, msg *transport.Message) {
	protoMsgID := crc32.ChecksumIEEE([]byte(msg.Name))
	fn, err := m.getRegProtoHandle(protoMsgID)
	if err != nil {
		logger.WithFields(logger.Fields{
			"message_id":   protoMsgID,
			"message_name": msg.Name,
			"error":        err,
		}).Warn("unregisted proto message received")
		return
	}

	// callback
	fn(sock, msg)
}

func (m *MsgParser) parserJsonBody(sock transport.Socket, msg *transport.Message) {
	/*msg, ok := msg.Body.(string)*/
	//if !ok {
	//logger.WithFields(logger.Fields{
	//"type": msg.Type,
	//"name": msg.Name,
	//"body": msg.Body,
	//}).Warn("parser message to json failed")
	//return
	//}

	//protoMsgID := crc32.ChecksumIEEE([]byte(msg.Name))
	//fn, err := m.getRegProtoHandle(protoMsgID)
	//if err != nil {
	//logger.WithFields(logger.Fields{
	//"message_id":   protoMsgID,
	//"message_name": name,
	//"error":        err,
	//}).Warn("unregisted proto message received")
	//return
	/*}*/

	// callback
	//fn(sock, msg)
}

func (m *MsgParser) handleClientLogon(sock transport.Socket, p *transport.Message) {
	msg, ok := p.Body.(*pbClient.MC_ClientLogon)
	if !ok {
		logger.Warn("Cannot assert value to message")
		return
	}

	client, err := m.g.cm.AddClient(msg.ClientId, msg.ClientName, sock)
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

func (m *MsgParser) handleHeartBeat(sock transport.Socket, p *transport.Message) {
	if client := m.g.cm.GetClientBySock(sock); client != nil {
		if t := int32(time.Now().Unix()); t == -1 {
			logger.Warn("Heart beat get time err")
			return
		}

		client.HeartBeat()
	}
}

func (m *MsgParser) handleClientConnected(sock transport.Socket, p *transport.Message) {
	if client := m.g.cm.GetClientBySock(sock); client != nil {
		clientID := p.Body.(*pbClient.MC_ClientConnected).ClientId
		logger.WithFields(logger.Fields{
			"client_id": clientID,
		}).Info("client connected")

		// todo after connected
	}
}

/*func (m *MsgParser) handleRequestPlayerInfo(con iface.ITCPConn, p proto.Message) {*/
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

//func (m *MsgParser) handleRequestGuildInfo(con iface.ITCPConn, p proto.Message) {
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

//func (m *MsgParser) handlePlayUltimateRecord(con iface.ITCPConn, p proto.Message) {
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

//func (m *MsgParser) handleRequestUltimatePlayer(con iface.ITCPConn, p proto.Message) {
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

//func (m *MsgParser) handleViewFormation(con iface.ITCPConn, p proto.Message) {
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
//func (m *MsgParser) handleArenaMatching(con iface.ITCPConn, p proto.Message) {
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

//func (m *MsgParser) handleArenaAddRecord(con iface.ITCPConn, p proto.Message) {
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

//func (m *MsgParser) handleArenaBattleResult(con iface.ITCPConn, p proto.Message) {
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

//func (m *MsgParser) handleReplacePlayerInfo(con iface.ITCPConn, p proto.Message) {
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

//func (m *MsgParser) handleReplaceGuildInfo(con iface.ITCPConn, p proto.Message) {
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

//func (m *MsgParser) handleRequestArenaRank(con iface.ITCPConn, p proto.Message) {
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

//func (m *MsgParser) handleAddInvite(con iface.ITCPConn, p proto.Message) {
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

//func (m *MsgParser) handleCheckInviteResult(con iface.ITCPConn, p proto.Message) {
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

//func (m *MsgParser) handleInviteRecharge(con iface.ITCPConn, p proto.Message) {
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

//func (m *MsgParser) handleArenaChampionOnline(con iface.ITCPConn, p proto.Message) {
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
