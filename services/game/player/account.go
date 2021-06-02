package player

import (
	"context"
	"errors"
	"time"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/iface"
	"bitbucket.org/funplus/server/transport"
	"bitbucket.org/funplus/server/utils"
	"github.com/hellodudu/task"
	log "github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

var (
	ErrAccountDisconnect       = errors.New("account disconnect") // handleSocket got this error will disconnect account
	ErrAccountKicked           = errors.New("account kickoff")
	ErrCreateMoreThanOnePlayer = errors.New("AccountManager.CreatePlayer failed: only can create one player") // only can create one player
	Account_MemExpire          = time.Hour * 2
	AccountTaskNum             = 100         // max account execute channel number
	AccountTaskTimeout         = time.Minute // 账号task超时
)

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

	tasker    *task.Tasker    `bson:"-" json:"-"`
	rpcCaller iface.RpcCaller `bson:"-" json:"-"`
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

}

func (a *Account) InitTask(fns ...task.StartFn) {
	fns = append(fns, a.onTaskStart)

	a.tasker = task.NewTasker(int32(AccountTaskNum))
	a.tasker.Init(
		task.WithContextDoneFn(func() {
			log.Info().
				Int64("account_id", a.GetId()).
				Str("socket_remote", a.sock.Remote()).
				Msg("account context done...")
		}),
		task.WithStartFns(fns...),
		task.WithUpdateFn(a.onTaskUpdate),
		task.WithTimeout(AccountTaskTimeout),
		task.WithSleep(time.Millisecond*100),
	)
}

func (a *Account) ResetTimeout() {
	a.tasker.ResetTimer()
}

func (a *Account) IsTaskRunning() bool {
	return a.tasker.IsRunning()
}

func (a *Account) GetId() int64 {
	return a.Id
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

func (a *Account) Stop() {
	a.tasker.Stop()
	a.sock.Close()

	// Pool.Put
	if a.GetPlayer() != nil {
		a.GetPlayer().Destroy()
	}
}

func (a *Account) AddTask(ctx context.Context, fn task.TaskHandler, m proto.Message) error {
	return a.tasker.Add(ctx, fn, a, m)
}

func (a *Account) TaskRun(ctx context.Context) error {
	return a.tasker.Run(ctx)
}

func (a *Account) onTaskStart() {
	if a.p != nil {
		a.p.onTaskStart()
	}
}

func (a *Account) onTaskUpdate() {
	if a.p != nil {
		a.p.onTaskUpdate()
	}
}

func (a *Account) LogonSucceed() {
	log.Info().Int64("account_id", a.Id).Msg("account logon success")

	// send logon success
	a.SendLogon()

	// sync time
	a.SendServerTime()

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
	msg.Name = string(p.ProtoReflect().Descriptor().Name())
	msg.Body = p

	err := a.sock.Send(&msg)
	_ = utils.ErrCheck(err, "Account.SendProtoMessage failed", a.Id, msg.Name)

	log.Info().Str("msg_name", msg.Name).Interface("msg_body", msg.Body).Msg("send message")
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
		reply.PlayerId = p.GetId()
		reply.PlayerName = p.GetName()
		reply.PlayerLevel = p.GetLevel()
	}

	a.SendProtoMessage(reply)
}

func (a *Account) SendServerTime() {
	reply := &pbGlobal.S2C_ServerTime{Timestamp: uint32(time.Now().Unix())}
	a.SendProtoMessage(reply)
}

func (a *Account) HeartBeat() {
	a.ResetTimeout()
}
