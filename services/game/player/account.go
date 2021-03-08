package player

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/east-eden/server/define"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/transport"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

var (
	ErrAccountDisconnect       = errors.New("account disconnect")                                             // handleSocket got this error will disconnect account
	ErrCreateMoreThanOnePlayer = errors.New("AccountManager.CreatePlayer failed: only can create one player") // only can create one player
	Account_MemExpire          = time.Hour * 2
)

// account delay handle func
type SlowHandleFunc func(context.Context, *Account, *transport.Message) error
type AccountSlowHandler struct {
	F SlowHandleFunc
	M *transport.Message
}

// full account info
type Account struct {
	ID        int64   `bson:"_id" json:"_id"`
	UserId    int64   `bson:"user_id" json:"user_id"`
	GameId    int16   `bson:"game_id" json:"game_id"`
	Name      string  `bson:"name" json:"name"`
	Level     int32   `bson:"level" json:"level"`
	PlayerIDs []int64 `bson:"player_id" json:"player_id"`

	sock transport.Socket `bson:"-" json:"-"`
	p    *Player          `bson:"-" json:"-"`

	timeOut *time.Timer `bson:"-" json:"-"`

	SlowHandler chan *AccountSlowHandler `bson:"-" json:"-"`
}

func NewAccount() interface{} {
	account := &Account{
		ID:        -1,
		Name:      "",
		Level:     1,
		PlayerIDs: []int64{},
		sock:      nil,
		p:         nil,
		timeOut:   time.NewTimer(define.Account_OnlineTimeout),
	}

	return account
}

func (a *Account) GetID() int64 {
	return a.ID
}

func (a *Account) SetID(id int64) {
	a.ID = id
}

func (a *Account) GetName() string {
	return a.Name
}

func (a *Account) SetName(name string) {
	a.Name = name
}

func (a *Account) GetLevel() int32 {
	return a.Level
}

func (a *Account) SetLevel(level int32) {
	a.Level = level
}

func (a *Account) AddPlayerID(playerID int64) {
	for _, value := range a.PlayerIDs {
		if value == playerID {
			return
		}
	}

	a.PlayerIDs = append(a.PlayerIDs, playerID)
}

func (a *Account) GetPlayerIDs() []int64 {
	return a.PlayerIDs
}

func (a *Account) GetSock() transport.Socket {
	return a.sock
}

func (a *Account) SetSock(s transport.Socket) {
	a.sock = s
}

func (a *Account) GetPlayer() *Player {
	return a.p
}

func (a *Account) SetPlayer(p *Player) {
	a.p = p
}

func (a *Account) Close() {
	close(a.SlowHandler)
	a.timeOut.Stop()
	a.sock.Close()

	// Pool.Put
	if a.GetPlayer() != nil {
		a.GetPlayer().Destroy()
	}
}

func (a *Account) Run(ctx context.Context) error {
	for {
		select {
		// context canceled
		case <-ctx.Done():
			log.Info().
				Int64("account_id", a.GetID()).
				Str("socket_remote", a.sock.Remote()).
				Msg("account context done...")
			return nil

		case handler, ok := <-a.SlowHandler:
			if !ok {
				log.Info().
					Int64("account_id", a.GetID()).
					Msg("delay handler channel closed")
				return nil
			} else {
				err := handler.F(ctx, a, handler.M)
				if err != nil && !errors.Is(err, ErrCreateMoreThanOnePlayer) {
					log.Warn().
						Int64("account_id", a.ID).
						Err(err).
						Msg("Account.Run execute failed")
				}
			}

		// lost connection
		case <-a.timeOut.C:
			return fmt.Errorf("account<%d> time out", a.GetID())

		// account update
		default:
			now := time.Now()
			a.update()
			d := time.Since(now)
			time.Sleep(time.Millisecond*100 - d)
		}
	}
}

func (a *Account) update() {
	if a.p != nil {
		a.p.update()
	}
}

/*
msg Example:
	Name: S2C_AccountLogon
	Body: protoBuf byte
*/
func (a *Account) SendProtoMessage(p proto.Message) {
	if a.sock == nil {
		return
	}

	var msg transport.Message
	// msg.Type = transport.BodyProtobuf
	msg.Name = string(proto.MessageReflect(p).Descriptor().Name())
	msg.Body = p

	if err := a.sock.Send(&msg); err != nil {
		log.Warn().
			Int64("account_id", a.ID).
			Str("msg_name", msg.Name).
			Err(err).
			Msg("Account.SendProtoMessage failed")
		return
	}
}

func (a *Account) HeartBeat() {
	a.timeOut.Reset(define.Account_OnlineTimeout)

	reply := &pbGlobal.S2C_HeartBeat{Timestamp: uint32(time.Now().Unix())}
	a.SendProtoMessage(reply)
}
