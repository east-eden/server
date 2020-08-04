package game

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/gammazero/workerpool"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/transport"
	"github.com/yokaiio/yokai_server/transport/codec"
)

var (
	tpcRecvInterval = time.Millisecond * 100 // tcp recv interval per connection
)

type TcpServer struct {
	tr    transport.Transport
	reg   transport.Register
	g     *Game
	wg    sync.WaitGroup
	mu    sync.Mutex
	wp    *workerpool.WorkerPool
	socks map[transport.Socket]struct{}

	accountConnectMax int
}

func NewTcpServer(g *Game, ctx *cli.Context) *TcpServer {
	s := &TcpServer{
		g:                 g,
		reg:               g.msgHandler.r,
		socks:             make(map[transport.Socket]struct{}),
		accountConnectMax: ctx.Int("account_connect_max"),
	}

	if s.accountConnectMax < 1 {
		s.accountConnectMax = 1
	}

	s.wp = workerpool.New(s.accountConnectMax)
	s.serve(ctx)

	return s
}

func (s *TcpServer) serve(ctx *cli.Context) error {
	s.tr = transport.NewTransport("tcp")

	s.tr.Init(
		transport.Timeout(transport.DefaultServeTimeout),
		transport.Codec(&codec.ProtoBufMarshaler{}),
	)

	go func() {
		err := s.tr.ListenAndServe(ctx, ctx.String("tcp_listen_addr"), s.handleSocket)
		if err != nil {
			logger.Warn("TcpServer serve error:", err)
		}
	}()

	logger.Info("TcpServer serve at:", ctx.String("tcp_listen_addr"))

	return nil
}

func (s *TcpServer) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			logger.Info("TcpServer context done...")
			return nil
		}
	}
}

func (s *TcpServer) Exit() {
	s.wg.Wait()
	logger.Info("tcp server exit...")
}

func (s *TcpServer) handleSocket(ctx context.Context, sock transport.Socket, closeHandler transport.SocketCloseHandler) {
	s.mu.Lock()
	sockNum := len(s.socks)
	if sockNum >= s.accountConnectMax {
		s.mu.Unlock()
		logger.WithFields(logger.Fields{
			"connections": sockNum,
		}).Warn("too many connections")
		return
	}

	s.socks[sock] = struct{}{}
	s.mu.Unlock()

	s.wg.Add(1)
	s.wp.Submit(func() {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				fmt.Printf("handleSocket panic recovered: %s\ncall stack: %s\n", r, buf)
			}

			s.mu.Lock()
			delete(s.socks, sock)
			s.mu.Unlock()

			// Socket close handler
			closeHandler()
			s.wg.Done()
		}()

		for {
			ct := time.Now()

			select {
			case <-ctx.Done():
				return
			default:
			}

			msg, h, err := sock.Recv(s.reg)
			if err != nil {
				logger.Warn("TcpServer.handleSocket error: ", err)
				return
			}

			if err := h.Fn(ctx, sock, msg); err != nil {
				// account need disconnect
				if errors.Is(err, ErrAccountDisconnect) {
					logger.Info("TcpServer.handleSocket account disconnect initiativly")
					return
				}

				logger.Warn("TcpServer.handleSocket callback error: ", err)
			}

			time.Sleep(tpcRecvInterval - time.Since(ct))
		}
	})
}
