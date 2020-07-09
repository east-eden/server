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

	OpsCreateHeroCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "create_hero_counter",
		Help: "创建英雄总数",
	})

	OpsCreateItemCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "create_item_counter",
		Help: "创建物品总数",
	})
)
