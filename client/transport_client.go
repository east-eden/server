package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
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

	heartBeatTimer    *time.Timer
	heartBeatDuration time.Duration

	gameInfo      *GameInfo
	gateEndpoints []string
	tlsConf       *tls.Config

	returnMsgName chan string
}

func NewTransportClient(c *Client, ctx *cli.Context) *TransportClient {

	t := &TransportClient{
		heartBeatDuration: ctx.Duration("heart_beat"),
		heartBeatTimer:    time.NewTimer(ctx.Duration("heart_beat")),
		gateEndpoints:     ctx.StringSlice("gate_endpoints"),
		reconn:            make(chan int, 1),
		returnMsgName:     make(chan string, 100),
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

	t.heartBeatTimer.Stop()

	return t
}

func (t *TransportClient) TcpConnect() error {
	t.tr = transport.NewTransport("tcp")
	t.tr.Init(
		transport.Timeout(transport.DefaultDialTimeout),
	)

	t.Disconnect()
	t.wg.Wait()

	logger.Info("start connect to server")
	t.wg.Wrap(func() {
		t.reconn <- 1
		t.onConnect()
	})

	t.wg.Wrap(func() {
		t.onRecv()
	})

	return fmt.Errorf("unexpected error")
}

func (t *TransportClient) WsConnect() error {
	t.tr = transport.NewTransport("ws")
	t.tr.Init(
		transport.Timeout(transport.DefaultDialTimeout),
		transport.TLSConfig(t.tlsConf),
	)

	return nil
}

func (t *TransportClient) Disconnect() {
	logger.WithFields(logger.Fields{
		"local": t.ts.Local(),
	}).Info("socket local disconnect")

	t.cancel()
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

func (t *TransportClient) onConnect() {
	for {
		select {
		case <-t.ctx.Done():
			logger.Info("transport client dial goroutine done...")
			return

		case <-t.heartBeatTimer.C:
			t.heartBeatTimer.Reset(t.heartBeatDuration)

			msg := &transport.Message{
				Type: transport.BodyJson,
				Name: "yokai_account.C2M_HeartBeat",
				Body: &pbAccount.C2M_HeartBeat{},
			}
			t.SendMessage(msg)

		case <-t.reconn:
			// close old connection
			if t.ts != nil {
				logger.WithFields(logger.Fields{
					"local": t.ts.Local(),
				}).Info("socket local reconnect")
				t.ts.Close()
			}

			// dial to server
			var err error
			addr := t.gameInfo.PublicTcpAddr
			if t.tr.Protocol() == "ws" {
				addr = "wss://" + t.gameInfo.PublicWsAddr
			}

			if t.ts, err = t.tr.Dial(addr); err != nil {
				logger.Warn("unexpected dial err:", err)
				continue
			}

			logger.WithFields(logger.Fields{
				"local":  t.ts.Local(),
				"remote": t.ts.Remote(),
			}).Info("tcp dial success")

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
			t.SendMessage(msg)

			logger.WithFields(logger.Fields{
				"user_id":    t.gameInfo.UserID,
				"account_id": t.gameInfo.AccountID,
				"local":      t.ts.Local(),
			}).Info("connect send message")
		}
	}
}

func (t *TransportClient) onRecv(ctx context.Context) {
	for {
		select {
		case <-t.ctx.Done():
			logger.Info("transport client recv goroutine done...")
			return

		default:
			// be called per 100ms
			ct := time.Now()
			defer func() {
				d := time.Since(ct)
				time.Sleep(100*time.Millisecond - d)
			}()

			if t.ts != nil {
				if msg, h, err := t.ts.Recv(t.c.msgHandler.r); err != nil {
					if errors.Is(err, io.EOF) {
						logger.Info("TransportClient.onRecv recv eof, close connection")
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
	t.heartBeatTimer.Stop()
	t.wg.Wait()
}

func (t *TransportClient) ReturnMsgName() <-chan string {
	return t.returnMsgName
}

func (t *TransportClient) GetGateEndPoints() []string {
	return t.gateEndpoints
}
