package prom

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	OpsLogonAccountCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logon_account_total",
		Help: "登录账号总数",
	})

	OpsOnlineAccountGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "online_account_gauge",
		Help: "在线账号数量",
	})
)
