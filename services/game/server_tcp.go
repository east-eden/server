package game

import (
	"context"
	"errors"
	"io"
	"net"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/east-eden/server/services/game/player"
	"github.com/east-eden/server/transport"
	"github.com/east-eden/server/transport/codec"
	"github.com/east-eden/server/utils"
	"github.com/panjf2000/ants/v2"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	tpcRecvInterval = time.Millisecond * 100 // tcp recv interval per connection
)

type TcpServer struct {
	tr   transport.Transport
	reg  transport.Register
	g    *Game
	wg   sync.WaitGroup
	pool *ants.Pool
}

func NewTcpServer(ctx *cli.Context, g *Game) *TcpServer {
	maxAccount := ctx.Int("account_connect_max")

	s := &TcpServer{
		g:   g,
		reg: g.msgRegister.r,
	}

	var err error
	s.pool, err = ants.NewPool(maxAccount, ants.WithExpiryDuration(10*time.Second))
	if !utils.ErrCheck(err, "NewPool failed when NewTcpServer", maxAccount) {
		return nil
	}

	err = s.serve(ctx)
	if !utils.ErrCheck(err, "serve failed when NewTcpServer", maxAccount) {
		return nil
	}

	return s
}

func (s *TcpServer) serve(ctx *cli.Context) error {
	s.tr = transport.NewTransport("tcp")

	err := s.tr.Init(
		transport.Timeout(transport.DefaultServeTimeout),
		transport.Codec(&codec.ProtoBufMarshaler{}),
	)

	if err != nil {
		return err
	}

	go func() {
		defer utils.CaptureException()

		err := s.tr.ListenAndServe(ctx.Context, ctx.String("tcp_listen_addr"), s)
		if err != nil {
			log.Warn().Err(err).Msg("tcp server ListenAndServe return with error")
			os.Exit(1)
		}
	}()

	log.Info().
		Str("addr", ctx.String("tcp_listen_addr")).
		Msg("tcp server serve at address")

	return nil
}

func (s *TcpServer) Run(ctx context.Context) error {
	<-ctx.Done()
	log.Info().Msg("tcp server context done...")
	return nil
}

func (s *TcpServer) Exit() {
	s.wg.Wait()
	log.Info().Msg("tcp server exit...")
}

func (s *TcpServer) HandleSocket(ctx context.Context, sock transport.Socket) {
	subCtx, cancel := context.WithCancel(ctx)
	s.wg.Add(1)
	err := s.pool.Submit(func() {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				log.Error().Msgf("catch exception:%v, panic recovered with stack:%s", err, stack)
			}

			sock.Close()
			cancel()
			s.wg.Done()
		}()

		for {
			ct := time.Now()

			select {
			case <-subCtx.Done():
				return
			default:
			}

			msg, h, err := sock.Recv(s.reg)
			if err != nil {
				if !errors.Is(err, io.EOF) && !errors.Is(err, net.ErrClosed) {
					log.Warn().Err(err).Msg("TcpServer.handleSocket error")
				}
				return
			}

			if err := h.Fn(subCtx, sock, msg); err != nil {
				// account need disconnect
				if errors.Is(err, player.ErrAccountDisconnect) {
					log.Info().Msg("TcpServer.handleSocket account disconnect initiativly")
					return
				}

				log.Warn().Caller().Err(err).Str("msg", msg.Name).Msg("TcpServer.handleSocket callback error")
			}

			time.Sleep(tpcRecvInterval - time.Since(ct))
		}
	})

	utils.ErrPrint(err, "Submit failed when handleSocket")
}
