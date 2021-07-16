package tcpkeepalive

import (
	"fmt"
	"net"
	"time"
)

func SetKeepAlive(conn net.Conn, idleTime time.Duration, count int, interval time.Duration) (err error) {
	c, ok := conn.(*net.TCPConn)
	if !ok {
		return fmt.Errorf("Bad connection type: %T", c)
	}

	if err := c.SetKeepAlivePeriod(idleTime); err != nil {
		return err
	}

	if err := c.SetKeepAlive(true); err != nil {
		return err
	}

	return nil
}
