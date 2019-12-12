package game

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/transport"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

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
	reply := &pbGame.MS_ItemList{Items: make([]*pbGame.Item, 0, len(list))}
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
	reply := &pbGame.MS_ItemList{Items: make([]*pbGame.Item, 0, len(list))}
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
	reply := &pbGame.MS_ItemList{Items: make([]*pbGame.Item, 0, len(list))}
	for _, v := range list {
		i := &pbGame.Item{
			Id:     v.GetID(),
			TypeId: v.GetTypeID(),
		}
		reply.Items = append(reply.Items, i)
	}
	cli.SendProtoMessage(reply)
}

func (m *MsgHandler) handlePutonEquip(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("puton equip failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_PutonEquip)
	if !ok {
		logger.Warn("Puton equip failed, recv message body error")
		return
	}

	cli.Player().ItemManager().GetItem(msg.EquipId)
	list := cli.Player().HeroManager().PutonEquip()
	reply := &pbGame.MS_ItemList{Items: make([]*pbGame.Item, 0, len(list))}
	for _, v := range list {
		i := &pbGame.Item{
			Id:     v.GetID(),
			TypeId: v.GetTypeID(),
		}
		reply.Items = append(reply.Items, i)
	}
	cli.SendProtoMessage(reply)
}

func (m *MsgHandler) handleTakeoffEquip(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("query items failed")
		return
	}

	list := cli.Player().ItemManager().GetItemList()
	reply := &pbGame.MS_ItemList{Items: make([]*pbGame.Item, 0, len(list))}
	for _, v := range list {
		i := &pbGame.Item{
			Id:     v.GetID(),
			TypeId: v.GetTypeID(),
		}
		reply.Items = append(reply.Items, i)
	}
	cli.SendProtoMessage(reply)
}
