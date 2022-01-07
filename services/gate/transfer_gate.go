package gate

import (
	"context"
	"errors"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/transport"
	"github.com/east-eden/server/utils"
	"github.com/panjf2000/ants/v2"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/proto"
)

var (
	TcpRecvInternal = 100 * time.Millisecond
)

type TransferGate struct {
	gate    *Gate
	wg      sync.WaitGroup
	pool    *ants.Pool
	Options *TransferGateOptions
	front   transport.Transport // transport with client
	reg     transport.Register
}

func NewTransferGate(ctx *cli.Context, gate *Gate) *TransferGate {
	maxClient := ctx.Int("max_client")
	tg := &TransferGate{
		gate: gate,
		reg:  transport.NewTransportRegister(),
	}

	var err error
	tg.pool, err = ants.NewPool(maxClient, ants.WithExpiryDuration(10*time.Second))
	if !utils.ErrCheck(err, "NewPool failed when NewTransferGate", maxClient) {
		return nil
	}

	err = tg.register()
	if !utils.ErrCheck(err, "register failed when NewTransferGate") {
		return nil
	}

	err = tg.serve(ctx)
	if !utils.ErrCheck(err, "serve failed when NewTransferGate") {
		return nil
	}

	return tg
}

func (tg *TransferGate) register() error {
	tg.reg.RegisterProtobufMessage(&pbGlobal.Handshake{}, tg.handleHandshake)
	return nil
}

func (tg *TransferGate) serve(ctx *cli.Context) error {
	tg.front = transport.NewTransport("tcp")
	tg.front.Init(
		transport.Timeout(transport.DefaultServeTimeout),
	)

	go func() {
		defer utils.CaptureException()

		err := tg.front.ListenAndServe(ctx.Context, ctx.String("tcp_listen_addr"), tg)
		if err != nil {
			log.Warn().Err(err).Msg("front ListenAndServe failed")
			os.Exit(1)
		}
	}()

	log.Info().
		Str("addr", ctx.String("tcp_listen_addr")).
		Msg("front server serve at address")

	return nil
}

func (tg *TransferGate) Run(ctx *cli.Context) error {
	<-ctx.Done()
	log.Info().Msg("transfer gate server context done...")
	return nil
}

func (tg *TransferGate) Exit() {
	tg.wg.Wait()
	log.Info().Msg("transfer gate exit...")
}

func (tg *TransferGate) HandleSocket(ctx context.Context, frontSock transport.Socket) {
	subCtx, cancel := context.WithCancel(ctx)
	tg.wg.Add(1)
	err := tg.pool.Submit(func() {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				log.Error().Msgf("catch exception:%v, panic recovered with stack:%s", err, stack)
			}

			frontSock.Close()
			cancel()
			tg.wg.Done()
		}()

		// handshake
		msg, h, err := frontSock.Recv(tg.reg)
		if err != nil {
			return
		}

		// validation
		if err := h.Fn(subCtx, frontSock, msg); err != nil {
			log.Warn().
				Caller().
				Err(err).
				Str("msg", string(msg.ProtoReflect().Descriptor().Name())).
				Msg("TransferGate.handleSocket callback error")
			return
		}

		handshake, ok := msg.(*pbGlobal.Handshake)
		if !ok {
			log.Warn().Caller().Msg("assert to Handshake failed")
			return
		}

		_, metadata := tg.gate.gs.SelectGame(handshake.UserId)
		backend := transport.NewTransport("tcp")
		backend.Init()
		backendSock, err := backend.Dial(metadata["publicTcpAddr"])
		if !utils.ErrCheck(err, "Dial failed when TransferGate.HandleSocket") {
			return
		}

		// begin transfer
		_, err = io.Copy(frontSock, backendSock)
		if err != nil {
			frontSock.Close()
			backendSock.Close()
		}
		go func() {
			_, err := io.Copy(backendSock, frontSock)
			if err != nil {
				frontSock.Close()
				backendSock.Close()
			}
		}()
	})

	utils.ErrPrint(err, "Submit failed when handleSocket")
}

func (tg *TransferGate) handleHandshake(ctx context.Context, sock transport.Socket, p proto.Message) error {
	_, ok := p.(*pbGlobal.Handshake)
	if !ok {
		return errors.New("handleHandshake failed: cannot assert value to message")
	}

	// todo check client version

	// todo validation

	return nil
}
