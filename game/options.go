package game

import (
	"time"
)

type Options struct {
	ConfigFile       string        `flag:"config_file"`
	GameID           int           `flag:"game_id"`
	ClientConnectMax int           `flag:"client_connect_max"`
	ClientTimeOut    time.Duration `flag:"client_timeout"`
	HeartBeat        time.Duration `flag:"heart_beat"`
	MysqlDSN         string        `flag:"mysql_dsn"`

	HTTPListenAddr string `flag:"http_listen_addr"`
	TCPListenAddr  string `flag:"tcp_listen_addr"`

	MicroRegistry  string `flag:"micro_registry"`
	MicroTransport string `flag:"micro_transport"`
	MicroBroker    string `flag:"micro_broker"`
}

func NewOptions() *Options {
	return &Options{
		ConfigFile:       "",
		GameID:           1001,
		ClientConnectMax: 5000,
		ClientTimeOut:    30 * time.Second,
		HeartBeat:        10 * time.Second,
		MysqlDSN:         "root:@(127.0.0.1:3306)/db_game",

		HTTPListenAddr: ":8080",
		TCPListenAddr:  ":7030",

		MicroRegistry:  "mdns",
		MicroTransport: "http",
		MicroBroker:    "http",
	}
}
