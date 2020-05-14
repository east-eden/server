package game

import (
	"hash/crc32"
	"testing"
)

var MsgName = "C2M_AccountLogon"

func TestMsgHandler(t *testing.T) {
	handler := NewMsgHandler(nil)
	h, err := handler.r.GetHandler(crc32.ChecksumIEEE([]byte(MsgName)))
	if err != nil {
		t.Errorf("MsgHandler get handler failed:%v", err)
	}

	if h.Name != MsgName {
		t.Errorf("MsgHandler name invalid:%s", h.Name)
	}
}
