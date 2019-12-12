package game

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/transport"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

func (m *MsgHandler) handleAddTalent(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("add talent failed")
		return
	}

	msg, ok := p.Body.(*pbGame.MC_AddTalent)
	if !ok {
		logger.Warn("Add Talent failed, recv message body error")
		return
	}

	err := cli.Player().TalentManager().AddTalent(msg.Id)
	if err != nil {
		logger.Warn("talent inc failed:", err)
	}

	cli.Player().TalentManager().Save()

	list := cli.Player().TalentManager().GetTalentList()
	reply := &pbGame.MS_TalentList{Talents: make([]*pbGame.Talent, 0, len(list))}
	for _, v := range list {
		reply.Talents = append(reply.Talents, &pbGame.Talent{
			Id: v.ID,
		})
	}

	cli.SendProtoMessage(reply)
}

func (m *MsgHandler) handleQueryTalents(sock transport.Socket, p *transport.Message) {
	cli := m.g.cm.GetClientBySock(sock)
	if cli == nil {
		logger.WithFields(logger.Fields{
			"client_id":   cli.ID(),
			"client_name": cli.Name(),
		}).Warn("query talents failed")
		return
	}

	list := cli.Player().TalentManager().GetTalentList()
	reply := &pbGame.MS_TalentList{Talents: make([]*pbGame.Talent, 0, len(list))}
	for _, v := range list {
		reply.Talents = append(reply.Talents, &pbGame.Talent{
			Id: v.ID,
		})
	}

	cli.SendProtoMessage(reply)
}
