package game

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"e.coding.net/mmstudio/blade/server/services/game/player"
	"e.coding.net/mmstudio/blade/server/transport"
	"e.coding.net/mmstudio/blade/server/transport/codec"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/panjf2000/ants/v2"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type WsServer struct {
	tr   transport.Transport
	reg  transport.Register
	g    *Game
	wg   sync.WaitGroup
	pool *ants.Pool
}

func NewWsServer(ctx *cli.Context, g *Game) *WsServer {
	maxAccount := ctx.Int("account_connect_max")

	s := &WsServer{
		g:   g,
		reg: g.msgRegister.r,
	}

	var err error
	s.pool, err = ants.NewPool(maxAccount, ants.WithExpiryDuration(10*time.Second))
	if !utils.ErrCheck(err, "NewPool failed when NewWsServer", maxAccount) {
		return nil
	}

	err = s.serve(ctx)
	if !utils.ErrCheck(err, "serve failed when NewWsServer") {
		return nil
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

func (s *WsServer) handleSocket(ctx context.Context, sock transport.Socket) {
	subCtx, cancel := context.WithCancel(ctx)
	s.wg.Add(1)
	err := s.pool.Submit(func() {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				log.Error().Caller().Msgf("catch exception:%v, panic recovered with stack:%s", err, stack)
			}

			cancel()
			s.wg.Done()
		}()

		for {
			select {
			case <-subCtx.Done():
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

	utils.ErrPrint(err, "Submit failed when handleSocket")
}
