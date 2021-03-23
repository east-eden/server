package transport

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime/debug"
	"sync"
	"testing"
	"time"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/transport/codec"
	"bitbucket.org/funplus/server/utils"
	"github.com/google/go-cmp/cmp"
)

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				fmt.Println(stack)
			}

			w.Done()
		}()

		cb()
	}()
}

var (
	trTcpSrv  = NewTransport("tcp")
	regTcpSrv = NewTransportRegister()
	trTcpCli  = NewTransport("tcp")
	regTcpCli = NewTransportRegister()
	wgTcp     WaitGroupWrapper
)

var (
	trWsSrv  = NewTransport("ws")
	regWsSrv = NewTransportRegister()
	trWsCli  = NewTransport("ws")
	regWsCli = NewTransportRegister()
	wgWs     WaitGroupWrapper
)

func handleTcpServerSocket(ctx context.Context, sock Socket, closeHandler SocketCloseHandler) {
	defer func() {
		sock.Close()
		closeHandler()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, h, err := sock.Recv(regTcpSrv)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("handleTcpServerSocket Recv eof, close connection")
				return
			}

			log.Printf("handleTcpServerSocket Recv failed: %v\n", err)
			return
		}

		_ = h.Fn(ctx, sock, msg)
	}
}

func handleTcpClientAccountLogon(ctx context.Context, sock Socket, p *Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_AccountLogon)
	if !ok {
		log.Fatalf("handleClient failed")
	}

	var sendMsg Message
	sendMsg.Name = "S2C_AccountLogon"
	sendMsg.Body = &pbGlobal.S2C_AccountLogon{
		PlayerName: msg.AccountName,
	}

	_ = sock.Send(&sendMsg)
	return nil
}

func handleTcpServerAccountLogon(ctx context.Context, sock Socket, p *Message) error {
	msg, ok := p.Body.(*pbGlobal.S2C_AccountLogon)
	if !ok {
		log.Fatalf("handleServer failed")
	}

	diff := cmp.Diff(msg.PlayerName, "test_name")
	if diff != "" {
		log.Fatalf("handleTcpServerAccountLogon failed: %s", diff)
	}

	return nil
}

func TestTransportTcp(t *testing.T) {

	// tcp server
	_ = trTcpSrv.Init(
		Timeout(DefaultServeTimeout),
		Codec(&codec.ProtoBufMarshaler{}),
	)

	_ = regTcpSrv.RegisterProtobufMessage(&pbGlobal.C2S_AccountLogon{}, handleTcpClientAccountLogon)

	ctxServ, cancelServ := context.WithCancel(context.Background())
	wgTcp.Wrap(func() {
		err := trTcpSrv.ListenAndServe(ctxServ, ":7030", handleTcpServerSocket)
		if err != nil {
			log.Fatalf("TcpServer ListenAndServe failed: %v", err)
		}
	})

	// tcp client
	_ = trTcpCli.Init(
		Timeout(DefaultDialTimeout),
	)

	_ = regTcpCli.RegisterProtobufMessage(&pbGlobal.S2C_AccountLogon{}, handleTcpServerAccountLogon)

	time.Sleep(time.Millisecond * 500)
	sockClient, err := trTcpCli.Dial("127.0.0.1:7030")
	if err != nil {
		log.Fatalf("unexpected tcp dial err: %v", err)
	}

	ctxCli, cancelCli := context.WithCancel(context.Background())
	wgTcp.Wrap(func() {
		defer func() {
			sockClient.Close()
		}()

		for {
			select {
			case <-ctxCli.Done():
				return
			default:
			}

			if msg, h, err := sockClient.Recv(regTcpCli); err != nil {
				log.Printf("Unexpected recv err: %v", err)
			} else {
				_ = h.Fn(ctxCli, sockClient, msg)
			}
		}
	})

	// send protobuf message
	msgProtobuf := &Message{
		Name: "C2S_AccountLogon",
		Body: &pbGlobal.C2S_AccountLogon{
			UserId:      1,
			AccountId:   1,
			AccountName: "test_name",
		},
	}

	wgTcp.Wrap(func() {
		if err := sockClient.Send(msgProtobuf); err != nil {
			log.Fatalf("client send socket failed: %v", err)
		}
	})

	time.Sleep(time.Second)
	cancelCli()
	cancelServ()
	wgTcp.Wait()
}

