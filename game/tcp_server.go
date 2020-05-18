package game

import (
	"context"
	"runtime"
	"sync"

	"github.com/gammazero/workerpool"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/transport"
	"github.com/yokaiio/yokai_server/transport/codec"
)

type TcpServer struct {
	tr     transport.Transport
	reg    transport.Register
	g      *Game
	wg     sync.WaitGroup
	mu     sync.Mutex
	socks  map[transport.Socket]struct{}
	wp     *workerpool.WorkerPool
	ctx    context.Context
	cancel context.CancelFunc

	accountConnectMax int
}

func NewTcpServer(g *Game, ctx *cli.Context) *TcpServer {
	s := &TcpServer{
		g:                 g,
		reg:               g.msgHandler.r,
		socks:             make(map[transport.Socket]struct{}),
		wp:                workerpool.New(runtime.GOMAXPROCS(runtime.NumCPU())),
		accountConnectMax: ctx.Int("account_connect_max"),
	}

	s.ctx, s.cancel = context.WithCancel(ctx)
	s.serve(ctx)
	return s
}

func (s *TcpServer) serve(ctx *cli.Context) error {
	s.tr = transport.NewTransport(
		"tcp",
		transport.Timeout(transport.DefaultDialTimeout),
		transport.Codec(&codec.ProtoBufMarshaler{}),
	)

	go func() {
		err := s.tr.ListenAndServe(ctx.String("tcp_listen_addr"), s.handleSocket)
		if err != nil {
			logger.Warn("TcpServer serve error:", err)
		}
	}()

	logger.Info("TcpServer serve at:", ctx.String("tcp_listen_addr"))

	return nil
}

func (s *TcpServer) Run() error {
	for {
		select {
		case <-s.ctx.Done():
			logger.Info("TcpServer context done...")
			return nil
		}
	}
}

func (s *TcpServer) Exit() {
	s.cancel()
	s.wg.Wait()
}

func (s *TcpServer) handleSocket(sock transport.Socket) {
	defer func() {
		sock.Close()
		s.wg.Done()
	}()

	s.wg.Add(1)
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

	for {
		select {
		case <-s.ctx.Done():
			break
		default:
		}

		msg, h, err := sock.Recv(s.reg)
		if err != nil {
			logger.Warn("tcp server handle socket error: ", err)
			return
		}

		sock := sock
		s.wp.Submit(func() {
			h.Fn(sock, msg)
		})
	}

	s.mu.Lock()
	delete(s.socks, sock)
	s.mu.Unlock()
}
