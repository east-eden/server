package client

import (
	"context"
	"time"

	"github.com/micro/go-micro/transport"
	"github.com/micro/go-plugins/transport/tcp"
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type TcpClient struct {
	tr        transport.Transport
	tc        transport.Client
	opts      *Options
	ctx       context.Context
	cancel    context.CancelFunc
	waitGroup utils.WaitGroupWrapper

	reconn    chan int
	sendQueue chan *transport.Message
	recvQueue chan *transport.Message
}

func NewTcpClient(opts *Options, ctx context.Context) *TcpClient {
	t := &TcpClient{
		tr:        tcp.NewTransport(transport.Timeout(time.Millisecond * 100)),
		opts:      opts,
		reconn:    make(chan int, 1),
		sendQueue: make(chan *transport.Message, 1000),
		recvQueue: make(chan *transport.Message, 1000),
	}

	t.ctx, t.cancel = context.WithCancel(ctx)

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

func (t *TcpClient) doConnect() {
	for {
		select {
		case <-t.ctx.Done():
			logger.Info("tcp client dial goroutine done...")
			return
		case <-t.reconn:
			// close old connection
			if t.tc != nil {
				t.tc.Close()
			}

			var err error
			if t.tc, err = t.tr.Dial(t.opts.TcpServerAddr); err != nil {
				logger.Warn("unexpected dial err:", err)
			}
		}
	}
}

func (t *TcpClient) doSend() {
	for {
		select {
		case <-t.ctx.Done():
			logger.Info("tcp client send goroutine done...")
			return

		default:

			func() {
				// be called per 100ms
				ct := time.Now()
				defer func() {
					d := time.Since(ct)
					time.Sleep(100*time.Millisecond - d)
				}()

				// make sure transport.Client existed
				if t.tc != nil {
					msg := <-t.sendQueue
					if err := t.tc.Send(msg); err != nil {
						logger.Warn("Unexpected send err", err)
						t.reconn <- 1
					}
				}
			}()

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

				if t.tc != nil {
					var msg transport.Message
					if err := t.tc.Recv(&msg); err != nil {
						logger.Warn("Unexpected recv err", err)
						t.reconn <- 1
					} else {
						t.recvQueue <- &msg
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
		}
	}
}

func (t *TcpClient) SendMessage(msg *transport.Message) {
	t.sendQueue <- msg
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
	t.tc.Close()
	t.waitGroup.Wait()
}
