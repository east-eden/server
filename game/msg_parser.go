package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hellodudu/Ultimate/game-service/define"
	"github.com/hellodudu/Ultimate/iface"
	pbArena "github.com/hellodudu/Ultimate/proto/arena"
	pbGame "github.com/hellodudu/Ultimate/proto/game"
	pbWorld "github.com/hellodudu/Ultimate/proto/world"
	"github.com/hellodudu/Ultimate/utils"
	logger "github.com/sirupsen/logrus"
)

// ProtoHandler handle function
type ProtoHandler func(iface.ITCPConn, proto.Message)

type MsgParser struct {
	protoHandler map[uint32]ProtoHandler
	gm           iface.IGameMgr
	wm           iface.IWorldMgr
}

func NewMsgParser(gm iface.IGameMgr, wm iface.IWorldMgr) *MsgParser {
	m := &MsgParser{
		protoHandler: make(map[uint32]ProtoHandler),
		gm:           gm,
		wm:           wm,
	}

	m.registerAllMessage()
	return m
}

func (m *MsgParser) registerAllMessage() {
	m.regProtoHandle("ultimate_service_world.MWU_WorldLogon", m.handleWorldLogon)
	m.regProtoHandle("ultimate_service_world.MWU_TestConnect", m.handleTestConnect)
	m.regProtoHandle("ultimate_service_world.MWU_HeartBeat", m.handleHeartBeat)
	m.regProtoHandle("ultimate_service_world.MWU_WorldConnected", m.handleWorldConnected)
	m.regProtoHandle("ultimate_service_game.MWU_RequestPlayerInfo", m.handleRequestPlayerInfo)
	m.regProtoHandle("ultimate_service_game.MWU_RequestGuildInfo", m.handleRequestGuildInfo)
	m.regProtoHandle("ultimate_service_game.MWU_PlayUltimateRecord", m.handlePlayUltimateRecord)
	m.regProtoHandle("ultimate_service_game.MWU_RequestUltimatePlayer", m.handleRequestUltimatePlayer)
	m.regProtoHandle("ultimate_service_game.MWU_ViewFormation", m.handleViewFormation)
	m.regProtoHandle("ultimate_service_game.MWU_AddInvite", m.handleAddInvite)
	m.regProtoHandle("ultimate_service_game.MWU_CheckInviteResult", m.handleCheckInviteResult)
	m.regProtoHandle("ultimate_service_game.MWU_InviteRecharge", m.handleInviteRecharge)
	m.regProtoHandle("ultimate_service_game.MWU_ReplacePlayerInfo", m.handleReplacePlayerInfo)
	m.regProtoHandle("ultimate_service_game.MWU_ReplaceGuildInfo", m.handleReplaceGuildInfo)
	m.regProtoHandle("ultimate_service_arena.MWU_ArenaMatching", m.handleArenaMatching)
	m.regProtoHandle("ultimate_service_arena.MWU_ArenaAddRecord", m.handleArenaAddRecord)
	m.regProtoHandle("ultimate_service_arena.MWU_ArenaBattleResult", m.handleArenaBattleResult)
	m.regProtoHandle("ultimate_service_arena.MWU_RequestArenaRank", m.handleRequestArenaRank)
	m.regProtoHandle("ultimate_service_arena.MWU_ArenaChampionOnline", m.handleArenaChampionOnline)

}

func (m *MsgParser) getRegProtoHandle(id uint32) (ProtoHandler, error) {
	v, ok := m.protoHandler[id]
	if ok {
		return v, nil
	}

	return nil, errors.New("cannot find proto type registered in msg_handle!")
}

