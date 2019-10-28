package game

import (
	"time"
)

type Options struct {
	ClientConnectMax int32         `flag:"client_connect_max"`
	ClientTimeOut    time.Duration `flag:"client_timeout"`
	HeartBeat        time.Duration `flag:"heart_beat"`
}

func NewOptions() *Options {
	return &Options{
		ClientConnectMax: 5000,
		ClientTimeOut:    30 * time.Second,
		HeartBeat:        10 * time.second,
	}
}
