package player

import (
	"context"
	"errors"
	"fmt"
	"time"

	"bitbucket.org/funplus/server/define"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/iface"
	"bitbucket.org/funplus/server/transport"
	"bitbucket.org/funplus/server/utils"
	"github.com/golang/protobuf/proto"
	log "github.com/rs/zerolog/log"
)

var (
	ErrAccountDisconnect       = errors.New("account disconnect") // handleSocket got this error will disconnect account
	ErrAccountKicked           = errors.New("account kickoff")
	ErrCreateMoreThanOnePlayer = errors.New("AccountManager.CreatePlayer failed: only can create one player") // only can create one player
	Account_MemExpire          = time.Hour * 2
	AccountTaskNum             = 100              // max account execute channel number
	AccountTaskTimeout         = 10 * time.Second // 账号task超时
)

// account delay handle func
type TaskHandler func(context.Context, *Account, *transport.Message) error
type AccountTasker struct {
	C context.Context
	E chan<- error
	F TaskHandler
	M *transport.Message
}

// full account info
type Account struct {
	Id             int64   `bson:"_id" json:"_id"`
	UserId         int64   `bson:"user_id" json:"user_id"`
	GameId         int16   `bson:"game_id" json:"game_id"` // 上次登陆的game节点
	Name           string  `bson:"name" json:"name"`
	Level          int32   `bson:"level" json:"level"`
	Privilege      int8    `bson:"privilege" json:"privilege"` // gm 权限
	PlayerIDs      []int64 `bson:"player_id" json:"player_id"`
	LastLogoffTime int32   `bson:"last_logoff_time" json:"last_logoff_time"` // 账号上次下线时间

	sock transport.Socket `bson:"-" json:"-"`
	p    *Player          `bson:"-" json:"-"`

	timeOut *time.Timer `bson:"-" json:"-"`

	TaskHandlers chan *AccountTasker `bson:"-" json:"-"`
	rpcCaller    iface.RpcCaller     `bson:"-" json:"-"`
}

func NewAccount() interface{} {
	return &Account{}
}

func (a *Account) Init() {
	a.Id = -1
	a.UserId = -1
	a.GameId = -1
	a.Name = ""
	a.Level = 1
	a.Privilege = 3
	a.PlayerIDs = make([]int64, 0, 5)
	a.sock = nil
	a.p = nil
	a.ResetTimeout()
	a.TaskHandlers = make(chan *AccountTasker, AccountTaskNum)
}

func (a *Account) ResetTimeout() {
	a.timeOut = time.NewTimer(define.Account_OnlineTimeout)
}

func (a *Account) GetID() int64 {
	return a.Id
}

func (a *Account) SetID(id int64) {
	a.Id = id
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

func (a *Account) SetRpcCaller(c iface.RpcCaller) {
	a.rpcCaller = c
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
	close(a.TaskHandlers)
	a.timeOut.Stop()
	a.sock.Close()

	// Pool.Put
	if a.GetPlayer() != nil {
		a.GetPlayer().Destroy()
	}
}

func (a *Account) AddTask(task *AccountTasker) error {
	subCtx, cancel := utils.WithTimeoutContext(task.C, AccountTaskTimeout)
	defer cancel()

	e := make(chan error, 1)
	task.E = e
	a.TaskHandlers <- task

	select {
	case err := <-e:
		return err
	case <-subCtx.Done():
		return subCtx.Err()
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

		case handler, ok := <-a.TaskHandlers:
			if !ok {
				log.Info().
					Int64("account_id", a.GetID()).
					Msg("delay handler channel closed")
				return nil
			} else {
				err := handler.F(ctx, a, handler.M)
				handler.E <- err
				if err == nil {
					continue
				}

				// 被踢下线
				if errors.Is(err, ErrAccountKicked) {
					return ErrAccountKicked
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

func (a *Account) LogonSucceed() {
	log.Info().Int64("account_id", a.Id).Msg("account logon success")

	// send logon success
	a.SendLogon()

	p := a.GetPlayer()
	if p == nil {
		return
	}

	// sync to client
	p.SendInitInfo()

	// 时间跨度检查
	p.CheckTimeChange()
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
			Int64("account_id", a.Id).
			Str("msg_name", msg.Name).
			Err(err).
			Msg("Account.SendProtoMessage failed")
		return
	}
}

func (a *Account) SendLogon() {
	reply := &pbGlobal.S2C_AccountLogon{
		UserId:      a.UserId,
		AccountId:   a.Id,
		PlayerId:    -1,
		PlayerName:  "",
		PlayerLevel: 0,
	}

	if p := a.GetPlayer(); p != nil {
		reply.PlayerId = p.GetID()
		reply.PlayerName = p.GetName()
		reply.PlayerLevel = p.GetLevel()
	}

	a.SendProtoMessage(reply)
}

func (a *Account) HeartBeat() {
	a.ResetTimeout()

	reply := &pbGlobal.S2C_HeartBeat{Timestamp: uint32(time.Now().Unix())}
	a.SendProtoMessage(reply)
}