func handleWsServerSocket(ctx context.Context, sock Socket, closeHandler SocketCloseHandler) {
	defer func() {
		sock.Close()
		closeHandler()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, h, err := sock.Recv(regWsSrv)
		if err != nil {
			log.Printf("ws server handle socket error: %v", err)
			return
		}

		_ = h.Fn(ctx, sock, msg)
	}
}

func handleWsClient(ctx context.Context, sock Socket, p *Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_AccountLogon)
	if !ok {
		log.Fatalf("handleClient failed")
	}

	var sendMsg Message
	sendMsg.Name = "S2C_AccountLogon"
	sendMsg.Body = &pbGlobal.S2C_AccountLogon{
		PlayerName: msg.AccountName,
	}

	_ = sock.Send(&sendMsg)
	return nil
}

func handleWsServer(ctx context.Context, sock Socket, p *Message) error {
	msg, ok := p.Body.(*pbGlobal.S2C_AccountLogon)
	if !ok {
		log.Fatalf("handleServer failed")
	}

	if msg.PlayerName != "test_name" {
		log.Fatalf("handleServer failed")
	}

	wgWs.Done()
	return nil
}

func TestTransportWs(t *testing.T) {

	// cert
	certPath := "../config/cert/localhost.crt"
	keyPath := "../config/cert/localhost.key"
	tlsConfServ := &tls.Config{}
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatalf("load certificates failed:%v", err)
	}

	tlsConfServ.Certificates = []tls.Certificate{cert}

	// ws server
	_ = trWsSrv.Init(
		Timeout(DefaultServeTimeout),
		Codec(&codec.ProtoBufMarshaler{}),
		TLSConfig(tlsConfServ),
	)

	_ = regWsSrv.RegisterProtobufMessage(&pbGlobal.C2S_AccountLogon{}, handleWsClient)

	go func() {
		defer utils.CaptureException()
		err := trWsSrv.ListenAndServe(context.Background(), ":4433", handleWsServerSocket)
		if err != nil {
			log.Fatalf("WsServer ListenAndServe failed: %v", err)
		}
	}()

	// ws client
	tlsConfCli := &tls.Config{InsecureSkipVerify: true}
	tlsConfCli.Certificates = []tls.Certificate{cert}
	_ = trWsCli.Init(
		Timeout(DefaultDialTimeout),
		TLSConfig(tlsConfCli),
	)

	_ = regWsCli.RegisterProtobufMessage(&pbGlobal.S2C_AccountLogon{}, handleWsServer)

	time.Sleep(time.Millisecond * 500)
	sockClient, err := trWsCli.Dial("wss://localhost:4433")
	if err != nil {
		log.Fatalf("unexpected web socket dial err: %v", err)
	}

	msg := &Message{
		Name: "C2S_AccountLogon",
		Body: &pbGlobal.C2S_AccountLogon{
			UserId:      1,
			AccountId:   1,
			AccountName: "test_name",
		},
	}

	ctxCli, cancelCli := context.WithCancel(context.Background())
	wgWs.Wrap(func() {
		defer sockClient.Close()

		for {
			select {
			case <-ctxCli.Done():
				return
			default:
			}

			if msg, h, err := sockClient.Recv(regWsCli); err != nil {
				log.Printf("Unexpected recv err: %v", err)
			} else {
				_ = h.Fn(ctxCli, sockClient, msg)
			}
		}
	})

	wgWs.Wrap(func() {
		if err := sockClient.Send(msg); err != nil {
			log.Fatalf("client send socket failed: %v", err)
		}
	})

	time.Sleep(time.Second * 2)
	cancelCli()
	wgWs.Wait()
}
