package game

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"

	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
	"e.coding.net/mmstudio/blade/server/services/game/player"
	"e.coding.net/mmstudio/blade/server/transport"
	"e.coding.net/mmstudio/blade/server/transport/codec"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/protobuf/proto"
)

var (
	ErrUnregistedMsgName = errors.New("unregisted message name")
)

func (m *MsgRegister) handleAccountTest(ctx context.Context, sock transport.Socket, p proto.Message) error {
	return nil
}

func (m *MsgRegister) handleWaitResponseMessage(ctx context.Context, sock transport.Socket, p proto.Message) error {
	msg, ok := p.(*pbGlobal.C2S_WaitResponseMessage)
	if !ok {
		return errors.New("handleWaitResponseMessage failed: cannot assert value to message")
	}

	handler, err := m.r.GetHandler(msg.GetInnerMsgCrc())
	if !utils.ErrCheck(err, "handleWaitResponseMessage GetHandler by MsgCrc failed", msg.GetInnerMsgCrc()) {
		return ErrUnregistedMsgName
	}

	codec := &codec.ProtoBufMarshaler{}
	innerMsg, err := codec.Unmarshal(msg.GetInnerMsgData(), handler.RType)
	if !utils.ErrCheck(err, "handleWaitResponseMessage protobuf Unmarshal failed") {
		return err
	}

	// direct handle inner message
	err = handler.Fn(ctx, sock, innerMsg.(proto.Message))
	if !utils.ErrCheck(err, "handle inner message failed", handler.Name) {
		return err
	}

	accountId, ok := m.am.GetAccountIdBySock(sock)
	if !ok {
		return ErrAccountNotFound
	}

	err = m.am.AddAccountTask(
		ctx,
		accountId,
		func(c context.Context, p ...interface{}) error {
			acct := p[0].(*player.Account)
			m := p[1].(*pbGlobal.C2S_WaitResponseMessage)
			reply := &pbGlobal.S2C_WaitResponseMessage{
				MsgId:   m.MsgId,
				ErrCode: 0,
			}

			acct.SendProtoMessage(reply)
			return nil
		},
		msg,
	)

	return err
}

func (m *MsgRegister) handleAccountPing(ctx context.Context, sock transport.Socket, p proto.Message) error {
	msg, ok := p.(*pbGlobal.C2S_Ping)
	if !ok {
		return errors.New("handleAccountLogon failed: cannot assert value to message")
	}

	reply := &pbGlobal.S2C_Pong{
		Pong: msg.Ping + 1,
	}

	return sock.Send(reply)
}

func (m *MsgRegister) handleAccountLogon(ctx context.Context, sock transport.Socket, p proto.Message) error {
	msg, ok := p.(*pbGlobal.C2S_AccountLogon)
	if !ok {
		return errors.New("handleAccountLogon failed: cannot assert value to message")
	}

	// todo userid暂时为crc32
	userId := crc32.ChecksumIEEE([]byte(msg.UserId))
	err := m.am.Logon(ctx, int64(userId), sock)
	if err != nil {
		return fmt.Errorf("handleAccountLogon failed: %w", err)
	}

	return err
}

func (m *MsgRegister) handleHeartBeat(ctx context.Context, sock transport.Socket, p proto.Message) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		m.timeHistogram.WithLabelValues("handleHeartBeat").Observe(v)
	}))

	_, ok := p.(*pbGlobal.C2S_HeartBeat)
	if !ok {
		return errors.New("handleHeartBeat failed: cannot assert value to message")
	}

	accountId, ok := m.am.GetAccountIdBySock(sock)
	if !ok {
		return fmt.Errorf("error: %w, sock.LocalAddr: %s, sock.RemoteAddr: %s", ErrAccountNotFound, sock.Local(), sock.Remote())
	}

	err := m.am.AddAccountTask(
		ctx,
		accountId,
		func(c context.Context, p ...interface{}) error {
			acct := p[0].(*player.Account)
			defer timer.ObserveDuration()
			acct.HeartBeat()
			return nil
		},
	)

	return err
}

// client disconnect
func (m *MsgRegister) handleAccountDisconnect(ctx context.Context, sock transport.Socket, p proto.Message) error {
	return player.ErrAccountDisconnect
}
