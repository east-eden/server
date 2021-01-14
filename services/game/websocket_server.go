package game

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"runtime"
	"sync"

	"e.coding.net/mmstudio/blade/server/services/game/player"
	"e.coding.net/mmstudio/blade/server/transport"
	"e.coding.net/mmstudio/blade/server/transport/codec"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/gammazero/workerpool"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
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

func NewWsServer(ctx *cli.Context, g *Game) *WsServer {
	s := &WsServer{
		g:                 g,
		reg:               g.msgHandler.r,
		socks:             make(map[transport.Socket]struct{}),
		wp:                workerpool.New(runtime.GOMAXPROCS(runtime.NumCPU())),
		accountConnectMax: ctx.Int("account_connect_max"),
	}

	if err := s.serve(ctx); err != nil {
		log.Warn().Err(err).Msg("websocket server return error")
	}
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

	err = s.tr.Init(
		transport.Timeout(transport.DefaultServeTimeout),
		transport.Codec(&codec.ProtoBufMarshaler{}),
		transport.TLSConfig(tlsConf),
	)

	if err != nil {
		log.Fatal().Err(err).Msg("websocket transport init failed")
	}

	go func() {
		defer utils.CaptureException()
		err := s.tr.ListenAndServe(ctx, ctx.String("websocket_listen_addr"), s.handleSocket)
		if err != nil {
			log.Warn().Err(err).Msg("web socket ListenAndServe return with error")
		}
	}()

	log.Info().Str("addr", ctx.String("websocket_listen_addr")).Msg("web socket serve at")

	return nil
}

func (s *WsServer) Run(ctx context.Context) error {
	<-ctx.Done()
	log.Info().Msg("web socket server context done...")
	return nil
}

func (s *WsServer) Exit() {
	s.wg.Wait()
	log.Info().Msg("web socket server exit...")
}

func (s *WsServer) handleSocket(ctx context.Context, sock transport.Socket, closeHandler transport.SocketCloseHandler) {

	s.wg.Add(1)
	s.mu.Lock()
	sockNum := len(s.socks)
	if sockNum >= s.accountConnectMax {
		s.mu.Unlock()
		log.Warn().
			Int("connection_num", sockNum).
			Msg("too many connections")
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
				return
			default:
			}

			msg, h, err := sock.Recv(s.reg)
			if err != nil {
				log.Warn().Err(err).Msg("WsServer.handleSocket error")
				return
			}

			if err := h.Fn(ctx, sock, msg); err != nil {
				// account need disconnect
				if errors.Is(err, player.ErrAccountDisconnect) {
					log.Info().Msg("WsServer.handleSocket account disconnect initiativly")
					break
				}

				log.Warn().Err(err).Msg("WsServer.handlerSocket callback error")
			}
		}
	})
}
