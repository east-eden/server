package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/transport"
	"bitbucket.org/funplus/server/utils"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type GameInfo struct {
	UserID        int64  `json:"userId"`
	UserName      string `json:"userName"`
	AccountID     int64  `json:"accountId"`
	GameID        string `json:"gameId"`
	PublicTcpAddr string `json:"publicTcpAddr"`
	PublicWsAddr  string `json:"publicWsAddr"`
}

type TransportClient struct {
	c       *Client
	tr      transport.Transport
	ts      transport.Socket
	wg      utils.WaitGroupWrapper
	wgRecon utils.WaitGroupWrapper

	gameInfo      *GameInfo
	gateEndpoints []string
	tlsConf       *tls.Config

	protocol       string
	connected      int32
	needReconnect  int32
	cancelRecvSend context.CancelFunc
	chDisconnect   chan int
	returnMsgName  chan string
	unProcedMsg    int32

	ticker *time.Ticker
	chSend chan *transport.Message
	sync.Mutex
}

func NewTransportClient(c *Client, ctx *cli.Context) *TransportClient {

	t := &TransportClient{
		c:              c,
		gateEndpoints:  ctx.StringSlice("gate_endpoints"),
		returnMsgName:  make(chan string, 100),
		ticker:         time.NewTicker(ctx.Duration("heart_beat")),
		chDisconnect:   make(chan int, 1),
		needReconnect:  0,
		connected:      0,
		cancelRecvSend: func() {},
		chSend:         make(chan *transport.Message, 100),
	}

	var certFile, keyFile string
	if ctx.Bool("debug") {
		certFile = ctx.String("cert_path_debug")
		keyFile = ctx.String("key_path_debug")
	} else {
		certFile = ctx.String("cert_path_release")
		keyFile = ctx.String("key_path_release")
	}

	t.tlsConf = &tls.Config{InsecureSkipVerify: true}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal().Err(err).Msg("load certificates failed")
	}

	t.tlsConf.Certificates = []tls.Certificate{cert}

	// toxi proxy
	// toxiClient := toxiproxy.NewClient("localhost:8474")
	// gateProxy, err := toxiClient.CreateProxy("gate_toxi_proxy", "[::]:20000", "host.docker.internal:443")
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("create gate proxy failed")
	// }

	// gameProxy, err := toxiClient.CreateProxy("game_toxi_proxy", "[::]:20001", "host.docker.internal:7030")
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("create game proxy failed")
	// }

	// gateProxy.AddToxic("latency_down", "latency", "downstream", 1.0, toxiproxy.Attributes{
	// 	"latency": 1000,
	// })

	// gameProxy.AddToxic("latency_down", "latency", "downstream", 1.0, toxiproxy.Attributes{
	// 	"latency": 1000,
	// })

	// t.gateEndpoints = []string{"https://localhost:20000/select_game_addr"}

	// timer heart beat
	go func() {
		defer utils.CaptureException()

		for {
			select {
			case <-t.ticker.C:
				if atomic.LoadInt32(&t.connected) == 0 {
					continue
				}

				msg := &transport.Message{
					Name: "C2S_HeartBeat",
					Body: &pbGlobal.C2S_HeartBeat{},
				}
				t.chSend <- msg

			default:
				time.Sleep(time.Millisecond * 100)
			}
		}

	}()

	return t
}

func (t *TransportClient) connect(ctx context.Context) error {
	// dial to server
	var err error
	addr := t.gameInfo.PublicTcpAddr
	if t.protocol == "ws" {
		addr = "wss://" + t.gameInfo.PublicWsAddr
	}

	// toxiproxy
	// addr = "127.0.0.1:20001"

	if t.ts, err = t.tr.Dial(addr); err != nil {
		return fmt.Errorf("TransportClient.Connect failed: %w", err)
	}

	atomic.StoreInt32(&t.connected, 1)

	log.Info().
		Str("local", t.ts.Local()).
		Str("remote", t.ts.Remote()).
		Msg("tcp dial success")

	// send logon
	msg := &transport.Message{
		Name: "C2S_AccountLogon",
		Body: &pbGlobal.C2S_AccountLogon{
			UserId:      t.gameInfo.UserID,
			AccountId:   t.gameInfo.AccountID,
			AccountName: t.gameInfo.UserName,
		},
	}
	t.chSend = make(chan *transport.Message, 100)
	t.chSend <- msg

	// goroutine to send and recv messages
	t.wg.Wrap(func() {
		err := t.onSend(ctx)
		if err != nil {
			log.Warn().
				Int64("client_id", t.c.Id).
				Err(err).
				Msg("TransportClient onSend finished")

			atomic.StoreInt32(&t.needReconnect, 1)
		}
	})

	t.wg.Wrap(func() {
		err := t.onRecv(ctx)
		if err != nil {
			log.Warn().
				Int64("client_id", t.c.Id).
				Err(err).
				Msg("TransportClient onRecv finished")

			atomic.StoreInt32(&t.needReconnect, 1)
		}
	})

	return nil
}

