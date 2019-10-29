package game

import (
	"time"
)

type Options struct {
	GameID           uint32        `flag:"game_id"`
	ClientConnectMax int32         `flag:"client_connect_max"`
	ClientTimeOut    time.Duration `flag:"client_timeout"`
	HeartBeat        time.Duration `flag:"heart_beat"`
	MysqlDSN         string        `flag:"mysql_dsn"`
}

func NewOptions() *Options {
	return &Options{
		GameID:           1001,
		ClientConnectMax: 5000,
		ClientTimeOut:    30 * time.Second,
		HeartBeat:        10 * time.second,
		MysqlDSN:         "root:@(127.0.0.1:3306)/db_game",
	}
}
