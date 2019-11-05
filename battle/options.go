package battle

type Options struct {
	ConfigFile string `flag:"config_file"`
	BattleID   int    `flag:"battle_id"`
	MysqlDSN   string `flag:"mysql_dsn"`

	HTTPListenAddr string `flag:"http_listen_addr"`

	MicroRegistry  string `flag:"micro_registry"`
	MicroTransport string `flag:"micro_transport"`
	MicroBroker    string `flag:"micro_broker"`
}

func NewOptions() *Options {
	return &Options{
		ConfigFile: "",
		BattleID:   2001,
		MysqlDSN:   "root:@(127.0.0.1:3306)/db_battle",

		HTTPListenAddr: ":8081",

		MicroRegistry:  "mdns",
		MicroTransport: "http",
		MicroBroker:    "http",
	}
}
