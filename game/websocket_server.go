package game

import (
	"context"
	"crypto/tls"
	"fmt"
	"runtime"
	"sync"

	"github.com/gammazero/workerpool"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/transport"
	"github.com/yokaiio/yokai_server/transport/codec"
)

type WsServer struct {
	tr    transport.Transport
	reg   transport.Register
	g     *Game
	wg    sync.WaitGroup
	mu    sync.Mutex
	socks map[transport.Socket]struct{}
	wp    *workerpool.WorkerPool

	accountConnectMax int
}

func NewWsServer(g *Game, ctx *cli.Context) *WsServer {
	s := &WsServer{
		g:                 g,
		reg:               g.msgHandler.r,
		socks:             make(map[transport.Socket]struct{}),
		wp:                workerpool.New(runtime.GOMAXPROCS(runtime.NumCPU())),
		accountConnectMax: ctx.Int("account_connect_max"),
	}

	s.serve(ctx)
	return s
}

func (s *WsServer) serve(ctx *cli.Context) error {
	// cert
	certPath := ctx.String("cert_path_release")
	keyPath := ctx.String("key_path_release")

	if ctx.Bool("debug") {
		certPath = ctx.String("cert_path_debug")
		keyPath = ctx.String("key_path_debug")
	}

	tlsConf := &tls.Config{}
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return fmt.Errorf("load certificates failed:%v", err)
	}
	tlsConf.Certificates = []tls.Certificate{cert}

	s.tr = transport.NewTransport(
		"ws",
		transport.Timeout(transport.DefaultDialTimeout),
		transport.Codec(&codec.ProtoBufMarshaler{}),
		transport.TLSConfig(tlsConf),
	)

	go func() {
		err := s.tr.ListenAndServe(ctx, ctx.String("websocket_listen_addr"), s.handleSocket)
		if err != nil {
			logger.Warn("WsServer serve error:", err)
		}
	}()

	logger.Info("WsServer serve at:", ctx.String("websocket_listen_addr"))

	return nil
}

func (s *WsServer) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			logger.Info("WsServer context done...")
			return nil
		}
	}
}

func (s *WsServer) Exit() {
	s.wg.Wait()
	logger.Info("web server exit...")
}

func (s *WsServer) handleSocket(ctx context.Context, sock transport.Socket) {
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
		case <-ctx.Done():
			break
		default:
		}

		msg, h, err := sock.Recv(s.reg)
		if err != nil {
			logger.Warn("websocket server handle socket error: ", err)
			return
		}

		sock := sock
		s.wp.Submit(func() {
			h.Fn(ctx, sock, msg)
		})
	}

	s.mu.Lock()
	delete(s.socks, sock)
	s.mu.Unlock()
}
