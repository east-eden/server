package client

import (
	"context"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/transport"
	"github.com/yokaiio/yokai_server/internal/utils"
	pbClient "github.com/yokaiio/yokai_server/proto/client"
)

type TcpClient struct {
	tr             transport.Transport
	ts             transport.Socket
	opts           *Options
	ctx            context.Context
	cancel         context.CancelFunc
	waitGroup      utils.WaitGroupWrapper
	heartBeatTimer *time.Timer

	reconn chan int
}

type MC_ClientTest struct {
	ClientId int64  `protobuf:"varint,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	Name     string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func NewTcpClient(opts *Options, ctx context.Context) *TcpClient {
	t := &TcpClient{
		tr:             transport.NewTransport(transport.Timeout(transport.DefaultDialTimeout)),
		opts:           opts,
		heartBeatTimer: time.NewTimer(opts.HeartBeat),
		reconn:         make(chan int, 1),
	}

	t.ctx, t.cancel = context.WithCancel(ctx)

	t.initSendMessage()

	t.waitGroup.Wrap(func() {
		t.doConnect()
	})

	t.waitGroup.Wrap(func() {
		t.doRecv()
	})

	t.reconn <- 1

	return t
}

func (t *TcpClient) initSendMessage() {

	transport.DefaultRegister.RegisterMessage("yokai_client.MS_ClientLogon", &pbClient.MS_ClientLogon{}, t.OnMS_ClientLogon)
	transport.DefaultRegister.RegisterMessage("yokai_client.MS_HeartBeat", &pbClient.MS_HeartBeat{}, t.OnMS_HeartBeat)
}

func (t *TcpClient) SendLogon() {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_client.MC_ClientLogon",
		Body: &pbClient.MC_ClientLogon{
			ClientId:   1,
			ClientName: "dudu",
		},
	}

	if err := t.ts.Send(msg); err != nil {
		logger.Warn("Unexpected send err", err)
		t.reconn <- 1
	}
}

func (t *TcpClient) SendHeartBeat() {
	msg := &transport.Message{
		Type: transport.BodyJson,
		Name: "yokai_client.MC_HeartBeat",
		Body: &pbClient.MC_HeartBeat{},
	}
	if err := t.ts.Send(msg); err != nil {
		logger.Warn("Unexpected send err", err)
		t.reconn <- 1
	}
}

func (t *TcpClient) SendConnected() {
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_client.MC_ClientConnected",
		Body: &pbClient.MC_ClientConnected{ClientId: 1, Name: "dudu"},
	}

	if err := t.ts.Send(msg); err != nil {
		logger.Warn("Unexpected send err", err)
		t.reconn <- 1
	}
}

func (t *TcpClient) SendTest() {
	msg := &transport.Message{
		Type: transport.BodyJson,
		Name: "MC_ClientTest",
		Body: &MC_ClientTest{ClientId: 1, Name: "test"},
	}

	if err := t.ts.Send(msg); err != nil {
		logger.Warn("Unexpected send err", err)
		t.reconn <- 1
	}
}

func (t *TcpClient) OnMS_ClientLogon(sock transport.Socket, msg *transport.Message) {
	logger.Info("recv MS_ClientLogon")

	t.SendConnected()
	t.SendTest()
}

func (t *TcpClient) OnMS_HeartBeat(sock transport.Socket, msg *transport.Message) {
	logger.Info("recv MS_HeartBeat")
}

func (t *TcpClient) doConnect() {
	for {
		select {
		case <-t.ctx.Done():
			logger.Info("tcp client dial goroutine done...")
			return

		case <-t.heartBeatTimer.C:
			t.heartBeatTimer.Reset(t.opts.HeartBeat)
			t.SendHeartBeat()

		case <-t.reconn:
			// close old connection
			if t.ts != nil {
				t.ts.Close()
			}

			var err error
			if t.ts, err = t.tr.Dial(t.opts.TcpServerAddr); err != nil {
				logger.Warn("unexpected dial err:", err)
				time.Sleep(time.Second * 3)
				t.reconn <- 1
				continue
			}

			logger.Info("tpc dial at remote:", t.ts.Remote())

			t.SendLogon()
		}
	}
}

func (t *TcpClient) doRecv() {
	for {
		select {
		case <-t.ctx.Done():
			logger.Info("tcp client recv goroutine done...")
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
