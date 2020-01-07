package define

import "time"

// base net message type define
type BaseNetMsg struct {
	ID   uint32 // message name crc32
	Size uint32 // message size
}

// transfer message type
type TransferNetMsg struct {
	BaseNetMsg
	WorldID  uint32 // world to recv message
	PlayerID int64  // player to recv message
}

// tcp read max buff size
var TCPReadBufMax = 1024 * 1024 * 2

// micro service ttl
var MicroServiceTTL = time.Second * 30

// micro service internal
var MicroServiceInternal = time.Second * 15