func (m *MsgParser) regProtoHandle(name string, fn ProtoHandler) {
	id := utils.Crc32(name)
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

// top 8 bytes are baseNetMsg
// if it is protobuf msg, then next 2 bytes are proto name length, the next is proto name, final is proto data.
// if it is transfer msg(transfer binarys to other world), then next are binarys to be transferred
func (m *MsgParser) ParserMessage(con iface.ITCPConn, data []byte) {
	if len(data) <= 8 {
		logger.WithFields(logger.Fields{
			"data": string(data),
		}).Warn("tcp recv data length <= 8")
		return
	}

	baseMsg := &define.BaseNetMsg{}
	byBaseMsg := make([]byte, binary.Size(baseMsg))

	copy(byBaseMsg, data[:binary.Size(baseMsg)])
	buf := &bytes.Buffer{}
	if _, err := buf.Write(byBaseMsg); err != nil {
		logger.WithFields(logger.Fields{
			"base_msg": byBaseMsg,
			"con":      con,
			"error":    err,
		}).Warn("cannot read message from connection")
		return
	}

	// get top 4 bytes messageid
	if err := binary.Read(buf, binary.LittleEndian, baseMsg); err != nil {
		logger.WithFields(logger.Fields{
			"base_msg": byBaseMsg,
			"con":      con,
			"error":    err,
		}).Warn("cannot read message from connection")
		return
	}

	// proto message
	if baseMsg.ID == utils.Crc32(string("MWU_DirectProtoMsg")) {
		newProto, err := m.decodeToProto(data)
		if err != nil {
			logger.Warn(err)
			return
		}

		protoMsgName := proto.MessageName(newProto)
		protoMsgID := utils.Crc32(protoMsgName)
		fn, err := m.getRegProtoHandle(protoMsgID)
		if err != nil {
			logger.WithFields(logger.Fields{
				"message_id":   protoMsgID,
				"message_name": protoMsgName,
				"error":        err,
			}).Warn("unregisted proto message received")
			return
		}

		// callback
		fn(con, newProto)

		// transfer message
	} else if baseMsg.ID == utils.Crc32(string("MWU_TransferMsg")) {
		transferMsg := &define.TransferNetMsg{}
		byTransferMsg := make([]byte, binary.Size(transferMsg))

		copy(byTransferMsg, data[:binary.Size(transferMsg)])
		buf := &bytes.Buffer{}
		if _, err := buf.Write(byTransferMsg); err != nil {
			logger.WithFields(logger.Fields{
				"transfer_msg": byTransferMsg,
				"con":          con,
				"error":        err,
			}).Warn("cannot read message from connection")
			return
		}

		// get top 4 bytes messageid
		if err := binary.Read(buf, binary.LittleEndian, transferMsg); err != nil {
			logger.WithFields(logger.Fields{
				"transfer_msg": byTransferMsg,
				"con":          con,
				"error":        err,
			}).Warn("cannot read message from connection")
			return
		}

		// send message to world
		sendWorld := m.wm.GetWorldByID(transferMsg.WorldID)
		if sendWorld == nil {
			logger.WithFields(logger.Fields{
				"world_id": transferMsg.WorldID,
			}).Warn("send transfer message to unconnected world")
			return
		}

		sendWorld.SendTransferMessage(data)
	}

}

func (m *MsgParser) handleWorldLogon(con iface.ITCPConn, p proto.Message) {
	msg, ok := p.(*pbWorld.MWU_WorldLogon)
	if !ok {
		logger.Warn("Cannot assert value to message")
		return
	}

	world, err := m.wm.AddWorld(msg.WorldId, msg.WorldName, con)
	if err != nil {
		logger.WithFields(logger.Fields{
			"id":   msg.WorldId,
			"name": msg.WorldName,
			"con":  con,
		}).Warn("add world failed")
		return
	}

	reply := &pbWorld.MUW_WorldLogon{}
	world.SendProtoMessage(reply)

}

func (m *MsgParser) handleTestConnect(con iface.ITCPConn, p proto.Message) {
	if world := m.wm.GetWorldByCon(con); world != nil {
		world.ResetTestConnect()
	}
}

func (m *MsgParser) handleHeartBeat(con iface.ITCPConn, p proto.Message) {
	if world := m.wm.GetWorldByCon(con); world != nil {
		if t := int32(time.Now().Unix()); t == -1 {
			logger.Warn("Heart beat get time err")
			return
		}

		reply := &pbWorld.MUW_HeartBeat{BattleTime: uint32(time.Now().Unix())}
		world.SendProtoMessage(reply)
	}
}

func (m *MsgParser) handleWorldConnected(con iface.ITCPConn, p proto.Message) {
	if world := m.wm.GetWorldByCon(con); world != nil {
		arrWorldID := p.(*pbWorld.MWU_WorldConnected).WorldId
		logger.WithFields(logger.Fields{
			"ref_id": arrWorldID,
		}).Info("world ref connected")

		// add reference world id
		m.wm.AddWorldRef(world.GetID(), arrWorldID)

		// request player info
		msgP := &pbGame.MUW_RequestPlayerInfo{MinLevel: 20}
		world.SendProtoMessage(msgP)

		// request guild info
		msgG := &pbGame.MUW_RequestGuildInfo{}
		world.SendProtoMessage(msgG)

		// sync arena data
		if season, seasonEndTime, err := m.gm.GetArenaSeasonData(); err == nil {
			logger.WithFields(logger.Fields{
				"season": season,
				"time":   seasonEndTime,
			}).Info("GetArenaSeasonData success")
			msgArena := &pbArena.MUW_SyncArenaSeason{
				Season:  season,
				EndTime: uint32(seasonEndTime),
			}
			world.SendProtoMessage(msgArena)
		}

		// 20s later sync arena champion
		t := time.NewTimer(20 * time.Second)
		go func(id uint32) {
			<-t.C
			w := m.wm.GetWorldByID(id)
			if w == nil {
				logger.WithFields(logger.Fields{
					"world_id": id,
				}).Warn("world disconnected, cannot sync arena champion")
				return
			}

			if championList, err := m.gm.GetArenaChampion(); err != nil {
				msg := &pbArena.MUW_ArenaChampion{
					Data: championList,
				}

				w.SendProtoMessage(msg)
				logger.WithFields(logger.Fields{
					"world_id":   w.GetID(),
					"world_name": w.GetName(),
				}).Info("sync arena champion to world")
			}
		}(world.GetID())
	}
}

func (m *MsgParser) handleRequestPlayerInfo(con iface.ITCPConn, p proto.Message) {
	if world := m.wm.GetWorldByCon(con); world != nil {
		msg, ok := p.(*pbGame.MWU_RequestPlayerInfo)
		if !ok {
			logger.WithFields(logger.Fields{
				"msg_name": proto.MessageName(p),
			}).Warn("parsing message name error")
			return
		}

		m.gm.AddPlayerInfoList(msg.Info)
	}
}

func (m *MsgParser) handleRequestGuildInfo(con iface.ITCPConn, p proto.Message) {
	if world := m.wm.GetWorldByCon(con); world != nil {
		msg, ok := p.(*pbGame.MWU_RequestGuildInfo)
		if !ok {
			logger.WithFields(logger.Fields{
				"msg_name": proto.MessageName(p),
			}).Warn("parsing message name error")
			return
		}

		m.gm.AddGuildInfoList(msg.Info)
	}
}

func (m *MsgParser) handlePlayUltimateRecord(con iface.ITCPConn, p proto.Message) {
	if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
		msg, ok := p.(*pbGame.MWU_PlayUltimateRecord)
		if !ok {
			logger.WithFields(logger.Fields{
				"msg_name": proto.MessageName(p),
			}).Warn("parsing message name error")
			return
		}

		dstWorld := m.wm.GetWorldByID(msg.DstServerId)
		if dstWorld == nil {
			return
		}

		msgSend := &pbGame.MUW_PlayUltimateRecord{
			SrcPlayerId: msg.SrcPlayerId,
			SrcServerId: msg.SrcServerId,
			RecordId:    msg.RecordId,
			DstServerId: msg.DstServerId,
		}
		dstWorld.SendProtoMessage(msgSend)
	}
}

