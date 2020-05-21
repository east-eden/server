package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/entries"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"github.com/yokaiio/yokai_server/transport"
	"github.com/yokaiio/yokai_server/utils"
)

type TransportClient struct {
	c         *Client
	tr        transport.Transport
	ts        transport.Socket
	r         transport.Register
	ctx       context.Context
	cancel    context.CancelFunc
	waitGroup utils.WaitGroupWrapper

	heartBeatTimer    *time.Timer
	heartBeatDuration time.Duration
	serverAddr        string
	gateEndpoints     []string
	certFile          string
	keyFile           string

	userID    int64
	userName  string
	accountID int64
	reconn    chan int
	connected bool
	recvCh    chan int
}

func NewTransportClient(c *Client, ctx *cli.Context) *TransportClient {

	t := &TransportClient{
		r:                 transport.NewTransportRegister(),
		heartBeatDuration: ctx.Duration("heart_beat"),
		heartBeatTimer:    time.NewTimer(ctx.Duration("heart_beat")),
		gateEndpoints:     ctx.StringSlice("gate_endpoints"),
		reconn:            make(chan int, 1),
		recvCh:            make(chan int, 100),
		connected:         false,
	}

	if ctx.Bool("debug") {
		t.certFile = ctx.String("cert_path_debug")
		t.keyFile = ctx.String("key_path_debug")
	} else {
		t.certFile = ctx.String("cert_path_release")
		t.keyFile = ctx.String("key_path_release")
	}

	t.ctx, t.cancel = context.WithCancel(ctx)
	t.heartBeatTimer.Stop()

	t.registerMessage()

	return t
}

func (t *TransportClient) registerMessage() {

	t.r.RegisterProtobufMessage(&pbAccount.M2C_AccountLogon{}, t.OnM2C_AccountLogon)
	t.r.RegisterProtobufMessage(&pbAccount.M2C_HeartBeat{}, t.OnM2C_HeartBeat)

	t.r.RegisterProtobufMessage(&pbGame.M2C_CreatePlayer{}, t.OnM2C_CreatePlayer)
	t.r.RegisterProtobufMessage(&pbGame.MS_SelectPlayer{}, t.OnMS_SelectPlayer)
	t.r.RegisterProtobufMessage(&pbGame.M2C_QueryPlayerInfo{}, t.OnM2C_QueryPlayerInfo)

	t.r.RegisterProtobufMessage(&pbGame.M2C_HeroList{}, t.OnM2C_HeroList)
	t.r.RegisterProtobufMessage(&pbGame.M2C_HeroInfo{}, t.OnM2C_HeroInfo)

	t.r.RegisterProtobufMessage(&pbGame.M2C_ItemList{}, t.OnM2C_ItemList)

	t.r.RegisterProtobufMessage(&pbGame.M2C_TokenList{}, t.OnM2C_TokenList)

	t.r.RegisterProtobufMessage(&pbGame.MS_TalentList{}, t.OnMS_TalentList)

	t.r.RegisterProtobufMessage(&pbGame.M2C_StartStageCombat{}, t.OnM2C_StartStageCombat)
}

func (t *TransportClient) SetTransportProtocol(protocol string) {
	tlsConf := &tls.Config{InsecureSkipVerify: true}
	cert, err := tls.LoadX509KeyPair(t.certFile, t.keyFile)
	if err != nil {
		logger.Fatal("load certificates failed:", err)
	}

	tlsConf.Certificates = []tls.Certificate{cert}

	t.tr = transport.NewTransport(
		protocol,
		transport.Timeout(transport.DefaultDialTimeout),
		//transport.TLSConfig(tlsConf),
	)
}

