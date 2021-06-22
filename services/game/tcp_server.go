package game

import (
	"context"
	"errors"
	"io"
	"net"
	"runtime/debug"
	"sync"
	"time"

	"e.coding.net/mmstudio/blade/server/services/game/player"
	"e.coding.net/mmstudio/blade/server/transport"
	"e.coding.net/mmstudio/blade/server/transport/codec"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/gammazero/workerpool"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	tpcRecvInterval = time.Millisecond * 100 // tcp recv interval per connection
)

type TcpServer struct {
	tr  transport.Transport
	reg transport.Register
	g   *Game
	wg  sync.WaitGroup
	wp  *workerpool.WorkerPool
}

func NewTcpServer(ctx *cli.Context, g *Game) *TcpServer {
	s := &TcpServer{
		g:   g,
		reg: g.msgRegister.r,
	}

	s.wp = workerpool.New(ctx.Int("account_connect_max"))
	err := s.serve(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("tcpserver serve return error")
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

		err := s.tr.ListenAndServe(ctx, ctx.String("tcp_listen_addr"), s.handleSocket)
		if err != nil {
			log.Warn().Err(err).Msg("tcp server ListenAndServe return with error")
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

func (s *TcpServer) handleSocket(ctx context.Context, sock transport.Socket) {
	subCtx, cancel := context.WithCancel(ctx)
	s.wg.Add(1)
	s.wp.Submit(func() {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				log.Error().Msgf("catch exception:%v, panic recovered with stack:%s", err, stack)
			}

			cancel()
			sock.Close()

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
}