func (m *MsgParser) handleRequestUltimatePlayer(con iface.ITCPConn, p proto.Message) {
	if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
		msg, ok := p.(*pbGame.MWU_RequestUltimatePlayer)
		if !ok {
			logger.WithFields(logger.Fields{
				"msg_name": proto.MessageName(p),
			}).Warn("parsing message name error")
			return
		}

		dstInfo, err := m.gm.GetPlayerInfoByID(msg.DstPlayerId)
		dstWorld := m.wm.GetWorldByID(msg.DstServerId)
		if err != nil {
			return
		}

		if int32(msg.DstServerId) == -1 {
			dstWorld = m.wm.GetWorldByID(dstInfo.ServerId)
		}

		if dstWorld == nil {
			return
		}

		msgSend := &pbGame.MUW_RequestUltimatePlayer{
			SrcPlayerId: msg.SrcPlayerId,
			SrcServerId: msg.SrcServerId,
			DstPlayerId: msg.DstPlayerId,
			DstServerId: dstWorld.GetID(),
		}
		dstWorld.SendProtoMessage(msgSend)
	}
}

func (m *MsgParser) handleViewFormation(con iface.ITCPConn, p proto.Message) {
	if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
		msg, ok := p.(*pbGame.MWU_ViewFormation)
		if !ok {
			logger.WithFields(logger.Fields{
				"msg_name": proto.MessageName(p),
			}).Warn("parsing message name error")
			return
		}

		dstInfo, err := m.gm.GetPlayerInfoByID(msg.DstPlayerId)
		dstWorld := m.wm.GetWorldByID(msg.DstServerId)
		if err != nil {
			return
		}

		if int32(msg.DstServerId) == -1 {
			dstWorld = m.wm.GetWorldByID(dstInfo.ServerId)
		}

		if dstWorld == nil {
			return
		}

		msgSend := &pbGame.MUW_ViewFormation{
			SrcPlayerId: msg.SrcPlayerId,
			SrcServerId: msg.SrcServerId,
			DstPlayerId: msg.DstPlayerId,
			DstServerId: dstWorld.GetID(),
		}
		dstWorld.SendProtoMessage(msgSend)
	}
}

///////////////////////////////
// arena battle
//////////////////////////////
func (m *MsgParser) handleArenaMatching(con iface.ITCPConn, p proto.Message) {
	if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
		msg, ok := p.(*pbArena.MWU_ArenaMatching)
		if !ok {
			logger.WithFields(logger.Fields{
				"msg_name": proto.MessageName(p),
			}).Warn("parsing message name error")
			return
		}

		m.gm.ArenaMatching(msg.PlayerId)
	}
}