func (t *TransportClient) Connect() error {
	if t.connected {
		t.Disconnect()
	}

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

func (t *TransportClient) Disconnect() {
	logger.WithFields(logger.Fields{
		"local": t.ts.Local(),
	}).Info("socket local disconnect")

	t.cancel()
}

func (t *TransportClient) SendMessage(msg *transport.Message) {
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

func (t *TransportClient) SetServerAddress(addr string) {
	t.serverAddr = addr
}

func (t *TransportClient) SetUserInfo(userID int64, accountID int64, userName string) {
	t.userID = userID
	t.userName = userName
	t.accountID = accountID
}

func (t *TransportClient) OnM2C_AccountLogon(sock transport.Socket, msg *transport.Message) {
	m := msg.Body.(*pbAccount.M2C_AccountLogon)

	logger.WithFields(logger.Fields{
		"local":        sock.Local(),
		"user_id":      m.UserId,
		"account_id":   m.AccountId,
		"player_id":    m.PlayerId,
		"player_name":  m.PlayerName,
		"player_level": m.PlayerLevel,
	}).Info("帐号登录成功")

	t.connected = true
	t.heartBeatTimer.Reset(t.heartBeatDuration)

	send := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_account.MC_AccountConnected",
		Body: &pbAccount.MC_AccountConnected{AccountId: m.AccountId, Name: m.PlayerName},
	}
	t.SendMessage(send)
}

func (t *TransportClient) OnM2C_HeartBeat(sock transport.Socket, msg *transport.Message) {
}

func (t *TransportClient) OnM2C_CreatePlayer(sock transport.Socket, msg *transport.Message) {
	m := msg.Body.(*pbGame.M2C_CreatePlayer)
	if m.Error == 0 {
		logger.WithFields(logger.Fields{
			"角色id":     m.Info.LiteInfo.Id,
			"角色名字":     m.Info.LiteInfo.Name,
			"角色经验":     m.Info.LiteInfo.Exp,
			"角色等级":     m.Info.LiteInfo.Level,
			"角色拥有英雄数量": m.Info.HeroNums,
			"角色拥有物品数量": m.Info.ItemNums,
		}).Info("角色创建成功：")
	} else {
		logger.Info("角色创建失败，error_code=", m.Error)
	}
}

func (t *TransportClient) OnMS_SelectPlayer(sock transport.Socket, msg *transport.Message) {
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

func (t *TransportClient) OnM2C_QueryPlayerInfo(sock transport.Socket, msg *transport.Message) {
	m := msg.Body.(*pbGame.M2C_QueryPlayerInfo)
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

func (t *TransportClient) OnM2C_HeroList(sock transport.Socket, msg *transport.Message) {
	m := msg.Body.(*pbGame.M2C_HeroList)
	fields := logger.Fields{}

	logger.Info("拥有英雄：")
	for k, v := range m.Heros {
		fields["id"] = v.Id
		fields["TypeID"] = v.TypeId
		fields["经验"] = v.Exp
		fields["等级"] = v.Level

		entry := entries.GetHeroEntry(v.TypeId)
		if entry != nil {
			fields["名字"] = entry.Name
		}

		logger.WithFields(fields).Info(fmt.Sprintf("英雄%d", k+1))
	}

}

func (t *TransportClient) OnM2C_HeroInfo(sock transport.Socket, msg *transport.Message) {
	m := msg.Body.(*pbGame.M2C_HeroInfo)

	entry := entries.GetHeroEntry(m.Info.TypeId)
	logger.WithFields(logger.Fields{
		"id":     m.Info.Id,
		"TypeID": m.Info.TypeId,
		"经验":     m.Info.Exp,
		"等级":     m.Info.Level,
		"名字":     entry.Name,
	}).Info("英雄信息：")
}

func (t *TransportClient) OnM2C_ItemList(sock transport.Socket, msg *transport.Message) {
	m := msg.Body.(*pbGame.M2C_ItemList)
	fields := logger.Fields{}

	logger.Info("拥有物品：")
	for k, v := range m.Items {
		fields["id"] = v.Id
		fields["type_id"] = v.TypeId

		entry := entries.GetItemEntry(v.TypeId)
		if entry != nil {
			fields["name"] = entry.Name
		}
		logger.WithFields(fields).Info(fmt.Sprintf("物品%d", k+1))
	}

}

func (t *TransportClient) OnM2C_TokenList(sock transport.Socket, msg *transport.Message) {
	m := msg.Body.(*pbGame.M2C_TokenList)
	fields := logger.Fields{}

	logger.Info("拥有代币：")
	for k, v := range m.Tokens {
		fields["type"] = v.Type
		fields["value"] = v.Value
		fields["max_hold"] = v.MaxHold

		entry := entries.GetTokenEntry(v.Type)
		if entry != nil {
			fields["name"] = entry.Name
		}
		logger.WithFields(fields).Info(fmt.Sprintf("代币%d", k+1))
	}

}

func (t *TransportClient) OnMS_TalentList(sock transport.Socket, msg *transport.Message) {
	m := msg.Body.(*pbGame.MS_TalentList)
	fields := logger.Fields{}

	logger.Info("已点击天赋：")
	for k, v := range m.Talents {
		fields["id"] = v.Id

		entry := entries.GetTalentEntry(v.Id)
		if entry != nil {
			fields["名字"] = entry.Name
			fields["描述"] = entry.Desc
		}

		logger.WithFields(fields).Info(fmt.Sprintf("天赋%d", k+1))
	}

}

func (t *TransportClient) OnM2C_StartStageCombat(sock transport.Socket, msg *transport.Message) {
	m := msg.Body.(*pbGame.M2C_StartStageCombat)

	logger.Info("战斗返回结果:", m)
}

func (t *TransportClient) doConnect() {
	for {
		select {
		case <-t.ctx.Done():
			logger.Info("transport client dial goroutine done...")
			return

		case <-t.heartBeatTimer.C:
			t.heartBeatTimer.Reset(t.heartBeatDuration)

			msg := &transport.Message{
				Type: transport.BodyJson,
				Name: "yokai_account.C2M_HeartBeat",
				Body: &pbAccount.C2M_HeartBeat{},
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
			if t.ts, err = t.tr.Dial(t.serverAddr); err != nil {
				logger.Warn("unexpected dial err:", err)
				continue
			}

			logger.WithFields(logger.Fields{
				"local":  t.ts.Local(),
				"remote": t.ts.Remote(),
			}).Info("tcp dial success")

			msg := &transport.Message{
				Type: transport.BodyProtobuf,
				Name: "yokai_account.C2M_AccountLogon",
				Body: &pbAccount.C2M_AccountLogon{
					RpcId:       1,
					UserId:      t.userID,
					AccountId:   t.accountID,
					AccountName: t.userName,
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

func (t *TransportClient) doRecv() {
	for {
		select {
		case <-t.ctx.Done():
			logger.Info("transport client recv goroutine done...")
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
					if msg, h, err := t.ts.Recv(t.r); err != nil {
						logger.Warn("Unexpected recv err:", err)
					} else {
						h.Fn(t.ts, msg)
						if msg.Name != "yokai_account.M2C_HeartBeat" {
							t.recvCh <- 1
						}
					}
				}
			}()
		}
	}
}

func (t *TransportClient) Run() error {
	for {
		select {
		case <-t.ctx.Done():
			logger.Info("transport client context done...")
			return nil
		}
	}

	return nil
}

func (t *TransportClient) Exit() {
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

func (t *TransportClient) WaitRecv() <-chan int {
	return t.recvCh
}

func (t *TransportClient) GetGateEndPoints() []string {
	return t.gateEndpoints
}
