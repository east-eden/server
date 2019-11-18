package game

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/micro/go-micro/transport"
	"github.com/micro/go-plugins/transport/tcp"
	logger "github.com/sirupsen/logrus"
)

// Context specifies a context for the service.
// Can be used to signal shutdown of the service.
// Can be used for extra option values.
func Context(ctx context.Context) transport.Option {
	return func(o *transport.Options) {
		o.Context = ctx
	}
}

type TcpServer struct {
	tr     transport.Transport
	ls     transport.Listener
	g      *Game
	mu     sync.Mutex
	parser *MsgParser
	socks  map[transport.Socket]struct{}
	wp     *workerpool.WorkerPool
	ctx    context.Context
	cancel context.CancelFunc
	errCh  chan error
}

func NewTcpServer(g *Game) *TcpServer {
	s := &TcpServer{
		g:      g,
		parser: NewMsgParser(g),
		socks:  make(map[transport.Socket]struct{}),
		wp:     workerpool.New(runtime.GOMAXPROCS(runtime.NumCPU())),
		errCh:  make(chan error, 0),
	}

	s.ctx, s.cancel = context.WithCancel(g.ctx)
	s.serve()
	return s
}

func (s *TcpServer) serve() error {
	s.tr = tcp.NewTransport(transport.Timeout(time.Millisecond * 100))
	var err error
	s.ls, err = s.tr.Listen(s.g.opts.TCPListenAddr)
	if err != nil {
		logger.Error("TcpServer listen error", err)
		return err
	}

	logger.Info("TcpServer listened at:", s.ls.Addr())

	go func() {
		if err := s.ls.Accept(s.handleSocket); err != nil {
			logger.Error("TcpServer accept error:", err)
			s.errCh <- err
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
		case err := <-s.errCh:
			logger.Error("TcpServer Run error:", err)
			return err
		}
	}
}

func (s *TcpServer) Exit() {
	s.cancel()
	s.ls.Close()
}

func (s *TcpServer) handleSocket(sock transport.Socket) {
	defer func() {
		sock.Close()
	}()

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
		var msg transport.Message
		if err := sock.Recv(&msg); err != nil {
			logger.Error("tcp server handle socket error", err)
			return
		}

		ctype := msg.Header["Content-Type"]
		name := msg.Header["Name"]

		// protobuf
		if ctype == "application/x-protobuf" && len(name) > 0 {
			p := s.parser
			sock := sock
			name := name
			s.wp.Submit(func() {
				p.ParserProtoMessage(sock, name, &msg)
			})
		} else {
			logger.WithFields(logger.Fields{
				"header": msg.Header,
				"body":   msg.Body,
			}).Warn("tcp server received invalid protobuf message")
		}

	}

	s.mu.Lock()
	delete(s.socks, sock)
	s.mu.Unlock()
}

/* m := transport.Message{*/
//Header: map[string]string{
//"Content-Type": "application/json",
//},
//Body: []byte(`{"message": "Hello World"}`),
//}

//if err := c.Send(&m); err != nil {
//t.Errorf("Unexpected send err: %v", err)
//}
