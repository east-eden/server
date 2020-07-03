package transport

import (
	"context"
	"crypto/tls"
	"log"
	"sync"
	"testing"

	pbAccount "github.com/yokaiio/yokai_server/proto/account"
	"github.com/yokaiio/yokai_server/transport/codec"
)

var (
	trWsSrv  = NewTransport("ws")
	regWsSrv = NewTransportRegister()
	trWsCli  = NewTransport("ws")
	regWsCli = NewTransportRegister()
	wgWs     sync.WaitGroup
)

func handleWsServerSocket(ctx context.Context, sock Socket, closeHandler SocketCloseHandler) {
	defer func() {
		sock.Close()
		closeHandler()
	}()

	for {
		select {
		case <-ctx.Done():
			break
		default:
		}

		msg, h, err := sock.Recv(regWsSrv)
		if err != nil {
			log.Printf("ws server handle socket error: %w", err)
			return
		}

		h.Fn(ctx, sock, msg)
	}
}

func handleWsClient(ctx context.Context, sock Socket, p *Message) {
	msg, ok := p.Body.(*pbAccount.C2M_AccountLogon)
	if !ok {
		log.Fatalf("handleClient failed")
	}

	var sendMsg Message
	sendMsg.Type = BodyProtobuf
	sendMsg.Name = "M2C_AccountLogon"
	sendMsg.Body = &pbAccount.M2C_AccountLogon{
		RpcId:      2,
		Error:      0,
		Message:    "OK",
		PlayerName: msg.AccountName,
	}

	sock.Send(&sendMsg)
	wgWs.Done()
}

func handleWsServer(ctx context.Context, sock Socket, p *Message) {
	msg, ok := p.Body.(*pbAccount.M2C_AccountLogon)
	if !ok {
		log.Fatalf("handleServer failed")
	}

	if msg.PlayerName != "test_name" {
		log.Fatalf("handleServer failed")
	}

	wgWs.Done()
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
		Timeout(DefaultDialTimeout),
		Codec(&codec.ProtoBufMarshaler{}),
		TLSConfig(tlsConfServ),
	)

	regWsSrv.RegisterProtobufMessage(&pbAccount.C2M_AccountLogon{}, handleWsClient)

	wgWs.Add(1)
	go func() {
		err := trWsSrv.ListenAndServe(context.TODO(), ":443", handleWsServerSocket)
		if err != nil {
			log.Fatalf("WsServer ListenAndServe failed:%w", err)
		}
	}()

	// ws client
	tlsConfCli := &tls.Config{InsecureSkipVerify: true}
	tlsConfCli.Certificates = []tls.Certificate{cert}
	trWsCli.Init(
		Timeout(DefaultDialTimeout),
		TLSConfig(tlsConfCli),
	)

	regWsCli.RegisterProtobufMessage(&pbAccount.M2C_AccountLogon{}, handleWsServer)

	sockClient, err := trWsCli.Dial("wss://localhost:443")
	if err != nil {
		log.Fatalf("unexpected web socket dial err:%w", err)
	}

	wgWs.Add(1)
	msg := &Message{
		Type: BodyProtobuf,
		Name: "yokai_account.C2M_AccountLogon",
		Body: &pbAccount.C2M_AccountLogon{
			RpcId:       1,
			UserId:      1,
			AccountId:   1,
			AccountName: "test_name",
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer func() {
			sockClient.Close()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			if msg, h, err := sockClient.Recv(regWsCli); err != nil {
				log.Fatalf("Unexpected recv err:%w", err)
			} else {
				h.Fn(ctx, sockClient, msg)
			}
		}
	}()

	if err := sockClient.Send(msg); err != nil {
		log.Fatalf("client send socket failed:%w", err)
	}

	wgWs.Wait()
	cancel()
}
