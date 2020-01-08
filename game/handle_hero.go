package game

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/transport"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

func (m *MsgHandler) handleAddHero(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("add hero failed")
		return
	}

	if acct.GetPlayer() == nil {
		return
	}

	msg, ok := p.Body.(*pbGame.MC_AddHero)
	if !ok {
		logger.Warn("Add Hero failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		acct.GetPlayer().HeroManager().AddHero(msg.TypeId)
		list := acct.GetPlayer().HeroManager().GetHeroList()
		reply := &pbGame.MS_HeroList{Heros: make([]*pbGame.Hero, 0, len(list))}
		for _, v := range list {
			h := &pbGame.Hero{
				Id:     v.GetID(),
				TypeId: v.GetTypeID(),
				Exp:    v.GetExp(),
				Level:  v.GetLevel(),
			}
			reply.Heros = append(reply.Heros, h)
		}
		acct.SendProtoMessage(reply)
	})

}

func (m *MsgHandler) handleDelHero(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("delete hero failed")
		return
	}

	if acct.GetPlayer() == nil {
		return
	}

	msg, ok := p.Body.(*pbGame.MC_DelHero)
	if !ok {
		logger.Warn("Delete Hero failed, recv message body error")
		return
	}

	acct.PushWrapHandler(func() {
		acct.GetPlayer().HeroManager().DelHero(msg.Id)
		list := acct.GetPlayer().HeroManager().GetHeroList()
		reply := &pbGame.MS_HeroList{Heros: make([]*pbGame.Hero, 0, len(list))}
		for _, v := range list {
			h := &pbGame.Hero{
				Id:     v.GetID(),
				TypeId: v.GetTypeID(),
				Exp:    v.GetExp(),
				Level:  v.GetLevel(),
			}
			reply.Heros = append(reply.Heros, h)
		}
		acct.SendProtoMessage(reply)
	})
}

func (m *MsgHandler) handleQueryHeros(sock transport.Socket, p *transport.Message) {
	acct := m.g.am.GetAccountBySock(sock)
	if acct == nil {
		logger.WithFields(logger.Fields{
			"account_id":   acct.GetID(),
			"account_name": acct.GetName(),
		}).Warn("query heros failed")
		return
	}

	if acct.GetPlayer() == nil {
		return
	}

	acct.PushWrapHandler(func() {
		list := acct.GetPlayer().HeroManager().GetHeroList()
		reply := &pbGame.MS_HeroList{Heros: make([]*pbGame.Hero, 0, len(list))}
		for _, v := range list {
			h := &pbGame.Hero{
				Id:     v.GetID(),
				TypeId: v.GetTypeID(),
				Exp:    v.GetExp(),
				Level:  v.GetLevel(),
			}
			reply.Heros = append(reply.Heros, h)
		}
		acct.SendProtoMessage(reply)
	})

}

//func (m *MsgHandler) handleHeroAddExp(sock transport.Socket, p *transport.Message) {
//cli := m.g.cm.GetClientBySock(sock)
//if cli == nil {
//logger.WithFields(logger.Fields{
//"client_id":   cli.ID(),
//"client_name": cli.Name(),
//}).Warn("hero add exp failed")
//return
//}

//msg, ok := p.Body.(*pbGame.MC_HeroAddExp)
//if !ok {
//logger.Warn("hero add exp failed, recv message body error")
//return
//}

//if cli.Player() == nil {
//logger.Warn("client has no player", cli.ID())
//return
//}

//cli.Player().HeroManager().HeroAddExp(msg.HeroId, msg.Exp)
//hero := cli.Player().HeroManager().GetHero(msg.HeroId)
//if hero == nil {
//logger.Warn("get hero by id error:", msg.HeroId)
//return
//}

//reply := &pbGame.MS_HeroInfo{
//Info: &pbGame.Hero{
//Id:     hero.GetID(),
//TypeId: hero.GetTypeID(),
//Exp:    hero.GetExp(),
//Level:  hero.GetLevel(),
//},
//}

//cli.SendProtoMessage(reply)
//}

//func (m *MsgHandler) handleHeroAddLevel(sock transport.Socket, p *transport.Message) {
//cli := m.g.cm.GetClientBySock(sock)
//if cli == nil {
//logger.WithFields(logger.Fields{
//"client_id":   cli.ID(),
//"client_name": cli.Name(),
//}).Warn("hero add level failed")
//return
//}

//msg, ok := p.Body.(*pbGame.MC_HeroAddLevel)
//if !ok {
//logger.Warn("hero add level failed, recv message body error")
//return
//}

//if cli.Player() == nil {
//logger.Warn("client has no player", cli.ID())
//return
//}

//cli.Player().HeroManager().HeroAddLevel(msg.HeroId, msg.Level)
//hero := cli.Player().HeroManager().GetHero(msg.HeroId)
//if hero == nil {
//logger.Warn("get hero by id error:", msg.HeroId)
//return
//}

//reply := &pbGame.MS_HeroInfo{
//Info: &pbGame.Hero{
//Id:     hero.GetID(),
//TypeId: hero.GetTypeID(),
//Exp:    hero.GetExp(),
//Level:  hero.GetLevel(),
//},
//}

//cli.SendProtoMessage(reply)
//}
