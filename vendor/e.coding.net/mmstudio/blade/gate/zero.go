package gate

import (
	"net"

	"e.coding.net/mmstudio/blade/golib/paniccatcher"
	"e.coding.net/mmstudio/blade/golib/zerocopy"
	"github.com/rs/zerolog/log"
)

func pipeZeroCopy(dstConn, srcConn net.Conn) {
	paniccatcher.Do(func() {
		defer func() { _ = dstConn.Close() }()
		_, _ = zerocopy.Transfer(dstConn, srcConn)
	}, func(p *paniccatcher.Panic) {
		log.Warn().Interface("reason", p.Reason).Msg("pipeZeroCopy panic")
	})
}