func (m *MsgParser) handleArenaAddRecord(con iface.ITCPConn, p proto.Message) {
	if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
		msg, ok := p.(*pbArena.MWU_ArenaAddRecord)
		if !ok {
			logger.WithFields(logger.Fields{
				"msg_name": proto.MessageName(p),
			}).Warn("parsing message name error")
			return
		}

		m.gm.ArenaAddRecord(msg.Record)
	}
}

func (m *MsgParser) handleArenaBattleResult(con iface.ITCPConn, p proto.Message) {
	if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
		msg, ok := p.(*pbArena.MWU_ArenaBattleResult)
		if !ok {
			logger.WithFields(logger.Fields{
				"msg_name": proto.MessageName(p),
			}).Warn("parsing message name error")
			return
		}

		m.gm.ArenaBattleResult(msg.AttackId, msg.TargetId, msg.AttackWin)
	}
}

func (m *MsgParser) handleReplacePlayerInfo(con iface.ITCPConn, p proto.Message) {
	if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
		msg, ok := p.(*pbGame.MWU_ReplacePlayerInfo)
		if !ok {
			logger.WithFields(logger.Fields{
				"msg_name": proto.MessageName(p),
			}).Warn("parsing message name error")
			return
		}

		m.gm.AddPlayerInfo(msg.Info)
	}
}

func (m *MsgParser) handleReplaceGuildInfo(con iface.ITCPConn, p proto.Message) {
	if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
		msg, ok := p.(*pbGame.MWU_ReplaceGuildInfo)
		if !ok {
			logger.WithFields(logger.Fields{
				"msg_name": proto.MessageName(p),
			}).Warn("parsing message name error")
			return
		}

		m.gm.AddGuildInfo(msg.Info)
	}
}

func (m *MsgParser) handleRequestArenaRank(con iface.ITCPConn, p proto.Message) {
	if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
		msg, ok := p.(*pbArena.MWU_RequestArenaRank)
		if !ok {
			logger.WithFields(logger.Fields{
				"msg_name": proto.MessageName(p),
			}).Warn("parsing message name error")
			return
		}

		m.gm.ArenaGetRank(msg.PlayerId, msg.Page)
	}
}

func (m *MsgParser) handleAddInvite(con iface.ITCPConn, p proto.Message) {
	if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
		msg, ok := p.(*pbGame.MWU_AddInvite)
		if !ok {
			logger.WithFields(logger.Fields{
				"msg_name": proto.MessageName(p),
			}).Warn("parsing message name error")
			return
		}

		ret := m.gm.Invite().AddInvite(msg.NewbieId, msg.InviterId)
		if ret != 0 {
			msgRet := &pbGame.MUW_AddInviteResult{
				NewbieId:  msg.NewbieId,
				InviterId: msg.InviterId,
				ErrorCode: ret,
			}

			srcWorld.SendProtoMessage(msgRet)
		}
	}
}

func (m *MsgParser) handleCheckInviteResult(con iface.ITCPConn, p proto.Message) {
	if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
		msg, ok := p.(*pbGame.MWU_CheckInviteResult)
		if !ok {
			logger.WithFields(logger.Fields{
				"msg_name": proto.MessageName(p),
			}).Warn("parsing message name error")
			return
		}

		m.gm.Invite().CheckInviteResult(msg.NewbieId, msg.InviterId, msg.ErrorCode)
	}
}

func (m *MsgParser) handleInviteRecharge(con iface.ITCPConn, p proto.Message) {
	if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
		msg, ok := p.(*pbGame.MWU_InviteRecharge)
		if !ok {
			logger.WithFields(logger.Fields{
				"msg_name": proto.MessageName(p),
			}).Warn("parsing message name error")
			return
		}

		m.gm.Invite().InviteRecharge(msg.NewbieId, msg.NewbieName, msg.InviterId, msg.DiamondGift)
	}
}

func (m *MsgParser) handleArenaChampionOnline(con iface.ITCPConn, p proto.Message) {
	if srcWorld := m.wm.GetWorldByCon(con); srcWorld != nil {
		msg, ok := p.(*pbArena.MWU_ArenaChampionOnline)
		if !ok {
			logger.WithFields(logger.Fields{
				"msg_name": proto.MessageName(p),
			}).Warn("parsing message name error")
			return
		}

		msgSend := &pbArena.MUW_ArenaChampionOnline{
			PlayerId:   msg.PlayerId,
			PlayerName: msg.PlayerName,
			ServerName: msg.ServerName,
		}

		m.wm.BroadCast(msgSend)
	}
}
