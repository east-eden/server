module github.com/yokaiio/yokai_server

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/gammazero/workerpool v0.0.0-20191005194639-971bc780f6d7
	github.com/go-delve/delve v1.3.2
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.2
	github.com/grafana/grafana v6.1.6+incompatible
	github.com/jinzhu/gorm v1.9.11
	github.com/judwhite/go-svc v1.1.2
	github.com/micro/cli v0.2.0
	github.com/micro/examples v0.2.0 // indirect
	github.com/micro/go-micro v1.15.1
	github.com/micro/go-plugins v1.4.0
	github.com/micro/micro v1.15.1 // indirect
	github.com/mreiferson/go-options v0.0.0-20190302064952-20ba7d382d05
	github.com/prometheus/client_golang v1.1.0
	github.com/sirupsen/logrus v1.4.2
)

replace github.com/gogo/protobuf v0.0.0-20190410021324-65acae22fc9 => github.com/gogo/protobuf v0.0.0-20190723190241-65acae22fc9d
