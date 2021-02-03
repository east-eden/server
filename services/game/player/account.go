package player

import (
	"context"
	"errors"
	"fmt"
	"time"

	"bitbucket.org/east-eden/server/define"
	pbGlobal "bitbucket.org/east-eden/server/proto/global"
	"bitbucket.org/east-eden/server/store"
	"bitbucket.org/east-eden/server/transport"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

var (
	ErrAccountDisconnect       = errors.New("account disconnect")                                             // handleSocket got this error will disconnect account
	ErrCreateMoreThanOnePlayer = errors.New("AccountManager.CreatePlayer failed: only can create one player") // only can create one player
	Account_MemExpire          = time.Hour * 2
)

// account delay handle func
type DelayHandleFunc func(*Account) error

// lite account info
type LiteAccount struct {
	store.StoreObjector `bson:"-" json:"-"`
	ID                  int64   `bson:"_id" json:"_id"`
	UserId              int64   `bson:"user_id" json:"user_id"`
	GameId              int16   `bson:"game_id" json:"game_id"`
	Name                string  `bson:"name" json:"name"`
	Level               int32   `bson:"level" json:"level"`
	PlayerIDs           []int64 `bson:"player_id" json:"player_id"`
}

func (la *LiteAccount) GetObjID() int64 {
	return la.ID
}

func (la *LiteAccount) GetStoreIndex() int64 {
	return -1
}

func (la *LiteAccount) AfterLoad() error {
	return nil
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
	LiteAccount `bson:"inline" json:",inline"`

	sock transport.Socket `bson:"-" json:"-"`
	p    *Player          `bson:"-" json:"-"`

	timeOut *time.Timer `bson:"-" json:"-"`

	DelayHandler chan DelayHandleFunc `bson:"-" json:"-"`
}

func NewLiteAccount() interface{} {
	return &LiteAccount{
		ID:        -1,
		Name:      "",
		Level:     1,
		PlayerIDs: []int64{},
	}
}

func NewAccount() interface{} {
	account := &Account{
		LiteAccount: *(NewLiteAccount().(*LiteAccount)),
		sock:        nil,
		p:           nil,
		timeOut:     time.NewTimer(define.Account_OnlineTimeout),
	}

	return account
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
	close(a.DelayHandler)
	a.timeOut.Stop()
	a.sock.Close()
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

		case fn, ok := <-a.DelayHandler:
			if !ok {
				log.Info().
					Int64("account_id", a.GetID()).
					Msg("delay handler channel closed")
				return nil
			} else {
				err := fn(a)
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
		}
	}
}

/*
msg Example:
	Type: transport.BodyProtobuf
	Name: M2C_AccountLogon
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

	reply := &pbGlobal.M2C_HeartBeat{Timestamp: uint32(time.Now().Unix())}
	a.SendProtoMessage(reply)
}
