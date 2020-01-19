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
	register  transport.Register
	ctx       context.Context
	cancel    context.CancelFunc
	waitGroup utils.WaitGroupWrapper

	heartBeatTimer    *time.Timer
	heartBeatDuration time.Duration
	tcpServerAddr     string
	gateEndpoints     []string

	userID    int64
	userName  string
	accountID int64
	reconn    chan int
	connected bool
	recvCh    chan int

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
		register:          transport.NewTransportRegister(),
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

	return t
}

func (t *TcpClient) registerMessage() {

	t.register.RegisterProtobufMessage(&pbAccount.MS_AccountLogon{}, t.OnMS_AccountLogon)
	t.register.RegisterProtobufMessage(&pbAccount.MS_HeartBeat{}, t.OnMS_HeartBeat)

	t.register.RegisterProtobufMessage(&pbGame.MS_CreatePlayer{}, t.OnMS_CreatePlayer)
	t.register.RegisterProtobufMessage(&pbGame.MS_SelectPlayer{}, t.OnMS_SelectPlayer)
	t.register.RegisterProtobufMessage(&pbGame.MS_QueryPlayerInfo{}, t.OnMS_QueryPlayerInfo)
	t.register.RegisterProtobufMessage(&pbGame.MS_QueryPlayerInfos{}, t.OnMS_QueryPlayerInfos)

	t.register.RegisterProtobufMessage(&pbGame.MS_HeroList{}, t.OnMS_HeroList)
	t.register.RegisterProtobufMessage(&pbGame.MS_HeroInfo{}, t.OnMS_HeroInfo)

	t.register.RegisterProtobufMessage(&pbGame.MS_ItemList{}, t.OnMS_ItemList)
	t.register.RegisterProtobufMessage(&pbGame.MS_HeroEquips{}, t.OnMS_HeroEquips)

	t.register.RegisterProtobufMessage(&pbGame.MS_TokenList{}, t.OnMS_TokenList)

	t.register.RegisterProtobufMessage(&pbGame.MS_TalentList{}, t.OnMS_TalentList)
}

func (t *TcpClient) Connect() error {
	if t.connected {
		t.Disconnect()
	}

	t.disconnectCtx, t.disconnectCancel = context.WithCancel(t.ctx)
	t.waitGroup.Wrap(func() {
		t.reconn <- 1
		t.doConnect()
	})

	t.waitGroup.Wrap(func() {
		t.doRecv()
	})

	timer := time.NewTimer(time.Second * 5)
	select {
	case <-timer.C:
		return fmt.Errorf("connect timeout")
	default:
		if t.connected {
			return nil
		}

		time.Sleep(time.Millisecond * 200)
	}

	return fmt.Errorf("unexpected error")
}

func (t *TcpClient) Disconnect() {
	logger.WithFields(logger.Fields{
		"local": t.ts.Local(),
	}).Info("socket local disconnect")

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
	}
}

func (t *TcpClient) SetTcpAddress(addr string) {
	t.tcpServerAddr = addr
}

func (t *TcpClient) SetUserInfo(userID int64, accountID int64, userName string) {
	t.userID = userID
	t.userName = userName
	t.accountID = accountID
}

func (t *TcpClient) OnMS_AccountLogon(sock transport.Socket, msg *transport.Message) {
	logger.WithFields(logger.Fields{
		"local": sock.Local(),
	}).Info("连接到服务器")

	t.connected = true
	t.heartBeatTimer.Reset(t.heartBeatDuration)

	send := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_account.MC_AccountConnected",
		Body: &pbAccount.MC_AccountConnected{AccountId: t.userID, Name: t.userName},
	}
	t.SendMessage(send)

	sendTest := &transport.Message{
		Type: transport.BodyJson,
		Name: "MC_AccountTest",
		Body: &MC_AccountTest{AccountId: t.userID, Name: t.userName},
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
				logger.WithFields(logger.Fields{
					"local": t.ts.Local(),
				}).Info("socket local reconnect")
				t.ts.Close()
			}

			var err error
			if t.ts, err = t.tr.Dial(t.tcpServerAddr); err != nil {
				logger.Warn("unexpected dial err:", err)
				continue
			}

			logger.WithFields(logger.Fields{
				"local":  t.ts.Local(),
				"remote": t.ts.Remote(),
				"conn":   t.ts.Conn(),
			}).Info("tcp dial success")

			msg := &transport.Message{
				Type: transport.BodyProtobuf,
				Name: "yokai_account.MC_AccountLogon",
				Body: &pbAccount.MC_AccountLogon{
					UserId:    t.userID,
					AccountId: t.accountID,
				},
			}
			t.SendMessage(msg)

			logger.WithFields(logger.Fields{
				"user_id":    t.userID,
				"account_id": t.accountID,
				"local":      t.ts.Local(),
			}).Info("connect send message")
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
					if msg, h, err := t.ts.Recv(t.register); err != nil {
						logger.Warn("Unexpected recv err:", err)
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
		logger.WithFields(logger.Fields{
			"local": t.ts.Local(),
		}).Info("socket exit close")

		t.ts.Close()
	}

	t.waitGroup.Wait()
}
