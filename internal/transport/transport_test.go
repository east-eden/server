package transport

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	pbClient "github.com/yokaiio/yokai_server/proto/client"
)

func init() {

}

func TestMarshal(t *testing.T) {

	var msg Message
	msg.Type = BodyProtobuf
	msg.Name = "yokai_client.MC_ClientLogon"

	protoMsg := &pbClient.MC_ClientLogon{
		ClientId:   1002,
		ClientName: "dudu",
	}

	data, err := proto.Marshal(protoMsg)
	if err != nil {
		t.Error("proto marshal error:", err)
	}

	fmt.Println("marshal success")

	var newMsg proto.Message
	if err := proto.Unmarshal(data, newMsg); err != nil {
		t.Error("proto unmarshal error:", err)
	}

	fmt.Println("unmarshal success")

	retMsg, ok := newMsg.(*pbClient.MC_ClientLogon)
	if !ok {
		t.Error("proto assert error")
	}

	fmt.Println("assert success:", retMsg)
}