package client

import (
	"context"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/internal/transport"
	"github.com/yokaiio/yokai_server/internal/utils"
	pbClient "github.com/yokaiio/yokai_server/proto/client"
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

	reconn    chan int
	connected bool

	disconnectCtx    context.Context
	disconnectCancel context.CancelFunc
}

type MC_ClientTest struct {
	ClientId int64  `protobuf:"varint,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	Name     string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func NewTcpClient(ctx *cli.Context) *TcpClient {
	t := &TcpClient{
		tr:                transport.NewTransport(transport.Timeout(transport.DefaultDialTimeout)),
		heartBeatDuration: ctx.Duration("heart_beat"),
		heartBeatTimer:    time.NewTimer(ctx.Duration("heart_beat")),
		tcpServerAddr:     ctx.String("tcp_server_addr"),
		reconn:            make(chan int, 1),
		connected:         false,
	}

	t.ctx, t.cancel = context.WithCancel(ctx)
	t.heartBeatTimer.Stop()

	t.initSendMessage()

	t.reconn <- 1

	return t
}

func (t *TcpClient) initSendMessage() {

	transport.DefaultRegister.RegisterMessage("yokai_client.MS_ClientLogon", &pbClient.MS_ClientLogon{}, t.OnMS_ClientLogon)
	transport.DefaultRegister.RegisterMessage("yokai_client.MS_HeartBeat", &pbClient.MS_HeartBeat{}, t.OnMS_HeartBeat)
}

func (t *TcpClient) Connect() {
	t.disconnectCtx, t.disconnectCancel = context.WithCancel(t.ctx)
	t.waitGroup.Wrap(func() {
		t.doConnect()
	})

	t.waitGroup.Wrap(func() {
		t.doRecv()
	})
}

func (t *TcpClient) Disconnect() {
	t.disconnectCancel()
}

func (t *TcpClient) SendMessage(msg *transport.Message) {
	if msg == nil {
		return
	}

	if err := t.ts.Send(msg); err != nil {
		logger.Warn("Unexpected send err", err)
		t.reconn <- 1
	}
}

func (t *TcpClient) OnMS_ClientLogon(sock transport.Socket, msg *transport.Message) {
	logger.Info("recv MS_ClientLogon")
	logger.Info("server connected")

	t.connected = true
	t.heartBeatTimer.Reset(t.heartBeatDuration)

	send := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_client.MC_ClientConnected",
		Body: &pbClient.MC_ClientConnected{ClientId: 1, Name: "dudu"},
	}
	t.SendMessage(send)

	sendTest := &transport.Message{
		Type: transport.BodyJson,
		Name: "MC_ClientTest",
		Body: &MC_ClientTest{ClientId: 1, Name: "test"},
	}
	t.SendMessage(sendTest)
}

func (t *TcpClient) OnMS_HeartBeat(sock transport.Socket, msg *transport.Message) {
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
				Name: "yokai_client.MC_HeartBeat",
				Body: &pbClient.MC_HeartBeat{},
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
				Name: "yokai_client.MC_ClientLogon",
				Body: &pbClient.MC_ClientLogon{
					ClientId:   1,
					ClientName: "dudu",
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
