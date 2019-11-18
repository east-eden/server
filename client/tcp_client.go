package client

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
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

	messages  map[int]*transport.Message
	reconn    chan int
	sendQueue chan *transport.Message
	recvQueue chan *transport.Message
}

func NewTcpClient(opts *Options, ctx context.Context) *TcpClient {
	t := &TcpClient{
		tr:        tcp.NewTransport(transport.Timeout(time.Millisecond * 100)),
		opts:      opts,
		messages:  make(map[int]*transport.Message),
		reconn:    make(chan int, 1),
		sendQueue: make(chan *transport.Message, 1000),
		recvQueue: make(chan *transport.Message, 1000),
	}

	t.ctx, t.cancel = context.WithCancel(ctx)

	t.initSendMessage()

	t.waitGroup.Wrap(func() {
		t.input()
	})

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
		Header: map[string]string{
			"Content-Type": "application/json",
		},
		Body: []byte(`{"message": "Hello World"}`),
	}
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

func (t *TcpClient) input() error {

	for {
		select {
		case <-t.ctx.Done():
			logger.Info("Client input context done...")
			return nil

		default:
			return func() error {
				// be called per 100ms
				ct := time.Now()
				defer func() {
					d := time.Since(ct)
					time.Sleep(100*time.Millisecond - d)
				}()

				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter send message number: ")

				text, err := reader.ReadString('\n')
				if err != nil {
					return err
				}

				number, err := strconv.Atoi(text)
				if err != nil {
					return err
				}

				msg, ok := t.messages[number]
				if !ok {
					logger.Warn("cannot find message template number:", number)
				} else {
					t.sendQueue <- msg
				}

				return nil
			}()
		}
	}

}
