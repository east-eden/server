package game

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/player"
	"github.com/yokaiio/yokai_server/internal/transport"
	"github.com/yokaiio/yokai_server/internal/utils"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
)

var WrapHandlerSize int = 100
var AsyncHandlerSize int = 100

type AccountInfo struct {
	ID   int64
	Name string
	sock transport.Socket
	p    player.Player
}

type Account struct {
	info *AccountInfo

	am        *AccountManager
	ctx       context.Context
	cancel    context.CancelFunc
	waitGroup utils.WaitGroupWrapper
	timeOut   *time.Timer

	wrapHandler  chan func()
	asyncHandler chan func()
}

func NewAccount(am *AccountManager, info *AccountInfo) *Account {
	account := &Account{
		am:           am,
		info:         info,
		timeOut:      time.NewTimer(am.accountTimeout),
		wrapHandler:  make(chan func(), WrapHandlerSize),
		asyncHandler: make(chan func(), AsyncHandlerSize),
	}

	account.ctx, account.cancel = context.WithCancel(am.ctx)

	return account
}

func (Account) TableName() string {
	return "Account"
}

func (a *Account) ID() int64 {
	return a.info.ID
}

func (a *Account) Name() string {
	return a.info.Name
}

func (a *Account) Sock() transport.Socket {
	return a.info.sock
}

func (a *Account) Player() player.Player {
	return a.info.p
}

func (a *Account) Main() error {

	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal("Account Main() error:", err)
			}
			exitCh <- err
		})
	}

	a.waitGroup.Wrap(func() {
		exitFunc(a.Run())
	})

	return <-exitCh
}

func (a *Account) Exit() {
	a.timeOut.Stop()
	a.info.sock.Close()
}

func (a *Account) Run() error {
	for {
		select {
		// context canceled
		case <-a.ctx.Done():
			logger.WithFields(logger.Fields{
				"id": a.ID(),
			}).Info("Account context done!")
			return nil

		// async handler
		case h := <-a.asyncHandler:
			h()

		// request handler
		case h := <-a.wrapHandler:
			t := time.Now()
			h()
			d := time.Since(t)
			time.Sleep(time.Millisecond*100 - d)

		// lost connection
		case <-a.timeOut.C:
			a.am.DisconnectAccount(a.info.sock, "timeout")
		}
	}
}

/*
msg Example:
	Type: transport.BodyProtobuf
	Name: yokai_account.MS_AccountLogon
	Body: protoBuf byte
*/
func (a *Account) SendProtoMessage(p proto.Message) {
	var msg transport.Message
	msg.Type = transport.BodyProtobuf
	msg.Name = proto.MessageName(p)
	msg.Body = p

	if err := a.info.sock.Send(&msg); err != nil {
		logger.Warn("send proto msg error:", err)
		return
	}
}

func (a *Account) HeartBeat() {
	a.timeOut.Reset(a.am.accountTimeout)

	reply := &pbAccount.MS_HeartBeat{Timestamp: uint32(time.Now().Unix())}
	a.SendProtoMessage(reply)
}

func (a *Account) PushWrapHandler(f func()) {
	if len(a.wrapHandler) >= WrapHandlerSize {
		logger.WithFields(logger.Fields{
			"account_id": a.ID(),
			"func":       f,
		}).Warn("wrap handler channel full, ignored.")
		return
	}

	a.wrapHandler <- f
}

func (a *Account) PushAsyncHandler(f func()) {
	a.asyncHandler <- f
}
