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

	messages  map[int]*transport.Message
	reconn    chan int
	sendQueue chan *transport.Message
	recvQueue chan *transport.Message
}

func NewTcpClient(opts *Options, ctx context.Context) *TcpClient {
	t := &TcpClient{
		tr:             transport.NewTransport(transport.Timeout(transport.DefaultDialTimeout)),
		opts:           opts,
		heartBeatTimer: time.NewTimer(opts.HeartBeat),
		messages:       make(map[int]*transport.Message),
		reconn:         make(chan int, 1),
		sendQueue:      make(chan *transport.Message, 1000),
		recvQueue:      make(chan *transport.Message, 1000),
	}

	t.ctx, t.cancel = context.WithCancel(ctx)

	t.initSendMessage()

	t.waitGroup.Wrap(func() {
		t.doConnect()
	})

	t.waitGroup.Wrap(func() {
		t.doSend()
	})

	t.waitGroup.Wrap(func() {
		t.doRecv()
	})

	t.waitGroup.Wrap(func() {
		t.handleRecv()
	})

	t.reconn <- 1

	return t
}

func (t *TcpClient) initSendMessage() {
	t.messages[1] = &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_client.MC_ClientLogon",
		Body: &pbClient.MC_ClientLogon{
			ClientId:   1,
			ClientName: "dudu",
		},
	}

	t.messages[2] = &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_client.MC_HeartBeat",
		Body: &pbClient.MC_HeartBeat{},
	}

	t.messages[3] = &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_client.MC_ClientConnected",
		Body: &pbClient.MC_ClientConnected{ClientId: 1},
	}
}

func (t *TcpClient) doConnect() {
	for {
		select {
		case <-t.ctx.Done():
			logger.Info("tcp client dial goroutine done...")
			return

		case <-t.heartBeatTimer.C:
			t.heartBeatTimer.Reset(t.opts.HeartBeat)
			t.SendMessage(t.messages[2])

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

			logger.Info("tpc dial at remote:%s, local:%s", t.ts.Remote(), t.ts.Local())

			t.SendMessage(t.messages[1])
		}
	}
}

func (t *TcpClient) doSend() {
	for {
		select {
		case <-t.ctx.Done():
			logger.Info("tcp client send goroutine done...")
			return

		case msg := <-t.sendQueue:
			logger.Info("begin send message:", msg)
			if err := t.ts.Send(msg); err != nil {
				logger.Warn("Unexpected send err", err)
				t.reconn <- 1
			}
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
					if msg, err := t.ts.Recv(); err != nil {
						logger.Warn("Unexpected recv err", err)
						t.reconn <- 1
					} else {
						t.recvQueue <- msg
					}
				}
			}()
		}
	}
}

func (t *TcpClient) handleRecv() {
	for {
		select {
		case <-t.ctx.Done():
			logger.Info("tcp handle recv context done...")
			return
		case msg := <-t.recvQueue:
			logger.Println("handle recv message:", msg)

			if msg.Name == "yokai_client.MS_ClientLogon" {
				t.SendMessage(t.messages[3])
			}
		}
	}
}

func (t *TcpClient) SendMessage(msg *transport.Message) {
	if t.ts != nil {
		t.sendQueue <- msg
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
}

func (t *TcpClient) Exit() {
	t.cancel()
	t.heartBeatTimer.Stop()

	if t.ts != nil {
		t.ts.Close()
	}

	t.waitGroup.Wait()
}
