package gate

import (
	"encoding/binary"
)

// TODO gate不需要解析所有协议
func onBytesInOut(byteOrder binary.ByteOrder, tag string, onFrame func([]byte, string)) func(buf []byte) {
	var localBytes []byte
	var localTag = tag
	var localByteOrder = byteOrder
	return func(buf []byte) {
		if onFrame == nil {
			return
		}
		localBytes = append(localBytes, buf...)
		for {
			if len(localBytes) <= 4 {
				return
			}
			lenBody := int(localByteOrder.Uint32(localBytes[:4]))
			if len(localBytes) < lenBody {
				return
			}
			localBytes = localBytes[4:]
			body := localBytes[:lenBody]
			localBytes = localBytes[lenBody:]
			onFrame(body, localTag)
		}
	}
}
