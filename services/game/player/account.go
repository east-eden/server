package player

import (
	"context"
	"errors"
	"time"

	"e.coding.net/mmstudio/blade/server/define"
	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
	"e.coding.net/mmstudio/blade/server/services/game/iface"
	"e.coding.net/mmstudio/blade/server/store"
	"e.coding.net/mmstudio/blade/server/transport"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/hellodudu/task"
	log "github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

var (
	ErrAccountDisconnect       = errors.New("account disconnect") // handleSocket got this error will disconnect account
	ErrAccountKicked           = errors.New("account kickoff")
	ErrCreateMoreThanOnePlayer = errors.New("AccountManager.CreatePlayer failed: only can create one player") // only can create one player
	Account_MemExpire          = time.Hour * 2
	AccountTaskNum             = 128         // max account execute channel number
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
	a.SetSock(nil)
	a.SetPlayer(nil)
}

func (a *Account) InitTask(startFn task.StartFn, stopFn task.StopFn) {
	startFns := make([]task.StartFn, 0, 8)
	startFns = append(startFns, startFn, a.onTaskStart)

	stopFns := make([]task.StopFn, 0, 8)
	stopFns = append(stopFns, stopFn, a.onTaskStop)

	a.tasker = task.NewTasker(int32(AccountTaskNum))
	a.tasker.Init(
		task.WithStartFns(startFns...),
		task.WithStopFns(stopFns...),
		task.WithUpdateFn(a.onTaskUpdate),
		task.WithTimeout(AccountTaskTimeout),
		task.WithSleep(time.Millisecond*100),
	)
}

func (a *Account) ResetTimeout() {
	a.tasker.ResetTimer()
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

func (a *Account) TaskStop() {
	a.tasker.Stop()
}

func (a *Account) AddWaitTask(ctx context.Context, fn task.TaskHandler, p ...interface{}) error {
	param := make([]interface{}, 0, len(p)+1)
	param = append(param, a)
	param = append(param, p...)
	return a.tasker.AddWait(ctx, fn, param...)
}

func (a *Account) AddTask(ctx context.Context, fn task.TaskHandler, p ...interface{}) {
	param := make([]interface{}, 0, len(p)+1)
	param = append(param, a)
	param = append(param, p...)
	a.tasker.Add(ctx, fn, param...)
}

func (a *Account) TaskRun(ctx context.Context) error {
	return a.tasker.Run(ctx)
}

func (a *Account) IsTaskRunning() bool {
	return a.tasker.IsRunning()
}

func (a *Account) onTaskStart() {
	if a.p != nil {
		a.p.onTaskStart()
	}
}

func (a *Account) onTaskStop() {
	log.Info().
		Caller().
		Int64("account_id", a.GetId()).
		// Str("socket_remote", a.GetSock().Remote()).
		Msg("account task stop...")

	// 记录下线时间
	a.saveLogoffTime()

	// 关闭socket
	a.GetSock().Close()
	a.SetSock(nil)
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

func (a *Account) HeartBeat() {
	a.ResetTimeout()
}

func (a *Account) SaveAccount() {
	// save account
	err := store.GetStore().UpdateOne(context.Background(), define.StoreType_Account, a.Id, a, true)
	utils.ErrPrint(err, "UpdateOne failed when Account.SaveAccount", a.Id, a.UserId)
}

// 记录下线时间
func (a *Account) saveLogoffTime() {
	a.LastLogoffTime = int32(time.Now().Unix())
	fields := map[string]interface{}{
		"last_logoff_time": a.LastLogoffTime,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Account, a.Id, fields)
	utils.ErrPrint(err, "account save last_logoff_time failed", a.Id, a.LastLogoffTime)
}

// 记录当前节点
func (a *Account) SaveGameNode(nodeId int16) {
	a.GameId = nodeId
	fields := map[string]interface{}{
		"game_id": a.GameId,
	}

	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Account, a.Id, fields, true)
	_ = utils.ErrCheck(err, "UpdateFields failed when Account.saveGameNode", a.Id, a.GameId)
}

/*
msg Example:
	Name: S2C_AccountLogon
	Body: protoBuf byte
*/
func (a *Account) SendProtoMessage(p proto.Message) {
	if a.GetSock() == nil {
		return
	}

	var msg transport.Message
	msg.Name = string(p.ProtoReflect().Descriptor().Name())
	msg.Body = p

	err := a.GetSock().Send(&msg)
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
