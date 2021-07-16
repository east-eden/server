package gate

import (
	"encoding/binary"
	"errors"
	"io"
	"time"

	"e.coding.net/mmstudio/blade/gate/selector"
	"e.coding.net/mmstudio/blade/golib/encoding2"
	"e.coding.net/mmstudio/blade/golib/encoding2/protobuf"
	"e.coding.net/mmstudio/blade/golib/net/link"
)

const (
	// max open file should at least be
	maxOpenfile   = uint64(1024 * 1024 * 1024)
	DefaultSecret = "A40h%vXGjT6XcL9jmnOv&DgEJoTA0M87"
)

//go:generate optiongen
func _SpecOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"ConnReadTimeout":         time.Duration(time.Duration(30) * time.Second),
		"ConnWriteTimeout":        time.Duration(time.Duration(30) * time.Second),
		"DispatcherDialTimeout":   time.Duration(time.Duration(30) * time.Second),
		"StableConnection":        false,
		"PortForClient":           8989,
		"PortForClientReuse":      true,
		"EnableWebSocket":         false,
		"PortForWebClient":        8990,
		"PathForWebClient":        "/gate",
		"EnableZeroCopy":          true,
		"EnableXListener":         true,
		"XListenerBacklogAccept":  10240,
		"XListenerTimeoutCanRead": time.Duration(time.Duration(30) * time.Second),
		"XListenerHandshake":      false,
		"OnFrameForLog":           func([]byte, string) {},
		"ByteOrderForLog":         binary.ByteOrder(binary.LittleEndian),
		"Selector":                selector.Selector(nil),
		"BackEndHandshake":        false,
		"TransferProvider": (func(io.ReadWriter) (link.Transporter, error))(func(_ io.ReadWriter) (link.Transporter, error) {
			return nil, errors.New("no transfer")
		}),
		"MessageCodec":       encoding2.Codec(protobuf.Codec),
		"PatchPath":          "./patch",
		"SelectorStrategy":   selector.Strategy(selector.RoundRobin),
		"Filter":             []FilterPlugin{},
		"HandshakeSecret":    false,                 // 是否开启握手验证签名
		"HandshakeCheckTime": false,                 // 是否验证时间戳
		"SecretKey":          string(DefaultSecret), // 握手验签secret
	}
}
