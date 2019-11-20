package game

import (
	"context"
	"runtime"
	"sync"

	"github.com/gammazero/workerpool"
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/codec"
	"github.com/yokaiio/yokai_server/internal/transport"
)

type TcpServer struct {
	tr     transport.Transport
	ls     transport.Listener
	g      *Game
	wg     sync.WaitGroup
	mu     sync.Mutex
	parser *MsgParser
	socks  map[transport.Socket]struct{}
	wp     *workerpool.WorkerPool
	ctx    context.Context
	cancel context.CancelFunc
}

func NewTcpServer(g *Game) *TcpServer {
	s := &TcpServer{
		g:      g,
		parser: NewMsgParser(g),
		socks:  make(map[transport.Socket]struct{}),
		wp:     workerpool.New(runtime.GOMAXPROCS(runtime.NumCPU())),
	}

	s.ctx, s.cancel = context.WithCancel(g.ctx)
	s.serve()
	return s
}

func (s *TcpServer) serve() error {
	s.tr = transport.NewTransport(
		transport.Timeout(transport.DefaultDialTimeout),
		transport.Codec(&codec.ProtoBufMarshaler{}),
	)

	var err error
	s.ls, err = s.tr.Listen(s.g.opts.TCPListenAddr)
	if err != nil {
		logger.Error("TcpServer listen error:", err)
		return err
	}

	logger.Info("TcpServer listened at:", s.ls.Addr())

	go func() {
		if err := s.ls.Accept(s.handleSocket); err != nil {
			logger.Warn("TcpServer accept error:", err)
		}
	}()

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
	s.ls.Close()
}

func (s *TcpServer) handleSocket(sock transport.Socket) {
	defer func() {
		sock.Close()
		s.wg.Done()
	}()

	s.wg.Add(1)
	s.mu.Lock()
	sockNum := len(s.socks)
	if sockNum >= s.g.opts.ClientConnectMax {
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

		msg, h, err := sock.Recv()
		if err != nil {
			logger.Warn("tcp server handle socket error", err)
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
