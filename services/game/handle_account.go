package game

import (
	"context"
	"errors"
	"fmt"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/player"
	"bitbucket.org/funplus/server/transport"
	"bitbucket.org/funplus/server/utils"
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ErrUnregistedMsgName = errors.New("unregisted message name")
)

func (m *MsgRegister) handleAccountTest(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	return nil
}

func (m *MsgRegister) handleWaitResponseMessage(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_WaitResponseMessage)
	if !ok {
		return errors.New("handleWaitResponseMessage failed: cannot assert value to message")
	}

	handler, err := m.r.GetHandler(msg.GetInnerMsgCrc())
	if !utils.ErrCheck(err, "handleWaitResponseMessage GetHandler by MsgCrc failed", msg.GetInnerMsgCrc()) {
		return ErrUnregistedMsgName
	}

	var innerMsg transport.Message
	innerMsg.Name = handler.Name
	innerMsg.Body, err = sock.PbMarshaler().Unmarshal(msg.GetInnerMsgData(), handler.RType)
	if !utils.ErrCheck(err, "handleWaitResponseMessage protobuf Unmarshal failed") {
		return err
	}

	// direct handle inner message
	err = handler.Fn(ctx, sock, &innerMsg)
	if !utils.ErrCheck(err, "handle inner message failed", handler.Name) {
		return err
	}

	err = m.am.AddAccountTask(
		m.am.GetAccountIdBySock(sock),
		&player.AccountTasker{
			C: ctx,
			F: func(ctx context.Context, acct *player.Account, _ *transport.Message) error {
				reply := &pbGlobal.S2C_WaitResponseMessage{
					MsgId:   msg.MsgId,
					ErrCode: 0,
				}

				acct.SendProtoMessage(reply)
				return nil
			},
			M: p,
		},
	)

	return err
}

func (m *MsgRegister) handleAccountPing(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_Ping)
	if !ok {
		return errors.New("handleAccountLogon failed: cannot assert value to message")
	}

	reply := &pbGlobal.S2C_Pong{
		Pong: msg.Ping + 1,
	}

	var send transport.Message
	send.Name = string(proto.MessageReflect(reply).Descriptor().Name())
	send.Body = reply

	return sock.Send(&send)
}

func (m *MsgRegister) handleAccountLogon(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_AccountLogon)
	if !ok {
		return errors.New("handleAccountLogon failed: cannot assert value to message")
	}

	err := m.am.Logon(ctx, msg.UserId, msg.AccountId, msg.AccountName, sock)
	if err != nil {
		return fmt.Errorf("handleAccountLogon failed: %w", err)
	}

	// err = m.am.AccountLazyHandle(
	// 	m.am.GetAccountIdBySock(sock),
	// 	&player.AccountLazyHandler{
	// 		F: func(ctx context.Context, acct *player.Account, _ *transport.Message) error {
	// 			reply := &pbGlobal.S2C_AccountLogon{
	// 				UserId:      acct.UserId,
	// 				AccountId:   acct.ID,
	// 				PlayerId:    -1,
	// 				PlayerName:  "",
	// 				PlayerLevel: 0,
	// 			}

	// 			if p := acct.GetPlayer(); p != nil {
	// 				reply.PlayerId = p.GetID()
	// 				reply.PlayerName = p.GetName()
	// 				reply.PlayerLevel = p.GetLevel()
	// 			}

	// 			acct.SendProtoMessage(reply)
	// 			return nil
	// 		},
	// 		M: p,
	// 	},
	// )

	return err
}

func (m *MsgRegister) handleHeartBeat(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		m.timeHistogram.WithLabelValues("handleHeartBeat").Observe(v)
	}))

	_, ok := p.Body.(*pbGlobal.C2S_HeartBeat)
	if !ok {
		return errors.New("handleHeartBeat failed: cannot assert value to message")
	}

	err := m.am.AddAccountTask(
		m.am.GetAccountIdBySock(sock),
		&player.AccountTasker{
			C: ctx,
			F: func(ctx context.Context, acct *player.Account, _ *transport.Message) error {
				defer timer.ObserveDuration()
				acct.HeartBeat()
				return nil
			},
			M: p,
		},
	)

	return err
}

// client disconnect
func (m *MsgRegister) handleAccountDisconnect(ctx context.Context, sock transport.Socket, p *transport.Message) error {
	return player.ErrAccountDisconnect
}
