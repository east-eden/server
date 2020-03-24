package transport

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	pbAccount "github.com/yokaiio/yokai_server/proto/account"
)

func init() {

}

func TestMarshal(t *testing.T) {

	var msg Message
	msg.Type = BodyProtobuf
	msg.Name = "yokai_account.C2M_AccountLogon"

	protoMsg := &pbAccount.C2M_AccountLogon{
		UserId:    1002,
		AccountId: 1002,
	}

	data, err := proto.Marshal(protoMsg)
	if err != nil {
		t.Error("proto marshal error:", err)
	}

	fmt.Println("marshal success")

	newMsg, ok := reflect.New(reflect.TypeOf(protoMsg).Elem()).Interface().(proto.Message)
	if !ok {
		t.Error("protobuf new elem interface failed")
	}

	if err := proto.Unmarshal(data, newMsg); err != nil {
		t.Error("proto unmarshal error:", err)
	}

	fmt.Println("unmarshal success")

	retMsg, ok := newMsg.(*pbAccount.C2M_AccountLogon)
	if !ok {
		t.Error("proto assert error")
	}

	fmt.Println("assert success:", retMsg)
}
