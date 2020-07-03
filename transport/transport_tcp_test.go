package transport

import (
	"context"
	"log"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
	"github.com/yokaiio/yokai_server/transport/codec"
)

var (
	trTcpSrv  = NewTransport("tcp")
	regTcpSrv = NewTransportRegister()
	trTcpCli  = NewTransport("tcp")
	regTcpCli = NewTransportRegister()
	wgTcp     sync.WaitGroup
)

func handleTcpServerSocket(ctx context.Context, sock Socket, closeHandler SocketCloseHandler) {
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

		msg, h, err := sock.Recv(regTcpSrv)
		if err != nil {
			log.Printf("tcp server handle socket error: %w", err)
			return
		}

		h.Fn(ctx, sock, msg)
	}
}

func handleTcpClientAccountLogon(ctx context.Context, sock Socket, p *Message) {
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
	wgTcp.Done()
}

func handleTcpClientAccountTest(ctx context.Context, sock Socket, p *Message) {
	msg, ok := p.Body.(*C2M_AccountTest)
	if !ok {
		log.Fatalf("handleClient failed")
	}

	var sendMsg Message
	sendMsg.Type = BodyJson
	sendMsg.Name = "M2C_AccountTest"
	sendMsg.Body = &M2C_AccountTest{
		AccountId: msg.AccountId,
		Name:      msg.Name,
	}

	sock.Send(&sendMsg)
	wgTcp.Done()
}

func handleTcpServerAccountLogon(ctx context.Context, sock Socket, p *Message) {
	msg, ok := p.Body.(*pbAccount.M2C_AccountLogon)
	if !ok {
		log.Fatalf("handleServer failed")
	}

	diff := cmp.Diff(msg.PlayerName, "test_name")
	if diff != "" {
		log.Fatalf("handleTcpServerAccountLogon failed: %s", diff)
	}

	wgTcp.Done()
}

func handleTcpServerAccountTest(ctx context.Context, sock Socket, p *Message) {
	msg, ok := p.Body.(*M2C_AccountTest)
	if !ok {
		log.Fatalf("handleServer json failed")
	}

	diff := cmp.Diff(msg.Name, "test_json")
	if diff != "" {
		log.Fatalf("handleTcpServerAccountTest failed: %s", diff)
	}

	wgTcp.Done()
}

// json message define
type C2M_AccountTest struct {
	AccountId int64  `protobuf:"varint,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (msg *C2M_AccountTest) GetName() string {
	return "C2M_AccountTest"
}

type M2C_AccountTest struct {
	AccountId int64  `protobuf:"varint,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (msg *M2C_AccountTest) GetName() string {
	return "M2C_AccountTest"
}

func TestTransportTcp(t *testing.T) {

	// tcp server
	trTcpSrv.Init(
		Timeout(DefaultDialTimeout),
		Codec(&codec.ProtoBufMarshaler{}),
	)

	regTcpSrv.RegisterProtobufMessage(&pbAccount.C2M_AccountLogon{}, handleTcpClientAccountLogon)
	regTcpSrv.RegisterJsonMessage(&C2M_AccountTest{}, handleTcpClientAccountTest)

	wgTcp.Add(2)
	go func() {
		err := trTcpSrv.ListenAndServe(context.TODO(), ":7030", handleTcpServerSocket)
		if err != nil {
			log.Fatal("TcpServer ListenAndServe failed:%w", err)
		}
	}()

	// tcp client
	trTcpCli.Init(
		Timeout(DefaultDialTimeout),
	)

	regTcpCli.RegisterProtobufMessage(&pbAccount.M2C_AccountLogon{}, handleTcpServerAccountLogon)
	regTcpCli.RegisterJsonMessage(&M2C_AccountTest{}, handleTcpServerAccountTest)

	sockClient, err := trTcpCli.Dial("127.0.0.1:7030")
	if err != nil {
		log.Fatalf("unexpected tcp dial err:%w", err)
	}

	wgTcp.Add(2)
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

			if msg, h, err := sockClient.Recv(regTcpCli); err != nil {
				log.Fatalf("Unexpected recv err:%w", err)
			} else {
				h.Fn(ctx, sockClient, msg)
			}
		}
	}()

	// send protobuf message
	msgProtobuf := &Message{
		Type: BodyProtobuf,
		Name: "yokai_account.C2M_AccountLogon",
		Body: &pbAccount.C2M_AccountLogon{
			RpcId:       1,
			UserId:      1,
			AccountId:   1,
			AccountName: "test_name",
		},
	}

	if err := sockClient.Send(msgProtobuf); err != nil {
		log.Fatalf("client send socket failed:%w", err)
	}

	// send json message
	msgJson := &Message{
		Type: BodyJson,
		Name: "C2M_AccountTest",
		Body: &C2M_AccountTest{
			AccountId: 1,
			Name:      "test_json",
		},
	}

	if err := sockClient.Send(msgJson); err != nil {
		log.Fatalf("client send socket failed:%w", err)
	}

	wgTcp.Wait()
	cancel()
}
