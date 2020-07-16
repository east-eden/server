package game

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"runtime"
	"sync"

	"github.com/gammazero/workerpool"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/transport"
	"github.com/yokaiio/yokai_server/transport/codec"
)

type WsServer struct {
	tr                transport.Transport
	reg               transport.Register
	g                 *Game
	wg                sync.WaitGroup
	mu                sync.Mutex
	wp                *workerpool.WorkerPool
	socks             map[transport.Socket]struct{}
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

	s.tr = transport.NewTransport("ws")

	s.tr.Init(
		transport.Timeout(transport.DefaultServeTimeout),
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

func (s *WsServer) handleSocket(ctx context.Context, sock transport.Socket, closeHandler transport.SocketCloseHandler) {

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

	s.wp.Submit(func() {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				fmt.Printf("handleSocket panic recovered: %s\ncall stack: %s\n", r, buf)
			}

			sock.Close()
			s.wg.Done()

			s.mu.Lock()
			delete(s.socks, sock)
			s.mu.Unlock()

			closeHandler()
		}()

		for {
			select {
			case <-ctx.Done():
				break
			default:
			}

			msg, h, err := sock.Recv(s.reg)
			if err != nil {
				if errors.Is(err, io.EOF) {
					logger.Info("WsServer.handleSocket Recv io.EOF, close connection :", err)
					return
				}

				logger.Warn("WsServer.handleSocket error: ", err)
				return
			}

			if err := h.Fn(ctx, sock, msg); err != nil {
				// account need disconnect
				if errors.Is(err, ErrAccountDisconnect) {
					logger.Info("WsServer.handleSocket account disconnect initiativly")
					break
				}

				logger.Warn("WsServer.handleSocket callback error: ", err)
			}
		}
	})
}
