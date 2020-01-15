package game

import (
	"context"
	"fmt"
	"math/rand"
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
	utils.CacheObjector `bson:"-"`
	ID                  int64       `bson:"_id"`
	UserID              int64       `bson:"user_id"`
	Name                string      `bson:"name"`
	Level               int32       `bson:"level"`
	PlayerIDs           []int64     `bson:"player_id"`
	Expire              *time.Timer `bson:"-"`
}

func (la *LiteAccount) GetObjID() interface{} {
	return la.ID
}

func (la *LiteAccount) GetExpire() *time.Timer {
	return la.Expire
}

func (la *LiteAccount) ResetExpire() {
	d := define.Account_MemExpire + time.Second*time.Duration(rand.Intn(60))
	la.Expire.Reset(d)
}

func (la *LiteAccount) StopExpire() {
	la.Expire.Stop()
}

func (la *LiteAccount) GetID() int64 {
	return la.ID
}

func (la *LiteAccount) SetID(id int64) {
	la.ID = id
}

func (la *LiteAccount) GetName() string {
	return la.Name
}

func (la *LiteAccount) SetName(name string) {
	la.Name = name
}

func (la *LiteAccount) GetLevel() int32 {
	return la.Level
}

func (la *LiteAccount) SetLevel(level int32) {
	la.Level = level
}

func (la *LiteAccount) AddPlayerID(playerID int64) {
	for _, value := range la.PlayerIDs {
		if value == playerID {
			return
		}
	}

	la.PlayerIDs = append(la.PlayerIDs, playerID)
}

func (la *LiteAccount) GetPlayerIDs() []int64 {
	return la.PlayerIDs
}

// full account info
type Account struct {
	*LiteAccount `bson:"inline"`

	sock transport.Socket `bson:"-"`
	p    *player.Player   `bson:"-"`

	ctx       context.Context        `bson:"-"`
	cancel    context.CancelFunc     `bson:"-"`
	waitGroup utils.WaitGroupWrapper `bson:"-"`
	timeOut   *time.Timer            `bson:"-"`

	wrapHandler  chan func() `bson:"-"`
	asyncHandler chan func() `bson:"-"`
}

func NewLiteAccount() interface{} {
	return &LiteAccount{
		ID:        -1,
		Name:      "",
		Level:     1,
		Expire:    time.NewTimer(define.Account_MemExpire + time.Second*time.Duration(rand.Intn(60))),
		PlayerIDs: make([]int64, 0),
	}
}

func NewAccount(ctx context.Context, la *LiteAccount, sock transport.Socket) *Account {
	account := &Account{
		LiteAccount: &LiteAccount{
			ID:        la.ID,
			Name:      la.Name,
			Level:     la.Level,
			Expire:    time.NewTimer(define.Account_MemExpire + time.Second*time.Duration(rand.Intn(60))),
			PlayerIDs: la.PlayerIDs,
		},

		sock:         sock,
		p:            nil,
		timeOut:      time.NewTimer(define.Account_OnlineTimeout),
		wrapHandler:  make(chan func(), WrapHandlerSize),
		asyncHandler: make(chan func(), AsyncHandlerSize),
	}

	account.ctx, account.cancel = context.WithCancel(ctx)

	return account
}

func (Account) TableName() string {
	return "account"
}

func (a *Account) GetSock() transport.Socket {
	return a.sock
}

func (a *Account) GetPlayer() *player.Player {
	return a.p
}

func (a *Account) SetPlayer(p *player.Player) {
	a.p = p
}

func (a *Account) Main() error {

	a.waitGroup.Wrap(func() {
		a.Run()
	})

	a.waitGroup.Wait()

	return nil
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
			return fmt.Errorf("account<%d> time out", a.GetID())
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
	a.timeOut.Reset(define.Account_OnlineTimeout)

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
