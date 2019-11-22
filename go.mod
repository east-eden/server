module github.com/yokaiio/yokai_server

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/gammazero/workerpool v0.0.0-20191005194639-971bc780f6d7
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.2
	github.com/jinzhu/gorm v1.9.11
	github.com/judwhite/go-svc v1.1.2
	github.com/micro/cli v0.2.0
	github.com/micro/go-micro v1.15.1
	github.com/micro/go-plugins v1.4.0
	github.com/miekg/dns v1.1.22 // indirect
	github.com/mreiferson/go-options v0.0.0-20190302064952-20ba7d382d05
	github.com/nats-io/nats.go v1.8.2-0.20190607221125-9f4d16fe7c2d // indirect
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/net v0.0.0-20191028085509-fe3aa8a45271 // indirect
)

replace github.com/gogo/protobuf v0.0.0-20190410021324-65acae22fc9 => github.com/gogo/protobuf v0.0.0-20190723190241-65acae22fc9d
