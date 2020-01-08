package game

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/player"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/transport"
	"github.com/yokaiio/yokai_server/internal/utils"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
)

var WrapHandlerSize int = 100
var AsyncHandlerSize int = 100

// lite account info
type LiteAccount struct {
	ID     int64       `bson:"account_id"`
	Name   string      `bson:"name"`
	Expire *time.Timer `bson:"-"`
}

func (la *LiteAccount) GetExpire() *time.Timer {
	return la.Expire
}

func (la *LiteAccount) ResetExpire() {
	d := define.Account_MemExpire + time.Second*time.Duration(rand.Intn(60))
	la.Expire.Reset(d)
}

// full account info
type Account struct {
	*LiteAccount `bson:"inline"`

	sock transport.Socket `bson:"-"`
	p    *player.Player   `bson:"-"`

	am        *AccountManager        `bson:"-"`
	ctx       context.Context        `bson:"-"`
	cancel    context.CancelFunc     `bson:"-"`
	waitGroup utils.WaitGroupWrapper `bson:"-"`
	timeOut   *time.Timer            `bson:"-"`
	Expire    *time.Timer            `bson:"-"`

	wrapHandler  chan func() `bson:"-"`
	asyncHandler chan func() `bson:"-"`
}

func NewLiteAccount(id int64) *LiteAccount {
	return &LiteAccount{ID: id}
}

func NewAccount(am *AccountManager, info *LiteAccount, sock transport.Socket, p *player.Player) *Account {
	account := &Account{
		LiteAccount: &LiteAccount{
			ID:     info.ID,
			Name:   info.Name,
			Expire: time.NewTimer(define.Account_MemExpire + time.Second*time.Duration(rand.Intn(60))),
		},

		sock:         sock,
		p:            p,
		am:           am,
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

func (a *Account) GetID() int64 {
	return a.ID
}

func (a *Account) GetName() string {
	return a.Name
}

func (a *Account) GetSock() transport.Socket {
	return a.sock
}

func (a *Account) GetPlayer() *player.Player {
	return a.p
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
	a.sock.Close()
}

func (a *Account) Run() error {
	for {
		select {
		// context canceled
		case <-a.ctx.Done():
			logger.WithFields(logger.Fields{
				"id": a.GetID(),
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
			a.am.DisconnectAccount(a, "timeout")
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

	if err := a.sock.Send(&msg); err != nil {
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
			"account_id": a.GetID(),
			"func":       f,
		}).Warn("wrap handler channel full, ignored.")
		return
	}

	a.wrapHandler <- f
}

func (a *Account) PushAsyncHandler(f func()) {
	a.asyncHandler <- f
}

func (a *Account) ResetExpire() {
}
