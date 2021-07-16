package link

import (
	"net"
	"regexp"
	"strings"

	"e.coding.net/mmstudio/blade/golib/reuseport"
)

func init() {
	RegisterMakeListener(ReusePort, reuseportNewListener)
}

func reuseportNewListener(s *Server, address string) (ln net.Listener, err error) {
	var network string
	if validIP4(address) {
		network = "tcp4"
	} else {
		network = "tcp6"
	}

	return reuseport.NewReusablePortListener(network, address)
}

var ip4Reg = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

func validIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")
	i := strings.LastIndex(ipAddress, ":")
	ipAddress = ipAddress[:i] //remove port

	return ip4Reg.MatchString(ipAddress)
}
