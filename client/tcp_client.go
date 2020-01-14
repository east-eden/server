package client

import (
	"context"
	"fmt"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/internal/global"
	"github.com/yokaiio/yokai_server/internal/transport"
	"github.com/yokaiio/yokai_server/internal/utils"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

type TcpClient struct {
	tr        transport.Transport
	ts        transport.Socket
	ctx       context.Context
	cancel    context.CancelFunc
	waitGroup utils.WaitGroupWrapper

	heartBeatTimer    *time.Timer
	heartBeatDuration time.Duration
	tcpServerAddr     string
	gateEndpoints     []string

	accountID   int64
	accountName string
	reconn      chan int
	connected   bool
	recvCh      chan int

	disconnectCtx    context.Context
	disconnectCancel context.CancelFunc
}

type MC_AccountTest struct {
	AccountId int64  `protobuf:"varint,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func NewTcpClient(ctx *cli.Context) *TcpClient {
	t := &TcpClient{
		tr:                transport.NewTransport(transport.Timeout(transport.DefaultDialTimeout)),
		heartBeatDuration: ctx.Duration("heart_beat"),
		heartBeatTimer:    time.NewTimer(ctx.Duration("heart_beat")),
		gateEndpoints:     ctx.StringSlice("gate_endpoints"),
		reconn:            make(chan int, 1),
		recvCh:            make(chan int, 100),
		connected:         false,
	}

	t.ctx, t.cancel = context.WithCancel(ctx)
	t.heartBeatTimer.Stop()

	t.registerMessage()

	t.reconn <- 1

	return t
}

func (t *TcpClient) registerMessage() {

	transport.DefaultRegister.RegisterProtobufMessage(&pbAccount.MS_AccountLogon{}, t.OnMS_AccountLogon)
	transport.DefaultRegister.RegisterProtobufMessage(&pbAccount.MS_HeartBeat{}, t.OnMS_HeartBeat)

	transport.DefaultRegister.RegisterProtobufMessage(&pbGame.MS_CreatePlayer{}, t.OnMS_CreatePlayer)
	transport.DefaultRegister.RegisterProtobufMessage(&pbGame.MS_SelectPlayer{}, t.OnMS_SelectPlayer)
	transport.DefaultRegister.RegisterProtobufMessage(&pbGame.MS_QueryPlayerInfo{}, t.OnMS_QueryPlayerInfo)
	transport.DefaultRegister.RegisterProtobufMessage(&pbGame.MS_QueryPlayerInfos{}, t.OnMS_QueryPlayerInfos)

	transport.DefaultRegister.RegisterProtobufMessage(&pbGame.MS_HeroList{}, t.OnMS_HeroList)
	transport.DefaultRegister.RegisterProtobufMessage(&pbGame.MS_HeroInfo{}, t.OnMS_HeroInfo)

	transport.DefaultRegister.RegisterProtobufMessage(&pbGame.MS_ItemList{}, t.OnMS_ItemList)
	transport.DefaultRegister.RegisterProtobufMessage(&pbGame.MS_HeroEquips{}, t.OnMS_HeroEquips)

	transport.DefaultRegister.RegisterProtobufMessage(&pbGame.MS_TokenList{}, t.OnMS_TokenList)

	transport.DefaultRegister.RegisterProtobufMessage(&pbGame.MS_TalentList{}, t.OnMS_TalentList)
}

func (t *TcpClient) Connect(accountID int64, accountName string) {
	if t.connected {
		t.Disconnect()
	}

	t.accountID = accountID
	t.accountName = accountName
	t.disconnectCtx, t.disconnectCancel = context.WithCancel(t.ctx)
	t.waitGroup.Wrap(func() {
		t.doConnect()
	})

	t.waitGroup.Wrap(func() {
		t.doRecv()
	})
}

func (t *TcpClient) Disconnect() {
	t.ts.Close()
	t.disconnectCancel()
	t.waitGroup.Wait()
}

func (t *TcpClient) SendMessage(msg *transport.Message) {
	if msg == nil {
		return
	}

	if t.ts == nil {
		logger.Warn("未连接到服务器")
		return
	}

	if err := t.ts.Send(msg); err != nil {
		logger.Warn("Unexpected send err", err)
		t.reconn <- 1
	}
}

func (t *TcpClient) SetTcpAddress(addr string) {
	t.tcpServerAddr = addr
}

func (t *TcpClient) OnMS_AccountLogon(sock transport.Socket, msg *transport.Message) {
	logger.Info("连接到服务器")

	t.connected = true
	t.heartBeatTimer.Reset(t.heartBeatDuration)

	send := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_account.MC_AccountConnected",
		Body: &pbAccount.MC_AccountConnected{AccountId: 1, Name: "dudu"},
	}
	t.SendMessage(send)

	sendTest := &transport.Message{
		Type: transport.BodyJson,
		Name: "MC_AccountTest",
		Body: &MC_AccountTest{AccountId: 1, Name: "test"},
	}
	t.SendMessage(sendTest)
}

func (t *TcpClient) OnMS_HeartBeat(sock transport.Socket, msg *transport.Message) {
}

func (t *TcpClient) OnMS_CreatePlayer(sock transport.Socket, msg *transport.Message) {
	m := msg.Body.(*pbGame.MS_CreatePlayer)
	if m.ErrorCode == 0 {
		logger.WithFields(logger.Fields{
			"角色id":     m.Info.LiteInfo.Id,
			"角色名字":     m.Info.LiteInfo.Name,
			"角色经验":     m.Info.LiteInfo.Exp,
			"角色等级":     m.Info.LiteInfo.Level,
			"角色拥有英雄数量": m.Info.HeroNums,
			"角色拥有物品数量": m.Info.ItemNums,
		}).Info("角色创建成功：")
	} else {
		logger.Info("角色创建失败，error_code=", m.ErrorCode)
	}
}

func (t *TcpClient) OnMS_SelectPlayer(sock transport.Socket, msg *transport.Message) {
	m := msg.Body.(*pbGame.MS_SelectPlayer)
	if m.ErrorCode == 0 {
		logger.WithFields(logger.Fields{
			"角色id":     m.Info.LiteInfo.Id,
			"角色名字":     m.Info.LiteInfo.Name,
			"角色经验":     m.Info.LiteInfo.Exp,
			"角色等级":     m.Info.LiteInfo.Level,
			"角色拥有英雄数量": m.Info.HeroNums,
			"角色拥有物品数量": m.Info.ItemNums,
		}).Info("使用此角色：")
	} else {
		logger.Info("选择角色失败，error_code=", m.ErrorCode)
	}
}

func (t *TcpClient) OnMS_QueryPlayerInfo(sock transport.Socket, msg *transport.Message) {
	m := msg.Body.(*pbGame.MS_QueryPlayerInfo)
	if m.Info == nil {
		logger.Info("该账号下还没有角色，请先创建一个角色")
		return
	}

	logger.WithFields(logger.Fields{
		"角色id":     m.Info.LiteInfo.Id,
		"角色名字":     m.Info.LiteInfo.Name,
		"角色经验":     m.Info.LiteInfo.Exp,
		"角色等级":     m.Info.LiteInfo.Level,
		"角色拥有英雄数量": m.Info.HeroNums,
		"角色拥有物品数量": m.Info.ItemNums,
	}).Info("角色信息：")
}

func (t *TcpClient) OnMS_QueryPlayerInfos(sock transport.Socket, msg *transport.Message) {
	m := msg.Body.(*pbGame.MS_QueryPlayerInfos)
	if len(m.Infos) == 0 {
		logger.Info("该账号下还没有角色，请先创建一个角色")
		return
	}

	logger.Info("所有角色信息：")
	for k, v := range m.Infos {
		fields := logger.Fields{}
		fields["id"] = v.LiteInfo.Id
		fields["名字"] = v.LiteInfo.Name
		fields["经验"] = v.LiteInfo.Exp
		fields["等级"] = v.LiteInfo.Level
		fields["拥有英雄数量"] = v.HeroNums
		fields["拥有物品数量"] = v.ItemNums
		logger.WithFields(fields).Info(fmt.Sprintf("角色%d", k+1))
	}

}

func (t *TcpClient) OnMS_HeroList(sock transport.Socket, msg *transport.Message) {
	m := msg.Body.(*pbGame.MS_HeroList)
	fields := logger.Fields{}

	logger.Info("拥有英雄：")
	for k, v := range m.Heros {
		fields["id"] = v.Id
		fields["TypeID"] = v.TypeId
		fields["经验"] = v.Exp
		fields["等级"] = v.Level

		entry := global.GetHeroEntry(v.TypeId)
		if entry != nil {
			fields["名字"] = entry.Name
		}

		logger.WithFields(fields).Info(fmt.Sprintf("英雄%d", k+1))
	}

}

func (t *TcpClient) OnMS_HeroInfo(sock transport.Socket, msg *transport.Message) {
	m := msg.Body.(*pbGame.MS_HeroInfo)

	entry := global.GetHeroEntry(m.Info.TypeId)
	logger.WithFields(logger.Fields{
		"id":     m.Info.Id,
		"TypeID": m.Info.TypeId,
		"经验":     m.Info.Exp,
		"等级":     m.Info.Level,
		"名字":     entry.Name,
	}).Info("英雄信息：")
}

func (t *TcpClient) OnMS_ItemList(sock transport.Socket, msg *transport.Message) {
	m := msg.Body.(*pbGame.MS_ItemList)
	fields := logger.Fields{}

	logger.Info("拥有物品：")
	for k, v := range m.Items {
		fields["id"] = v.Id
		fields["type_id"] = v.TypeId

		entry := global.GetItemEntry(v.TypeId)
		if entry != nil {
			fields["name"] = entry.Name
		}
		logger.WithFields(fields).Info(fmt.Sprintf("物品%d", k+1))
	}

}

func (t *TcpClient) OnMS_HeroEquips(sock transport.Socket, msg *transport.Message) {
	m := msg.Body.(*pbGame.MS_HeroEquips)
	fields := logger.Fields{}

	logger.Info("此英雄穿有装备：")
	for k, v := range m.Equips {
		fields["id"] = v.Id
		fields["type_id"] = v.TypeId

		entry := global.GetItemEntry(v.TypeId)
		if entry != nil {
			fields["name"] = entry.Name
		}
		logger.WithFields(fields).Info(fmt.Sprintf("装备%d", k+1))
	}

}

func (t *TcpClient) OnMS_TokenList(sock transport.Socket, msg *transport.Message) {
	m := msg.Body.(*pbGame.MS_TokenList)
	fields := logger.Fields{}

	logger.Info("拥有代币：")
	for k, v := range m.Tokens {
		fields["type"] = v.Type
		fields["value"] = v.Value
		fields["max_hold"] = v.MaxHold

		entry := global.GetTokenEntry(v.Type)
		if entry != nil {
			fields["name"] = entry.Name
		}
		logger.WithFields(fields).Info(fmt.Sprintf("代币%d", k+1))
	}

}

func (t *TcpClient) OnMS_TalentList(sock transport.Socket, msg *transport.Message) {
	m := msg.Body.(*pbGame.MS_TalentList)
	fields := logger.Fields{}

	logger.Info("已点击天赋：")
	for k, v := range m.Talents {
		fields["id"] = v.Id

		entry := global.GetTalentEntry(v.Id)
		if entry != nil {
			fields["名字"] = entry.Name
			fields["描述"] = entry.Desc
		}

		logger.WithFields(fields).Info(fmt.Sprintf("天赋%d", k+1))
	}

}

func (t *TcpClient) doConnect() {
	for {
		select {
		case <-t.ctx.Done():
			logger.Info("tcp client dial goroutine done...")
			return

		case <-t.disconnectCtx.Done():
			logger.Info("connect goroutine context down...")
			return

		case <-t.heartBeatTimer.C:
			t.heartBeatTimer.Reset(t.heartBeatDuration)

			msg := &transport.Message{
				Type: transport.BodyJson,
				Name: "yokai_account.MC_HeartBeat",
				Body: &pbAccount.MC_HeartBeat{},
			}
			t.SendMessage(msg)

		case <-t.reconn:
			// close old connection
			if t.ts != nil {
				t.ts.Close()
			}

			var err error
			if t.ts, err = t.tr.Dial(t.tcpServerAddr); err != nil {
				logger.Warn("unexpected dial err:", err)
				time.Sleep(time.Second * 3)
				t.reconn <- 1
				continue
			}

			logger.Info("tpc dial at remote:", t.ts.Remote())

			msg := &transport.Message{
				Type: transport.BodyProtobuf,
				Name: "yokai_account.MC_AccountLogon",
				Body: &pbAccount.MC_AccountLogon{
					AccountId:   t.accountID,
					AccountName: t.accountName,
				},
			}
			t.SendMessage(msg)

		}
	}
}

func (t *TcpClient) doRecv() {
	for {
		select {
		case <-t.ctx.Done():
			logger.Info("tcp client recv goroutine done...")
			return

		case <-t.disconnectCtx.Done():
			logger.Info("recv goroutine context down...")
			return
		default:

			func() {
				// be called per 100ms
				ct := time.Now()
				defer func() {
					d := time.Since(ct)
					time.Sleep(100*time.Millisecond - d)
				}()

				if t.ts != nil {
					if msg, h, err := t.ts.Recv(); err != nil {
						logger.Warn("Unexpected recv err:", err)
						t.reconn <- 1
					} else {
						h.Fn(t.ts, msg)
						if msg.Name != "yokai_account.MS_HeartBeat" {
							t.recvCh <- 1
						}
					}
				}
			}()
		}
	}
}

func (t *TcpClient) Run() error {
	for {
		select {
		case <-t.ctx.Done():
			logger.Info("tcp client context done...")
			return nil
		}
	}

	return nil
}

func (t *TcpClient) Exit() {
	t.cancel()
	t.heartBeatTimer.Stop()

	if t.ts != nil {
		t.ts.Close()
	}

	t.waitGroup.Wait()
}
