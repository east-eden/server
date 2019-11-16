package client

import "time"

type Options struct {
	ConfigFile string        `flag:"config_file"`
	ClientID   int           `flag:"client_id"`
	ClientName string        `flag:"client_name"`
	HeartBeat  time.Duration `flag:"heart_beat"`

	TcpServerAddr string `flag:"tcp_server_addr"`
}

func NewOptions() *Options {
	return &Options{
		ConfigFile:    "../../config/client/config.toml",
		ClientID:      8001,
		ClientName:    "Anonymous",
		HeartBeat:     5 * time.Second,
		TcpServerAddr: "127.0.0.1:7030",
	}
}