func (t *TransportClient) StartConnect(ctx context.Context) error {
	if t.tr != nil {
		return errors.New("TransportClient.StartConnect failed: connection existed")
	}

	if t.protocol == "tcp" {
		t.tr = transport.NewTransport("tcp")
		err := t.tr.Init(
			transport.Timeout(transport.DefaultDialTimeout),
		)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	} else {
		t.tr = transport.NewTransport("ws")
		err := t.tr.Init(
			transport.Timeout(transport.DefaultDialTimeout),
			transport.TLSConfig(t.tlsConf),
		)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	}

	t.wgRecon.Wrap(func() {
		t.onReconnect(ctx)
	})

	atomic.StoreInt32(&t.needReconnect, 1)

	return nil
}

// disconnect send cancel signal, and wait onRecv and onSend goroutine's context done
func (t *TransportClient) disconnect() {
	log.Info().Int64("client_id", t.c.Id).Msg("transport client disconnect")

	close(t.chSend)
	t.cancelRecvSend()
	atomic.StoreInt32(&t.connected, 0)
	t.wg.Wait()

	if t.ts != nil {
		t.ts.Close()
	}
}

func (t *TransportClient) StartDisconnect() {
	t.chDisconnect <- 1
}

func (t *TransportClient) SendMessage(msg *transport.Message) {
	if msg == nil {
		return
	}

	if t.ts == nil {
		log.Warn().Msg("未连接到服务器")
		return
	}

	t.chSend <- msg
}

func (t *TransportClient) SetGameInfo(info *GameInfo) {
	t.gameInfo = info
}

func (t *TransportClient) SetProtocol(p string) {
	t.protocol = p
}

func (t *TransportClient) onSend(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Info().Int64("client_id", t.c.Id).Msg("transport client send goroutine done...")
			return nil

		case msg := <-t.chSend:
			if atomic.LoadInt32(&t.connected) == 0 {
				log.Warn().Msg("TransportClient.onSend failed: unconnected to server")
				continue
			}

			if err := t.ts.Send(msg); err != nil {
				return fmt.Errorf("TransportClient.OnSend failed: %w", err)
			}
		}
	}
}

func (t *TransportClient) onRecv(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Info().Int64("client_id", t.c.Id).Msg("transport client recv goroutine done...")
			return nil

		default:
			// be called per 100ms
			// ct := time.Now()
			// defer func() {
			// d := time.Since(ct)
			// time.Sleep(100*time.Millisecond - d)
			// }()

			if atomic.LoadInt32(&t.connected) == 0 {
				log.Warn().Msg("TransportClient.onRecv failed: unconnected to server")
				continue
			}

			if msg, h, err := t.ts.Recv(t.c.msgHandler.r); err != nil {
				return fmt.Errorf("TransportClient.onRecv failed: %w", err)

			} else {
				err := h.Fn(ctx, t.ts, msg)
				if err != nil {
					return fmt.Errorf("TransportClient.onRecv failed: %w", err)
				}

				if msg.Name != "S2C_HeartBeat" {
					t.returnMsgName <- msg.Name
					atomic.AddInt32(&t.unProcedMsg, 1)
					num := atomic.LoadInt32(&t.unProcedMsg)
					if num >= 90 {
						log.Warn().Int64("client_id", t.c.Id).Int32("unproc", num).Msg("return msg name ")
					}
				}
			}
		}
	}
}

func (t *TransportClient) onReconnect(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Info().Int64("client_id", t.c.Id).Msg("transport client reconnect goroutine done...")
			return

		case <-t.chDisconnect:
			log.Info().Int64("client_id", t.c.Id).Msg("transport client disconnected, please rerun to start connection to server again")
			return

		default:
			func() {
				ct := time.Now()
				defer func() {
					d := time.Since(ct)
					time.Sleep(2*time.Second - d)
				}()

				// reconnect
				re := atomic.LoadInt32(&t.needReconnect)
				if re > 0 {
					t.disconnect()
					log.Info().Msg("start reconnect...")

					subCtx, subCancel := context.WithCancel(ctx)
					t.cancelRecvSend = subCancel
					err := t.connect(subCtx)
					if err != nil {
						log.Warn().Err(err).Msg("TransportClient.onReconnect failed")
					} else {
						atomic.StoreInt32(&t.needReconnect, 0)
					}
				}
			}()
		}
	}
}

func (t *TransportClient) Run(ctx *cli.Context) error {
	<-ctx.Done()
	log.Info().Int64("client_id", t.c.Id).Msg("transport client context done...")
	return nil
}

func (t *TransportClient) Exit(ctx *cli.Context) {
	if t.ts != nil {
		t.ts.Close()
	}

	// wait for onRecv and onSend context done
	t.wg.Wait()

	// wait for onReconnect context done
	t.wgRecon.Wait()
}

func (t *TransportClient) ReturnMsgName() <-chan string {
	return t.returnMsgName
}

func (t *TransportClient) GetGateEndPoints() []string {
	return t.gateEndpoints
}
