package transport

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"log"
	"sync"
	"testing"
	"time"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/transport/codec"
	"github.com/east-eden/server/utils"
	"github.com/google/go-cmp/cmp"
)

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		defer func() {
			utils.CaptureException()
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

		h.Fn(ctx, sock, msg)
	}
}

func handleTcpClientAccountLogon(ctx context.Context, sock Socket, p *Message) error {
	msg, ok := p.Body.(*pbGlobal.C2S_AccountLogon)
	if !ok {
		log.Fatalf("handleClient failed")
	}

	var sendMsg Message
	// sendMsg.Type = BodyProtobuf
	sendMsg.Name = "S2C_AccountLogon"
	sendMsg.Body = &pbGlobal.S2C_AccountLogon{
		PlayerName: msg.AccountName,
	}

	sock.Send(&sendMsg)
	return nil
}

func handleTcpClientAccountTest(ctx context.Context, sock Socket, p *Message) error {
	msg, ok := p.Body.(*C2S_AccountTest)
	if !ok {
		log.Fatalf("handleClient failed")
	}

	var sendMsg Message
	// sendMsg.Type = BodyJson
	sendMsg.Name = "S2C_AccountTest"
	sendMsg.Body = &S2C_AccountTest{
		AccountId: msg.AccountId,
		Name:      msg.Name,
	}

	sock.Send(&sendMsg)
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

	wgTcp.Done()
	return nil
}

func handleTcpServerAccountTest(ctx context.Context, sock Socket, p *Message) error {
	msg, ok := p.Body.(*S2C_AccountTest)
	if !ok {
		log.Fatalf("handleServer json failed")
	}

	diff := cmp.Diff(msg.Name, "test_json")
	if diff != "" {
		log.Fatalf("handleTcpServerAccountTest failed: %s", diff)
	}

	wgTcp.Done()
	return nil
}

// json message define
type C2S_AccountTest struct {
	AccountId int64  `protobuf:"varint,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (msg *C2S_AccountTest) GetName() string {
	return "C2S_AccountTest"
}

type S2C_AccountTest struct {
	AccountId int64  `protobuf:"varint,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (msg *S2C_AccountTest) GetName() string {
	return "S2C_AccountTest"
}

func TestTransportTcp(t *testing.T) {

	// tcp server
	trTcpSrv.Init(
		Timeout(DefaultServeTimeout),
		Codec(&codec.ProtoBufMarshaler{}),
	)

	regTcpSrv.RegisterProtobufMessage(&pbGlobal.C2S_AccountLogon{}, handleTcpClientAccountLogon)
	regTcpSrv.RegisterJsonMessage(&C2S_AccountTest{}, handleTcpClientAccountTest)

	ctxServ, cancelServ := context.WithCancel(context.Background())
	wgTcp.Wrap(func() {
		err := trTcpSrv.ListenAndServe(ctxServ, ":7030", handleTcpServerSocket)
		if err != nil {
			log.Fatalf("TcpServer ListenAndServe failed: %v", err)
		}
	})

	// tcp client
	trTcpCli.Init(
		Timeout(DefaultDialTimeout),
	)

	regTcpCli.RegisterProtobufMessage(&pbGlobal.S2C_AccountLogon{}, handleTcpServerAccountLogon)
	regTcpCli.RegisterJsonMessage(&S2C_AccountTest{}, handleTcpServerAccountTest)

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
				log.Fatalf("Unexpected recv err: %v", err)
			} else {
				h.Fn(ctxCli, sockClient, msg)
			}
		}
	})

	// send protobuf message
	msgProtobuf := &Message{
		// Type: BodyProtobuf,
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

	// send json message
	msgJson := &Message{
		// Type: BodyJson,
		Name: "C2S_AccountTest",
		Body: &C2S_AccountTest{
			AccountId: 1,
			Name:      "test_json",
		},
	}

	wgTcp.Wrap(func() {
		if err := sockClient.Send(msgJson); err != nil {
			log.Fatalf("client send socket failed: %v", err)
		}
	})

	time.Sleep(time.Second)
	cancelServ()
	cancelCli()
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

		h.Fn(ctx, sock, msg)
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

	sock.Send(&sendMsg)
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
	trWsSrv.Init(
		Timeout(DefaultServeTimeout),
		Codec(&codec.ProtoBufMarshaler{}),
		TLSConfig(tlsConfServ),
	)

	regWsSrv.RegisterProtobufMessage(&pbGlobal.C2S_AccountLogon{}, handleWsClient)

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
	trWsCli.Init(
		Timeout(DefaultDialTimeout),
		TLSConfig(tlsConfCli),
	)

	regWsCli.RegisterProtobufMessage(&pbGlobal.S2C_AccountLogon{}, handleWsServer)

	time.Sleep(time.Millisecond * 500)
	sockClient, err := trWsCli.Dial("wss://localhost:4433")
	if err != nil {
		log.Fatalf("unexpected web socket dial err: %v", err)
	}

	msg := &Message{
		// Type: BodyProtobuf,
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
				log.Fatalf("Unexpected recv err: %v", err)
			} else {
				h.Fn(ctxCli, sockClient, msg)
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
