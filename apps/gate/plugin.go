package main

import (
	_ "github.com/micro/go-plugins/broker/nsq"
	_ "github.com/micro/go-plugins/registry/consul"
	_ "github.com/micro/go-plugins/store/consul"
	_ "github.com/micro/go-plugins/transport/tcp"
)
