package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
	"github.com/yokaiio/yokai_server/transport"
	"github.com/yokaiio/yokai_server/utils"
)

type GameInfo struct {
	UserID        int64  `json:"userId"`
	UserName      string `json:"userName"`
	AccountID     int64  `json:"accountId"`
	GameID        string `json:"gameId"`
	PublicTcpAddr string `json:"publicTcpAddr"`
	PublicWsAddr  string `json:"publicWsAddr"`
	Section       string `json:"section"`
}

type TransportClient struct {
	c  *Client
	tr transport.Transport
	ts transport.Socket
	wg utils.WaitGroupWrapper

	gameInfo      *GameInfo
	gateEndpoints []string
	tlsConf       *tls.Config

	connected     atomic.Value
	cancel        context.CancelFunc
	returnMsgName chan string

	ticker *time.Ticker
	chSend chan *transport.Message
	sync.Mutex
}

func NewTransportClient(c *Client, ctx *cli.Context) *TransportClient {

	t := &TransportClient{
		c:             c,
		gateEndpoints: ctx.StringSlice("gate_endpoints"),
		returnMsgName: make(chan string, 100),
		ticker:        time.NewTicker(ctx.Duration("heart_beat")),
		chSend:        make(chan *transport.Message, 100),
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
		logger.Fatal("load certificates failed:", err)
	}

	t.tlsConf.Certificates = []tls.Certificate{cert}
	t.connected.Store(false)

	// timer heart beat
	go func() {
		for {
			select {
			case <-t.ticker.C:
				if !t.connected.Load().(bool) {
					return
				}

				msg := &transport.Message{
					Type: transport.BodyJson,
					Name: "yokai_account.C2M_HeartBeat",
					Body: &pbAccount.C2M_HeartBeat{},
				}
				t.chSend <- msg
			}
		}

	}()

	return t
}

func (t *TransportClient) Connect(ctx context.Context, protocol string) error {
	if t.connected.Load().(bool) {
		t.Disconnect(ctx)
	}

	if protocol == "tcp" {
		t.tr = transport.NewTransport("tcp")
		t.tr.Init(
			transport.Timeout(transport.DefaultDialTimeout),
		)
	} else {
		t.tr = transport.NewTransport("ws")
		t.tr.Init(
			transport.Timeout(transport.DefaultDialTimeout),
			transport.TLSConfig(t.tlsConf),
		)
	}

	// dial to server
	var err error
	addr := t.gameInfo.PublicTcpAddr
	if protocol == "ws" {
		addr = "wss://" + t.gameInfo.PublicWsAddr
	}

	if t.ts, err = t.tr.Dial(addr); err != nil {
		return fmt.Errorf("TransportClient.Connect failed: %w", err)
	}

	t.connected.Store(true)

	logger.WithFields(logger.Fields{
		"local":  t.ts.Local(),
		"remote": t.ts.Remote(),
	}).Info("tcp dial success")

	// send logon
	msg := &transport.Message{
		Type: transport.BodyProtobuf,
		Name: "yokai_account.C2M_AccountLogon",
		Body: &pbAccount.C2M_AccountLogon{
			RpcId:       1,
			UserId:      t.gameInfo.UserID,
			AccountId:   t.gameInfo.AccountID,
			AccountName: t.gameInfo.UserName,
		},
	}
	t.chSend <- msg

	// goroutine to send and recv messages
	subCtx, cancel := context.WithCancel(ctx)
	t.cancel = cancel
	t.wg.Wrap(func() {
		t.onSend(subCtx)
	})

	t.wg.Wrap(func() {
		t.onRecv(subCtx)
	})

	return nil
}

func (t *TransportClient) Disconnect(ctx context.Context) {
	logger.Info("transport client disconnect")
	t.cancel()
	t.connected.Store(false)
	t.wg.Wait()
}

func (t *TransportClient) SendMessage(msg *transport.Message) {
	if msg == nil {
		return
	}

	if t.ts == nil {
		logger.Warn("未连接到服务器")
		return
	}

	if err := t.ts.Send(msg); err != nil {
		logger.Warn("Unexpected send err", err)
	}
}

func (t *TransportClient) SetGameInfo(info *GameInfo) {
	t.gameInfo = info
}

func (t *TransportClient) onSend(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("transport client send goroutine done...")
			return

		case msg := <-t.chSend:
			if !t.connected.Load().(bool) {
				continue
			}

			logger.Info("transport clinet send message: ", msg.Name)
			if err := t.ts.Send(msg); err != nil {
				logger.Warn("Unexpected send err", err)
			}
		}
	}
}

func (t *TransportClient) onRecv(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("transport client recv goroutine done...")
			return

		default:
			// be called per 100ms
			ct := time.Now()
			defer func() {
				d := time.Since(ct)
				time.Sleep(100*time.Millisecond - d)
			}()

			if !t.connected.Load().(bool) {
				continue
			}

			if msg, h, err := t.ts.Recv(t.c.msgHandler.r); err != nil {
				if errors.Is(err, io.EOF) {
					logger.Info("TransportClient.onRecv recv io.EOF, close connection: ", err)
					return
				}

				logger.Warn("Unexpected recv err:", err)

			} else {
				h.Fn(ctx, t.ts, msg)
				if msg.Name != "yokai_account.M2C_HeartBeat" {
					t.returnMsgName <- msg.Name
				}
			}
		}
	}
}

func (t *TransportClient) Run(ctx *cli.Context) error {
	for {
		select {
		case <-ctx.Done():
			logger.Info("transport client context done...")
			return nil
		}
	}

	return nil
}

func (t *TransportClient) Exit(ctx *cli.Context) {
	t.wg.Wait()
}

func (t *TransportClient) ReturnMsgName() <-chan string {
	return t.returnMsgName
}

func (t *TransportClient) GetGateEndPoints() []string {
	return t.gateEndpoints
}
